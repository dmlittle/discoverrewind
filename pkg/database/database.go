package database

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/dmlittle/discoverrewind/pkg/config"
	"github.com/dmlittle/discoverrewind/pkg/logger"
	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
)

type logQueryHook struct {
	log logger.Logger
}

func (logQueryHook) BeforeQuery(ctx context.Context, event *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

func (qh logQueryHook) AfterQuery(ctx context.Context, event *pg.QueryEvent) error {
	query, err := event.FormattedQuery()
	if err != nil {
		return errors.WithStack(err)
	}

	qh.log.Debug(string(query))

	return nil
}

// New initializes a new database struct.
func New(cfg *config.Config) (*pg.DB, error) {
	addr := fmt.Sprintf("%s:%d", cfg.DatabaseHost, cfg.DatabasePort)
	opts := &pg.Options{
		Addr:     addr,
		User:     cfg.DatabaseUser,
		Password: cfg.DatabasePassword,
		Database: cfg.DatabaseName,
	}

	if cfg.DatabaseSSLMode != "disable" {
		opts.TLSConfig = &tls.Config{ServerName: cfg.DatabaseHost}
	}

	db := pg.Connect(opts)

	// print out all queries in debug mode
	if cfg.DatabaseDebug {
		db.AddQueryHook(logQueryHook{logger.NewWithLevel("debug")})
	}

	// retry up to 5 times to ensure that the database can connect
	var err error
	for i := 0; i < cfg.DatabaseConnectRetryCount; i++ {
		_, err = db.Exec("SELECT 1")
		if err != nil {
			time.Sleep(cfg.DatabaseConnectRetryDelay)
			continue
		}
		// successfully connected
		break
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return db, nil
}
