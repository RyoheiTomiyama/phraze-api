package model

import (
	"context"
	"errors"

	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/go-playground/validator/v10"
)

func (i *ReviewCardInput) Validate(ctx context.Context) error {
	v := validate()

	type input struct {
		CardID int64 `json:"cardId" validate:"required"`
		Grade  int   `json:"grade" validate:"min=1,max=5"`
	}
	err := v.StructCtx(ctx, input(*i))
	if err == nil {
		return nil
	}

	var errs validator.ValidationErrors
	if ok := errors.As(err, &errs); !ok {
		return errutil.Wrap(err)
	}
	if len(errs) == 0 {
		return nil
	}

	return errutil.New(errutil.CodeBadRequest, "%s", translateValidateError(errs[0]))
}
