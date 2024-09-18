package db

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
)

func (c *client) DeleteDeck(ctx context.Context, id int64) (int64, error) {
	e := c.execerFrom(ctx)

	query := "DELETE FROM decks WHERE id=$1"

	result, err := e.ExecContext(ctx, query, id)
	if err != nil {
		return 0, errutil.Wrap(err)
	}

	ar, err := result.RowsAffected()
	if err != nil {
		return 0, errutil.Wrap(err)
	}

	return ar, nil
}
