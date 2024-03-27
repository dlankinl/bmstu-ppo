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
	Create(ctx context.Context, skill *Skill) error
	GetById(ctx context.Context, id uuid.UUID) (*Skill, error)
	GetAll(ctx context.Context) ([]*Skill, error)
	Update(ctx context.Context, skill *Skill) error
	DeleteById(ctx context.Context, id uuid.UUID) error
}

type ISkillService interface {
	Create(skill *Skill) error
	GetById(id uuid.UUID) (*Skill, error)
	GetAll() ([]*Skill, error)
	Update(skill *Skill) error
	DeleteById(id uuid.UUID) error
}
