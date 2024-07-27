package db

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/infra/db/model"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
)

func (c *client) UpsertCardSchedule(ctx context.Context, schedule *domain.CardSchedule) (*domain.CardSchedule, error) {
	e := c.execerFrom(ctx)

	query := `
		INSERT INTO card_schedules (card_id, schedule_at, interval, efactor)
		VALUES (:card_id, :schedule_at, :interval, :efactor)
		ON CONFLICT (card_id)
		DO UPDATE SET
			schedule_at=EXCLUDED.schedule_at,
			interval=EXCLUDED.interval,
			efactor=EXCLUDED.efactor
		RETURNING *
	`

	d := model.CardSchedule{
		CardID:     schedule.CardID,
		ScheduleAt: schedule.ScheduleAt,
		Interval:   schedule.Interval,
		Efactor:    schedule.Efactor,
	}

	query, args, err := e.BindNamed(query, d)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	if err = e.QueryRowxContext(ctx, query, args...).StructScan(&d); err != nil {
		return nil, errutil.Wrap(err)
	}

	return d.ToDomain(), nil
}
