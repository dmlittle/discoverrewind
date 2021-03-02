package config

import (
	"os"
	"strconv"
)

func loadDevelopmentConfig(cfg *Config) {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err == nil {
		cfg.Port = port
	}

	cfg.BaseURL = "http://localhost:5000"
	cfg.DatabaseHost = "postgres"
	cfg.DatabaseName = "discoverrewind"
	cfg.DatabaseSSLMode = "disable"
	cfg.DatabaseUser = "discoverrewind_admin"
	cfg.DatabasePassword = "discoverrewind"
	cfg.Environment = "development"
	cfg.Version = "development"

	// This EncryptionKey is only to be used in development mode. It is hard-coded because
	// the server will panic if one is not set. It is OK for it to be checked into code.
	cfg.EncodedEncryptionKey = "9da920d43b7d18668edebe0a2b11cf23fae90ea48bb00842ad087af2eab3a13a"

}
