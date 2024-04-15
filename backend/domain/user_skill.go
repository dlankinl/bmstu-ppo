package domain

import (
	"context"
	"github.com/google/uuid"
)

type UserSkill struct {
	Username string
	SkillId  uuid.UUID
}

type IUserSkillRepository interface {
	Create(ctx context.Context, pair *UserSkill) error
	Delete(ctx context.Context, pair *UserSkill) error
	GetUserSkillsByUsername(ctx context.Context, username string) ([]*UserSkill, error)
	GetUserSkillsBySkillId(ctx context.Context, skillId uuid.UUID) ([]*UserSkill, error)
}

type IUserSkillService interface {
	Create(pair *UserSkill) error
	Delete(pair *UserSkill) error
	GetSkillsForUser(username string) ([]*Skill, error)
	GetUsersForSkill(skillId uuid.UUID) ([]*User, error)
	DeleteSkillsForUser(userId uuid.UUID) error
}
