package domain

import (
	"context"
	"github.com/google/uuid"
)

type UserAuth struct {
	ID         uuid.UUID
	Username   string
	Password   string
	HashedPass string
}

type IAuthRepository interface {
	Register(ctx context.Context, authInfo *UserAuth) (err error)
	GetByUsername(ctx context.Context, username string) (*UserAuth, error)
}

type IAuthService interface {
	Login(authInfo *UserAuth) (string, error)
	Register(authInfo *UserAuth) (err error)
}
