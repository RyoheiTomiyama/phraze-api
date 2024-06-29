package db

import (
	"context"
	"database/sql"
)

type txCtxKey struct{}

func (c *client) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := c.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return err
	}

	ctxWithTx := context.WithValue(ctx, txCtxKey{}, tx)

	if err := fn(ctxWithTx); err != nil {
		if innerErr := tx.Rollback(); innerErr != nil {
			return err
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
