package health

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

var resp = []byte(`{"healthy":true}`)

// RegisterRoutes takes in an Echo router and registers routes onto it.
func RegisterRoutes(e *echo.Echo) {
	e.GET("/health", healthHandler)
}

func healthHandler(c echo.Context) error {
	return errors.WithStack(c.JSONBlob(http.StatusOK, resp))
}
