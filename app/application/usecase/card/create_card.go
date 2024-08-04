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
		answer, err := u.genemiClient.GenAnswer(ctx, card.Question)
		if err != nil {
			log.Error(err, "card", card)
		}

		log.Debug(answer)
	}()

	return card, nil
}
