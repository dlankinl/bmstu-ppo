package activity_field

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"math"
	"ppo/domain"
)

type Service struct {
	repo        domain.IActivityFieldRepository
	companyRepo domain.ICompanyRepository
}

func NewService(repo domain.IActivityFieldRepository) domain.IActivityFieldService {
	return &Service{
		repo: repo,
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

	err = s.repo.Create(ctx, data)
	if err != nil {
		return fmt.Errorf("создание сферы деятельности: %w", err)
	}

	return nil
}

func (s *Service) DeleteById(id uuid.UUID) (err error) {
	ctx := context.Background()

	err = s.repo.DeleteById(ctx, id)
	if err != nil {
		return fmt.Errorf("удаление сферы деятельности по id: %w", err)
	}

	return nil
}

func (s *Service) Update(data *domain.ActivityField) (err error) {
	ctx := context.Background()

	err = s.repo.Update(ctx, data)
	if err != nil {
		return fmt.Errorf("обновление информации о cфере деятельности: %w", err)
	}

	return nil
}

func (s *Service) GetById(id uuid.UUID) (data *domain.ActivityField, err error) {
	ctx := context.Background()

	data, err = s.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("получение сферы деятельности по id: %w", err)
	}

	return data, nil
}

func (s *Service) GetCostByCompanyId(id uuid.UUID) (cost float32, err error) {
	ctx := context.Background()

	cost, err = s.repo.GetByCompanyId(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("получение веса сферы деятельности по id компании: %w", err)
	}

	return cost, nil
}

func (s *Service) GetMaxCost() (maxCost float32, err error) {
	ctx := context.Background()

	maxCost, err = s.repo.GetMaxCost(ctx)
	if err != nil {
		return 0, fmt.Errorf("получение максимального веса сфееры деятельности: %w", err)
	}

	return maxCost, nil
}
