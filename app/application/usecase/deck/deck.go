package deck

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/infra/db"
	"github.com/RyoheiTomiyama/phraze-api/util/auth"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
)

type IUsecase interface {
	GetDeck(ctx context.Context, id int64) (*domain.Deck, error)
	GetDecks(ctx context.Context) ([]*domain.Deck, error)
}

type usecase struct {
	dbClient db.IClient
}

func New(dbClient db.IClient) IUsecase {
	return &usecase{dbClient}
}

func (u *usecase) GetDeck(ctx context.Context, id int64) (*domain.Deck, error) {
	deck, err := u.dbClient.GetDeck(ctx, id)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	return deck, nil
}

func (u *usecase) GetDecks(ctx context.Context) ([]*domain.Deck, error) {
	user := auth.FromCtx(ctx)

	decks, err := u.dbClient.GetDecks(ctx, user.ID)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	return decks, nil
}
