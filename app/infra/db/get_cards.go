package db

import (
	"context"
	"fmt"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/infra/db/builder"
	"github.com/RyoheiTomiyama/phraze-api/infra/db/model"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
)

func (c *client) GetCards(ctx context.Context, where *domain.CardsWhere, limit, offset int) ([]*domain.Card, error) {
	e := c.execerFrom(ctx)

	query := `SELECT * FROM cards`
	arg := map[string]interface{}{}

	if where != nil {
		b := builder.CardsWhere(*where)
		query, arg = b.BuildNamedWhere(ctx, query, arg)
	}

	query += " ORDER BY %s LIMIT %d OFFSET %d"
	query = fmt.Sprintf(query, "updated_at DESC", limit, offset)

	query, args, err := e.BindNamed(query, arg)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	var cards []*model.Card
	if err = sqlx.SelectContext(ctx, e, &cards, query, args...); err != nil {
		return nil, errutil.Wrap(err)
	}

	return lo.Map(cards, func(c *model.Card, _ int) *domain.Card {
		return c.ToDomain()
	}), nil
}
