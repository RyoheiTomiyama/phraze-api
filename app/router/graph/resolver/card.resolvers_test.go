package resolver_test

import (
	"context"
	"testing"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/infra/db/fixture"
	dbModel "github.com/RyoheiTomiyama/phraze-api/infra/db/model"
	"github.com/RyoheiTomiyama/phraze-api/router/graph/model"
	"github.com/RyoheiTomiyama/phraze-api/test/assertion"
	"github.com/RyoheiTomiyama/phraze-api/util/auth"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func assertCard(t *testing.T, expect *domain.Card, actual *model.Card) {
	t.Helper()

	em := &model.Card{
		ID:        expect.ID,
		DeckID:    expect.DeckID,
		Question:  expect.Question,
		Answer:    expect.Answer,
		AiAnswer:  expect.AIAnswer,
		CreatedAt: expect.CreatedAt,
		UpdatedAt: expect.UpdatedAt,
	}
	assert.Equal(t, em, actual)
}

func (s *resolverSuite) TestCardResolverSchedule() {
	userID := "test_user"
	ctx := context.Background()
	ctx = auth.New(&domain.User{ID: userID}).WithCtx(ctx)

	fx := fixture.New(s.dbx)
	decks := fx.CreateDeck(s.T(), &fixture.DeckInput{UserID: lo.ToPtr(userID)})
	cards := fx.CreateCard(s.T(), decks[0].ID, make([]fixture.CardInput, 2)...)
	schedules := fx.CreateCardSchedule(s.T(), []fixture.CardScheduleInput{
		{CardID: cards[0].ID, ScheduleAt: time.Now().Add(-10 * time.Hour)},
	}...)

	s.T().Run("スケジュールがある場合", func(t *testing.T) {
		result, err := s.resolver.Card().Schedule(ctx, &model.Card{ID: cards[0].ID})

		assert.NoError(t, err)
		assert.Equal(t, schedules[0].ID, result.ID)
	})
	s.T().Run("スケジュールがない場合", func(t *testing.T) {
		result, err := s.resolver.Card().Schedule(ctx, &model.Card{ID: cards[1].ID})

		assert.NoError(t, err)
		assert.Nil(t, result)
	})
}

func (s *resolverSuite) TestCreateCard() {
	userID := "test_user"
	ctx := context.Background()
	ctx = auth.New(&domain.User{ID: userID}).WithCtx(ctx)

	fx := fixture.New(s.dbx)
	decks := fx.CreateDeck(s.T(), &fixture.DeckInput{UserID: lo.ToPtr(userID)})

	s.geminiClient.On("GenAnswer", mock.Anything, "question").Return("ai-answer", nil)

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

		t.Run("非同期でAIAnswerが保存されていること", func(t *testing.T) {
			time.Sleep(100 * time.Millisecond)
			var card dbModel.Card
			query := "SELECT * FROM cards WHERE id=$1"
			if err = s.dbx.Get(&card, query, result.Card.ID); err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, "ai-answer", card.AIAnswer)
		})
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
		decks := fx.CreateDeck(s.T(), &fixture.DeckInput{UserID: lo.ToPtr("other_user")})

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
		&fixture.DeckInput{UserID: lo.ToPtr(userID)},
		&fixture.DeckInput{UserID: lo.ToPtr("other_user")},
	)
	cardsInput := make([]fixture.CardInput, 9)
	cardsInput = append(cardsInput, fixture.CardInput{Question: lo.ToPtr("qwertyuiop")})
	cards := fx.CreateCard(s.T(), decks[0].ID, cardsInput...)
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
	s.T().Run("Cardsを検索できること", func(t *testing.T) {
		result, err := s.resolver.Query().Cards(ctx, &model.CardsInput{
			Where: &model.CardsWhere{
				DeckID: decks[0].ID,
				Q:      lo.ToPtr("ertyuio"),
			},
			Limit:  lo.ToPtr(2),
			Offset: lo.ToPtr(0),
		})
		if assert.Nil(t, err) {
			assert.Len(t, result.Cards, 1)
			assertCard(t, cards[9].ToDomain(), result.Cards[0])
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
		&fixture.DeckInput{UserID: lo.ToPtr(userID)},
		&fixture.DeckInput{UserID: lo.ToPtr("other_user")},
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

func (s *resolverSuite) TestPendingCards() {
	userID := "test_user"
	ctx := context.Background()
	ctx = auth.New(&domain.User{ID: userID}).WithCtx(ctx)

	fx := fixture.New(s.dbx)
	decks := fx.CreateDeck(s.T(),
		&fixture.DeckInput{UserID: lo.ToPtr(userID)},
	)
	cards := fx.CreateCard(s.T(), decks[0].ID, make([]fixture.CardInput, 5)...)
	fx.CreateCardSchedule(s.T(), []fixture.CardScheduleInput{
		{CardID: cards[2].ID, ScheduleAt: time.Now().Add(-10 * time.Hour)},
		{CardID: cards[1].ID, ScheduleAt: time.Now().Add(-1 * time.Hour)},
		// 学習済みのカード
		{CardID: cards[4].ID, ScheduleAt: time.Now().Add(10 * time.Hour)},
	}...)

	s.T().Run("Validationエラー", func(t *testing.T) {
		testCases := []struct {
			name   string
			input  model.PendingCardsInput
			assert func(err error)
		}{
			{
				name:  "DeckIDが0値の場合",
				input: model.PendingCardsInput{Where: &model.PendingCardsWhere{}},
				assert: func(err error) {
					assertion.AssertError(t, "DeckIDは必須項目です", errutil.CodeBadRequest, err)
				},
			},
			{
				name:  "Limitが100より大きい場合",
				input: model.PendingCardsInput{Where: &model.PendingCardsWhere{DeckID: 100}, Limit: lo.ToPtr(101)},
				assert: func(err error) {
					assertion.AssertError(t, "Limitは100が最大です", errutil.CodeBadRequest, err)
				},
			},
		}

		for _, tc := range testCases {
			result, err := s.resolver.Query().PendingCards(ctx, &tc.input)
			assert.Nil(t, result)
			tc.assert(err)
		}
	})

	s.T().Run("未学習のCardsがSchedule古い順に取得できること", func(t *testing.T) {
		result, err := s.resolver.Query().PendingCards(ctx, &model.PendingCardsInput{
			Where: &model.PendingCardsWhere{
				DeckID: decks[0].ID,
			},
			Limit:  lo.ToPtr(100),
			Offset: lo.ToPtr(0),
		})
		if assert.Nil(t, err) {
			assert.Len(t, result.Cards, 4)
			assertCard(t, cards[0].ToDomain(), result.Cards[0])
			assertCard(t, cards[3].ToDomain(), result.Cards[1])
			assertCard(t, cards[2].ToDomain(), result.Cards[2])
			assertCard(t, cards[1].ToDomain(), result.Cards[3])
		}
	})
}
