package card

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/util/auth"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
)

func (u *usecase) CreateCard(ctx context.Context, card *domain.Card) (*domain.Card, error) {
	user := auth.FromCtx(ctx)

	deck, err := u.dbClient.GetDeck(ctx, card.DeckID)
	if err != nil {
		return nil, errutil.Wrap(err)
	}
	if deck == nil || deck.UserID != user.ID {
		return nil, errutil.New(errutil.CodeBadRequest, "指定されたDeckのCardは作成できません")
	}

	card, err = u.dbClient.CreateCard(ctx, card)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	return card, nil
}
