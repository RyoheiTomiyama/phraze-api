package card

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/util/auth"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/samber/lo"
)

type CardsWhere struct {
	DeckID *int
}

type GetCardsOutput struct {
	Cards      []*domain.Card
	TotalCount int
}

func (u *usecase) GetCards(ctx context.Context, input domain.GetCardsInput) (*GetCardsOutput, error) {
	user := auth.FromCtx(ctx)

	deck, err := u.dbClient.GetDeck(ctx, input.Where.DeckID)
	if err != nil {
		return nil, errutil.Wrap(err)
	}
	if deck == nil || deck.UserID != user.ID {
		return nil, errutil.New(errutil.CodeForbidden, "指定されたDeckのCardは取得できません")
	}

	input.Where.UserID = lo.ToPtr(user.ID)

	cards, err := u.dbClient.GetCards(ctx, input.Where, *input.Limit, *input.Offset)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	count, err := u.dbClient.CountCards(ctx, input.Where)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	return &GetCardsOutput{
		Cards:      cards,
		TotalCount: count,
	}, nil
}
