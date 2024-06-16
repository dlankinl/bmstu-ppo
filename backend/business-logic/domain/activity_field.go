package domain

import (
	"context"

	"github.com/google/uuid"
)

//go:generate mockgen -source=activity_field.go -destination=mocks/activity_field.go -package=mocks

type ActivityField struct {
	ID          uuid.UUID
	Name        string
	Description string
	Cost        float32
}

type IActivityFieldRepository interface {
	Create(context.Context, *ActivityField) error
	DeleteById(context.Context, uuid.UUID) error
	Update(context.Context, *ActivityField) error
	GetById(context.Context, uuid.UUID) (*ActivityField, error)
	GetMaxCost(context.Context) (float32, error)
	GetAll(context.Context, int) ([]*ActivityField, error)
}

type IActivityFieldService interface {
	Create(*ActivityField) error
	DeleteById(uuid.UUID) error
	Update(*ActivityField) error
	GetById(uuid.UUID) (*ActivityField, error)
	GetCostByCompanyId(uuid.UUID) (float32, error)
	GetMaxCost() (float32, error)
	GetAll(page int) ([]*ActivityField, error)
}
