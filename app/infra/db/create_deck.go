package db

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/infra/db/model"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/jmoiron/sqlx"
)

func (c *client) CreateDeck(ctx context.Context, deck domain.Deck) (int64, error) {
	e := c.execerFrom(ctx)

	query := `INSERT INTO decks (user_id, name, created_at) VALUES (:user_id, :name, :created_at)`
	d := model.Deck{UserID: deck.UserID, Name: deck.Name}

	result, err := sqlx.NamedExecContext(ctx, e, query, d)
	if err != nil {
		return 0, errutil.Wrap(err)
	}

	return result.RowsAffected()
}
