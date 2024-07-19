package review

import (
	"context"
	"fmt"
	"ppo/domain"
	"ppo/pkg/logger"

	"github.com/google/uuid"
)

type Service struct {
	revRepo domain.IReviewRepository
	logger  logger.ILogger
}

func NewService(revRepo domain.IReviewRepository, logger logger.ILogger) domain.IReviewService {
	return &Service{
		revRepo: revRepo,
		logger:  logger,
	}
}

func (s *Service) Create(ctx context.Context, rev *domain.Review) (err error) {
	prompt := "ReviewCreate"

	if rev.Rating <= 0 || rev.Rating > 5 {
		s.logger.Infof("%s: оценка должна быть целым числом от 1 до 5", prompt)
		return fmt.Errorf("оценка должна быть целым числом от 1 до 5")
	}

	if rev.Pros == "" {
		s.logger.Infof("%s: описание преимуществ не должно быть пустым", prompt)
		return fmt.Errorf("описание преимуществ не должно быть пустым")
	}

	if rev.Cons == "" {
		s.logger.Infof("%s: описание недостатков не должно быть пустым", prompt)
		return fmt.Errorf("описание недостатков не должно быть пустым")
	}

	err = s.revRepo.Create(ctx, rev)
	if err != nil {
		s.logger.Infof("%s: создание отзыва: %v", prompt, err)
		return fmt.Errorf("создание отзыва: %w", err)
	}

	return nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (rev *domain.Review, err error) {
	prompt := "ReviewGet"

	rev, err = s.revRepo.Get(ctx, id)
	if err != nil {
		s.logger.Infof("%s: получение отзыва: %v", prompt, err)
		return nil, fmt.Errorf("получение отзыва по id: %w", err)
	}

	return rev, nil
}

func (s *Service) GetAllForReviewer(ctx context.Context, id uuid.UUID, page int) (revs []*domain.Review, numPages int, err error) {
	prompt := "ReviewGetAllForReviewer"

	revs, numPages, err = s.revRepo.GetAllForReviewer(ctx, id, page)
	if err != nil {
		s.logger.Infof("%s: получение всех отзывов ревьювера: %v", prompt, err)
		return nil, 0, fmt.Errorf("получение всех отзывов ревьювера: %w", err)
	}

	return revs, numPages, nil
}

func (s *Service) GetAllForTarget(ctx context.Context, id uuid.UUID, page int) (revs []*domain.Review, numPages int, err error) {
	prompt := "ReviewGetAllForTarget"

	revs, numPages, err = s.revRepo.GetAllForTarget(ctx, id, page)
	if err != nil {
		s.logger.Infof("%s: получение всех отзывов объекта: %v", prompt, err)
		return nil, 0, fmt.Errorf("получение всех отзывов объекта: %w", err)
	}

	return revs, numPages, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) (err error) {
	prompt := "ReviewDelete"

	err = s.revRepo.Delete(ctx, id)
	if err != nil {
		s.logger.Infof("%s: удаление отзыва: %v", prompt, err)
		return fmt.Errorf("удаление отзыва: %w", err)
	}

	return nil
}
