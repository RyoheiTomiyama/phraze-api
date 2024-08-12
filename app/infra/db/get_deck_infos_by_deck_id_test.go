package db

import (
	"context"
	"testing"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/infra/db/fixture"
	db_test "github.com/RyoheiTomiyama/phraze-api/test/db"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestGetDeckInfosByDeckID(t *testing.T) {
	ctx := context.Background()
	db := db_test.GetDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	fx := fixture.New(db)
	var dis []*fixture.DeckInput
	for range 6 {
		dis = append(dis, &fixture.DeckInput{UserID: lo.ToPtr("own")})
	}
	decks := fx.CreateDeck(t, dis...)
	cards1 := fx.CreateCard(t, decks[0].ID, make([]fixture.CardInput, 2)...)
	cards2 := fx.CreateCard(t, decks[1].ID, make([]fixture.CardInput, 2)...)
	cards3 := fx.CreateCard(t, decks[2].ID, make([]fixture.CardInput, 2)...)
	cards4 := fx.CreateCard(t, decks[3].ID, make([]fixture.CardInput, 3)...)
	cards5 := fx.CreateCard(t, decks[4].ID, []fixture.CardInput{
		{Answer: lo.ToPtr(""), AIAnswer: lo.ToPtr("")},
		{},
	}...)

	schedules := fx.CreateCardSchedule(t, []fixture.CardScheduleInput{
		// 過去日のみ
		{CardID: cards1[0].ID, ScheduleAt: time.Now().Add(-3 * time.Hour)},
		{CardID: cards1[1].ID, ScheduleAt: time.Now().Add(-10 * time.Hour)},
		// 未来日のみ
		{CardID: cards2[0].ID, ScheduleAt: time.Now().Add(3 * time.Hour)},
		{CardID: cards2[1].ID, ScheduleAt: time.Now().Add(10 * time.Hour)},
		// 過去日＋未来日
		{CardID: cards3[0].ID, ScheduleAt: time.Now().Add(-3 * time.Hour)},
		{CardID: cards3[1].ID, ScheduleAt: time.Now().Add(10 * time.Hour)},
		// 過去日＋未来日＋スケジュールなし
		{CardID: cards4[0].ID, ScheduleAt: time.Now().Add(-3 * time.Hour)},
		{CardID: cards4[1].ID, ScheduleAt: time.Now().Add(10 * time.Hour)},
		// 未来日（解答なし）
		{CardID: cards5[0].ID, ScheduleAt: time.Now().Add(10 * time.Hour)},
		{CardID: cards5[1].ID, ScheduleAt: time.Now().Add(11 * time.Hour)},
	}...)

	t.Run("正常系", func(t *testing.T) {
		client := NewTestClient(t, db)
		result, err := client.GetDeckInfosByDeckID(ctx,
			[]int64{decks[0].ID, decks[1].ID, decks[2].ID, decks[3].ID, decks[4].ID, decks[5].ID},
		)
		assert.NoError(t, err)

		assert.Equal(t, &domain.DeckInfo{
			TotalCardCount: 2, PendingCardCount: 2, LearnedCardCount: 0, ScheduleAt: nil,
		}, result[decks[0].ID])
		assert.Equal(t, &domain.DeckInfo{
			TotalCardCount: 2, PendingCardCount: 0, LearnedCardCount: 2, ScheduleAt: &schedules[2].ScheduleAt,
		}, result[decks[1].ID])
		assert.Equal(t, &domain.DeckInfo{
			TotalCardCount: 2, PendingCardCount: 1, LearnedCardCount: 1, ScheduleAt: nil,
		}, result[decks[2].ID])
		assert.Equal(t, &domain.DeckInfo{
			TotalCardCount: 3, PendingCardCount: 2, LearnedCardCount: 1, ScheduleAt: nil,
		}, result[decks[3].ID])
		assert.Equal(t, &domain.DeckInfo{
			TotalCardCount: 1, PendingCardCount: 0, LearnedCardCount: 1, ScheduleAt: &schedules[9].ScheduleAt,
		}, result[decks[4].ID])
		assert.Equal(t, &domain.DeckInfo{
			TotalCardCount: 0, PendingCardCount: 0, LearnedCardCount: 0, ScheduleAt: nil,
		}, result[decks[5].ID])
	})
}
