package pages

import (
	"github.com/dmlittle/discoverrewind/pkg/auth"
	"github.com/dmlittle/discoverrewind/pkg/spotify"
	"github.com/dmlittle/discoverrewind/pkg/user"
	"github.com/foolin/goview"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	spotifyLib "github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"net/http"
)

var errDiscoverWeeklyNotFound = errors.New("unable to find discover weekly playlist")

func (h *handler) indexHandler(c echo.Context) error {
	return errors.WithStack(c.Render(http.StatusOK, "index", nil))
}

func (h *handler) setupHandler(c echo.Context) error {
	ctx := c.Request().Context()

	authDetails, ok := auth.FromEchoContext(c)
	if !ok {
		return errors.WithStack(errors.New("unable to fetch authentication details"))
	}

	// If user already has linked their Discover Weekly playlist
	// redirect them to /home
	if authDetails.User.DiscoverWeeklyPlaylistID != "" {
		return errors.WithStack(c.Redirect(http.StatusFound, "/home"))
	}

	spotifyClient := spotifyLib.NewAuthenticator("").NewClient(&oauth2.Token{
		AccessToken:  authDetails.User.AccessToken,
		TokenType:    authDetails.User.TokenType,
		RefreshToken: authDetails.User.RefreshToken,
	})

	dw, err := findDiscoverWeeklyPlaylist(spotifyClient)
	if err == errDiscoverWeeklyNotFound {
		return errors.WithStack(c.Render(http.StatusOK, "setup", nil))
	} else if err != nil {
		return errors.WithStack(err)
	}

	err = h.spotifySvc.CreatePlaylistSnapshot(ctx, dw)
	if err != nil {
		return errors.WithStack(err)
	}

	err = h.userSvc.UpdateUser(ctx, &user.User{ID: authDetails.UserID, DiscoverWeeklyPlaylistID: dw.ID.String()}, []string{"discover_weekly_playlist_id"})
	if err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(c.Redirect(http.StatusFound, "/home"))
}

func (h *handler) homeHandler(c echo.Context) error {
	ctx := c.Request().Context()

	authDetails, ok := auth.FromEchoContext(c)
	if !ok {
		return errors.WithStack(c.Redirect(http.StatusFound, "/"))
	} else if authDetails.User.DiscoverWeeklyPlaylistID == "" {
		return errors.WithStack(c.Redirect(http.StatusFound, "/setup"))
	}

	u, err := h.userSvc.FetchUser(ctx, user.FetchUserInput{ID: authDetails.UserID})
	if err != nil {
		return errors.WithStack(err)
	}

	snapshots, err := h.spotifySvc.FetchPlaylistSnapshots(ctx, u.DiscoverWeeklyPlaylistID)
	if err != nil {
		return errors.WithStack(err)
	}

	querySnapshotID := c.QueryParam("snapshot")
	var currentSnapshot *spotify.PlaylistSnapshot
	querySnapshotNotFound := false

	// Attempt to find the snapshot specified in the query. If we cannot find one we'll
	// default to the latest snapshot.
	if querySnapshotID != "" {
		for _, s := range snapshots {
			if s.ID == querySnapshotID {
				currentSnapshot = s
				break
			}
		}
	}

	if currentSnapshot == nil {
		currentSnapshot = snapshots[0]
	}

	if querySnapshotID != "" && querySnapshotID != currentSnapshot.ID {
		querySnapshotNotFound = true
	}

	currentSnapshotDetails, err := h.spotifySvc.FetchPlaylistSnapshotDetails(ctx, currentSnapshot.ID)
	if err != nil {
		return errors.WithStack(err)
	}

	totalDurationMS := 0
	for _, detail := range currentSnapshotDetails {
		totalDurationMS += detail.Track.DurationMS
	}

	return errors.WithStack(c.Render(http.StatusOK, "home", goview.M{
		"snapshots":             snapshots,
		"currentSnapshot":       currentSnapshot,
		"currentSnapshotTracks": currentSnapshotDetails,
		"querySnapshotNotFound": querySnapshotNotFound,
		"totalDurationMS":       totalDurationMS,
		"displayName":           u.DisplayName,
		"profileImageURL":       u.ProfileImageURL,
	}))
}

func findDiscoverWeeklyPlaylist(c spotifyLib.Client) (*spotifyLib.FullPlaylist, error) {
	limit := 50

	playlists, err := c.CurrentUsersPlaylistsOpt(&spotifyLib.Options{
		Limit: &limit,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for {
		for _, p := range playlists.Playlists {
			if p.Name == "Discover Weekly" && p.Owner.ID == "spotify" {
				return c.GetPlaylist(p.ID)
			}
		}

		err = c.NextPage(playlists)
		if err == spotifyLib.ErrNoMorePages {
			break
		}
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return nil, errDiscoverWeeklyNotFound
}
