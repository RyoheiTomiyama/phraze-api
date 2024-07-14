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

func TestUpdateCardByID(t *testing.T) {
	db := db_test.GetDB(t)
	defer db.Close()

	fx := fixture.New(db)
	decks := fx.CreateDeck(t, &fixture.DeckInput{UserID: "own"}, &fixture.DeckInput{UserID: "own"})
	cards := fx.CreateCard(t, decks[0].ID, fixture.CardInput{})

	t.Run("正常系", func(t *testing.T) {
		client := NewTestClient(t, db)

		testCases := []struct {
			name    string
			arrange func() (int64, *domain.UpdateCardInput)
			assert  func(result *domain.Card)
		}{
			{
				name: "where句なしの場合",
				arrange: func() (int64, *domain.UpdateCardInput) {
					return cards[0].ID, &domain.UpdateCardInput{
						Field: domain.UpdateCardField{
							DeckID:   lo.ToPtr(decks[1].ID),
							Question: lo.ToPtr("question-updated-full"),
							Answer:   lo.ToPtr("answer-updated-full"),
						},
					}
				},
				assert: func(result *domain.Card) {
					t.Run("更新できること", func(t *testing.T) {
						assert.NotEqual(t, cards[0].UpdatedAt.UnixMilli(), result.UpdatedAt.UnixMilli())

						expect := cards[0]
						expect.DeckID = decks[1].ID
						expect.Question = "question-updated-full"
						expect.Answer = "answer-updated-full"
						expect.UpdatedAt = result.UpdatedAt

						assert.Equal(t, expect.ToDomain(), result)
					})
				},
			},
		}

		for _, tc := range testCases {
			id, input := tc.arrange()
			result, err := client.UpdateCardByID(context.Background(), id, input)
			if assert.NoError(t, err) {
				tc.assert(result)
			}
		}
	})
}
