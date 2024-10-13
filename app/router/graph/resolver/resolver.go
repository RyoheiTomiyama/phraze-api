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
	cardLoader loader.ICardLoader
	deckLoader loader.IDeckLoader

	/* usecase */
	cardUsecase card.IUsecase
	deckUsecase deck.IUsecase
}

func New(
	cardUsecase card.IUsecase,
	deckUsecase deck.IUsecase,
) *Resolver {
	cardLoader := loader.NewCardLoader(cardUsecase)
	deckLoader := loader.NewDeckLoader(deckUsecase)

	return &Resolver{
		cardLoader,
		deckLoader,
		cardUsecase,
		deckUsecase,
	}
}
