package db

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/infra/db/model"
	"github.com/jmoiron/sqlx"
)

func (c *client) GetDeck(ctx context.Context, id int64) (*domain.Deck, error) {
	e := c.execerFrom(ctx)

	var deck *model.Deck

	sqlx.SelectContext(ctx, e, deck, "SELECT * FROM decks WHERE id=$1", id)

	return &domain.Deck{
		ID:        deck.ID,
		UserID:    deck.UserID,
		Name:      deck.Name,
		CreateAt:  deck.CreatedAt.Time,
		UpdatedAt: deck.UpdatedAt.Time,
	}, nil
}
