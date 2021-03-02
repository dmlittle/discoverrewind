package session

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

type Session struct {
	ID     string `pg:",pk"`
	UserID string
}

type Service struct {
	db *pg.DB
}

func NewService(db *pg.DB) *Service {
	return &Service{db}
}

func (svc *Service) CreateSession(ctx context.Context, sess *Session) error {
	if sess.ID == "" {
		id, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		sess.ID = id.String()
	}

	_, err := svc.db.ModelContext(ctx, sess).Insert()

	return err
}

func (svc *Service) FetchSession(ctx context.Context, ID string) (*Session, error) {
	sess := &Session{}

	err := svc.db.ModelContext(ctx, sess).Where("id = ?", ID).First()

	return sess, err
}
