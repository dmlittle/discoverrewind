// language=PostgreSQL prefix="db.Exec(" suffix=")"
package main

import (
	"github.com/go-pg/pg/v10/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec("ALTER TABLE users ADD COLUMN token_expiration TIMESTAMPTZ;")
		if err != nil {
			return err
		}

		_, err = db.Exec("UPDATE users SET token_expiration = NOW() - INTERVAL '1 HOUR';")
		if err != nil {
			return err
		}

		_, err = db.Exec("ALTER TABLE users ALTER COLUMN token_expiration SET NOT NULL;")
		if err != nil {
			return err
		}

		return nil
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("ALTER TABLE users DROP COLUMN token_expiration;")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20210303072421_add_token_expiration_to_users_table", up, down, opts)
}
