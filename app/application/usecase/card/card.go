package card

import (
	"context"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/infra/db"
	"github.com/RyoheiTomiyama/phraze-api/service/card"
	"github.com/RyoheiTomiyama/phraze-api/util/auth"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
)

type IUsecase interface {
	CreateCard(ctx context.Context, card *domain.Card) (*domain.Card, error)
	GetCard(ctx context.Context, id int64) (*domain.Card, error)
	GetCards(ctx context.Context, input domain.GetCardsInput) (*GetCardsOutput, error)
	GetPendingCards(ctx context.Context, input domain.GetPendingCardsInput) (*GetPendingCardsOutput, error)
	ReviewCard(ctx context.Context, id int64, grade int) error
	UpdateCard(ctx context.Context, id int64, input domain.UpdateCardInput) (*domain.Card, error)

	ReadScheduleAt(ctx context.Context, deckIDs []int64) ([]*time.Time, error)
}

type usecase struct {
	dbClient    db.IClient
	cardService card.ICardService
}

func New(dbClient db.IClient, cardService card.ICardService) IUsecase {
	return &usecase{dbClient, cardService}
}

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
