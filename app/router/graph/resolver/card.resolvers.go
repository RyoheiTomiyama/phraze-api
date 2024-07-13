package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"fmt"

	"github.com/RyoheiTomiyama/phraze-api/router/graph/model"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
)

// CreateCard is the resolver for the createCard field.
func (r *mutationResolver) CreateCard(ctx context.Context, input *model.CreateCardInput) (*model.CreateCardOutput, error) {
	if err := input.Validate(ctx); err != nil {
		return nil, errutil.Wrap(err)
	}

	return nil, nil
}

// Cards is the resolver for the cards field.
func (r *queryResolver) Cards(ctx context.Context) ([]*model.Card, error) {
	panic(fmt.Errorf("not implemented: Cards - cards"))
}

// Card is the resolver for the card field.
func (r *queryResolver) Card(ctx context.Context) (*model.Card, error) {
	panic(fmt.Errorf("not implemented: Card - card"))
}
