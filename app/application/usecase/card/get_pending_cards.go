package card

import (
	"context"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/util/auth"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
)

type GetPendingCardsOutput struct {
	Cards []*domain.Card
}

func (u *usecase) GetPendingCards(ctx context.Context, input domain.GetPendingCardsInput) (*GetPendingCardsOutput, error) {
	user := auth.FromCtx(ctx)

	deck, err := u.dbClient.GetDeck(ctx, input.Where.DeckID)
	if err != nil {
		return nil, errutil.Wrap(err)
	}
	if deck == nil || deck.UserID != user.ID {
		return nil, errutil.New(errutil.CodeForbidden, "指定されたDeckのCardは取得できません")
	}

	cards, err := u.dbClient.GetPendingCards(ctx, input.Where.DeckID, time.Now(), *input.Limit, *input.Offset)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	return &GetPendingCardsOutput{
		Cards: cards,
	}, nil
}
