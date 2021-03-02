package spotify

import (
	"github.com/dmlittle/discoverrewind/pkg/config"
	"github.com/dmlittle/discoverrewind/pkg/session"
	"github.com/dmlittle/discoverrewind/pkg/user"
	"github.com/go-pg/pg/v10"
	"github.com/labstack/echo/v4"
	"github.com/zmb3/spotify"
)

type handler struct {
	authenticator spotify.Authenticator
	state         string

	userSvc    *user.Service
	sessionSvc *session.Service
	spotifySvc *Service
}

// RegisterRoutes takes in an Echo router and registers routes onto it.
func RegisterRoutes(e *echo.Echo, cfg *config.Config, db *pg.DB) {
	h := &handler{
		authenticator: spotify.NewAuthenticator(
			cfg.BaseURL+"/callback",
			spotify.ScopePlaylistReadPrivate,
			spotify.ScopeUserReadEmail,
			spotify.ScopeUserReadPrivate,
			spotify.ScopeStreaming,
		),
		state:      cfg.OAuthStateCode,
		userSvc:    user.New(db),
		sessionSvc: session.NewService(db),
		spotifySvc: New(db),
	}

	e.GET("/login", h.login)

	e.GET("/logout", h.logout)

	e.GET("/callback", h.callbackHandler)
}
