package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upInitUserTokensStorage, downInitUserTokensStorage)
}

func upInitUserTokensStorage(_ context.Context, tx *sql.Tx) error {
	if _, err := tx.Exec(`
	CREATE TABLE IF NOT EXISTS user_tokens (
	    id SERIAL PRIMARY KEY,
	    wallet_id INTEGER NOT NULL,
	    tokens JSONB NOT NULL,
	    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    	deleted_at TIMESTAMP NULL,
	    FOREIGN KEY (wallet_id) REFERENCES wallets(id) ON DELETE CASCADE
	)
`); err != nil {
		return err
	}

	return nil
}

func downInitUserTokensStorage(_ context.Context, tx *sql.Tx) error {
	if _, err := tx.Exec(`DROP TABLE IF EXISTS user_tokens`); err != nil {
		return err
	}

	return nil
}
