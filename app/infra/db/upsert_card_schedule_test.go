package db

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/infra/db/fixture"
	db_test "github.com/RyoheiTomiyama/phraze-api/test/db"
	"github.com/stretchr/testify/assert"
)

func assertUpsertedCardSchedule(t *testing.T, expect *domain.CardSchedule, actual *domain.CardSchedule) {
	t.Helper()

	// 環境によってtime.TimeがLocal/UTC、nano秒/micro秒のズレが起こるのでFormatして比較
	assert.Equal(t, expect.ScheduleAt.Format(time.StampMicro), actual.ScheduleAt.Format(time.StampMicro))
	expect.ID = actual.ID
	expect.ScheduleAt = actual.ScheduleAt
	assert.Equal(t, expect, actual)
}

func TestUpsertCardSchedule(t *testing.T) {
	if err := os.Setenv("TZ", "UTC"); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	db := db_test.GetDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	fx := fixture.New(db)
	decks := fx.CreateDeck(t, &fixture.DeckInput{UserID: "own"}, &fixture.DeckInput{UserID: "own"})
	cards := fx.CreateCard(t, decks[0].ID, make([]fixture.CardInput, 1)...)

	client := NewTestClient(t, db)
	t.Run("スケジュールの作成・更新ができる", func(t *testing.T) {
		schedule := &domain.CardSchedule{
			CardID: cards[0].ID, ScheduleAt: time.Now().Add(time.Hour),
			Interval: 20, Efactor: 1.1,
		}

		// 新規
		result, err := client.UpsertCardSchedule(ctx, schedule)
		assert.NoError(t, err)
		assertUpsertedCardSchedule(t, schedule, result)

		schedule = result
		schedule.ScheduleAt = time.Now().Add(10 * time.Hour)

		// 更新
		result, err = client.UpsertCardSchedule(ctx, schedule)
		assert.NoError(t, err)
		assertUpsertedCardSchedule(t, schedule, result)
	})

	t.Run("異常系", func(t *testing.T) {
		cases := []struct {
			name    string
			arrange func() (review *domain.CardReview)
			assert  func(err error)
		}{
			{
				name: "存在しないCardIDが入力された場合",
				arrange: func() (review *domain.CardReview) {
					review = &domain.CardReview{CardID: 0, Grade: 3}

					return review
				},
				assert: func(err error) {
					assert.Error(t, err)
					assert.ErrorContains(t, err, "card_reviews_card_id_fkey")
				},
			},
		}

		for _, tc := range cases {
			// DBエラー起こすと、txdbのトランザクション内でエラーになるので別でコネクション作る
			db := db_test.GetDB(t)
			defer func() {
				err := db.Close()
				t.Fatal(err)
			}()

			client := NewTestClient(t, db)

			review := tc.arrange()

			result, err := client.UpsertCardReview(ctx, review)
			assert.Nil(t, result)
			tc.assert(err)
		}
	})
}
