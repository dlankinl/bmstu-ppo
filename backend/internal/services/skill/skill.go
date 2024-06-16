package skill

import (
	"context"
	"fmt"
	"ppo/domain"

	"github.com/google/uuid"
)

type Service struct {
	skillRepo domain.ISkillRepository
}

func NewService(skillRepo domain.ISkillRepository) domain.ISkillService {
	return &Service{
		skillRepo: skillRepo,
	}
}

func (s *Service) Create(ctx context.Context, skill *domain.Skill) (err error) {
	if skill.Name == "" {
		return fmt.Errorf("должно быть указано название навыка")
	}

	if skill.Description == "" {
		return fmt.Errorf("должно быть указано описание навыка")
	}

	err = s.skillRepo.Create(ctx, skill)
	if err != nil {
		return fmt.Errorf("добавление навыка: %w", err)
	}

	return nil
}

func (s *Service) GetById(ctx context.Context, id uuid.UUID) (skill *domain.Skill, err error) {
	skill, err = s.skillRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("получение навыка по id: %w", err)
	}

	return skill, nil
}

func (s *Service) GetAll(ctx context.Context, page int) (skills []*domain.Skill, numPages int, err error) {
	skills, numPages, err = s.skillRepo.GetAll(ctx, page)
	if err != nil {
		return nil, 0, fmt.Errorf("получение списка всех навыков: %w", err)
	}

	return skills, numPages, nil
}

func (s *Service) Update(ctx context.Context, skill *domain.Skill) (err error) {
	err = s.skillRepo.Update(ctx, skill)
	if err != nil {
		return fmt.Errorf("обновление информации о навыке: %w", err)
	}

	return nil
}

func (s *Service) DeleteById(ctx context.Context, id uuid.UUID) (err error) {
	err = s.skillRepo.DeleteById(ctx, id)
	if err != nil {
		return fmt.Errorf("удаление навыка по id: %w", err)
	}

	return nil
}
