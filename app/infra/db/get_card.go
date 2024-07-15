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

func (c *client) GetCard(ctx context.Context, id int64) (*domain.Card, error) {
	e := c.execerFrom(ctx)

	var card model.Card

	if err := sqlx.GetContext(ctx, e, &card, "SELECT * FROM cards WHERE id=$1", id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, errutil.Wrap(err)
	}

	return card.ToDomain(), nil
}
