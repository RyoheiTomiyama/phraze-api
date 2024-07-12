package resolver

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

func (s *resolverSuite) TestDeck() {
	userID := "test_user"
	ctx := context.Background()
	ctx = auth.New(&domain.User{ID: userID}).WithCtx(ctx)

	fx := fixture.New(s.dbx)

	s.T().Run("Deckが取得できること", func(t *testing.T) {
		decks := fx.CreateDeck(t, &fixture.DeckInput{UserID: userID})

		deck, err := s.resolver.Query().Deck(ctx, decks[0].ID)
		assert.NoError(t, err)

		var m model.Deck
		if err = model.FromDomain(ctx, decks[0].ToDomain(), &m); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, m, *deck)
	})

	s.T().Run("ログインユーザが他ユーザ所有のDeckを取得しようとした場合", func(t *testing.T) {
		decks := fx.CreateDeck(t, &fixture.DeckInput{})

		deck, err := s.resolver.Query().Deck(ctx, decks[0].ID)
		assert.Nil(t, deck)
		assertion.AssertError(t, "取得する権限がありません", errutil.CodeBadRequest, err)
	})
}
