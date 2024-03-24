package domain

import (
	"context"
	"github.com/google/uuid"
)

type Company struct {
	ID              uuid.UUID
	Name            string
	FieldOfActivity string
	City            string
}

type ICompanyRepository interface { // FIXME: разобраться с сущностями
	Create(ctx context.Context, company *Company) error
	GetById(ctx context.Context, id uuid.UUID) (*Company, error)
	GetByOwnerId(ctx context.Context, id uuid.UUID) ([]*Company, error)
	GetAll(ctx context.Context) ([]*Company, error)
	Update(ctx context.Context, company *Company) error
	DeleteById(ctx context.Context, id uuid.UUID) error
}

type ICompanyService interface {
	Create(ctx context.Context, company *Company) error
	GetById(ctx context.Context, id uuid.UUID) (*Company, error)
	GetByOwnerId(ctx context.Context, id uuid.UUID) ([]*Company, error)
	GetAll(ctx context.Context) ([]*Company, error)
	Update(ctx context.Context, company *Company) error
	DeleteById(ctx context.Context, id uuid.UUID) error
}
