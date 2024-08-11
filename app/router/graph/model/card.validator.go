package model

import (
	"context"
	"errors"

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

	var errs validator.ValidationErrors
	if ok := errors.As(err, &errs); !ok {
		return errutil.Wrap(err)
	}
	if len(errs) == 0 {
		return nil
	}

	return errutil.New(errutil.CodeBadRequest, translateValidateError(errs[0]))
}

func (i *CardsInput) Validate(ctx context.Context) error {
	v := validate()

	type where struct {
		DeckID int64 `json:"deckId" validate:"required"`
	}
	type input struct {
		Where  *where `json:"where" validate:"required"`
		Limit  *int   `json:"limit,omitempty" validate:"omitempty,max=100"`
		Offset *int   `json:"offset,omitempty"`
	}

	w := where(*i.Where)
	err := v.StructCtx(ctx, input{
		Where:  &w,
		Limit:  i.Limit,
		Offset: i.Offset,
	})
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

	return errutil.New(errutil.CodeBadRequest, translateValidateError(errs[0]))
}

func (i *UpdateCardInput) Validate(ctx context.Context) error {
	v := validate()

	if i.Answer == nil && i.Question == nil {
		return errutil.New(errutil.CodeBadRequest, "更新するフィールドを指定してください")
	}

	type input struct {
		ID       int64   `json:"ID" validate:"required"`
		Question *string `json:"question" validate:"omitempty,max=1000"`
		Answer   *string `json:"answer,omitempty" validate:"omitempty,max=10000"`
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

	return errutil.New(errutil.CodeBadRequest, translateValidateError(errs[0]))
}
func (i *UpdateCardWithGenAnswerInput) Validate(ctx context.Context) error {
	v := validate()

	type input struct {
		ID       int64  `json:"ID" validate:"required"`
		Question string `json:"question" validate:"required,max=1000"`
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

	return errutil.New(errutil.CodeBadRequest, translateValidateError(errs[0]))
}

func (i *PendingCardsInput) Validate(ctx context.Context) error {
	v := validate()

	type where struct {
		DeckID int64 `json:"deckId" validate:"required"`
	}
	type input struct {
		Where  *where `json:"where" validate:"required"`
		Limit  *int   `json:"limit,omitempty" validate:"omitempty,max=100"`
		Offset *int   `json:"offset,omitempty"`
	}

	w := where(*i.Where)
	err := v.StructCtx(ctx, input{
		Where:  &w,
		Limit:  i.Limit,
		Offset: i.Offset,
	})
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

	return errutil.New(errutil.CodeBadRequest, translateValidateError(errs[0]))
}
