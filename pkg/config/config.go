package config

import (
	"encoding/hex"
	"os"
	"time"

	"github.com/pkg/errors"
)

// Config contains the environment specific configuration values needed by the
// application.
type Config struct {
	BaseURL                   string
	DatabaseConnectRetryCount int
	DatabaseConnectRetryDelay time.Duration
	DatabaseDebug             bool
	DatabaseHost              string
	DatabaseName              string
	DatabasePassword          string
	DatabasePort              int
	DatabaseSSLMode           string
	DatabaseUser              string
	Environment               string
	Hostname                  string
	Port                      int
	Version                   string
	OAuthStateCode            string
	EncodedEncryptionKey      string
	EncryptionKey             [32]byte
}

const environmentENV = "ENVIRONMENT"

// New returns an instance of Config based on the "ENVIRONMENT" environment
// variable.
func New() (*Config, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cfg := &Config{
		DatabaseConnectRetryCount: 10,
		DatabaseConnectRetryDelay: 2 * time.Second,
		DatabaseDebug:             os.Getenv("DATABASE_DEBUG") == "true",
		DatabasePort:              5432,
		Hostname:                  hostname,
		Port:                      5000,
		Version:                   os.Getenv("VERSION"),
		OAuthStateCode:            "development",
		EncryptionKey:             [32]byte{},
	}

	switch os.Getenv(environmentENV) {
	case "development", "":
		loadDevelopmentConfig(cfg)
	case "production":
		loadProductionConfig(cfg)
	}

	if cfg.EncodedEncryptionKey != "" {
		key, err := hex.DecodeString(cfg.EncodedEncryptionKey)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		for i := range cfg.EncryptionKey {
			cfg.EncryptionKey[i] = key[i]
		}
	}

	return cfg, nil
}
