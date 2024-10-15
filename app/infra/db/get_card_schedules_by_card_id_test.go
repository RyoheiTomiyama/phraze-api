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

func TestGetCardSchedulesByCardID(t *testing.T) {
	ctx := context.Background()
	db := db_test.GetDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	fx := fixture.New(db)
	var dis []*fixture.DeckInput
	for range 2 {
		dis = append(dis, &fixture.DeckInput{UserID: lo.ToPtr("own")})
	}
	decks := fx.CreateDeck(t, dis...)
	cards1 := fx.CreateCard(t, decks[0].ID, make([]fixture.CardInput, 2)...)
	cards2 := fx.CreateCard(t, decks[1].ID, make([]fixture.CardInput, 2)...)
	schedules := fx.CreateCardSchedule(t, []fixture.CardScheduleInput{
		{CardID: cards1[0].ID, ScheduleAt: time.Now().Add(-3 * time.Hour)},
		{CardID: cards1[1].ID, ScheduleAt: time.Now().Add(-10 * time.Hour)},
		{CardID: cards2[0].ID, ScheduleAt: time.Now().Add(3 * time.Hour)},
		{CardID: cards2[1].ID, ScheduleAt: time.Now().Add(10 * time.Hour)},
	}...)

	t.Run("正常系", func(t *testing.T) {
		client := NewTestClient(t, db)
		result, err := client.GetCardSchedulesByCardID(ctx,
			// 順不同で重複もさせている
			[]int64{cards1[1].ID, cards2[0].ID, cards1[0].ID, cards2[1].ID, cards2[1].ID},
		)

		assert.Nil(t, err)
		assert.Len(t, result, 4)
		assert.Equal(t, schedules[0].ID, result[cards1[0].ID].ID)
		assert.Equal(t, schedules[1].ID, result[cards1[1].ID].ID)
		assert.Equal(t, schedules[2].ID, result[cards2[0].ID].ID)
		assert.Equal(t, schedules[3].ID, result[cards2[1].ID].ID)
	})
}
