package card

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/domain/infra/db"
	"github.com/RyoheiTomiyama/phraze-api/domain/infra/gemini"
	"github.com/RyoheiTomiyama/phraze-api/service/card"
)

type IUsecase interface {
	CreateCard(ctx context.Context, card *domain.Card) (*domain.Card, error)
	CreateCardWithGenAnswer(ctx context.Context, card *domain.Card) (*domain.Card, error)
	GetCard(ctx context.Context, id int64) (*domain.Card, error)
	GetCards(ctx context.Context, input domain.GetCardsInput) (*GetCardsOutput, error)
	GetPendingCards(ctx context.Context, input domain.GetPendingCardsInput) (*GetPendingCardsOutput, error)
	ReviewCard(ctx context.Context, id int64, grade int) error
	UpdateCard(ctx context.Context, id int64, input domain.UpdateCardInput) (*domain.Card, error)
	UpdateCardWithGendAnswer(ctx context.Context, id int64, input domain.UpdateCardInput) (*domain.Card, error)
	DeleteCard(ctx context.Context, id int64) (int64, error)

	ReadSchedules(ctx context.Context, cardIDs []int64) (map[int64]*domain.CardSchedule, error)
}

type usecase struct {
	dbClient     db.IClient
	geminiClient gemini.IClient
	cardService  card.ICardService
}

func New(dbClient db.IClient, geminiClient gemini.IClient, cardService card.ICardService) IUsecase {
	return &usecase{dbClient, geminiClient, cardService}
}
