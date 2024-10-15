package db

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/infra/db/model"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/jmoiron/sqlx"
)

func (c *client) GetCardSchedulesByCardID(ctx context.Context, cardIDs []int64) (map[int64]*domain.CardSchedule, error) {
	e := c.execerFrom(ctx)

	query := `
		SELECT * FROM card_schedules
		WHERE card_schedules.card_id IN (:card_ids)
	`
	arg := map[string]interface{}{
		"card_ids": cardIDs,
	}

	query, args, err := sqlx.Named(query, arg)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	query = e.Rebind(query)

	var schedules []*model.CardSchedule
	if err = sqlx.SelectContext(ctx, e, &schedules, query, args...); err != nil {
		return nil, errutil.Wrap(err)
	}

	smap := make(map[int64]*domain.CardSchedule)
	for _, s := range schedules {
		if s == nil {
			continue
		}
		cs := s.ToDomain()
		if cs != nil {
			smap[s.CardID] = cs
		}
	}

	return smap, nil
}
