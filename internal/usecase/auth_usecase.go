package usecase

import (
	"context"
	"fmt"
	"my-clean-architecture-template/config"
	"my-clean-architecture-template/internal/gateway/messaging"
	"my-clean-architecture-template/internal/model"
	"my-clean-architecture-template/pkg/jwt"
	"net/http"
)

type TranslationUseCase struct {
	repo AuthRepo
	cfg  config.Config
	// pg   *postgres.Postgres
}

func New(cfg config.Config, repo AuthRepo) *TranslationUseCase {
	return &TranslationUseCase{
		cfg:  cfg,
		repo: repo,
		// pg:   pg,
	}
}

func (uc *TranslationUseCase) Login(ctx context.Context, payload model.LoginUserRequest) (model.UserResponse, int, []error) {
	var res model.UserResponse
	var errors []error
	// var p *messaging.Producer

	// tx, err := uc.pg.Pool.Begin(ctx)
	// if err != nil {
	// 	return res, http.StatusInternalServerError, append(errors, fmt.Errorf("internal server error"))
	// }
	// defer tx.Rollback(ctx)

	ss, err := uc.repo.Login(context.Background(), payload)
	if err != nil {
		return res, http.StatusInternalServerError, append(errors, fmt.Errorf("internal server error"))
	}

	token, err := jwt.GenerateToken(uc.cfg.JWT.SecretKey)
	if err != nil {
		return res, http.StatusInternalServerError, append(errors, fmt.Errorf("internal server error"))
	}

	res.ID = "1"
	res.Name = ss
	res.Token = token

	err = messaging.LoginProducer(payload)
	if err != nil {
		return res, http.StatusInternalServerError, append(errors, fmt.Errorf("internal server error"))
	}

	// err = tx.Commit(ctx)
	// if err != nil {
	// 	return res, http.StatusInternalServerError, append(errors, fmt.Errorf("internal server error"))
	// }

	return res, http.StatusOK, nil
}

func (uc *TranslationUseCase) Users(ctx context.Context) (string, error) {
	return "masok", nil
}
