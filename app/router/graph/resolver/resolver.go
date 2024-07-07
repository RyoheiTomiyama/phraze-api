package resolver

import (
	"github.com/RyoheiTomiyama/phraze-api/application/usecase/deck"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	deckUsecase deck.IUsecase
}

func New(
	deckUsecase deck.IUsecase,
) *Resolver {
	return &Resolver{
		deckUsecase,
	}
}
