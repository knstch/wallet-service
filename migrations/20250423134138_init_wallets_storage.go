package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upInitWalletsStorage, downInitWalletsStorage)
}

func upInitWalletsStorage(_ context.Context, tx *sql.Tx) error {
	if _, err := tx.Exec(`
	CREATE TABLE IF NOT EXISTS wallets (
	    id SERIAL PRIMARY KEY,
	    user_id INTEGER UNIQUE NOT NULL,
	    public_key VARCHAR(255) NOT NULL,
	    private_key BYTEA NOT NULL,
	    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
	    deleted_at TIMESTAMP WITH TIME ZONE NULL
	)
`); err != nil {
		return err
	}

	return nil
}

func downInitWalletsStorage(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.Exec(`DROP TABLE IF EXISTS wallets`); err != nil {
		return err
	}

	return nil
}
