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

func (c *client) GetLatestCardReview(ctx context.Context, cardID int64) (*domain.CardReview, error) {
	e := c.execerFrom(ctx)

	var review model.CardReview

	if err := sqlx.GetContext(ctx, e, &review, "SELECT * FROM card_reviews WHERE card_id=$1", cardID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, errutil.Wrap(err)
	}

	return review.ToDomain(), nil
}
