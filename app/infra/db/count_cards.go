package db

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/infra/db/builder"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
)

func (c *client) CountCards(ctx context.Context, where *domain.CardsWhere) (int, error) {
	e := c.execerFrom(ctx)

	query := `SELECT COUNT(DISTINCT id) FROM cards`
	arg := map[string]interface{}{}

	if where != nil {
		b := builder.CardsWhere(*where)
		query, arg = b.BuildNamedWhere(ctx, query, arg)
	}

	query, args, err := e.BindNamed(query, arg)
	if err != nil {
		return 0, errutil.Wrap(err)
	}

	type Result struct {
		Count int
	}
	var result Result
	if err = e.QueryRowxContext(ctx, query, args...).StructScan(&result); err != nil {
		// if errors.Is(err, sql.ErrNoRows) {
		// 	return lo.ToPtr(0), nil
		// }

		return 0, errutil.Wrap(err)
	}

	return result.Count, nil
}
