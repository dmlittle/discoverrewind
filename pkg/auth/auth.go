package auth

import (
	"context"
	"github.com/dmlittle/discoverrewind/pkg/user"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type key int

type Details struct {
	UserID string
	User   *user.User
}

const (
	cookieKey     = "_session_store"
	ctxKey    key = 0
)

func SetSessionCookie(c echo.Context, sid string) {
	setCookieHelper(c, sid)
}

func DeleteSessionCookie(c echo.Context) {
	setCookieHelper(c, "")
}

func setCookieHelper(c echo.Context, value string) {
	cookie := &http.Cookie{
		Name:     cookieKey,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   !strings.HasPrefix(c.Request().Host, "localhost:"),
	}

	c.SetCookie(cookie)
}

func (d *Details) WithEchoContext(c echo.Context) echo.Context {
	ctx := d.WithContext(c.Request().Context())
	c.SetRequest(c.Request().WithContext(ctx))

	return c
}

func (d *Details) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey, d)
}

// FromEchoContext returns an Details from the given echo.Context, if any.
// It fetches the Details on the underlying context.Context.
func FromEchoContext(c echo.Context) (*Details, bool) {
	return FromContext(c.Request().Context())
}

// FromContext returns an Details from the given context.Context, if any.
func FromContext(ctx context.Context) (*Details, bool) {
	d, ok := ctx.Value(ctxKey).(*Details)
	return d, ok
}
