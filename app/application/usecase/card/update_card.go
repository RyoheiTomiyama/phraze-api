package card

import (
	"context"
	"fmt"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/util/auth"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/RyoheiTomiyama/phraze-api/util/logger"
)

func (u *usecase) UpdateCard(ctx context.Context, id int64, input domain.UpdateCardInput) (*domain.Card, error) {
	_, _, err := u.getCardWithRoleCheck(ctx, id)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	card, err := u.dbClient.UpdateCardByID(ctx, id, &input)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	return card, nil
}

func (u *usecase) UpdateCardWithGendAnswer(ctx context.Context, id int64, input domain.UpdateCardInput) (*domain.Card, error) {
	log := logger.FromCtx(ctx)

	_, _, err := u.getCardWithRoleCheck(ctx, id)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	if input.Field.Question == nil || *input.Field.Question == "" {
		err = fmt.Errorf("Questionが見つからずAI生成ができない")
		log.Error(err, "id", id, "input", input)
		return nil, errutil.New(errutil.CodeBadRequest, err.Error())
	}

	ans, err := u.genemiClient.GenAnswer(ctx, *input.Field.Question)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	input.Field.Answer = &ans
	input.Field.AIAnswer = &ans

	card, err := u.dbClient.UpdateCardByID(ctx, id, &input)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	return card, nil
}

// idのカードが更新可能なものかチェックする
// ついでにCard,Deckを返す
func (u *usecase) getCardWithRoleCheck(ctx context.Context, id int64) (*domain.Card, *domain.Deck, error) {
	user := auth.FromCtx(ctx)

	card, err := u.dbClient.GetCard(ctx, id)
	if err != nil {
		return nil, nil, errutil.Wrap(err)
	}
	if card == nil {
		return nil, nil, errutil.New(errutil.CodeBadRequest, "指定されたカードが存在しません")
	}
	deck, err := u.dbClient.GetDeck(ctx, card.DeckID)
	if err != nil {
		return nil, nil, errutil.Wrap(err)
	}
	if deck == nil || deck.UserID != user.ID {
		return nil, nil, errutil.New(errutil.CodeBadRequest, "指定されたCardは更新できません")
	}

	return card, deck, nil
}
