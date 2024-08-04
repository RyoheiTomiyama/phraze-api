package card

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/util/auth"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
)

func (u *usecase) UpdateCard(ctx context.Context, id int64, input domain.UpdateCardInput) (*domain.Card, error) {
	user := auth.FromCtx(ctx)

	card, err := u.dbClient.GetCard(ctx, id)
	if err != nil {
		return nil, errutil.Wrap(err)
	}
	if card == nil {
		return nil, errutil.New(errutil.CodeBadRequest, "指定されたカードが存在しません")
	}
	deck, err := u.dbClient.GetDeck(ctx, card.DeckID)
	if err != nil {
		return nil, errutil.Wrap(err)
	}
	if deck == nil || deck.UserID != user.ID {
		return nil, errutil.New(errutil.CodeBadRequest, "指定されたCardは更新できません")
	}

	card, err = u.dbClient.UpdateCardByID(ctx, id, &input)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	return card, nil
}
