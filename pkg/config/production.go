package config

import (
	"github.com/labstack/gommon/random"
	"os"
	"strconv"
)

func loadProductionConfig(cfg *Config) {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err == nil {
		cfg.Port = port
	}

	cfg.BaseURL = "https://www.discoverrewind.com"
	cfg.DatabaseHost = os.Getenv("DATABASE_HOST")
	cfg.DatabaseName = os.Getenv("DATABASE_NAME")
	cfg.DatabaseUser = os.Getenv("DATABASE_USER")
	cfg.DatabasePassword = os.Getenv("DATABASE_PASSWORD")
	cfg.Environment = "production"
	cfg.OAuthStateCode = random.String(8, random.Alphanumeric)
	cfg.EncodedEncryptionKey = os.Getenv("ENCODED_ENCRYPTION_KEY")
}
