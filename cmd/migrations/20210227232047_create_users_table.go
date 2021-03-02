// language=PostgreSQL prefix="db.Exec(" suffix=")"
package main

import (
	"github.com/go-pg/pg/v10/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec(`
			CREATE TABLE users (
				id TEXT PRIMARY KEY,
				spotify_id TEXT NOT NULL,
				display_name TEXT NOT NULL,
				profile_image_url TEXT,
				product TEXT NOT NULL,
				token_type TEXT NOT NULL,
				access_token TEXT NOT NULL,
				refresh_token TEXT NOT NULL,
				discover_weekly_playlist_id TEXT,
				created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
				deleted_at TIMESTAMPTZ,
				UNIQUE(spotify_id)
			);
		`)
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("DROP TABLE IF EXISTS users")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20210227232047_create_users_table", up, down, opts)
}
