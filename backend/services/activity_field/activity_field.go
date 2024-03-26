package activity_field

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"math"
	"ppo/domain"
)

type Service struct {
	repo domain.IActivityFieldRepository
}

func NewService(repo domain.IActivityFieldRepository) domain.IActivityFieldService {
	return &Service{
		repo: repo,
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

	err = s.repo.Create(ctx, data)
	if err != nil {
		return fmt.Errorf("создание сферы деятельности: %w", err)
	}

	return nil
}

func (s *Service) DeleteById(ctx context.Context, id uuid.UUID) (err error) {
	err = s.repo.DeleteById(ctx, id)
	if err != nil {
		return fmt.Errorf("удаление сферы деятельности по id: %w", err)
	}

	return nil
}

func (s *Service) Update(ctx context.Context, data *domain.ActivityField) (err error) {
	err = s.repo.Update(ctx, data)
	if err != nil {
		return fmt.Errorf("обновление информации о cфере деятельности: %w", err)
	}

	return nil
}

func (s *Service) GetById(ctx context.Context, id uuid.UUID) (data *domain.ActivityField, err error) {
	data, err = s.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("получение сферы деятельности по id: %w", err)
	}

	return data, nil
}
