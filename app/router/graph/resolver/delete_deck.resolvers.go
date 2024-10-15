package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/router/graph/model"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
)

// DeleteDeck is the resolver for the deleteDeck field.
func (r *mutationResolver) DeleteDeck(ctx context.Context, input model.DeleteDeckInput) (*model.DeleteDeckOutput, error) {
	ar, err := r.deckUsecase.DeleteDeck(ctx, input.ID)
	if err != nil {
		return nil, errutil.Wrap(err)
	}

	return &model.DeleteDeckOutput{AffectedRows: int(ar)}, nil
}
