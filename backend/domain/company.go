package domain

import "github.com/google/uuid"

type Company struct {
	ID              uuid.UUID
	Name            string
	FieldOfActivity string
	City            string
}

type ICompanyRepository interface { // FIXME: разобраться с сущностями
	Create(company *Company) error
	GetById(id uuid.UUID) (*Company, error)
	GetByOwnerId(id uuid.UUID) ([]*Company, error)
	GetAll() ([]*Company, error)
	Update(company *Company) error
	DeleteById(id uuid.UUID) error
}

type ICompanyService interface {
	Create(company *Company) error
	GetById(id uuid.UUID) (*Company, error)
	GetByOwnerId(id uuid.UUID) ([]*Company, error)
	GetAll() ([]*Company, error)
	Update(company *Company) error
	DeleteById(id uuid.UUID) error
	GetFinancialReport(id uuid.UUID, period Period) (*FinancialReport, error)
}
