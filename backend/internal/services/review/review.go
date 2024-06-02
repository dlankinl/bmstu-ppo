package review

import (
	"context"
	"fmt"
	"ppo/domain"

	"github.com/google/uuid"
)

type Service struct {
	revRepo domain.IReviewRepository
}

func NewService(revRepo domain.IReviewRepository) domain.IReviewService {
	return &Service{
		revRepo: revRepo,
	}
}

func (s *Service) Create(ctx context.Context, rev *domain.Review) (err error) {
	if rev.Rating <= 0 || rev.Rating > 5 {
		return fmt.Errorf("оценка должна быть целым числом от 1 до 5")
	}

	if rev.Pros == "" {
		return fmt.Errorf("описание преимуществ не должно быть пустым")
	}

	if rev.Cons == "" {
		return fmt.Errorf("описание недостатков не должно быть пустым")
	}

	err = s.revRepo.Create(ctx, rev)
	if err != nil {
		return fmt.Errorf("создание отзыва: %w", err)
	}

	return nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (rev *domain.Review, err error) {
	rev, err = s.revRepo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("получение отзыва по id: %w", err)
	}

	return rev, nil
}

func (s *Service) GetAllForReviewer(ctx context.Context, id uuid.UUID, page int) (revs []*domain.Review, numPages int, err error) {
	revs, numPages, err = s.revRepo.GetAllForReviewer(ctx, id, page)
	if err != nil {
		return nil, 0, fmt.Errorf("получение всех отзывов ревьювера: %w", err)
	}

	return revs, numPages, nil
}

func (s *Service) GetAllForTarget(ctx context.Context, id uuid.UUID, page int) (revs []*domain.Review, numPages int, err error) {
	revs, numPages, err = s.revRepo.GetAllForTarget(ctx, id, page)
	if err != nil {
		return nil, 0, fmt.Errorf("получение всех отзывов объекта: %w", err)
	}

	return revs, numPages, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) (err error) {
	err = s.revRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("удаление отзыва: %w", err)
	}

	return nil
}
