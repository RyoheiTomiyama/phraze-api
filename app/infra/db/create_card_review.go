package db

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/infra/db/model"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
)

func (c *client) CreateCardReview(ctx context.Context, review *domain.CardReview) (*domain.CardReview, error) {
	e := c.execerFrom(ctx)

	query := `
		INSERT INTO card_reviews (card_id, grade, user_id)
		VALUES (:card_id, :grade, :user_id)
		RETURNING *
	`

	d := model.CardReview{CardID: review.CardID, Grade: review.Grade, UserID: review.UserID}

	query, args, err := e.BindNamed(query, d)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	if err = e.QueryRowxContext(ctx, query, args...).StructScan(&d); err != nil {
		return nil, errutil.Wrap(err)
	}

	return d.ToDomain(), nil
}
