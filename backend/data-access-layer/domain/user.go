package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

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
