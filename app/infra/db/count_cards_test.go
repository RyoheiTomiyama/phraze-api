package db

import (
	"context"
	"testing"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/infra/db/fixture"
	db_test "github.com/RyoheiTomiyama/phraze-api/test/db"
	"github.com/stretchr/testify/assert"
)

func TestCountCards(t *testing.T) {
	db := db_test.GetDB(t)
	defer func() {
		err := db.Close()
		t.Fatal(err)
	}()

	fx := fixture.New(db)
	decks := fx.CreateDeck(
		t,
		&fixture.DeckInput{UserID: "own"},
		&fixture.DeckInput{UserID: "another user"},
	)

	fx.CreateCard(t, decks[0].ID, make([]fixture.CardInput, 10)...)
	fx.CreateCard(t, decks[1].ID, make([]fixture.CardInput, 5)...)

	t.Run("正常系", func(t *testing.T) {
		client := NewTestClient(t, db)

		testCases := []struct {
			name    string
			arrange func() (where *domain.CardsWhere)
			assert  func(result int)
		}{
			{
				name: "where句なしの場合",
				arrange: func() (where *domain.CardsWhere) {
					return nil
				},
				assert: func(result int) {
					t.Run("全件の数であること", func(t *testing.T) {
						assert.Equal(t, 15, result)
					})
				},
			},
			{
				name: "DeckIDで絞り込んだ場合",
				arrange: func() (where *domain.CardsWhere) {
					return &domain.CardsWhere{
						DeckID: decks[0].ID,
					}
				},
				assert: func(result int) {
					t.Run("Deckに紐づくCardの数であること", func(t *testing.T) {
						assert.Equal(t, 10, result)
					})
				},
			},
			{
				name: "絞り込み結果が0件の場合",
				arrange: func() (where *domain.CardsWhere) {
					return &domain.CardsWhere{
						DeckID: -1,
					}
				},
				assert: func(result int) {
					t.Run("0件が取得できること", func(t *testing.T) {
						assert.Equal(t, 0, result)
					})
				},
			},
		}

		for _, tc := range testCases {
			where := tc.arrange()
			result, err := client.CountCards(context.Background(), where)
			assert.NoError(t, err)
			tc.assert(result)
		}
	})
}
