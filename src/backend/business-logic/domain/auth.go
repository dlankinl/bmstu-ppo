package domain

import (
	"context"
	"github.com/google/uuid"
)

//go:generate mockgen -source=auth.go -destination=mocks/auth.go -package=mocks

type UserAuth struct {
	ID         uuid.UUID
	Username   string
	Password   string
	HashedPass string
	Role       string
}

type IAuthRepository interface {
	Register(ctx context.Context, authInfo *UserAuth) (err error)
	GetByUsername(ctx context.Context, username string) (*UserAuth, error)
}

type IAuthService interface {
	Login(authInfo *UserAuth) (string, error)
	Register(authInfo *UserAuth) (err error)
}
