package user

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
)

type Service struct {
	db *pg.DB
}

func New(db *pg.DB) *Service {
	return &Service{db}
}

func (s *Service) CreateUser(ctx context.Context, u *User) error {
	if u.ID == "" {
		id, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		u.ID = id.String()
	}

	_, err := s.db.ModelContext(ctx, u).
		OnConflict("(spotify_id) DO UPDATE").
		Set("access_token = EXCLUDED.access_token, refresh_token = EXCLUDED.refresh_token, updated_at = NOW()").
		Insert()

	return err
}

func (s *Service) UpdateUser(ctx context.Context, u *User, columns []string) error {
	now := time.Now().UTC()
	u.UpdatedAt = &now

	columns = append(columns, "updated_at")

	_, err := s.db.ModelContext(ctx, u).
		Column(columns...).
		WherePK().
		Update()

	return err
}

type FetchUserInput struct {
	ID        string
	SpotifyID string
}

func (s *Service) FetchUser(ctx context.Context, input FetchUserInput) (*User, error) {
	u := &User{}

	q := s.db.ModelContext(ctx, u)

	if input.ID != "" {
		q = q.Where("id = ?", input.ID)
	}

	if input.SpotifyID != "" {
		q = q.Where("spotify_id = ?", input.SpotifyID)
	}

	err := q.First()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return u, nil
}
