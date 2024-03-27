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
	Create(user *User) error
	GetById(id uuid.UUID) (*User, error)
	GetAll(filters utils.Filters) ([]*User, error)
	Update(user *User) error
	DeleteById(id uuid.UUID) error
	GetFinancialReport(companies []*Company, period *Period) (finReports []*FinancialReportByPeriod, err error)
	//GetFinancialReport(ctx context.Context, id uuid.UUID, period *Period) ([]*FinancialReportByPeriod, error)
	//CalculateRating(ctx context.Context, id uuid.UUID, mainFieldWeight float32) (rating float32, err error)
}
