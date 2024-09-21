package db

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/jmoiron/sqlx"
)

type execer interface {
	sqlx.QueryerContext
	sqlx.ExecerContext
	sqlx.ExtContext
}

// トランザクション内でも、各メソッドが正常に実行できるように、execerで共通化する
func (c *client) execerFrom(ctx context.Context) execer {
	tx, ok := ctx.Value(txCtxKey{}).(*sqlx.Tx)
	if !ok {
		return c.db
	}

	return tx
}

// トランザクション内で絶対に実行した場合はこちらを使う
func (c *client) txFrom(ctx context.Context) (execer, error) {
	tx, ok := ctx.Value(txCtxKey{}).(*sqlx.Tx)
	if !ok {
		return nil, errutil.New(errutil.CodeInternalError, "transaction内で実行してください")
	}

	return tx, nil
}
