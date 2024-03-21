package domain

import "github.com/google/uuid"

type Skill struct {
	ID          uuid.UUID
	Name        string
	Description string
}

type ISkillRepository interface {
	Create(skill *Skill) error
	GetById(id uuid.UUID) (*Skill, error)
	GetAll() ([]*Skill, error)
	Update(skill *Skill) error
	DeleteById(id uuid.UUID) error
}

type ISkillService interface {
	Create(skill *Skill) error
	GetById(id uuid.UUID) (*Skill, error)
	GetAll() ([]*Skill, error)
	Update(skill *Skill) error
	DeleteById(id uuid.UUID) error
}
