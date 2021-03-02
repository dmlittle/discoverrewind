// language=PostgreSQL prefix="db.Exec(" suffix=")"
package main

import (
	"github.com/go-pg/pg/v10/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec(`
			CREATE TABLE artists (
				id TEXT PRIMARY KEY,
				name TEXT NOT NULL
			);

		`)
		if err != nil {
			return err
		}

		_, err = db.Exec(`
			CREATE TABLE albums (
				id TEXT PRIMARY KEY,
				name TEXT NOT NULL
			);

		`)
		if err != nil {
			return err
		}

		_, err = db.Exec(`
			CREATE TABLE album_images (
				album_id TEXT NOT NULL REFERENCES albums(id),
				height INTEGER,
				width INTEGER,
				url TEXT NOT NULL,
				UNIQUE(album_id, height, width)
			);

		`)
		if err != nil {
			return err
		}

		_, err = db.Exec(`
			CREATE TABLE tracks (
				id TEXT PRIMARY KEY,
				name TEXT NOT NULL,
				album_id TEXT NOT NULL REFERENCES albums(id),
				disc_number INTEGER NOT NULL,
				track_number INTEGER NOT NULL,
				explicit BOOLEAN NOT NULL DEFAULT false,
				duration_ms INTEGER NOT NULL
			);

		`)
		if err != nil {
			return err
		}

		_, err = db.Exec(`
			CREATE TABLE track_artists (
				track_id TEXT NOT NULL REFERENCES tracks(id),
				artist_id TEXT NOT NULL REFERENCES artists(id),
				PRIMARY KEY(track_id, artist_id)
			);

		`)
		if err != nil {
			return err
		}

		return nil
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("DROP TABLE IF EXISTS track_artists;")
		if err != nil {
			return err
		}

		_, err = db.Exec("DROP TABLE IF EXISTS tracks;")
		if err != nil {
			return err
		}

		_, err = db.Exec("DROP TABLE IF EXISTS album_images;")
		if err != nil {
			return err
		}

		_, err = db.Exec("DROP TABLE IF EXISTS albums;")
		if err != nil {
			return err
		}

		_, err = db.Exec("DROP TABLE IF EXISTS artists;")
		if err != nil {
			return err
		}

		return nil
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20210228051839_create_tracks_artists_albums_tables", up, down, opts)
}
