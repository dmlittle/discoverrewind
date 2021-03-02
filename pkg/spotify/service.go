package spotify

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/zmb3/spotify"
	"time"
)

type Service struct {
	db *pg.DB
}

func New(db *pg.DB) *Service {
	return &Service{db}
}

func (s *Service) CreatePlaylistSnapshot(ctx context.Context, playlist *spotify.FullPlaylist) error {
	playlistSnapshotID, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	createdAt, err := time.Parse(spotify.TimestampLayout, playlist.Tracks.Tracks[0].AddedAt)
	if err != nil {
		return err
	}

	playlistSnapshot := &PlaylistSnapshot{
		ID:         playlistSnapshotID.String(),
		PlaylistID: playlist.ID.String(),
		SnapshotID: playlist.SnapshotID,
		CreatedAt:  &createdAt,
	}

	if len(playlist.Images) > 0 {
		playlistSnapshot.ImageURL = playlist.Images[0].URL
	}

	var artists []*Artist
	var albums []*Album
	var albumImages []*AlbumImage
	var tracks []*Track
	var trackArtists []*TrackArtists
	var playlistSnapshotDetails []*PlaylistSnapshotDetail

	for i, t := range playlist.Tracks.Tracks {
		playlistSnapshotDetails = append(playlistSnapshotDetails, &PlaylistSnapshotDetail{
			PlaylistSnapshotID: playlistSnapshotID.String(),
			Rank:               i,
			TrackID:            t.Track.ID.String(),
		})

		tracks = append(tracks, &Track{
			ID:          t.Track.ID.String(),
			Name:        t.Track.Name,
			AlbumID:     t.Track.Album.ID.String(),
			DiscNumber:  t.Track.DiscNumber,
			TrackNumber: t.Track.TrackNumber,
			Explicit:    t.Track.Explicit,
			DurationMS:  t.Track.Duration,
		})

		for _, artist := range t.Track.Artists {
			artists = append(artists, &Artist{
				ID:   artist.ID.String(),
				Name: artist.Name,
			})

			trackArtists = append(trackArtists, &TrackArtists{
				TrackID:  t.Track.ID.String(),
				ArtistID: artist.ID.String(),
			})
		}

		albums = append(albums, &Album{
			ID:   t.Track.Album.ID.String(),
			Name: t.Track.Album.Name,
		})

		for _, img := range t.Track.Album.Images {
			albumImages = append(albumImages, &AlbumImage{
				AlbumID: t.Track.Album.ID.String(),
				Height:  img.Height,
				Width:   img.Width,
				URL:     img.URL,
			})
		}
	}

	err = s.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		_, err := tx.ModelContext(ctx, &artists).OnConflict("DO NOTHING").Insert()
		if err != nil {
			return errors.WithStack(err)
		}

		_, err = tx.ModelContext(ctx, &albums).OnConflict("DO NOTHING").Insert()
		if err != nil {
			return errors.WithStack(err)
		}

		_, err = tx.ModelContext(ctx, &albumImages).OnConflict("DO NOTHING").Insert()
		if err != nil {
			return errors.WithStack(err)
		}

		_, err = tx.ModelContext(ctx, &tracks).OnConflict("DO NOTHING").Insert()
		if err != nil {
			return errors.WithStack(err)
		}

		_, err = tx.ModelContext(ctx, &trackArtists).OnConflict("DO NOTHING").Insert()
		if err != nil {
			return errors.WithStack(err)
		}

		_, err = tx.ModelContext(ctx, playlistSnapshot).OnConflict("DO NOTHING").Insert()
		if err != nil {
			return errors.WithStack(err)
		}

		_, err = tx.ModelContext(ctx, &playlistSnapshotDetails).OnConflict("DO NOTHING").Insert()
		if err != nil {
			return errors.WithStack(err)
		}

		return nil
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Service) FetchPlaylistSnapshots(ctx context.Context, playlistID string) ([]*PlaylistSnapshot, error) {
	snapshots := []*PlaylistSnapshot{}

	err := s.db.ModelContext(ctx, &snapshots).
		Where("playlist_id = ?", playlistID).
		Order("created_at DESC").
		Select()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return snapshots, nil
}

func (s *Service) FetchPlaylistSnapshotDetails(ctx context.Context, playlistSnapshotID string) ([]*PlaylistSnapshotDetail, error) {
	details := []*PlaylistSnapshotDetail{}

	err := s.db.ModelContext(ctx, &details).
		Where("playlist_snapshot_id = ?", playlistSnapshotID).
		Relation("Track", func(q *orm.Query) (*orm.Query, error) {
			return q.
				Relation("Track.Artists").
				Relation("Track.Album", func(q *orm.Query) (*orm.Query, error) {
					return q.Relation("Track.Album.Images"), nil
				}), nil

		}).
		OrderExpr("rank ASC").
		Select()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return details, nil
}
