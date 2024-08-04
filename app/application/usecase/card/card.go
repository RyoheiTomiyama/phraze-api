package card

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/infra/db"
	"github.com/RyoheiTomiyama/phraze-api/service/card"
)

type IUsecase interface {
	CreateCard(ctx context.Context, card *domain.Card) (*domain.Card, error)
	GetCard(ctx context.Context, id int64) (*domain.Card, error)
	GetCards(ctx context.Context, input domain.GetCardsInput) (*GetCardsOutput, error)
	GetPendingCards(ctx context.Context, input domain.GetPendingCardsInput) (*GetPendingCardsOutput, error)
	ReviewCard(ctx context.Context, id int64, grade int) error
	UpdateCard(ctx context.Context, id int64, input domain.UpdateCardInput) (*domain.Card, error)
}

type usecase struct {
	dbClient    db.IClient
	cardService card.ICardService
}

func New(dbClient db.IClient, cardService card.ICardService) IUsecase {
	return &usecase{dbClient, cardService}
}
