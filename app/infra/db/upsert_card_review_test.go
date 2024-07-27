package db

import (
	"context"
	"testing"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/infra/db/fixture"
	db_test "github.com/RyoheiTomiyama/phraze-api/test/db"
	"github.com/stretchr/testify/assert"
)

func assertUpsertedCardReview(t *testing.T, expect *domain.CardReview, actual *domain.CardReview) {
	t.Helper()
	assert.NotEqual(t, expect.ReviewedAt.Format(time.RFC3339Nano), actual.ReviewedAt.Format(time.RFC3339Nano))

	expect.ID = actual.ID
	expect.ReviewedAt = actual.ReviewedAt
	assert.Equal(t, expect, actual)
}

func TestUpsertCardReview(t *testing.T) {
	ctx := context.Background()

	db := db_test.GetDB(t)
	defer db.Close()

	fx := fixture.New(db)
	decks := fx.CreateDeck(t, &fixture.DeckInput{UserID: "own"}, &fixture.DeckInput{UserID: "own"})
	cards := fx.CreateCard(t, decks[0].ID, make([]fixture.CardInput, 1)...)

	client := NewTestClient(t, db)
	t.Run("はじめてのレビューの場合", func(t *testing.T) {
		review := &domain.CardReview{CardID: cards[0].ID, Grade: 5}
		result, err := client.UpsertCardReview(ctx, review)
		assert.NoError(t, err)
		assertUpsertedCardReview(t, review, result)
		// assert.Equal(t, re)

		review = result
	})

	t.Run("レビュー更新の場合", func(t *testing.T) {
		review := &domain.CardReview{CardID: cards[0].ID, Grade: 5}

		// 1回目
		_, err := client.UpsertCardReview(ctx, review)
		assert.NoError(t, err)

		time.Sleep(time.Millisecond)

		// 2回目
		review.Grade = 3
		result, err := client.UpsertCardReview(ctx, review)
		assert.NoError(t, err)
		assertUpsertedCardReview(t, review, result)
	})

	t.Run("異常系", func(t *testing.T) {
		cases := []struct {
			name    string
			arrange func() (review *domain.CardReview)
			assert  func(err error)
		}{
			{
				name: "Gradeが入力されなかった場合",
				arrange: func() (review *domain.CardReview) {
					review = &domain.CardReview{CardID: cards[0].ID}

					return review
				},
				assert: func(err error) {
					assert.Error(t, err)
					assert.ErrorContains(t, err, "card_reviews_grade_check")
				},
			},
			{
				name: "5より大きいGradeが入力された場合",
				arrange: func() (review *domain.CardReview) {
					review = &domain.CardReview{CardID: cards[0].ID, Grade: 6}

					return review
				},
				assert: func(err error) {
					assert.Error(t, err)
					assert.ErrorContains(t, err, "card_reviews_grade_check")
				},
			},
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
			defer db.Close()
			client := NewTestClient(t, db)

			review := tc.arrange()

			result, err := client.UpsertCardReview(ctx, review)
			assert.Nil(t, result)
			tc.assert(err)
		}
	})
}
