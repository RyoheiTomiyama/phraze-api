package resolver_test

import (
	"context"
	"testing"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/infra/db/fixture"
	"github.com/RyoheiTomiyama/phraze-api/router/graph/model"
	"github.com/RyoheiTomiyama/phraze-api/test/assertion"
	"github.com/RyoheiTomiyama/phraze-api/util/auth"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func assertDeck(t *testing.T, expect *domain.Deck, actual *model.Deck) {
	t.Helper()

	em := &model.Deck{
		ID:        expect.ID,
		UserID:    expect.UserID,
		Name:      expect.Name,
		CreatedAt: expect.CreatedAt,
		UpdatedAt: expect.UpdatedAt,
	}
	assert.Equal(t, em, actual)
}

func (s *resolverSuite) TestCreateDeck() {
	userID := "test_user"
	ctx := context.Background()
	ctx = auth.New(&domain.User{ID: userID}).WithCtx(ctx)

	s.T().Run("Deckが作成できること", func(t *testing.T) {
		input := model.CreateDeckInput{
			Name: "created_test",
		}
		result, err := s.resolver.Mutation().CreateDeck(ctx, input)
		assert.NoError(t, err)
		assert.Equal(t, input.Name, result.Deck.Name)
		assert.Equal(t, userID, result.Deck.UserID)
	})
}

func (s *resolverSuite) TestDecks() {
	userID := "test_user"
	ctx := context.Background()
	ctx = auth.New(&domain.User{ID: userID}).WithCtx(ctx)

	fx := fixture.New(s.dbx)

	decks := fx.CreateDeck(s.T(),
		&fixture.DeckInput{UserID: lo.ToPtr(userID)},
		&fixture.DeckInput{UserID: lo.ToPtr(userID)},
		&fixture.DeckInput{UserID: lo.ToPtr("other_user")},
	)

	s.T().Run("Decksが取得できること", func(t *testing.T) {
		result, err := s.resolver.Query().Decks(ctx)
		assert.NoError(t, err)

		t.Run("ログインユーザのDeckが取得できていること", func(t *testing.T) {
			assert.Len(t, result.Decks, 2)
			assertDeck(t, decks[0].ToDomain(), result.Decks[0])
			assertDeck(t, decks[1].ToDomain(), result.Decks[1])
		})

		t.Run("他ユーザのDeckが取得されていないこと", func(t *testing.T) {
			for _, d := range result.Decks {
				assert.NotEqual(t, decks[2].ID, d.ID)
			}
		})
	})

	s.T().Run("ログインユーザが他ユーザ所有のDeckを取得しようとした場合", func(t *testing.T) {
		decks := fx.CreateDeck(t, &fixture.DeckInput{})

		deck, err := s.resolver.Query().Deck(ctx, decks[0].ID)
		assert.Nil(t, deck)
		assertion.AssertError(t, "取得する権限がありません", errutil.CodeBadRequest, err)
	})
}

func (s *resolverSuite) TestDeck() {
	userID := "test_user"
	ctx := context.Background()
	ctx = auth.New(&domain.User{ID: userID}).WithCtx(ctx)

	fx := fixture.New(s.dbx)

	s.T().Run("Deckが取得できること", func(t *testing.T) {
		decks := fx.CreateDeck(t, &fixture.DeckInput{UserID: lo.ToPtr(userID)})

		deck, err := s.resolver.Query().Deck(ctx, decks[0].ID)
		assert.NoError(t, err)

		assertDeck(t, decks[0].ToDomain(), deck)
	})

	s.T().Run("ログインユーザが他ユーザ所有のDeckを取得しようとした場合", func(t *testing.T) {
		decks := fx.CreateDeck(t, &fixture.DeckInput{})

		deck, err := s.resolver.Query().Deck(ctx, decks[0].ID)
		assert.Nil(t, deck)
		assertion.AssertError(t, "取得する権限がありません", errutil.CodeBadRequest, err)
	})
}

func (s *resolverSuite) TestDeckInfo() {
	userID := "test_user"
	ctx := context.Background()
	ctx = auth.New(&domain.User{ID: userID}).WithCtx(ctx)

	fx := fixture.New(s.dbx)
	decks := fx.CreateDeck(s.T(), &fixture.DeckInput{UserID: lo.ToPtr(userID)}, &fixture.DeckInput{UserID: lo.ToPtr(userID)})
	cards1 := fx.CreateCard(s.T(), decks[0].ID, make([]fixture.CardInput, 2)...)
	cards2 := fx.CreateCard(s.T(), decks[1].ID, make([]fixture.CardInput, 2)...)

	fx.CreateCardSchedule(s.T(), []fixture.CardScheduleInput{
		// 過去日のみ
		{CardID: cards1[0].ID, ScheduleAt: time.Now().Add(-3 * time.Hour)},
		{CardID: cards1[1].ID, ScheduleAt: time.Now().Add(-10 * time.Hour)},
	}...)
	schedules2 := fx.CreateCardSchedule(s.T(), []fixture.CardScheduleInput{
		// 未来日のみ
		{CardID: cards2[0].ID, ScheduleAt: time.Now().Add(3 * time.Hour)},
		{CardID: cards2[1].ID, ScheduleAt: time.Now().Add(10 * time.Hour)},
	}...)

	s.T().Run("正常系", func(t *testing.T) {
		t.Run("学習待ちデッキの場合", func(t *testing.T) {
			di, err := s.resolver.Deck().DeckInfo(ctx, &model.Deck{ID: decks[0].ID})
			assert.Nil(t, err)
			assert.Equal(t, &model.DeckInfo{
				TotalCardCount:   2,
				PendingCardCount: 2,
				LearnedCardCount: 0,
				ScheduleAt:       nil,
			}, di)
		})
		t.Run("学習待ちデッキの場合", func(t *testing.T) {
			di, err := s.resolver.Deck().DeckInfo(ctx, &model.Deck{ID: decks[1].ID})
			assert.Nil(t, err)
			assert.Equal(t, &model.DeckInfo{
				TotalCardCount:   2,
				PendingCardCount: 0,
				LearnedCardCount: 2,
				ScheduleAt:       &schedules2[0].ScheduleAt,
			}, di)
		})
	})
}
