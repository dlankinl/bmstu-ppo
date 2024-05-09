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
	Create(ctx context.Context, pair *UserSkill) error
	Delete(ctx context.Context, pair *UserSkill) error
	GetUserSkillsByUserId(ctx context.Context, userId uuid.UUID, page int, isPaginated bool) ([]*UserSkill, error)
	GetUserSkillsBySkillId(ctx context.Context, skillId uuid.UUID, page int) ([]*UserSkill, error)
}

type IUserSkillService interface {
	Create(pair *UserSkill) error
	Delete(pair *UserSkill) error
	GetSkillsForUser(userId uuid.UUID, page int, isPaginated bool) ([]*Skill, error)
	GetUsersForSkill(skillId uuid.UUID, page int) ([]*User, error)
	DeleteSkillsForUser(userId uuid.UUID) error
}
