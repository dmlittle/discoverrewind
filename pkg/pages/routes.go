package pages

import (
	"github.com/dmlittle/discoverrewind/pkg/auth"
	"github.com/dmlittle/discoverrewind/pkg/spotify"
	"github.com/dmlittle/discoverrewind/pkg/user"
	"github.com/go-pg/pg/v10"
	"github.com/labstack/echo/v4"
)

type handler struct {
	userSvc    *user.Service
	spotifySvc *spotify.Service
}

// RegisterRoutes takes in an Echo router and registers routes onto it.
func RegisterRoutes(e *echo.Echo, db *pg.DB) {
	h := &handler{
		userSvc:    user.New(db),
		spotifySvc: spotify.New(db),
	}

	e.GET("/", h.indexHandler)

	e.GET("/setup", h.setupHandler, auth.Middleware(db))

	e.GET("/home", h.homeHandler, auth.Middleware(db))

	e.POST("/saveTrack", h.modifyLibraryTracks(true), auth.Middleware(db))

	e.POST("/removeTrack", h.modifyLibraryTracks(false), auth.Middleware(db))

}
