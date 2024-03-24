package domain

import (
	"context"
	"github.com/google/uuid"
)

type UserAuth struct {
	ID           uuid.UUID
	Username     string
	PasswordHash string
}

type IAuthRepository interface {
	Register(ctx context.Context, username, password string) (err error)
	GetByUsername(ctx context.Context, username string) (*UserAuth, error)
}

type IAuthService interface {
	Login(ctx context.Context, username, password string) (string, error)
	Register(ctx context.Context, username, password string) (err error)
}
