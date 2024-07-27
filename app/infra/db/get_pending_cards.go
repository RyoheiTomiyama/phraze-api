package db

import (
	"context"
	"fmt"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/infra/db/model"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
)

func (c *client) GetPendingCards(ctx context.Context, deckID int64, to time.Time, limit, offset int) ([]*domain.Card, error) {
	e := c.execerFrom(ctx)

	query := `
		SELECT cards.* FROM cards
		LEFT JOIN card_schedules ON card_schedules.card_id = cards.id
		WHERE cards.deck_id=:deck_id 
			AND (card_schedules.schedule_at IS NULL OR card_schedules.schedule_at < :schedule_at)
	`
	arg := map[string]interface{}{
		"deck_id":     deckID,
		"schedule_at": to,
	}

	query = query + " ORDER BY %s LIMIT %d OFFSET %d"
	query = fmt.Sprintf(query, "card_schedules.schedule_at ASC NULLS FIRST", limit, offset)

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
