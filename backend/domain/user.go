package domain

import (
	"context"
	"github.com/google/uuid"
	"ppo/pkg/utils"
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
	Create(ctx context.Context, user *User) error
	GetById(ctx context.Context, id uuid.UUID) (*User, error)
	GetAll(ctx context.Context, filters utils.Filters) ([]*User, error)
	Update(ctx context.Context, user *User) error
	DeleteById(ctx context.Context, id uuid.UUID) error
}

type IUserService interface {
	Create(ctx context.Context, user *User) error
	GetById(ctx context.Context, id uuid.UUID) (*User, error)
	GetAll(ctx context.Context, filters utils.Filters) ([]*User, error)
	Update(ctx context.Context, user *User) error
	DeleteById(ctx context.Context, id uuid.UUID) error
	GetFinancialReport(ctx context.Context, id uuid.UUID, period Period) ([]*FinancialReportByPeriod, error)
	CalculateRating(ctx context.Context, id uuid.UUID) (float32, error)
}
