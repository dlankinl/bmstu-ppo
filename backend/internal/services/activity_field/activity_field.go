package activity_field

import (
	"context"
	"fmt"
	"math"
	"ppo/domain"
	"ppo/pkg/logger"

	"github.com/google/uuid"
)

type Service struct {
	actFieldRepo domain.IActivityFieldRepository
	compRepo     domain.ICompanyRepository
	logger       logger.ILogger
}

func NewService(
	actFieldRepo domain.IActivityFieldRepository,
	compRepo domain.ICompanyRepository,
	logger logger.ILogger,
) domain.IActivityFieldService {
	return &Service{
		actFieldRepo: actFieldRepo,
		compRepo:     compRepo,
		logger:       logger,
	}
}

func (s *Service) Create(ctx context.Context, data *domain.ActivityField) (err error) {
	prompt := "ActivityFieldCreate"
	if data.Name == "" {
		s.logger.Infof("%s: должно быть указано название сферы деятельности", prompt)
		return fmt.Errorf("должно быть указано название сферы деятельности")
	}

	if data.Description == "" {
		s.logger.Infof("%s: должно быть указано описание сферы деятельности", prompt)
		return fmt.Errorf("должно быть указано описание сферы деятельности")
	}

	if math.Abs(float64(data.Cost)) < 1e-7 {
		s.logger.Infof("%s: вес сферы деятельности не может быть равен 0", prompt)
		return fmt.Errorf("вес сферы деятельности не может быть равен 0")
	}

	err = s.actFieldRepo.Create(ctx, data)
	if err != nil {
		s.logger.Infof("%s: создание сферы деятельности: %v", prompt, err)
		return fmt.Errorf("создание сферы деятельности: %w", err)
	}

	return nil
}

func (s *Service) DeleteById(ctx context.Context, id uuid.UUID) (err error) {
	prompt := "ActivityFieldDeleteById"

	err = s.actFieldRepo.DeleteById(ctx, id)
	if err != nil {
		s.logger.Infof("%s: удаление сферы деятельности по id: %v", prompt, err)
		return fmt.Errorf("удаление сферы деятельности по id: %w", err)
	}

	return nil
}

func (s *Service) Update(ctx context.Context, data *domain.ActivityField) (err error) {
	prompt := "ActivityFieldUpdate"

	err = s.actFieldRepo.Update(ctx, data)
	if err != nil {
		s.logger.Infof("%s: обновление информации о cфере деятельности: %v", prompt, err)
		return fmt.Errorf("обновление информации о cфере деятельности: %w", err)
	}

	return nil
}

func (s *Service) GetById(ctx context.Context, id uuid.UUID) (data *domain.ActivityField, err error) {
	prompt := "ActivityFieldGetById"

	data, err = s.actFieldRepo.GetById(ctx, id)
	if err != nil {
		s.logger.Infof("%s: получение сферы деятельности по id: %v", prompt, err)
		return nil, fmt.Errorf("получение сферы деятельности по id: %w", err)
	}

	return data, nil
}

func (s *Service) GetCostByCompanyId(ctx context.Context, companyId uuid.UUID) (cost float32, err error) {
	prompt := "ActivityFieldGetCostByCompanyId"

	company, err := s.compRepo.GetById(ctx, companyId)
	if err != nil {
		s.logger.Infof("%s: получение компании по id: %v", prompt, err)
		return 0, fmt.Errorf("получение компании по id: %w", err)
	}

	field, err := s.actFieldRepo.GetById(ctx, company.ActivityFieldId)
	if err != nil {
		s.logger.Infof("%s: получение сферы деятельности по id: %v", prompt, err)
		return 0, fmt.Errorf("получение сферы деятельности по id: %w", err)
	}
	cost = field.Cost

	return cost, nil
}

func (s *Service) GetMaxCost(ctx context.Context) (maxCost float32, err error) {
	prompt := "ActivityFieldGetMaxCost"

	maxCost, err = s.actFieldRepo.GetMaxCost(ctx)
	if err != nil {
		s.logger.Infof("%s: получение максимального веса сферы деятельности: %v", prompt, err)
		return 0, fmt.Errorf("получение максимального веса сферы деятельности: %w", err)
	}

	return maxCost, nil
}

func (s *Service) GetAll(ctx context.Context, page int, isPaginated bool) (fields []*domain.ActivityField, numPages int, err error) {
	prompt := "ActivityFieldGetAll"

	fields, numPages, err = s.actFieldRepo.GetAll(ctx, page, isPaginated)
	if err != nil {
		s.logger.Infof("%s: получение списка всех сфер деятельности: %v", prompt, err)
		return nil, 0, fmt.Errorf("получение списка всех сфер деятельности: %w", err)
	}

	return fields, numPages, nil
}
