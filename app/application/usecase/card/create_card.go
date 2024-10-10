package card

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/util/auth"
	ctxutil "github.com/RyoheiTomiyama/phraze-api/util/context"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/RyoheiTomiyama/phraze-api/util/logger"
)

func (u *usecase) CreateCard(ctx context.Context, card *domain.Card) (*domain.Card, error) {
	user := auth.FromCtx(ctx)

	permissions, err := u.dbClient.GetPermissionsByUserID(ctx, user.ID)
	if err != nil {
		return nil, errutil.Wrap(err)
	}
	// 無制限の権限ない場合は、合計1000個までしかカードを作成できない
	if !domain.Permissions(permissions).HasKey(ctx, domain.PermissionUnlimitedCardCreation) {
		c, err := u.dbClient.CountCards(ctx, &domain.CardsWhere{
			UserID: &user.ID,
		})
		if err != nil {
			return nil, errutil.Wrap(err)
		}
		if c > 999 {
			return nil, errutil.New(errutil.CodeForbidden, "カードは合計1000件まで作成可能です")
		}
	}

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

func (u *usecase) CreateCardWithGenAnswer(ctx context.Context, card *domain.Card) (*domain.Card, error) {
	card, err := u.CreateCard(ctx, card)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	ctx, _ = ctxutil.AsyncContext(ctx)
	go func() {
		log := logger.FromCtx(ctx)
		if card.Question == "" {
			return
		}
		answer, err := u.geminiClient.GenAnswer(ctx, card.Question)
		if err != nil {
			log.ErrorWithNotify(ctx, err, "card", card)
		}

		if _, err = u.dbClient.UpdateCardByID(ctx, card.ID, &domain.UpdateCardInput{
			Field: domain.UpdateCardField{
				AIAnswer: &answer,
			},
		}); err != nil {
			log.ErrorWithNotify(ctx, err)
		}
	}()

	return card, nil
}
