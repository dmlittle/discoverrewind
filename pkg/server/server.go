package server

import (
	"context"
	"fmt"
	"github.com/dmlittle/discoverrewind/pkg/config"
	"github.com/dmlittle/discoverrewind/pkg/crypto"
	"github.com/dmlittle/discoverrewind/pkg/errcodes"
	"github.com/dmlittle/discoverrewind/pkg/health"
	"github.com/dmlittle/discoverrewind/pkg/logger"
	"github.com/dmlittle/discoverrewind/pkg/pages"
	"github.com/dmlittle/discoverrewind/pkg/recovery"
	"github.com/dmlittle/discoverrewind/pkg/spotify"
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/echoview-v4"
	"github.com/go-pg/pg/v10"
	"github.com/labstack/echo/v4"
	elog "github.com/labstack/gommon/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// New returns a new HTTP server with the registered routes.
func New(cfg *config.Config, db *pg.DB) (*http.Server, error) {
	log := logger.New()
	e := echo.New()

	e.Logger.SetLevel(elog.OFF)

	e.Static("/assets", "assets")

	e.Use(logger.Middleware())
	e.Use(recovery.Middleware())
	e.Use(crypto.Middleware(cfg))

	//Set Renderer
	e.Renderer = echoview.New(goview.Config{
		Root:      "views",
		Extension: ".html",
		Master:    "layouts/master",
		Partials:  []string{},
		Funcs: map[string]interface{}{
			"add": func(a int, b int) int {
				return a + b
			},
			"date": func(t time.Time) string {
				suffix := "th"
				switch t.Day() {
				case 1, 21, 31:
					suffix = "st"
				case 2, 22:
					suffix = "nd"
				case 3, 23:
					suffix = "rd"
				}

				return t.Format("January 2" + suffix + ", 2006")
			},
			"duration": func(durationMS int) string {
				minutes := durationMS / 1000 / 60
				seconds := durationMS/1000 - minutes*60

				return fmt.Sprintf("%d:%02d", minutes, seconds)
			},
			"playlistDuration": func(totalMS int) string {
				seconds := totalMS / 1000

				hours := seconds / 3600                // 3600 seconds = 1 hour
				minutes := (seconds - hours*3600) / 60 // 60 seconds = 1 minute

				return fmt.Sprintf("%dhr %dmin", hours, minutes)
			},
		},
		DisableCache: cfg.Environment == "development",
		Delims:       goview.Delims{Left: "{{", Right: "}}"},
	})

	health.RegisterRoutes(e)
	spotify.RegisterRoutes(e, cfg, db)
	pages.RegisterRoutes(e, db)

	echo.NotFoundHandler = notFoundHandler
	e.HTTPErrorHandler = errcodes.NewHandler().Handle

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: e,
	}

	graceful := signalsSetup()

	go func() {
		<-graceful
		ctx := context.Background()
		err := srv.Shutdown(ctx)
		if err != nil {
			log.Err(err).Error("server shutdown error")
		}
	}()

	return srv, nil
}

func notFoundHandler(c echo.Context) error {
	return c.Render(http.StatusNotFound, "error", goview.M{
		"code":    http.StatusNotFound,
		"message": "Sorry, the page you are looking for could not be found.",
	})
}

// signalsSetup registered for SIGTERM and SIGINT. A stop channel is returned
// which is closed on one of these signals. If a second signal is caught, the program
// is terminated with exit code 1.
func signalsSetup() (stopCh <-chan struct{}) {
	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1)
	}()

	return stop
}
