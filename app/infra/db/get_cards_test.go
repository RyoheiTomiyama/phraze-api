package db

import (
	"context"
	"testing"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/infra/db/fixture"
	db_test "github.com/RyoheiTomiyama/phraze-api/test/db"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestGetCards(t *testing.T) {
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
		&fixture.DeckInput{UserID: lo.ToPtr("another user")},
	)

	cards := fx.CreateCard(t, decks[0].ID, make([]fixture.CardInput, 10)...)
	cards2 := fx.CreateCard(t, decks[1].ID, make([]fixture.CardInput, 10)...)

	t.Run("正常系", func(t *testing.T) {
		client := NewTestClient(t, db)

		testCases := []struct {
			name    string
			arrange func() (where *domain.CardsWhere, limit int, offset int)
			assert  func(result []*domain.Card)
		}{
			{
				name: "where句なしの場合",
				arrange: func() (where *domain.CardsWhere, limit int, offset int) {
					return nil, 2, 0
				},
				assert: func(result []*domain.Card) {
					t.Run("更新日降順で取得できること", func(t *testing.T) {
						assert.Len(t, result, 2)
						assert.Equal(t, cards2[len(cards2)-1].ToDomain(), result[0])
						assert.Equal(t, cards2[len(cards2)-2].ToDomain(), result[1])
					})
				},
			},
			{
				name: "DeckIDで絞り込んだ場合",
				arrange: func() (where *domain.CardsWhere, limit int, offset int) {
					return &domain.CardsWhere{
						DeckID: lo.ToPtr(decks[0].ID),
					}, 2, 0
				},
				assert: func(result []*domain.Card) {
					t.Run("更新日降順で取得できること", func(t *testing.T) {
						assert.Len(t, result, 2)
						assert.Equal(t, cards[len(cards)-1].ToDomain(), result[0])
						assert.Equal(t, cards[len(cards)-2].ToDomain(), result[1])
					})
				},
			},
			{
				name: "UserIDで絞り込んだ場合",
				arrange: func() (where *domain.CardsWhere, limit int, offset int) {
					return &domain.CardsWhere{
						UserID: lo.ToPtr("own"),
					}, 2, 0
				},
				assert: func(result []*domain.Card) {
					t.Run("更新日降順で取得できること", func(t *testing.T) {
						assert.Len(t, result, 2)
						assert.Equal(t, cards[len(cards)-1].ToDomain(), result[0])
						assert.Equal(t, cards[len(cards)-2].ToDomain(), result[1])
					})
				},
			},
			{
				name: "offsetが取得できるデータ数より大きい場合",
				arrange: func() (where *domain.CardsWhere, limit int, offset int) {
					return nil, 2, 9999
				},
				assert: func(result []*domain.Card) {
					t.Run("0件が取得できること", func(t *testing.T) {
						assert.Len(t, result, 0)
					})
				},
			},
		}

		for _, tc := range testCases {
			where, limit, offset := tc.arrange()
			result, err := client.GetCards(context.Background(), where, limit, offset)
			assert.NoError(t, err)
			tc.assert(result)
		}
	})
}
