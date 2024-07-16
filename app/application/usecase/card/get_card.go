package card

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/util/auth"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
)

func (u *usecase) GetCard(ctx context.Context, id int64) (*domain.Card, error) {
	user := auth.FromCtx(ctx)

	card, err := u.dbClient.GetCard(ctx, id)
	if err != nil {
		return nil, errutil.Wrap(err)
	}
	if card == nil {
		return nil, nil
	}

	deck, err := u.dbClient.GetDeck(ctx, card.DeckID)
	if err != nil {
		return nil, errutil.Wrap(err)
	}
	if deck != nil && deck.UserID != user.ID {
		return nil, errutil.New(errutil.CodeForbidden, "権限がありません")
	}

	return card, nil
}
