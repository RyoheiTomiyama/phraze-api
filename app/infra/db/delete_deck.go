package db

import (
	"context"
	"fmt"

	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
)

func (c *client) DeleteDeck(ctx context.Context, id int64) (int64, error) {
	e, err := c.txFrom(ctx)
	if err != nil {
		return 0, errutil.Wrap(err)
	}

	ar := int64(0)

	query := "SELECT id FROM cards WHERE deck_id=$1"
	var cards []struct {
		ID int64
	}
	if err = sqlx.SelectContext(ctx, e, &cards, query, id); err != nil {
		return 0, errutil.Wrap(err)
	}

	if len(cards) > 0 {
		cardIDs := lo.Map(cards, func(card struct{ ID int64 }, _ int) int64 {
			return card.ID
		})

		query = "DELETE FROM card_reviews WHERE card_id IN (?)"
		query, args, err := sqlx.In(query, cardIDs)
		if err != nil {
			return 0, errutil.Wrap(err)
		}
		fmt.Println(query, args, err)
		query = e.Rebind(query)
		result, err := e.ExecContext(ctx, query, args...)
		if err != nil {
			return 0, errutil.Wrap(err)
		}
		ra, err := result.RowsAffected()
		if err != nil {
			return 0, errutil.Wrap(err)
		}
		ar += ra

		query = "DELETE FROM card_schedules WHERE card_id IN (?)"
		query, args, err = sqlx.In(query, cardIDs)
		if err != nil {
			return 0, errutil.Wrap(err)
		}
		fmt.Println(query, args, err)
		query = e.Rebind(query)
		result, err = e.ExecContext(ctx, query, args...)
		if err != nil {
			return 0, errutil.Wrap(err)
		}
		ra, err = result.RowsAffected()
		if err != nil {
			return 0, errutil.Wrap(err)
		}
		ar += ra

		query = "DELETE FROM cards WHERE deck_id=$1"
		result, err = e.ExecContext(ctx, query, id)
		if err != nil {
			return 0, errutil.Wrap(err)
		}
		fmt.Println(query, err)
		ra, err = result.RowsAffected()
		if err != nil {
			return 0, errutil.Wrap(err)
		}
		ar += ra
	}

	query = "DELETE FROM decks WHERE id=$1"
	result, err := e.ExecContext(ctx, query, id)
	if err != nil {
		return 0, errutil.Wrap(err)
	}
	fmt.Println(query, err)
	ra, err := result.RowsAffected()
	if err != nil {
		return 0, errutil.Wrap(err)
	}
	ar += ra

	return ar, nil
}
