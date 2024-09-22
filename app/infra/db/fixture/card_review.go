package fixture

import (
	"fmt"
	"testing"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/infra/db/model"
	"github.com/samber/lo"
)

type CardReviewInput struct {
	CardID     int64
	UserID     string
	ReviewedAt time.Time
	Grade      int
}

func (f *fixture) CreateCardReview(t *testing.T, cardReviews ...CardReviewInput) []*model.CardReview {
	var list []*model.CardReview
	for _, d := range cardReviews {
		list = append(list, &model.CardReview{
			CardID:     d.CardID,
			UserID:     d.UserID,
			Grade:      lo.Ternary(d.Grade == 0, 1, d.Grade),
			ReviewedAt: lo.Ternary(d.ReviewedAt.IsZero(), time.Now(), d.ReviewedAt),
		})
		// 日時の作成順を担保するためスリープする
		time.Sleep(time.Millisecond)
	}

	query := `
		INSERT INTO card_reviews (card_id, grade, reviewed_at) 
		VALUES (:card_id, :grade, :reviewed_at)
		RETURNING *
	`

	tx := f.db.MustBegin()
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		t.Fatal(err)

		return nil
	}

	var insertedCards []*model.CardReview

	for _, l := range list {
		var result model.CardReview

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

	f.CardReviews = append(f.CardReviews, insertedCards...)

	return insertedCards
}
