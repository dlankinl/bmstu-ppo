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
}

type IUserRepository interface {
	GetById(ctx context.Context, id uuid.UUID) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
	Update(ctx context.Context, user *User) error
	DeleteById(ctx context.Context, id uuid.UUID) error
}

type IUserService interface {
	GetById(ctx context.Context, id uuid.UUID) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
	Update(ctx context.Context, user *User) error
	DeleteById(ctx context.Context, id uuid.UUID) error
	GetFinancialReport(ctx context.Context, id uuid.UUID, period Period) (*FinancialReport, error)
	CalculateRating(ctx context.Context, id uuid.UUID) (float32, error)
}
