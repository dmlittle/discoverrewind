package spotify

import (
	"github.com/dmlittle/discoverrewind/pkg/auth"
	"github.com/dmlittle/discoverrewind/pkg/errcodes"
	"github.com/dmlittle/discoverrewind/pkg/session"
	"github.com/dmlittle/discoverrewind/pkg/user"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
)

func (h *handler) login(c echo.Context) error {
	authURL := h.authenticator.AuthURL(h.state)

	return errors.WithStack(c.Redirect(http.StatusFound, authURL))
}

func (h *handler) logout(c echo.Context) error {
	auth.DeleteSessionCookie(c)

	return errors.WithStack(c.Redirect(http.StatusFound, "/"))
}

func (h *handler) callbackHandler(c echo.Context) error {
	ctx := c.Request().Context()

	tok, err := h.authenticator.Token(h.state, c.Request())
	if err != nil {
		return errcodes.ValidationError(err.Error())
	}

	spotifyClient := h.authenticator.NewClient(tok)

	spotifyUser, err := spotifyClient.CurrentUser()
	if err != nil {
		return errors.WithStack(err)
	}

	u := &user.User{
		SpotifyID:       spotifyUser.ID,
		DisplayName:     spotifyUser.DisplayName,
		Product:         spotifyUser.Product,
		TokenType:       tok.TokenType,
		AccessToken:     tok.AccessToken,
		RefreshToken:    tok.RefreshToken,
		TokenExpiration: tok.Expiry,
	}

	if len(spotifyUser.Images) > 0 {
		u.ProfileImageURL = spotifyUser.Images[0].URL
	}

	err = h.userSvc.CreateUser(ctx, u)
	if err != nil {
		return errors.WithStack(err)
	}

	// Fetch user again. If a user has already created an account the user.ID will not
	// match the one stored on the database.
	u, err = h.userSvc.FetchUser(ctx, user.FetchUserInput{SpotifyID: spotifyUser.ID})
	if err != nil {
		return errors.WithStack(err)
	}

	sess := &session.Session{UserID: u.ID}
	err = h.sessionSvc.CreateSession(ctx, sess)
	if err != nil {
		return errors.WithStack(err)
	}

	auth.SetSessionCookie(c, sess.ID)

	return errors.WithStack(c.Redirect(http.StatusTemporaryRedirect, "/setup"))
}
