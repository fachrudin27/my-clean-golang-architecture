package v1

import (
	"my-clean-architecture-template/internal/model"
)

func NewErrors(err []error) *[]model.ErrorData {
	out := make([]model.ErrorData, len(err))
	for i, v := range err {
		out[i] = model.ErrorData{
			Message: v.Error(),
		}
	}
	return &out
}
