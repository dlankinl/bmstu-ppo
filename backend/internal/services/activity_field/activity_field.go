package activity_field

import (
	"context"
	"fmt"
	"math"
	"ppo/domain"

	"github.com/google/uuid"
)

type Service struct {
	actFieldRepo domain.IActivityFieldRepository
	compRepo     domain.ICompanyRepository
}

func NewService(
	actFieldRepo domain.IActivityFieldRepository,
	compRepo domain.ICompanyRepository,
) domain.IActivityFieldService {
	return &Service{
		actFieldRepo: actFieldRepo,
		compRepo:     compRepo,
	}
}

func (s *Service) Create(ctx context.Context, data *domain.ActivityField) (err error) {
	if data.Name == "" {
		return fmt.Errorf("должно быть указано название сферы деятельности")
	}

	if data.Description == "" {
		return fmt.Errorf("должно быть указано описание сферы деятельности")
	}

	if math.Abs(float64(data.Cost)) < 1e-7 {
		return fmt.Errorf("вес сферы деятельности не может быть равен 0")
	}

	err = s.actFieldRepo.Create(ctx, data)
	if err != nil {
		return fmt.Errorf("создание сферы деятельности: %w", err)
	}

	return nil
}

func (s *Service) DeleteById(ctx context.Context, id uuid.UUID) (err error) {
	err = s.actFieldRepo.DeleteById(ctx, id)
	if err != nil {
		return fmt.Errorf("удаление сферы деятельности по id: %w", err)
	}

	return nil
}

func (s *Service) Update(ctx context.Context, data *domain.ActivityField) (err error) {
	err = s.actFieldRepo.Update(ctx, data)
	if err != nil {
		return fmt.Errorf("обновление информации о cфере деятельности: %w", err)
	}

	return nil
}

func (s *Service) GetById(ctx context.Context, id uuid.UUID) (data *domain.ActivityField, err error) {

	data, err = s.actFieldRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("получение сферы деятельности по id: %w", err)
	}

	return data, nil
}

func (s *Service) GetCostByCompanyId(ctx context.Context, companyId uuid.UUID) (cost float32, err error) {
	company, err := s.compRepo.GetById(ctx, companyId)
	if err != nil {
		return 0, fmt.Errorf("получение компании по id: %w", err)
	}

	field, err := s.actFieldRepo.GetById(ctx, company.ActivityFieldId)
	if err != nil {
		return 0, fmt.Errorf("получение сферы деятельности по id: %w", err)
	}
	cost = field.Cost

	return cost, nil
}

func (s *Service) GetMaxCost(ctx context.Context) (maxCost float32, err error) {
	maxCost, err = s.actFieldRepo.GetMaxCost(ctx)
	if err != nil {
		return 0, fmt.Errorf("получение максимального веса сферы деятельности: %w", err)
	}

	return maxCost, nil
}

func (s *Service) GetAll(ctx context.Context, page int, isPaginated bool) (fields []*domain.ActivityField, numPages int, err error) {
	fields, numPages, err = s.actFieldRepo.GetAll(ctx, page, isPaginated)
	if err != nil {
		return nil, 0, fmt.Errorf("получение списка всех сфер деятельности: %w", err)
	}

	return fields, numPages, nil
}
