package card

import (
	"context"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/util/auth"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
)

func (u *usecase) ReviewCard(ctx context.Context, id int64, grade int) error {
	user := auth.FromCtx(ctx)

	card, err := u.dbClient.GetCard(ctx, id)
	if err != nil {
		return errutil.Wrap(err)
	}
	if card == nil {
		return errutil.New(errutil.CodeBadRequest, "Cardが見つかりませんでした")
	}

	deck, err := u.dbClient.GetDeck(ctx, card.DeckID)
	if err != nil {
		return errutil.Wrap(err)
	}
	if deck == nil || deck.UserID != user.ID {
		return errutil.New(errutil.CodeBadRequest, "指定されたCardのレビューはできません")
	}

	schedule, err := u.dbClient.GetCardSchedule(ctx, id)
	if err != nil {
		return errutil.Wrap(err)
	}

	schedule = u.cardService.EvalSchedule(ctx, grade, schedule)

	if err = u.dbClient.Tx(ctx, func(ctx context.Context) error {
		if _, err := u.dbClient.UpsertCardReview(ctx, &domain.CardReview{
			CardID:     id,
			ReviewedAt: time.Now(),
			Grade:      grade,
		}); err != nil {
			return nil
		}

		if _, err = u.dbClient.UpsertCardSchedule(ctx, schedule); err != nil {
			return nil
		}

		return nil
	}); err != nil {
		return errutil.Wrap(err)
	}

	return nil
}
