package model

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var val *validator.Validate

func validate() *validator.Validate {
	if val == nil {
		val = validator.New()
	}

	return val
}

func translateValidateError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%sは必須項目です", e.StructField())
	case "max":
		return fmt.Sprintf("%sは%sが最大です", e.StructField(), e.Param())
	}

	return e.Error()
}
