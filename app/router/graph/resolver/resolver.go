package resolver

import (
	"github.com/RyoheiTomiyama/phraze-api/application/usecase/card"
	"github.com/RyoheiTomiyama/phraze-api/application/usecase/deck"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	cardUsecase card.IUsecase
	deckUsecase deck.IUsecase
}

func New(
	cardUsecase card.IUsecase,
	deckUsecase deck.IUsecase,
) *Resolver {
	return &Resolver{
		cardUsecase,
		deckUsecase,
	}
}
