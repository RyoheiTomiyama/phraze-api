package deck

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/util/auth"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
)

func (u *usecase) DeleteDeck(ctx context.Context, id int64) (int64, error) {
	user := auth.FromCtx(ctx)

	deck, err := u.dbClient.GetDeck(ctx, id)
	if err != nil {
		return 0, errutil.Wrap(err)
	}
	if deck == nil || deck.UserID != user.ID {
		return 0, errutil.New(errutil.CodeForbidden, "指定されたDeckは取得できません")
	}

	var ar int64
	if err := u.dbClient.Tx(ctx, func(ctx context.Context) error {
		ar, err = u.dbClient.DeleteDeck(ctx, deck.ID)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return 0, errutil.Wrap(err)
	}

	return ar, nil
}
