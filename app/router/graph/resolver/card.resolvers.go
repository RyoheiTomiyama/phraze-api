package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"fmt"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/router/graph/generated"
	"github.com/RyoheiTomiyama/phraze-api/router/graph/model"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/samber/lo"
)

// Schedule is the resolver for the schedule field.
func (r *cardResolver) Schedule(ctx context.Context, obj *model.Card) (*model.CardSchedule, error) {
	panic(fmt.Errorf("not implemented: Schedule - schedule"))
}

// CreateCard is the resolver for the createCard field.
func (r *mutationResolver) CreateCard(ctx context.Context, input model.CreateCardInput) (*model.CreateCardOutput, error) {
	if err := input.Validate(ctx); err != nil {
		return nil, errutil.Wrap(err)
	}

	card, err := r.cardUsecase.CreateCardWithGenAnswer(ctx, &domain.Card{
		DeckID:   input.DeckID,
		Question: input.Question,
		Answer: lo.TernaryF(
			input.Answer != nil,
			func() string { return *input.Answer },
			func() string { return "" },
		),
	})
	if err != nil {
		return nil, errutil.Wrap(err)
	}
	if card == nil {
		return nil, errutil.New(errutil.CodeInternalError, "cardの作成に失敗: deck_id: %d", input.DeckID)
	}

	var m model.Card
	if err = model.FromDomain(ctx, card, &m); err != nil {
		return nil, errutil.Wrap(err)
	}

	return &model.CreateCardOutput{
		Card: &m,
	}, nil
}

// Cards is the resolver for the cards field.
func (r *queryResolver) Cards(ctx context.Context, input *model.CardsInput) (*model.CardsOutput, error) {
	if err := input.Validate(ctx); err != nil {
		return nil, errutil.Wrap(err)
	}

	output, err := r.cardUsecase.GetCards(ctx, domain.GetCardsInput{
		Where: &domain.CardsWhere{
			DeckID: lo.ToPtr(input.Where.DeckID),
		},
		Limit:  input.Limit,
		Offset: input.Offset,
	})
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	var cards []*model.Card
	for _, item := range output.Cards {
		var m model.Card
		if err = model.FromDomain(ctx, item, &m); err != nil {
			return nil, errutil.Wrap(err)
		}

		cards = append(cards, &m)
	}

	return &model.CardsOutput{
		Cards: cards,
		PageInfo: &model.PageInfo{
			TotalCount: output.TotalCount,
		},
	}, nil
}

// Card is the resolver for the card field.
func (r *queryResolver) Card(ctx context.Context, id int64) (*model.Card, error) {
	card, err := r.cardUsecase.GetCard(ctx, id)
	if err != nil {
		return nil, errutil.Wrap(err)
	}
	if card == nil {
		return nil, errutil.New(errutil.CodeNotFound, "Cardが見つかりませんでした")
	}

	var m model.Card
	if err = model.FromDomain(ctx, card, &m); err != nil {
		return nil, errutil.Wrap(err)
	}

	return &m, nil
}

// PendingCards is the resolver for the pendingCards field.
func (r *queryResolver) PendingCards(ctx context.Context, input *model.PendingCardsInput) (*model.PendingCardsOutput, error) {
	if err := input.Validate(ctx); err != nil {
		return nil, errutil.Wrap(err)
	}

	output, err := r.cardUsecase.GetPendingCards(ctx, domain.GetPendingCardsInput{
		Where: &domain.CardsWhere{
			DeckID: lo.ToPtr(input.Where.DeckID),
		},
		Limit:  input.Limit,
		Offset: input.Offset,
	})
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	var cards []*model.Card
	for _, item := range output.Cards {
		var m model.Card
		if err = model.FromDomain(ctx, item, &m); err != nil {
			return nil, errutil.Wrap(err)
		}

		cards = append(cards, &m)
	}

	return &model.PendingCardsOutput{
		Cards: cards,
	}, nil
}

// Card returns generated.CardResolver implementation.
func (r *Resolver) Card() generated.CardResolver { return &cardResolver{r} }

type cardResolver struct{ *Resolver }
