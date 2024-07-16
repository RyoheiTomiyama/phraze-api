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
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func assertCard(t *testing.T, expect *domain.Card, actual *model.Card) {
	t.Helper()

	em := &model.Card{
		ID:        expect.ID,
		DeckID:    expect.DeckID,
		Question:  expect.Question,
		Answer:    expect.Answer,
		CreatedAt: expect.CreateAt,
		UpdatedAt: expect.UpdatedAt,
	}
	assert.Equal(t, em, actual)
}

func (s *resolverSuite) TestCreateCard() {
	userID := "test_user"
	ctx := context.Background()
	ctx = auth.New(&domain.User{ID: userID}).WithCtx(ctx)

	fx := fixture.New(s.dbx)
	decks := fx.CreateDeck(s.T(), &fixture.DeckInput{UserID: userID})

	s.T().Run("Cardが作成できること", func(t *testing.T) {
		input := model.CreateCardInput{
			DeckID:   decks[0].ID,
			Question: "question",
			Answer:   lo.ToPtr("answer"),
		}
		result, err := s.resolver.Mutation().CreateCard(ctx, input)
		assert.NoError(t, err)
		assert.Equal(t, input.DeckID, result.Card.DeckID)
		assert.Equal(t, input.Question, result.Card.Question)
		assert.Equal(t, lo.FromPtr(input.Answer), result.Card.Answer)
	})

	s.T().Run("Validationエラー", func(t *testing.T) {
		testCases := []struct {
			name   string
			input  model.CreateCardInput
			assert func(err error)
		}{
			{
				name:  "DeckIDが0値",
				input: model.CreateCardInput{Question: "question", Answer: lo.ToPtr("answer")},
				assert: func(err error) {
					assertion.AssertError(t, "DeckIDは必須項目です", errutil.CodeBadRequest, err)
				},
			},
			{
				name:  "Questionが0値",
				input: model.CreateCardInput{DeckID: 10, Answer: lo.ToPtr("answer")},
				assert: func(err error) {
					assertion.AssertError(t, "Questionは必須項目です", errutil.CodeBadRequest, err)
				},
			},
		}

		for _, tc := range testCases {
			result, err := s.resolver.Mutation().CreateCard(ctx, tc.input)
			assert.Nil(t, result)
			tc.assert(err)
		}
	})

	s.T().Run("存在しないDeckのCardを作成する場合", func(t *testing.T) {
		input := model.CreateCardInput{
			DeckID:   -1,
			Question: "question",
			Answer:   lo.ToPtr("answer"),
		}
		result, err := s.resolver.Mutation().CreateCard(ctx, input)
		assert.Nil(t, result)
		assertion.AssertError(t, "指定されたDeckのCardは作成できません", errutil.CodeBadRequest, err)
	})

	s.T().Run("他ユーザのDeckのCardを作成する場合", func(t *testing.T) {
		decks := fx.CreateDeck(s.T(), &fixture.DeckInput{UserID: "other_user"})

		input := model.CreateCardInput{
			DeckID:   decks[0].ID,
			Question: "question",
			Answer:   lo.ToPtr("answer"),
		}
		result, err := s.resolver.Mutation().CreateCard(ctx, input)
		assert.Nil(t, result)
		assertion.AssertError(t, "指定されたDeckのCardは作成できません", errutil.CodeBadRequest, err)
	})
}

func (s *resolverSuite) TestCards() {
	userID := "test_user"
	ctx := context.Background()
	ctx = auth.New(&domain.User{ID: userID}).WithCtx(ctx)

	fx := fixture.New(s.dbx)
	decks := fx.CreateDeck(s.T(),
		&fixture.DeckInput{UserID: userID},
		&fixture.DeckInput{UserID: "other_user"},
	)
	cards := fx.CreateCard(s.T(), decks[0].ID, make([]fixture.CardInput, 10)...)
	fx.CreateCard(s.T(), decks[1].ID, make([]fixture.CardInput, 10)...)

	s.T().Run("Validationエラー", func(t *testing.T) {
		testCases := []struct {
			name   string
			input  model.CardsInput
			assert func(err error)
		}{
			{
				name:  "DeckIDが0値の場合",
				input: model.CardsInput{Where: &model.CardsWhere{}},
				assert: func(err error) {
					assertion.AssertError(t, "DeckIDは必須項目です", errutil.CodeBadRequest, err)
				},
			},
			{
				name:  "Limitが100より大きい場合",
				input: model.CardsInput{Where: &model.CardsWhere{DeckID: 100}, Limit: lo.ToPtr(101)},
				assert: func(err error) {
					assertion.AssertError(t, "Limitは100が最大です", errutil.CodeBadRequest, err)
				},
			},
		}

		for _, tc := range testCases {
			result, err := s.resolver.Query().Cards(ctx, &tc.input)
			assert.Nil(t, result)
			tc.assert(err)
		}
	})

	s.T().Run("Cardsが取得できること", func(t *testing.T) {
		result, err := s.resolver.Query().Cards(ctx, &model.CardsInput{
			Where: &model.CardsWhere{
				DeckID: decks[0].ID,
			},
			Limit:  lo.ToPtr(2),
			Offset: lo.ToPtr(0),
		})
		if assert.Nil(t, err) {
			assert.Len(t, result.Cards, 2)
			assertCard(t, cards[len(cards)-1].ToDomain(), result.Cards[0])
			assertCard(t, cards[len(cards)-2].ToDomain(), result.Cards[1])
		}
	})

	s.T().Run("他ユーザのDeckのCardsを取得しようとした場合", func(t *testing.T) {
		result, err := s.resolver.Query().Cards(ctx, &model.CardsInput{
			Where: &model.CardsWhere{
				DeckID: decks[1].ID,
			},
			Limit:  lo.ToPtr(2),
			Offset: lo.ToPtr(0),
		})
		assert.Nil(t, result)
		assertion.AssertError(t, "指定されたDeckのCardは取得できません", errutil.CodeForbidden, err)
	})
}

func (s *resolverSuite) TestCard() {
	userID := "test_user"
	ctx := context.Background()
	ctx = auth.New(&domain.User{ID: userID}).WithCtx(ctx)

	fx := fixture.New(s.dbx)
	decks := fx.CreateDeck(s.T(),
		&fixture.DeckInput{UserID: userID},
		&fixture.DeckInput{UserID: "other_user"},
	)
	cards := fx.CreateCard(s.T(), decks[0].ID, make([]fixture.CardInput, 1)...)
	cards2 := fx.CreateCard(s.T(), decks[1].ID, make([]fixture.CardInput, 1)...)

	s.T().Run("Cardsが取得できること", func(t *testing.T) {
		result, err := s.resolver.Query().Card(ctx, cards[0].ID)
		if assert.Nil(t, err) {
			assertCard(t, cards[0].ToDomain(), result)
		}
	})

	s.T().Run("存在しないカードの場合", func(t *testing.T) {
		result, err := s.resolver.Query().Card(ctx, -1)
		assert.Nil(t, result)
		assertion.AssertError(t, "Cardが見つかりませんでした", errutil.CodeNotFound, err)

	})

	s.T().Run("他ユーザのCardを取得しようとした場合", func(t *testing.T) {
		otherCard := cards2[0]
		result, err := s.resolver.Query().Card(ctx, otherCard.ID)
		assert.Nil(t, result)
		assertion.AssertError(t, "指定されたCardは取得できません", errutil.CodeForbidden, err)
	})
}
