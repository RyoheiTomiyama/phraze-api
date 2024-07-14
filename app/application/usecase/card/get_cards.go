package card

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
)

type CardsWhere struct {
	DeckID *int
}

type GetCardsOutput struct {
	Cards []*domain.Card
}

func (u *usecase) GetCards(ctx context.Context, input domain.GetCardsInput) (*GetCardsOutput, error) {
	return nil, nil
}
