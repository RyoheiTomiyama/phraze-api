package db

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
)

func (c *client) DeleteCard(ctx context.Context, id int64) (int64, error) {
	e, err := c.txFrom(ctx)
	if err != nil {
		return 0, errutil.Wrap(err)
	}

	ar := int64(0)

	query := "DELETE FROM card_reviews WHERE card_id=$1"
	result, err := e.ExecContext(ctx, query, id)
	if err != nil {
		return 0, errutil.Wrap(err)
	}
	ra, err := result.RowsAffected()
	if err != nil {
		return 0, errutil.Wrap(err)
	}
	ar += ra

	query = "DELETE FROM card_schedules WHERE card_id=$1"
	result, err = e.ExecContext(ctx, query, id)
	if err != nil {
		return 0, errutil.Wrap(err)
	}
	ra, err = result.RowsAffected()
	if err != nil {
		return 0, errutil.Wrap(err)
	}
	ar += ra

	query = "DELETE FROM cards WHERE id=$1"
	result, err = e.ExecContext(ctx, query, id)
	if err != nil {
		return 0, errutil.Wrap(err)
	}
	ra, err = result.RowsAffected()
	if err != nil {
		return 0, errutil.Wrap(err)
	}
	ar += ra

	return ar, nil
}
