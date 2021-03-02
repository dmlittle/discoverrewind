package logger

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type key int

const (
	levelHeader       = "log-level"
	echoLoggerKey     = "logger"
	echoIDKey         = "id"
	ctxKey        key = 0
)

// Middleware attaches a Logger instance with a request ID onto the context. It
// also logs every request along with metadata about the request.
func Middleware() func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			level := c.Request().Header.Get(levelHeader)
			l := NewWithLevel(level)
			t1 := time.Now()
			id, err := uuid.NewRandom()
			if err != nil {
				return errors.WithStack(err)
			}
			idStr := id.String()
			c.Set(echoIDKey, idStr)
			log := l.ID(idStr)
			c.Set(echoLoggerKey, log)
			if err := next(c); err != nil {
				c.Error(err)
			}
			t2 := time.Now()

			// If request was a successful health check do not log it.
			// These logs are too noisy.
			if c.Request().Method == "GET" && c.Request().URL.Path == "/health" && c.Response().Status == http.StatusOK {
				return nil
			}

			log.Root(Data{
				"status_code": c.Response().Status,
				"method":      c.Request().Method,
				"path":        c.Request().URL.Path,
				"route":       c.Path(),
				"duration":    t2.Sub(t1).Seconds() * 1000,
				"referer":     c.Request().Referer(),
				"user_agent":  c.Request().UserAgent(),
			}).Info("request handled")
			return nil
		}
	}
}

// IDFromEchoContext returns the request ID from the given echo.Context. If
// there is no request ID, then this will just return the empty string.
func IDFromEchoContext(c echo.Context) string {
	id, ok := c.Get(echoIDKey).(string)
	if !ok {
		return ""
	}
	return id
}

// FromEchoContext returns a Logger from the given echo.Context. If there is no
// attached logger, then this will just return a new Logger instance.
func FromEchoContext(c echo.Context) Logger {
	var log Logger
	log, ok := c.Get(echoLoggerKey).(Logger)
	if !ok {
		log = New()
	}
	return log
}

// FromContext returns a Logger from the given context.Context. If there is no
// attached logger, then this will just return a new Logger instance.
func FromContext(ctx context.Context) Logger {
	var log Logger
	log, ok := ctx.Value(ctxKey).(Logger)
	if !ok {
		log = New()
	}
	return log
}
