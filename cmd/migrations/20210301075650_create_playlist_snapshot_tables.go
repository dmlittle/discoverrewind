// language=PostgreSQL prefix="db.Exec(" suffix=")"
package main

import (
	"github.com/go-pg/pg/v10/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec(`
			CREATE TABLE playlist_snapshots (
				id TEXT PRIMARY KEY,
				playlist_id TEXT NOT NULL,
				snapshot_id TEXT NOT NULL,
				image_url TEXT,
				created_at DATE NOT NULL DEFAULT CURRENT_DATE,
				UNIQUE (playlist_id, created_at)
			);
		`)
		if err != nil {
			return err
		}

		_, err = db.Exec(`
			CREATE TABLE playlist_snapshot_details (
				playlist_snapshot_id TEXT NOT NULL REFERENCES playlist_snapshots(id),
				rank INT NOT NULL,
				track_id TEXT NOT NULL REFERENCES tracks(id),
				UNIQUE(playlist_snapshot_id, rank)
			);
		`)
		if err != nil {
			return err
		}

		return nil
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("DROP TABLE IF EXISTS playlist_snapshot_details;")
		if err != nil {
			return err
		}

		_, err = db.Exec("DROP TABLE IF EXISTS playlist_snapshots;")
		if err != nil {
			return err
		}

		return nil
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20210301075650_create_playlist_snapshot_tables", up, down, opts)
}
