package card

import (
	"context"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/util/auth"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/RyoheiTomiyama/phraze-api/util/logger"
	"github.com/samber/lo"
)

type GetPendingCardsOutput struct {
	Cards []*domain.Card
}

func (u *usecase) GetPendingCards(ctx context.Context, input domain.GetPendingCardsInput) (*GetPendingCardsOutput, error) {
	user := auth.FromCtx(ctx)
	log := logger.FromCtx(ctx)

	if input.Where.DeckID == nil {
		err := errutil.New(errutil.CodeBadRequest, "DeckIDを指定してください")
		log.ErrorWithNotify(ctx, err, "input", input)

		return nil, err
	}

	deckID := lo.FromPtr(input.Where.DeckID)

	deck, err := u.dbClient.GetDeck(ctx, deckID)
	if err != nil {
		return nil, errutil.Wrap(err)
	}
	if deck == nil || deck.UserID != user.ID {
		return nil, errutil.New(errutil.CodeForbidden, "指定されたDeckのCardは取得できません")
	}

	cards, err := u.dbClient.GetPendingCards(ctx, deckID, time.Now(), *input.Limit, *input.Offset)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	return &GetPendingCardsOutput{
		Cards: cards,
	}, nil
}
