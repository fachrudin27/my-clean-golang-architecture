package usecase

import (
	"context"
	"my-clean-architecture-template/internal/model"
)

type (
	Auth interface {
		Login(context.Context, model.LoginUserRequest) (model.UserResponse, int, []error)
		Users(context.Context) (string, error)
	}

	AuthRepo interface {
		Login(context.Context, model.LoginUserRequest) (string, error)
	}
)
