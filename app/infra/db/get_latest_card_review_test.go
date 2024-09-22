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

func assertCardReview(t *testing.T, expect *domain.CardReview, actual *domain.CardReview) {
	t.Helper()
	assert.Equal(t, expect.ReviewedAt.Format(time.RFC3339Nano), actual.ReviewedAt.Format(time.RFC3339Nano))

	expect.ReviewedAt = actual.ReviewedAt
	assert.Equal(t, expect, actual)
}

func TestGetLatestCardReview(t *testing.T) {
	ctx := context.Background()
	db := db_test.GetDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	fx := fixture.New(db)
	decks := fx.CreateDeck(t, &fixture.DeckInput{UserID: lo.ToPtr("own")})
	cards := fx.CreateCard(t, decks[0].ID, make([]fixture.CardInput, 3)...)

	t.Run("正常系", func(t *testing.T) {
		cases := []struct {
			name   string
			assert func(exec func(cardID int64) (*domain.CardReview, error))
		}{
			{
				name: "単一レビュー持つカードの場合",
				assert: func(exec func(cardID int64) (*domain.CardReview, error)) {
					rs := fx.CreateCardReview(t, []fixture.CardReviewInput{
						{CardID: cards[0].ID, ReviewedAt: time.Now().Add(-3 * time.Hour), Grade: 3, UserID: "me"},
					}...)

					cr, err := exec(cards[0].ID)
					assert.NoError(t, err)
					assertCardReview(t, rs[0].ToDomain(), cr)
				},
			},
			{
				name: "複数レビュー持つカードの場合",
				assert: func(exec func(cardID int64) (*domain.CardReview, error)) {
					rs := fx.CreateCardReview(t, []fixture.CardReviewInput{
						{CardID: cards[1].ID, ReviewedAt: time.Now().Add(-10 * time.Hour), Grade: 1, UserID: "me"},
						{CardID: cards[1].ID, ReviewedAt: time.Now().Add(-3 * time.Hour), Grade: 5, UserID: "me"},
						{CardID: cards[1].ID, ReviewedAt: time.Now().Add(-5 * time.Hour), Grade: 1, UserID: "me"},
					}...)

					cr, err := exec(cards[1].ID)
					assert.NoError(t, err)
					assertCardReview(t, rs[1].ToDomain(), cr)
				},
			},
			{
				name: "レビュー持たないカードの場合",
				assert: func(exec func(cardID int64) (*domain.CardReview, error)) {
					cr, err := exec(cards[2].ID)
					assert.NoError(t, err)
					assert.Nil(t, cr)
				},
			},
		}

		client := NewTestClient(t, db)

		for _, cs := range cases {
			t.Run(cs.name, func(t *testing.T) {
				cs.assert(func(cardID int64) (*domain.CardReview, error) {
					return client.GetLatestCardReview(ctx, cardID)
				})
			})
		}
	})
}
