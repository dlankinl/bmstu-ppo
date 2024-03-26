package domain

import (
	"context"
	"github.com/google/uuid"
)

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
}

type IActivityFieldService interface {
	Create(context.Context, *ActivityField) error
	DeleteById(context.Context, uuid.UUID) error
	Update(context.Context, *ActivityField) error
	GetById(context.Context, uuid.UUID) (*ActivityField, error)
}
