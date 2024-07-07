package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/infra/db/model"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
)

func (c *client) GetDecks(ctx context.Context, userID string) ([]*domain.Deck, error) {
	e := c.execerFrom(ctx)

	var decks []*model.Deck

	if err := sqlx.SelectContext(ctx, e, &decks, "SELECT * FROM decks WHERE user_id=$1", userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, errutil.Wrap(err)
	}

	d := lo.Map(decks, func(item *model.Deck, _ int) *domain.Deck {
		return item.ToDomain()
	})

	return d, nil
}
