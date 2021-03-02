package main

import (
	"os"

	"github.com/dmlittle/discoverrewind/pkg/config"
	"github.com/dmlittle/discoverrewind/pkg/database"
	"github.com/dmlittle/discoverrewind/pkg/logger"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

const directory = "./cmd/migrations"

func main() {
	log := logger.New()

	cfg, err := config.New()
	if err != nil {
		log.Err(err).Fatal("config error")
	}
	db, err := database.New(cfg)
	if err != nil {
		log.Err(err).Fatal("database error")
	}

	err = migrations.Run(db, directory, os.Args)
	if err != nil {
		log.Err(err).Fatal("migration error")
	}
}
