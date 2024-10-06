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
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (s *resolverSuite) TestUpdateCardWithGenAnswer() {
	userID := "test_user"
	ctx := context.Background()
	ctx = auth.New(&domain.User{ID: userID}).WithCtx(ctx)

	fx := fixture.New(s.dbx)
	decks := fx.CreateDeck(s.T(), &fixture.DeckInput{UserID: lo.ToPtr(userID)})
	cards := fx.CreateCard(s.T(), decks[0].ID, fixture.CardInput{Question: lo.ToPtr("before question"), Answer: lo.ToPtr("before answer")})

	s.geminiClient.On("GenAnswer", mock.Anything, "updated question").Return("ai-answer", nil)

	s.T().Run("Cardが更新できること", func(t *testing.T) {
		input := model.UpdateCardWithGenAnswerInput{
			ID:       cards[0].ID,
			Question: "updated question",
		}

		result, err := s.resolver.Mutation().UpdateCardWithGenAnswer(ctx, input)
		assert.NoError(t, err)
		assertCard(t, &domain.Card{
			ID:        cards[0].ID,
			DeckID:    cards[0].DeckID,
			Question:  input.Question,
			Answer:    "ai-answer",
			AIAnswer:  "ai-answer",
			CreatedAt: result.Card.CreatedAt,
			UpdatedAt: result.Card.UpdatedAt,
		}, result.Card)
	})

	s.T().Run("Validationエラー", func(t *testing.T) {
		testCases := []struct {
			name   string
			input  model.UpdateCardWithGenAnswerInput
			assert func(err error)
		}{
			{
				name:  "IDが0値",
				input: model.UpdateCardWithGenAnswerInput{Question: "question"},
				assert: func(err error) {
					assertion.AssertError(t, "IDは必須項目です", errutil.CodeBadRequest, err)
				},
			},
			{
				name:  "Questionが0値",
				input: model.UpdateCardWithGenAnswerInput{ID: 10, Question: ""},
				assert: func(err error) {
					assertion.AssertError(t, "Questionは必須項目です", errutil.CodeBadRequest, err)
				},
			},
		}

		for _, tc := range testCases {
			result, err := s.resolver.Mutation().UpdateCardWithGenAnswer(ctx, tc.input)
			assert.Nil(t, result)
			tc.assert(err)
		}
	})
}
