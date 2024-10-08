package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/router/graph/model"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
)

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
func (r *mutationResolver) UpdateCardWithGenAnswer(ctx context.Context, input model.UpdateCardWithGenAnswerInput) (*model.UpdateCardWithGenAnswerOutput, error) {
	if err := input.Validate(ctx); err != nil {
		return nil, errutil.Wrap(err)
	}

	card, err := r.cardUsecase.UpdateCardWithGendAnswer(ctx, input.ID, domain.UpdateCardInput{
		Field: domain.UpdateCardField{
			Question: &input.Question,
		},
	})
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	var m model.Card
	if err = model.FromDomain(ctx, card, &m); err != nil {
		return nil, errutil.Wrap(err)
	}

	return &model.UpdateCardWithGenAnswerOutput{
		Card: &m,
	}, nil
}
