package fixture

import (
	"fmt"
	"testing"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/infra/db/model"
	"github.com/samber/lo"
)

type CardScheduleInput struct {
	CardID     int64
	ScheduleAt time.Time
	Interval   int
	Efactor    float64
}

func (f *fixture) CreateCardSchedule(t *testing.T, cardSchedules ...CardScheduleInput) []*model.CardSchedule {
	var list []*model.CardSchedule
	for _, d := range cardSchedules {
		list = append(list, &model.CardSchedule{
			CardID:     d.CardID,
			Interval:   lo.Ternary(d.Interval == 0, 20, d.Interval),
			Efactor:    lo.Ternary(d.Efactor == 0, 1.0, d.Efactor),
			ScheduleAt: lo.Ternary(d.ScheduleAt.IsZero(), time.Now(), d.ScheduleAt),
		})
		// 日時の作成順を担保するためスリープする
		time.Sleep(time.Millisecond)
	}

	query := `
		INSERT INTO card_schedules (card_id, interval, efactor) 
		VALUES (:card_id, :interval, :efactor)
		RETURNING *
	`

	tx := f.db.MustBegin()
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		t.Fatal(err)

		return nil
	}

	var insertedCards []*model.CardSchedule

	for _, l := range list {
		var result model.CardSchedule

		if err = stmt.QueryRowx(l).StructScan(&result); err != nil {
			fmt.Print(fmt.Errorf("%w", err))
			if inerr := tx.Rollback(); inerr != nil {
				t.Fatal(inerr)

				return nil
			}

			t.Fatal(err)

			return nil
		}

		insertedCards = append(insertedCards, &result)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)

		return nil
	}

	f.CardSchedules = append(f.CardSchedules, insertedCards...)

	return insertedCards
}
