package skill

import (
	"context"
	"fmt"
	"ppo/domain"
	"ppo/pkg/logger"

	"github.com/google/uuid"
)

type Service struct {
	skillRepo domain.ISkillRepository
	logger    logger.ILogger
}

func NewService(skillRepo domain.ISkillRepository, logger logger.ILogger) domain.ISkillService {
	return &Service{
		skillRepo: skillRepo,
		logger:    logger,
	}
}

func (s *Service) Create(ctx context.Context, skill *domain.Skill) (err error) {
	prompt := "SkillCreate"

	if skill.Name == "" {
		s.logger.Infof("%s: должно быть указано название навыка", prompt)
		return fmt.Errorf("должно быть указано название навыка")
	}

	if skill.Description == "" {
		s.logger.Infof("%s: должно быть указано описание навыка", prompt)
		return fmt.Errorf("должно быть указано описание навыка")
	}

	err = s.skillRepo.Create(ctx, skill)
	if err != nil {
		s.logger.Infof("%s: добавление навыка: %v", prompt, err)
		return fmt.Errorf("добавление навыка: %w", err)
	}

	return nil
}

func (s *Service) GetById(ctx context.Context, id uuid.UUID) (skill *domain.Skill, err error) {
	prompt := "SkillGetById"

	skill, err = s.skillRepo.GetById(ctx, id)
	if err != nil {
		s.logger.Infof("%s: получение навыка по id: %v", prompt, err)
		return nil, fmt.Errorf("получение навыка по id: %w", err)
	}

	return skill, nil
}

func (s *Service) GetAll(ctx context.Context, page int) (skills []*domain.Skill, numPages int, err error) {
	prompt := "SkillGetAll"

	skills, numPages, err = s.skillRepo.GetAll(ctx, page)
	if err != nil {
		s.logger.Infof("%s: получение списка всех навыков: %v", prompt, err)
		return nil, 0, fmt.Errorf("получение списка всех навыков: %w", err)
	}

	return skills, numPages, nil
}

func (s *Service) Update(ctx context.Context, skill *domain.Skill) (err error) {
	prompt := "SkillUpdate"

	err = s.skillRepo.Update(ctx, skill)
	if err != nil {
		s.logger.Infof("%s: обновление информации о навыке: %v", prompt, err)
		return fmt.Errorf("обновление информации о навыке: %w", err)
	}

	return nil
}

func (s *Service) DeleteById(ctx context.Context, id uuid.UUID) (err error) {
	prompt := "SkillDeleteById"

	err = s.skillRepo.DeleteById(ctx, id)
	if err != nil {
		s.logger.Infof("%s: удаление навыка по id: %v", prompt, err)
		return fmt.Errorf("удаление навыка по id: %w", err)
	}

	return nil
}
