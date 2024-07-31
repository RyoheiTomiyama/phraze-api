package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/infra/db/fixture"
	db_test "github.com/RyoheiTomiyama/phraze-api/test/db"
	"github.com/stretchr/testify/assert"
)

func TestGetLatestCardSchedulesByDeckID(t *testing.T) {
	ctx := context.Background()
	db := db_test.GetDB(t)
	defer db.Close()

	fx := fixture.New(db)
	var dis []*fixture.DeckInput
	for range 5 {
		dis = append(dis, &fixture.DeckInput{UserID: "own"})
	}
	decks := fx.CreateDeck(t, dis...)
	cards1 := fx.CreateCard(t, decks[0].ID, make([]fixture.CardInput, 2)...)
	cards2 := fx.CreateCard(t, decks[1].ID, make([]fixture.CardInput, 2)...)
	cards3 := fx.CreateCard(t, decks[2].ID, make([]fixture.CardInput, 2)...)
	cards4 := fx.CreateCard(t, decks[3].ID, make([]fixture.CardInput, 3)...)

	schedules := fx.CreateCardSchedule(t, []fixture.CardScheduleInput{
		//過去日のみ
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
	}...)

	t.Run("正常系", func(t *testing.T) {
		client := NewTestClient(t, db)
		result, err := client.GetLatestCardSchedulesByDeckID(ctx,
			[]int64{decks[0].ID, decks[1].ID, decks[2].ID, decks[3].ID, decks[4].ID},
		)
		fmt.Println(err)
		t.Log(result, err)

		assert.Nil(t, result[decks[0].ID])
		assert.Equal(t, schedules[2].ID, result[decks[1].ID].ID)
		assert.Equal(t, schedules[5].ID, result[decks[2].ID].ID)
		assert.Equal(t, schedules[7].ID, result[decks[3].ID].ID)
		assert.Nil(t, result[decks[4].ID])
	})
}
