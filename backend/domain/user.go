package domain

import (
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
	Create(user *User) error
	GetById(id uuid.UUID) (*User, error)
	GetAll() ([]*User, error)
	Update(user *User) error
	DeleteById(id uuid.UUID) error
}

type IUserService interface {
	Create(user *User) error
	GetById(id uuid.UUID) (*User, error)
	GetAll() ([]*User, error)
	Update(user *User) error
	DeleteById(id uuid.UUID) error
	GetFinancialReport(id uuid.UUID, period Period) (*FinancialReport, error)
	CalculateRating(id uuid.UUID) (float32, error)
}
