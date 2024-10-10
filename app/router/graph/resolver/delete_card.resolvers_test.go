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

func (s *resolverSuite) TestDeleteCard() {
	userID := "test_user"
	ctx := context.Background()
	ctx = auth.New(&domain.User{ID: userID}).WithCtx(ctx)

	fx := fixture.New(s.dbx)
	decks := fx.CreateDeck(s.T(),
		&fixture.DeckInput{UserID: lo.ToPtr(userID)},
		&fixture.DeckInput{UserID: lo.ToPtr("other user")},
	)
	cards := fx.CreateCard(s.T(), decks[0].ID, fixture.CardInput{})
	otherCards := fx.CreateCard(s.T(), decks[1].ID, fixture.CardInput{})

	s.geminiClient.On("GenAnswer", mock.Anything, "updated question").Return("ai-answer", nil)

	s.T().Run("Cardが削除できること", func(t *testing.T) {
		input := model.DeleteCardInput{
			ID: cards[0].ID,
		}

		result, err := s.resolver.Mutation().DeleteCard(ctx, input)
		assert.NoError(t, err)
		assert.Equal(t, &model.DeleteCardOutput{AffectedRows: 1}, result)
	})

	s.T().Run("Validationエラー", func(t *testing.T) {
		testCases := []struct {
			name   string
			input  model.DeleteCardInput
			assert func(err error)
		}{
			{
				name:  "Cardが存在しない場合",
				input: model.DeleteCardInput{ID: -1},
				assert: func(err error) {
					assertion.AssertError(t, "指定されたCardは取得できません", errutil.CodeForbidden, err)
				},
			},
			{
				name:  "自分のCardでない場合",
				input: model.DeleteCardInput{ID: otherCards[0].ID},
				assert: func(err error) {
					assertion.AssertError(t, "指定されたCardは取得できません", errutil.CodeForbidden, err)
				},
			},
		}

		for _, tc := range testCases {
			result, err := s.resolver.Mutation().DeleteCard(ctx, tc.input)
			assert.Nil(t, result)
			tc.assert(err)
		}
	})
}
