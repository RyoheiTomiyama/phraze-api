package model

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
)

func (i *CreateCardInput) Validate(ctx context.Context) error {
	v := validate()

	if err := v.VarCtx(ctx, i.DeckID, "required"); err != nil {
		return errutil.New(errutil.CodeBadRequest, err.Error())
	}
	if err := v.VarCtx(ctx, i.Question, "required,max=1000"); err != nil {
		return errutil.New(errutil.CodeBadRequest, err.Error())
	}
	if err := v.VarCtx(ctx, i.Answer, "required,max=10000"); err != nil {
		return errutil.New(errutil.CodeBadRequest, err.Error())
	}

	return nil
}
