package model

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/go-playground/validator/v10"
)

func (i *ReviewCardInput) Validate(ctx context.Context) error {
	v := validate()

	type input struct {
		CardID int64 `json:"cardId" validate:"required"`
		Grade  int   `json:"grade" validate:"required,min=1,max=5"`
	}
	err := v.StructCtx(ctx, input(*i))
	if err == nil {
		return nil
	}

	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return errutil.Wrap(err)
	}

	return errutil.New(errutil.CodeBadRequest, translateValidateError(errs[0]))
}
