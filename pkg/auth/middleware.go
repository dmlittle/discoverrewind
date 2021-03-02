package auth

import (
	"github.com/dmlittle/discoverrewind/pkg/user"
	"github.com/go-pg/pg/v10"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
)

// Middleware returns a Session middleware.
func Middleware(db *pg.DB) echo.MiddlewareFunc {
	userSvc := user.New(db)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		type session struct {
			ID     string `pg:",pk"`
			UserID string
		}

		return func(c echo.Context) error {
			ctx := c.Request().Context()

			cookie, err := c.Cookie(cookieKey)
			if err != nil {
				return c.Redirect(http.StatusFound, "/")
			}

			sess := &session{}
			err = db.ModelContext(ctx, sess).Where("id = ?", cookie.Value).First()
			if err != nil {
				if err == pg.ErrNoRows {
					return c.Redirect(http.StatusFound, "/")
				}
				return errors.WithStack(err)
			}

			u, err := userSvc.FetchUser(ctx, user.FetchUserInput{ID: sess.UserID})
			if err != nil {
				return errors.WithStack(err)
			}

			d := &Details{UserID: sess.UserID, User: u}

			return next(d.WithEchoContext(c))
		}
	}

}
