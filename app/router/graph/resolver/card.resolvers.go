package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"fmt"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/router/graph/model"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/samber/lo"
)

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

// UpdateCard is the resolver for the updateCard field.
func (r *mutationResolver) UpdateCard(ctx context.Context, input model.UpdateCardInput) (*model.UpdateCardOutput, error) {
	if err := input.Validate(ctx); err != nil {
		return nil, errutil.Wrap(err)
	}

	card, err := r.cardUsecase.UpdateCard(ctx, input.ID, domain.UpdateCardInput{
		Field: domain.UpdateCardField{
			Question: input.Question,
			Answer:   input.Answer,
		},
	})
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	var m model.Card
	if err = model.FromDomain(ctx, card, &m); err != nil {
		return nil, errutil.Wrap(err)
	}

	return &model.UpdateCardOutput{
		Card: &m,
	}, nil
}

// UpdateCardWithGenAnswer is the resolver for the updateCardWithGenAnswer field.
func (r *mutationResolver) UpdateCardWithGenAnswer(ctx context.Context, input model.UpdateCardInput) (*model.UpdateCardOutput, error) {
	panic(fmt.Errorf("not implemented: UpdateCardWithGenAnswer - updateCardWithGenAnswer"))
}

// Cards is the resolver for the cards field.
func (r *queryResolver) Cards(ctx context.Context, input *model.CardsInput) (*model.CardsOutput, error) {
	if err := input.Validate(ctx); err != nil {
		return nil, errutil.Wrap(err)
	}

	output, err := r.cardUsecase.GetCards(ctx, domain.GetCardsInput{
		Where: &domain.CardsWhere{
			DeckID: input.Where.DeckID,
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
			DeckID: input.Where.DeckID,
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
