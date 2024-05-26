package domain

import (
	"context"
	"github.com/google/uuid"
)

type UserSkill struct {
	UserId  uuid.UUID
	SkillId uuid.UUID
}

type IUserSkillRepository interface {
	Create(context.Context, *UserSkill) error
	Delete(context.Context, *UserSkill) error
	GetUserSkillsByUserId(context.Context, uuid.UUID, int, bool) ([]*UserSkill, int, error)
	GetUserSkillsBySkillId(context.Context, uuid.UUID, int) ([]*UserSkill, error)
}

type IUserSkillService interface {
	Create(context.Context, *UserSkill) error
	Delete(context.Context, *UserSkill) error
	GetSkillsForUser(context.Context, uuid.UUID, int, bool) ([]*Skill, int, error)
	GetUsersForSkill(context.Context, uuid.UUID, int) ([]*User, error)
	DeleteSkillsForUser(context.Context, uuid.UUID) error
}
