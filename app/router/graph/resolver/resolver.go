package resolver

import (
	"github.com/RyoheiTomiyama/phraze-api/application/usecase/card"
	"github.com/RyoheiTomiyama/phraze-api/application/usecase/deck"
	"github.com/RyoheiTomiyama/phraze-api/router/graph/loader"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	/* loader */
	deckLoader loader.IDeckLoader

	/* usecase */
	cardUsecase card.IUsecase
	deckUsecase deck.IUsecase
}

func New(
	cardUsecase card.IUsecase,
	deckUsecase deck.IUsecase,
) *Resolver {
	deckLoader := loader.NewDeckLoader(cardUsecase)

	return &Resolver{
		deckLoader,
		cardUsecase,
		deckUsecase,
	}
}
