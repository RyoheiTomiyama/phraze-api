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

func (c *client) GetCardSchedule(ctx context.Context, cardID int64) (*domain.CardSchedule, error) {
	e := c.execerFrom(ctx)

	var schedule model.CardSchedule

	if err := sqlx.GetContext(ctx, e, &schedule, "SELECT * FROM card_schedules WHERE card_id=$1", cardID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, errutil.Wrap(err)
	}

	return schedule.ToDomain(), nil
}
