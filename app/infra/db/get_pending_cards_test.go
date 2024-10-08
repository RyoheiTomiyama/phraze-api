package db

import (
	"context"
	"testing"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/infra/db/fixture"
	db_test "github.com/RyoheiTomiyama/phraze-api/test/db"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestGetPendingCards(t *testing.T) {
	db := db_test.GetDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	fx := fixture.New(db)
	decks := fx.CreateDeck(
		t,
		&fixture.DeckInput{UserID: lo.ToPtr("own")},
		&fixture.DeckInput{UserID: lo.ToPtr("own")},
		&fixture.DeckInput{UserID: lo.ToPtr("own")},
		&fixture.DeckInput{UserID: lo.ToPtr("own")},
	)

	cards := fx.CreateCard(t, decks[0].ID, make([]fixture.CardInput, 6)...)
	// 未来日しか持たないデッキカード
	cards2 := fx.CreateCard(t, decks[1].ID, make([]fixture.CardInput, 1)...)
	// 解答なしを持つデッキカード
	cards3 := fx.CreateCard(t, decks[2].ID,
		fixture.CardInput{Answer: lo.ToPtr(""), AIAnswer: lo.ToPtr("")},
		fixture.CardInput{Answer: lo.ToPtr(""), AIAnswer: lo.ToPtr("ai-answer")},
		fixture.CardInput{Answer: lo.ToPtr("answer")},
	)

	fx.CreateCardSchedule(t, []fixture.CardScheduleInput{
		{
			CardID:     cards[2].ID,
			ScheduleAt: time.Now().Add(-1 * time.Hour),
		},
		{
			CardID:     cards[1].ID,
			ScheduleAt: time.Now().Add(-10 * time.Hour),
		},
		{
			CardID:     cards[3].ID,
			ScheduleAt: time.Now().Add(-100 * time.Hour),
		},
		{
			CardID:     cards[4].ID,
			ScheduleAt: time.Now().Add(-1000 * time.Hour),
		},
	}...)

	fx.CreateCardSchedule(t, fixture.CardScheduleInput{
		CardID:     cards2[0].ID,
		ScheduleAt: time.Now().Add(1 * time.Hour),
	})

	fx.CreateCardSchedule(t, []fixture.CardScheduleInput{
		{
			CardID:     cards3[0].ID,
			ScheduleAt: time.Now().Add(-1 * time.Hour),
		},
		{
			CardID:     cards3[1].ID,
			ScheduleAt: time.Now().Add(-2 * time.Hour),
		},
		{
			CardID:     cards3[2].ID,
			ScheduleAt: time.Now().Add(-3 * time.Hour),
		},
	}...)

	t.Run("正常系", func(t *testing.T) {
		client := NewTestClient(t, db)

		result, err := client.GetPendingCards(
			context.Background(), decks[0].ID, time.Now(), 100, 0,
		)
		assert.NoError(t, err)

		t.Run("Schedule古い順で取得できること", func(t *testing.T) {
			assert.Len(t, result, 6)
			assert.Equal(t, cards[0].ToDomain(), result[0])
			assert.Equal(t, cards[5].ToDomain(), result[1])
			assert.Equal(t, cards[4].ToDomain(), result[2])
			assert.Equal(t, cards[3].ToDomain(), result[3])
			assert.Equal(t, cards[1].ToDomain(), result[4])
			assert.Equal(t, cards[2].ToDomain(), result[5])
		})
	})

	t.Run("すべてのカードがScheduleが未来日の場合", func(t *testing.T) {
		client := NewTestClient(t, db)

		result, err := client.GetPendingCards(
			context.Background(), decks[1].ID, time.Now(), 100, 0,
		)
		assert.NoError(t, err)
		assert.Len(t, result, 0)
	})

	t.Run("解答がないカードを含むデッキの場合", func(t *testing.T) {
		client := NewTestClient(t, db)

		result, err := client.GetPendingCards(
			context.Background(), decks[2].ID, time.Now(), 100, 0,
		)
		assert.NoError(t, err)

		t.Run("解答のあるカードのみ取得できること", func(t *testing.T) {
			assert.Len(t, result, 2)
			assert.Equal(t, cards3[2].ToDomain(), result[0])
			assert.Equal(t, cards3[1].ToDomain(), result[1])
		})
	})

	t.Run("カードを持たないデッキの場合", func(t *testing.T) {
		client := NewTestClient(t, db)

		result, err := client.GetPendingCards(
			context.Background(), decks[3].ID, time.Now(), 100, 0,
		)
		assert.NoError(t, err)
		assert.Len(t, result, 0)
	})
}
