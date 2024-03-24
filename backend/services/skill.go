package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"ppo/domain"
)

type SkillService struct {
	skillRepo domain.ISkillRepository
}

func (s SkillService) Create(ctx context.Context, skill *domain.Skill) (err error) {
	if skill.Name == "" {
		return fmt.Errorf("должно быть указано название навыка")
	}

	if skill.Description == "" {
		return fmt.Errorf("должно быть указано описание навыка")
	}

	err = s.skillRepo.Create(ctx, skill)
	if err != nil {
		return fmt.Errorf("создание навыка: %w", err)
	}

	return nil
}

func (s SkillService) GetById(ctx context.Context, id uuid.UUID) (skill *domain.Skill, err error) {
	skill, err = s.skillRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("получение навыка по id: %w", err)
	}

	return skill, nil
}

func (s SkillService) GetAll(ctx context.Context) (skills []*domain.Skill, err error) {
	skills, err = s.skillRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("получение всех навыков: %w", err)
	}

	return skills, nil
}

func (s SkillService) Update(ctx context.Context, skill *domain.Skill) (err error) {
	err = s.skillRepo.Update(ctx, skill)
	if err != nil {
		return fmt.Errorf("обновление информации о навыке: %w", err)
	}

	return nil
}

func (s SkillService) DeleteById(ctx context.Context, id uuid.UUID) (err error) {
	err = s.skillRepo.DeleteById(ctx, id)
	if err != nil {
		return fmt.Errorf("удаление навыка по id: %w", err)
	}

	return nil
}
