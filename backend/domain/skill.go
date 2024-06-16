package domain

import (
	"context"

	"github.com/google/uuid"
)

type Skill struct {
	ID          uuid.UUID
	Name        string
	Description string
}

type ISkillRepository interface {
	Create(context.Context, *Skill) error
	GetById(context.Context, uuid.UUID) (*Skill, error)
	GetAll(context.Context, int) ([]*Skill, int, error)
	Update(context.Context, *Skill) error
	DeleteById(context.Context, uuid.UUID) error
}

type ISkillService interface {
	Create(context.Context, *Skill) error
	GetById(context.Context, uuid.UUID) (*Skill, error)
	GetAll(context.Context, int) ([]*Skill, int, error)
	Update(context.Context, *Skill) error
	DeleteById(context.Context, uuid.UUID) error
}
