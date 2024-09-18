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
)

func (s *resolverSuite) TestDeleteDeck() {
	userID := "test_user"
	ctx := context.Background()
	ctx = auth.New(&domain.User{ID: userID}).WithCtx(ctx)

	fx := fixture.New(s.dbx)
	decks := fx.CreateDeck(s.T(),
		&fixture.DeckInput{UserID: lo.ToPtr(userID)},
		&fixture.DeckInput{UserID: lo.ToPtr("other user")},
	)

	s.T().Run("Deckを削除できること", func(t *testing.T) {
		input := model.DeleteDeckInput{
			ID: decks[0].ID,
		}

		result, err := s.resolver.Mutation().DeleteDeck(ctx, input)
		assert.NoError(t, err)
		assert.Equal(t, &model.DeleteDeckOutput{AffectedRows: 1}, result)
	})

	s.T().Run("Validationエラー", func(t *testing.T) {
		testCases := []struct {
			name   string
			input  model.DeleteDeckInput
			assert func(err error)
		}{
			{
				name:  "Deckが存在しない場合",
				input: model.DeleteDeckInput{ID: -1},
				assert: func(err error) {
					assertion.AssertError(t, "指定されたDeckは取得できません", errutil.CodeForbidden, err)
				},
			},
			{
				name:  "自分のCardでない場合",
				input: model.DeleteDeckInput{ID: decks[1].ID},
				assert: func(err error) {
					assertion.AssertError(t, "指定されたDeckは取得できません", errutil.CodeForbidden, err)
				},
			},
		}

		for _, tc := range testCases {
			result, err := s.resolver.Mutation().DeleteDeck(ctx, tc.input)
			assert.Nil(t, result)
			tc.assert(err)
		}
	})
}
