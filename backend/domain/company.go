package domain

import (
	"context"
	"github.com/google/uuid"
)

type Company struct {
	ID              uuid.UUID
	OwnerID         uuid.UUID
	ActivityFieldId uuid.UUID
	Name            string
	City            string
}

type ICompanyRepository interface {
	Create(context.Context, *Company) error
	GetById(context.Context, uuid.UUID) (*Company, error)
	GetByOwnerId(context.Context, uuid.UUID, int, bool) ([]*Company, int, error)
	GetAll(context.Context, int) ([]*Company, error)
	Update(context.Context, *Company) error
	DeleteById(context.Context, uuid.UUID) error
}

type ICompanyService interface {
	Create(context.Context, *Company) error
	GetById(context.Context, uuid.UUID) (*Company, error)
	GetByOwnerId(context.Context, uuid.UUID, int, bool) ([]*Company, int, error)
	GetAll(context.Context, int) ([]*Company, error)
	Update(context.Context, *Company) error
	DeleteById(context.Context, uuid.UUID) error
}
