package model

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/go-playground/validator/v10"
)

func (i *CreateCardInput) Validate(ctx context.Context) error {
	v := validate()

	type input struct {
		DeckID   int64   `json:"deckId" validate:"required"`
		Question string  `json:"question" validate:"required,max=1000"`
		Answer   *string `json:"answer,omitempty" validate:"omitempty,max=10000"`
	}
	err := v.StructCtx(ctx, input{
		DeckID:   i.DeckID,
		Question: i.Question,
		Answer:   i.Answer,
	})
	if err == nil {
		return nil
	}

	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return errutil.Wrap(err)
	}

	return errutil.New(errutil.CodeBadRequest, translateValidateError(errs[0]))
}
