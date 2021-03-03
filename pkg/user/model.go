package user

import (
	"context"
	"github.com/dmlittle/discoverrewind/pkg/crypto"
	"github.com/go-pg/pg/v10/orm"
	"github.com/pkg/errors"
	"time"
)

type User struct {
	ID                       string
	SpotifyID                string
	DisplayName              string
	ProfileImageURL          string
	Product                  string
	TokenType                string
	AccessToken              string
	RefreshToken             string
	TokenExpiration          time.Time
	DiscoverWeeklyPlaylistID string
	CreatedAt                *time.Time
	UpdatedAt                *time.Time
	DeletedAt                *time.Time
}

var _ orm.BeforeInsertHook = (*User)(nil)

func (u *User) BeforeInsert(ctx context.Context) (context.Context, error) {
	cryptoSvc, ok := crypto.FromContext(ctx)
	if !ok {
		panic("crypto context has not been set")
	}

	encryptedAccessToken, err := cryptoSvc.Encrypt(u.AccessToken)
	if err != nil {
		return ctx, errors.WithStack(err)
	}
	u.AccessToken = encryptedAccessToken

	encryptedRefreshToken, err := cryptoSvc.Encrypt(u.RefreshToken)
	if err != nil {
		return ctx, errors.WithStack(err)
	}
	u.RefreshToken = encryptedRefreshToken

	return ctx, nil
}

var _ orm.AfterSelectHook = (*User)(nil)

func (u *User) AfterSelect(ctx context.Context) error {
	cryptoSvc, ok := crypto.FromContext(ctx)
	if !ok {
		panic("crypto context has not been set")
	}

	decryptedAccessToken, err := cryptoSvc.Decrypt(u.AccessToken)
	if err != nil {
		return errors.WithStack(err)
	}
	u.AccessToken = decryptedAccessToken

	decryptedRefreshToken, err := cryptoSvc.Decrypt(u.RefreshToken)
	if err != nil {
		return errors.WithStack(err)
	}
	u.RefreshToken = decryptedRefreshToken

	return nil
}
