package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/infra/db/model"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
)

func (c *client) GetCards(ctx context.Context, where *domain.CardsWhere, limit, offset int) ([]*domain.Card, error) {
	e := c.execerFrom(ctx)

	query := `SELECT * FROM cards WHERE %s ORDER BY %s LIMIT %d OFFSET %d`

	var wheres []string
	arg := map[string]interface{}{}
	if where != nil {
		if where.DeckID != nil {
			wheres = append(wheres, "deck_id=:deck_id")
			arg["deck_id"] = *where.DeckID
		}
	}

	query = fmt.Sprintf(query, strings.Join(wheres, " AND "), "updated_at DESC", limit, offset)

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
