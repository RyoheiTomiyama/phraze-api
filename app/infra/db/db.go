package db

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type execer interface {
	sqlx.QueryerContext
	sqlx.ExecerContext
}

// トランザクション内でも、各メソッドが正常に実行できるように、execerで共通化する
func (c *client) execerFrom(ctx context.Context) execer {
	tx, ok := ctx.Value(txCtxKey{}).(*sqlx.Tx)
	if !ok {
		return c.db
	}

	return tx
}
