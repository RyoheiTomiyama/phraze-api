package resolver_test

import (
	"context"
	"testing"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/infra/db/fixture"
	"github.com/RyoheiTomiyama/phraze-api/router/graph/model"
	"github.com/RyoheiTomiyama/phraze-api/test/assertion"
	"github.com/RyoheiTomiyama/phraze-api/util/auth"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/stretchr/testify/assert"
)

func (s *resolverSuite) TestReviewCard() {
	userID := "test_user"
	ctx := context.Background()
	ctx = auth.New(&domain.User{ID: userID}).WithCtx(ctx)

	fx := fixture.New(s.dbx)
	decks := fx.CreateDeck(s.T(),
		&fixture.DeckInput{UserID: userID},
	)

	s.T().Run("Schedule, Reviewが更新されること", func(t *testing.T) {
		cards := fx.CreateCard(t, decks[0].ID, make([]fixture.CardInput, 1)...)

		result, err := s.resolver.Mutation().ReviewCard(ctx, model.ReviewCardInput{
			CardID: cards[0].ID,
			Grade:  3,
		})
		assert.NoError(t, err)
		assert.Equal(t, &model.ReviewCardOutput{CardID: cards[0].ID}, result)

		schedule, err := s.dbClient.GetCardSchedule(ctx, cards[0].ID)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, schedule)

		review, err := s.dbClient.GetCardReview(ctx, cards[0].ID)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, review)
		assert.Equal(t, 3, review.Grade)
	})

	s.T().Run("Validationエラー", func(t *testing.T) {
		testCases := []struct {
			name   string
			input  model.ReviewCardInput
			assert func(err error)
		}{
			{
				name:  "CardIDが0値の場合",
				input: model.ReviewCardInput{Grade: 3},
				assert: func(err error) {
					assertion.AssertError(t, "CardIDは必須項目です", errutil.CodeBadRequest, err)
				},
			},
			{
				name:  "Gradeが0の場合",
				input: model.ReviewCardInput{CardID: 100, Grade: 0},
				assert: func(err error) {
					assertion.AssertError(t, "Gradeは1が最小です", errutil.CodeBadRequest, err)
				},
			},
			{
				name:  "Gradeが5より大きい場合",
				input: model.ReviewCardInput{CardID: 100, Grade: 6},
				assert: func(err error) {
					assertion.AssertError(t, "Gradeは5が最大です", errutil.CodeBadRequest, err)
				},
			},
		}

		for _, tc := range testCases {
			result, err := s.resolver.Mutation().ReviewCard(ctx, tc.input)
			assert.Nil(t, result)
			tc.assert(err)
		}
	})
}
