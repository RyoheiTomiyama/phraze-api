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

func TestUpdateCardByID(t *testing.T) {
	db := db_test.GetDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	fx := fixture.New(db)
	decks := fx.CreateDeck(t,
		&fixture.DeckInput{UserID: lo.ToPtr("own")},
		&fixture.DeckInput{UserID: lo.ToPtr("own")},
	)

	cards := fx.CreateCard(t, decks[0].ID, lo.Map(make([]fixture.CardInput, 2), func(ci fixture.CardInput, i int) fixture.CardInput {
		ci.UpdatedAt = lo.ToPtr(time.Now().Add(-time.Hour))
		return ci
	})...)

	t.Run("正常系", func(t *testing.T) {
		client := NewTestClient(t, db)

		testCases := []struct {
			name    string
			arrange func() (int64, *domain.UpdateCardInput)
			assert  func(result *domain.Card)
		}{
			{
				name: "全更新の場合",
				arrange: func() (int64, *domain.UpdateCardInput) {
					return cards[0].ID, &domain.UpdateCardInput{
						Field: domain.UpdateCardField{
							DeckID:   lo.ToPtr(decks[1].ID),
							Question: lo.ToPtr("question-updated-full"),
							Answer:   lo.ToPtr("answer-updated-full"),
							AIAnswer: lo.ToPtr("ai-answer-updated-full"),
						},
					}
				},
				assert: func(result *domain.Card) {
					t.Run("更新できること", func(t *testing.T) {
						t.Log(time.Now().Format(time.RFC3339Nano))
						t.Log(cards[0].UpdatedAt.Format(time.RFC3339Nano))
						t.Log(result.UpdatedAt.Format(time.RFC3339Nano))
						dbCard, err := client.GetCard(context.Background(), cards[0].ID)
						assert.NoError(t, err)
						t.Log(dbCard.UpdatedAt.Format(time.RFC3339Nano))

						assert.NotEqual(t, cards[0].UpdatedAt.UnixMilli(), result.UpdatedAt.UnixMilli())

						expect := *cards[0]
						expect.DeckID = decks[1].ID
						expect.Question = "question-updated-full"
						expect.Answer = "answer-updated-full"
						expect.AIAnswer = "ai-answer-updated-full"
						expect.UpdatedAt = result.UpdatedAt

						assert.Equal(t, expect.ToDomain(), result)
					})
				},
			},
			{
				name: "Questionのみ更新の場合",
				arrange: func() (int64, *domain.UpdateCardInput) {
					return cards[1].ID, &domain.UpdateCardInput{
						Field: domain.UpdateCardField{
							Question: lo.ToPtr("question-updated-only-question"),
						},
					}
				},
				assert: func(result *domain.Card) {
					t.Run("更新できること", func(t *testing.T) {
						assert.NotEqual(t, cards[1].UpdatedAt.UnixMilli(), result.UpdatedAt.UnixMilli())

						expect := *cards[1]
						expect.Question = "question-updated-only-question"
						expect.UpdatedAt = result.UpdatedAt

						assert.Equal(t, expect.ToDomain(), result)
					})
				},
			},
		}

		for _, tc := range testCases {
			id, input := tc.arrange()
			result, err := client.UpdateCardByID(context.TODO(), id, input)
			if assert.NoError(t, err) {
				tc.assert(result)
			}
		}
	})
}
