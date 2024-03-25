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
	GetUserSkillsByUserId(ctx context.Context, userId uuid.UUID) ([]*UserSkill, error)
	GetUserSkillsBySkillId(ctx context.Context, skillId uuid.UUID) ([]*UserSkill, error)
}

type IUserSkillService interface {
	Create(ctx context.Context, pair *UserSkill) error
	Delete(ctx context.Context, pair *UserSkill) error
	GetSkillsForUser(ctx context.Context, userId uuid.UUID) ([]*Skill, error)
	GetUsersForSkill(ctx context.Context, skillId uuid.UUID) ([]*User, error)
	DeleteSkillsForUser(ctx context.Context, userId uuid.UUID) error
}
