package validator

import (
	"errors"
	"my-clean-architecture-template/internal/model"

	"github.com/go-playground/validator/v10"
)

func NewValidator() *validator.Validate {
	return validator.New()
}

func msgForTag(field string, fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return field + " field is required"
	case "email":
		return "invalid email"
	}
	return fe.Error()
}

func CustomError(err error) []model.ErrorData {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]model.ErrorData, len(ve))
		for i, v := range ve {
			out[i] = model.ErrorData{msgForTag(v.Field(), v)}
		}
		return out
	}
	return nil
}
