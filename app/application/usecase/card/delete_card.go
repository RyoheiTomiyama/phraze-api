package card

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/util/auth"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
)

func (u *usecase) DeleteCard(ctx context.Context, id int64) (int, error) {
	user := auth.FromCtx(ctx)

	card, err := u.dbClient.GetCard(ctx, id)
	if err != nil {
		return 0, errutil.Wrap(err)
	}
	if card == nil {
		return 0, errutil.New(errutil.CodeForbidden, "指定されたCardは取得できません")
	}

	deck, err := u.dbClient.GetDeck(ctx, card.DeckID)
	if err != nil {
		return 0, errutil.Wrap(err)
	}
	if deck != nil && deck.UserID != user.ID {
		return 0, errutil.New(errutil.CodeForbidden, "指定されたCardは取得できません")
	}

	ar, err := u.dbClient.DeleteCard(ctx, card.ID)
	if err != nil {
		return 0, errutil.Wrap(err)
	}

	return ar, nil
}
