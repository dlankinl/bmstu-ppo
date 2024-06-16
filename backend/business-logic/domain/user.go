package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

//go:generate mockgen -source=user.go -destination=mocks/user.go -package=mocks

type User struct {
	ID       uuid.UUID
	Username string
	FullName string
	Gender   string
	Birthday time.Time
	City     string
	Role     string
}

type IUserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetById(ctx context.Context, userId uuid.UUID) (*User, error)
	GetAll(ctx context.Context, page int) ([]*User, error)
	Update(ctx context.Context, user *User) error
	DeleteById(ctx context.Context, id uuid.UUID) error
}

type IUserService interface {
	Create(user *User) error
	GetByUsername(username string) (*User, error)
	GetById(userId uuid.UUID) (*User, error)
	GetAll(page int) ([]*User, error)
	Update(user *User) error
	DeleteById(id uuid.UUID) error
}
