package db

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/infra/db/model"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/jmoiron/sqlx"
)

func (c *client) GetLatestCardSchedulesByDeckID(ctx context.Context, deckIDs []int64) (map[int64]*domain.CardSchedule, error) {
	e := c.execerFrom(ctx)

	query := `
		WITH ranked_schedules AS (
			SELECT 
				cs.*,
				c.deck_id as deck_id,
				ROW_NUMBER() OVER (PARTITION BY c.deck_id ORDER BY cs.schedule_at ASC) AS rn
			FROM 
				card_schedules cs
			JOIN 
				cards c ON cs.card_id = c.id
			WHERE 
				c.deck_id IN (:deck_ids)
		)
		SELECT 
			id, card_id, schedule_at, interval, efactor, deck_id
		FROM 
			ranked_schedules
		WHERE 
			rn = 1;
	`
	arg := map[string]interface{}{
		"deck_ids": deckIDs,
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

	type scheduleWithDeckID struct {
		model.CardSchedule
		DeckID int64 `db:"deck_id"`
	}
	var schedulesWithDeckID []*scheduleWithDeckID
	if err = sqlx.SelectContext(ctx, e, &schedulesWithDeckID, query, args...); err != nil {
		return nil, errutil.Wrap(err)
	}

	smap := make(map[int64]*domain.CardSchedule)
	for _, s := range schedulesWithDeckID {
		if s == nil {
			continue
		}
		cs := s.ToDomain()
		if cs != nil {
			smap[s.DeckID] = cs
		}
	}

	return smap, nil
}
