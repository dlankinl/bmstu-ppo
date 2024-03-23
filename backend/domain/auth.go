package domain

import (
	"github.com/google/uuid"
)

type UserAuth struct {
	ID           uuid.UUID
	Username     string
	PasswordHash string
}

type IAuthRepository interface {
	Login()
	Register(*UserAuth) error
	GetByUsername(username string) (*UserAuth, error)
}

type IAuthService interface {
	Login()
	Register(user *UserAuth) error
}
