package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/dmlittle/discoverrewind/pkg/config"
	"github.com/dmlittle/discoverrewind/pkg/database"
	"github.com/dmlittle/discoverrewind/pkg/logger"
	"github.com/dmlittle/discoverrewind/pkg/server"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	log := logger.New()

	cfg, err := config.New()
	if err != nil {
		log.Err(err).Fatal("config error")
	}
	db, err := database.New(cfg)
	if err != nil {
		log.Err(err).Fatal("database error")
	}

	srv, err := server.New(cfg, db)
	if err != nil {
		log.Err(err).Fatal("server error")
	}

	log.Info("server started", logger.Data{"port": cfg.Port})
	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Err(err).Fatal("server stopped")
	}
	log.Info("server stopped")
}
