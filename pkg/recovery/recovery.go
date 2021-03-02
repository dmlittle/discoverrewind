package recovery

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// Middleware returns recovers from any possible panics in subsequent handlers
// and funnels it to the error handler to be returned as a 500.
func Middleware() func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					var err error
					err, ok := r.(error)
					if !ok {
						err = errors.Errorf("%v", r)
					}
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}
