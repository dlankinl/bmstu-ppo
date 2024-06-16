package user_skill

import (
	"business-logic/domain"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type Service struct {
	userSkillRepo domain.IUserSkillRepository
	userRepo      domain.IUserRepository
	skillRepo     domain.ISkillRepository
}

func NewService(
	userSkillRepo domain.IUserSkillRepository,
	userRepo domain.IUserRepository,
	skillRepo domain.ISkillRepository,
) domain.IUserSkillService {
	return &Service{
		userSkillRepo: userSkillRepo,
		userRepo:      userRepo,
		skillRepo:     skillRepo,
	}
}

func (s *Service) Create(pair *domain.UserSkill) (err error) {
	ctx := context.Background()

	err = s.userSkillRepo.Create(ctx, pair)
	if err != nil {
		return fmt.Errorf("связывание пользователя и навыка: %w", err)
	}

	return nil
}

func (s *Service) Delete(pair *domain.UserSkill) (err error) {
	ctx := context.Background()

	err = s.userSkillRepo.Delete(ctx, pair)
	if err != nil {
		return fmt.Errorf("удаление связи пользователь-навык: %w", err)
	}

	return nil
}

func (s *Service) GetSkillsForUser(userId uuid.UUID, page int) (skills []*domain.Skill, err error) {
	ctx := context.Background()

	userSkills, err := s.userSkillRepo.GetUserSkillsByUserId(ctx, userId, page)
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

func (s *Service) GetUsersForSkill(skillId uuid.UUID, page int) (users []*domain.User, err error) {
	ctx := context.Background()

	userSkills, err := s.userSkillRepo.GetUserSkillsBySkillId(ctx, skillId, page)
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

func (s *Service) DeleteSkillsForUser(userId uuid.UUID) (err error) {
	ctx := context.Background()

	userSkills, err := s.userSkillRepo.GetUserSkillsByUserId(ctx, userId, 0)
	if err != nil {
		return fmt.Errorf("получение связок пользователь-навык по userId: %w", err)
	}

	for _, userSkill := range userSkills {
		err = s.userSkillRepo.Delete(ctx, userSkill)
		if err != nil {
			return fmt.Errorf("удаление пары пользователь-навык: %w", err)
		}
	}

	return nil
}
