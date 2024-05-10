package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
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
	Create(context.Context, *User) error
	GetByUsername(context.Context, string) (*User, error)
	GetById(context.Context, uuid.UUID) (*User, error)
	GetAll(context.Context, int) ([]*User, error)
	Update(context.Context, *User) error
	DeleteById(context.Context, uuid.UUID) error
}

type IUserService interface {
	Create(context.Context, *User) error
	GetByUsername(context.Context, string) (*User, error)
	GetById(context.Context, uuid.UUID) (*User, error)
	GetAll(context.Context, int) ([]*User, error)
	Update(context.Context, *User) error
	DeleteById(context.Context, uuid.UUID) error
}
