package model

import "github.com/go-playground/validator/v10"

var val *validator.Validate

func validate() *validator.Validate {
	if val == nil {
		val = validator.New()
	}

	return val
}
