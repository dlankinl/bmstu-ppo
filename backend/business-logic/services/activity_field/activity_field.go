package activity_field

import (
	"business-logic/domain"
	"context"
	"fmt"
	"github.com/google/uuid"
	"math"
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

func (s *Service) Create(data *domain.ActivityField) (err error) {
	if data.Name == "" {
		return fmt.Errorf("должно быть указано название сферы деятельности")
	}

	if data.Description == "" {
		return fmt.Errorf("должно быть указано описание сферы деятельности")
	}

	if math.Abs(float64(data.Cost)) < 1e-7 {
		return fmt.Errorf("вес сферы деятельности не может быть равен 0")
	}

	ctx := context.Background()

	err = s.actFieldRepo.Create(ctx, data)
	if err != nil {
		return fmt.Errorf("создание сферы деятельности: %w", err)
	}

	return nil
}

func (s *Service) DeleteById(id uuid.UUID) (err error) {
	ctx := context.Background()

	err = s.actFieldRepo.DeleteById(ctx, id)
	if err != nil {
		return fmt.Errorf("удаление сферы деятельности по id: %w", err)
	}

	return nil
}

func (s *Service) Update(data *domain.ActivityField) (err error) {
	ctx := context.Background()

	err = s.actFieldRepo.Update(ctx, data)
	if err != nil {
		return fmt.Errorf("обновление информации о cфере деятельности: %w", err)
	}

	return nil
}

func (s *Service) GetById(id uuid.UUID) (data *domain.ActivityField, err error) {
	ctx := context.Background()

	data, err = s.actFieldRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("получение сферы деятельности по id: %w", err)
	}

	return data, nil
}

func (s *Service) GetCostByCompanyId(companyId uuid.UUID) (cost float32, err error) {
	ctx := context.Background()

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

func (s *Service) GetMaxCost() (maxCost float32, err error) {
	ctx := context.Background()

	maxCost, err = s.actFieldRepo.GetMaxCost(ctx)
	if err != nil {
		return 0, fmt.Errorf("получение максимального веса сферы деятельности: %w", err)
	}

	return maxCost, nil
}

func (s *Service) GetAll(page int) (fields []*domain.ActivityField, err error) {
	ctx := context.Background()

	fields, err = s.actFieldRepo.GetAll(ctx, page)
	if err != nil {
		return nil, fmt.Errorf("получение списка всех сфер деятельности: %w", err)
	}

	return fields, nil
}
