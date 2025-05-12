package repository

import (
	"context"
	"fmt"
	"my-clean-architecture-template/internal/model"
	"my-clean-architecture-template/pkg/postgres"
)

type UserRepo struct {
	*postgres.Postgres
}

// New -.
func New(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

// Login
func (r *UserRepo) Login(ctx context.Context, t model.LoginUserRequest) (string, error) {
	var name string

	sql := `select name from users where name=$1 limit 1`
	rows, err := r.Pool.Query(ctx, sql, t.Username)
	if err != nil {
		return "", fmt.Errorf("UserRepo - Login - r.Pool.QueryRow: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&name)
		if err != nil {
			return "", err
		}
	}

	return name, nil
}
