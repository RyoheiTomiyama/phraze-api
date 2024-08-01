package db

import (
	"context"
	"fmt"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/jmoiron/sqlx"
)

func (c *client) GetDeckInfosByDeckID(ctx context.Context, deckIDs []int64) (map[int64]*domain.DeckInfo, error) {
	e := c.execerFrom(ctx)

	query := `
		-- 未来日の最新のスケジュール
		SELECT
			c.deck_id,
			COUNT(c.id) AS total_card_count,
			COUNT(cs.card_id) FILTER(WHERE cs.schedule_at > NOW()) AS schedule_card_count,
			MIN(cs.schedule_at) FILTER(WHERE cs.schedule_at > NOW()) AS schedule_at
		FROM
			cards c
		LEFT JOIN
			card_schedules cs ON c.id = cs.card_id
		WHERE
			c.deck_id IN (:deck_ids)
		GROUP BY
			c.deck_id;
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
		DeckID            int64      `db:"deck_id"`
		TotalCardCount    int        `db:"total_card_count"`
		ScheduleCardCount int        `db:"schedule_card_count"`
		SchedulaAt        *time.Time `db:"schedule_at"`
	}
	var schedulesWithDeckID []*scheduleWithDeckID
	if err = sqlx.SelectContext(ctx, e, &schedulesWithDeckID, query, args...); err != nil {
		return nil, errutil.Wrap(err)
	}

	fmt.Println(schedulesWithDeckID)

	smap := make(map[int64]*domain.DeckInfo)
	for _, s := range schedulesWithDeckID {
		if s == nil {
			continue
		}
		deckInfo := &domain.DeckInfo{
			TotalCardCount:    s.TotalCardCount,
			ScheduleCardCount: s.ScheduleCardCount,
			ScheduleAt:        s.SchedulaAt,
		}
		smap[s.DeckID] = deckInfo
	}

	return smap, nil
}
