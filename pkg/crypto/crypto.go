package crypto

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"github.com/dmlittle/discoverrewind/pkg/config"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"io"
)

// The Encrypt and Decrypt functions have been adapted from
// https://github.com/gtank/cryptopasta/blob/1f550f6f2f69009f6ae57347c188e0a67cd4e500/encrypt.go

type key int

const ctxKey key = 0

func Middleware(cfg *config.Config) func(echo.HandlerFunc) echo.HandlerFunc {
	svc := &service{&cfg.EncryptionKey}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(svc.WithEchoContext(c))
		}
	}
}

type service struct {
	key *[32]byte
}

// Encrypt encrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Output takes the
// form nonce|ciphertext|tag where '|' indicates concatenation.
func (s *service) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(s.key[:])
	if err != nil {
		return "", errors.WithStack(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.WithStack(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return hex.EncodeToString(gcm.Seal(nonce, nonce, []byte(plaintext), nil)), nil
}

// Decrypt decrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Expects input
// form nonce|ciphertext|tag where '|' indicates concatenation.
func (s *service) Decrypt(ciphertext string) (string, error) {
	block, err := aes.NewCipher(s.key[:])
	if err != nil {
		return "", errors.WithStack(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.WithStack(err)
	}

	if len(ciphertext) < gcm.NonceSize() {
		return "", errors.WithStack(errors.New("malformed ciphertext"))
	}

	cypherData, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", errors.WithStack(err)
	}

	plainData, err := gcm.Open(nil,
		cypherData[:gcm.NonceSize()],
		cypherData[gcm.NonceSize():],
		nil,
	)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return string(plainData), nil
}

func (s *service) WithEchoContext(c echo.Context) echo.Context {
	ctx := s.WithContext(c.Request().Context())
	c.SetRequest(c.Request().WithContext(ctx))

	return c
}

func (s *service) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey, s)
}

// FromContext returns a Logger from the given context.Context. If there is no
// attached logger, then this will just return a new Logger instance.
func FromContext(ctx context.Context) (*service, bool) {
	var svc *service
	svc, ok := ctx.Value(ctxKey).(*service)
	if !ok {
		return nil, false
	}
	return svc, true
}
