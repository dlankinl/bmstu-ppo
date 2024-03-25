package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"ppo/domain"
)

type UserSkillService struct {
	userSkillRepo domain.IUserSkillRepository
	userRepo      domain.IUserRepository
	skillRepo     domain.ISkillRepository
}

func NewUserSkillService(
	userSkillRepo domain.IUserSkillRepository,
	userRepo domain.IUserRepository,
	skillRepo domain.ISkillRepository,
) domain.IUserSkillService {
	return &UserSkillService{
		userSkillRepo: userSkillRepo,
		userRepo:      userRepo,
		skillRepo:     skillRepo,
	}
}

func (s UserSkillService) Create(ctx context.Context, pair *domain.UserSkill) (err error) {
	err = s.userSkillRepo.Create(ctx, pair)
	if err != nil {
		return fmt.Errorf("связывание пользователя и навыка: %w", err)
	}

	return nil
}

func (s UserSkillService) Delete(ctx context.Context, pair *domain.UserSkill) (err error) {
	err = s.userSkillRepo.Delete(ctx, pair)
	if err != nil {
		return fmt.Errorf("удаление связи пользователь-навык: %w", err)
	}

	return nil
}

func (s UserSkillService) GetSkillsForUser(ctx context.Context, userId uuid.UUID) (skills []*domain.Skill, err error) {
	userSkills, err := s.userSkillRepo.GetUserSkillsByUserId(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("получение связок пользователь-навык по userId: %w", err)
	}

	skills = make([]*domain.Skill, len(userSkills))
	for i, userSkill := range userSkills {
		skill, err := s.skillRepo.GetById(ctx, userSkill.SkillId)
		if err != nil {
			return nil, fmt.Errorf("получение скилла по skillId: %w", err)
		}

		skills[i] = skill
	}

	return skills, nil
}

func (s UserSkillService) GetUsersForSkill(ctx context.Context, skillId uuid.UUID) (users []*domain.User, err error) {
	userSkills, err := s.userSkillRepo.GetUserSkillsBySkillId(ctx, skillId)
	if err != nil {
		return nil, fmt.Errorf("получение связок пользователь-навык по skillId: %w", err)
	}

	users = make([]*domain.User, len(userSkills))
	for i, userSkill := range userSkills {
		user, err := s.userRepo.GetById(ctx, userSkill.UserId)
		if err != nil {
			return nil, fmt.Errorf("получение пользователя по userId: %w", err)
		}

		users[i] = user
	}

	return users, nil
}

func (s UserSkillService) DeleteSkillsForUser(ctx context.Context, userId uuid.UUID) (err error) {
	userSkills, err := s.userSkillRepo.GetUserSkillsByUserId(ctx, userId)
	if err != nil {
		return fmt.Errorf("получение связок пользователь-навык по userId: %w", err)
	}

	for _, userSkill := range userSkills {
		err = s.skillRepo.DeleteById(ctx, userSkill.SkillId)
		if err != nil {
			return fmt.Errorf("удаление скилла по skillId: %w", err)
		}
	}

	// todo: додумать

	return nil
}
