package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/infra/db/model"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/jmoiron/sqlx"
)

func (c *client) GetDeck(ctx context.Context, id int64) (*domain.Deck, error) {
	e := c.execerFrom(ctx)

	var deck *model.Deck

	if err := sqlx.GetContext(ctx, e, &deck, "SELECT * FROM decks WHERE id=$1", id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, errutil.Wrap(err)
	}

	return &domain.Deck{
		ID:        deck.ID,
		UserID:    deck.UserID,
		Name:      deck.Name,
		CreateAt:  deck.CreatedAt.Time,
		UpdatedAt: deck.UpdatedAt.Time,
	}, nil
}
