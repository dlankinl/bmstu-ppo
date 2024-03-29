package domain

import (
	"context"
	"github.com/google/uuid"
	"ppo/pkg/utils"
)

type Company struct {
	ID              uuid.UUID
	OwnerID         uuid.UUID
	ActivityFieldId uuid.UUID
	Name            string
	City            string
}

type ICompanyRepository interface { // FIXME: разобраться с сущностями
	Create(ctx context.Context, company *Company) error
	GetById(ctx context.Context, id uuid.UUID) (*Company, error)
	GetByOwnerId(ctx context.Context, id uuid.UUID) ([]*Company, error)
	GetAll(ctx context.Context, filters utils.Filters, page int) ([]*Company, error)
	Update(ctx context.Context, company *Company) error
	DeleteById(ctx context.Context, id uuid.UUID) error
}

type ICompanyService interface {
	Create(company *Company) error
	GetById(id uuid.UUID) (*Company, error)
	GetByOwnerId(id uuid.UUID) ([]*Company, error)
	GetAll(filters utils.Filters, page int) ([]*Company, error)
	Update(company *Company) error
	DeleteById(id uuid.UUID) error
}

// разобраться с фильтрами
