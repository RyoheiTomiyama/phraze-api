package db

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/infra/db/model"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
)

func (c *client) CreateDeck(ctx context.Context, deck domain.Deck) (*domain.Deck, error) {
	e := c.execerFrom(ctx)

	query := `INSERT INTO decks (user_id, name) VALUES (:user_id, :name) RETURNING *`
	d := model.Deck{UserID: deck.UserID, Name: deck.Name}

	query, args, err := e.BindNamed(query, d)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	if err = e.QueryRowxContext(ctx, query, args...).StructScan(&d); err != nil {
		return nil, errutil.Wrap(err)
	}

	return d.ToDomain(), nil
}
