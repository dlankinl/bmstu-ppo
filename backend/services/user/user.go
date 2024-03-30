package user

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"ppo/domain"
	"ppo/pkg/utils"
	"strings"
)

type Service struct {
	userRepo     domain.IUserRepository
	companyRepo  domain.ICompanyRepository
	finRepo      domain.IFinancialReportRepository
	actFieldRepo domain.IActivityFieldRepository
}

func NewService(
	userRepo domain.IUserRepository,
	companyRepo domain.ICompanyRepository,
	finRepo domain.IFinancialReportRepository,
	actFieldRepo domain.IActivityFieldRepository,
) domain.IUserService {
	return &Service{
		userRepo:     userRepo,
		companyRepo:  companyRepo,
		finRepo:      finRepo,
		actFieldRepo: actFieldRepo,
	}
}

func (s *Service) Create(user *domain.User) (err error) {
	if user.Gender != "m" && user.Gender != "w" {
		return fmt.Errorf("неизвестный пол")
	}

	if user.City == "" {
		return fmt.Errorf("должно быть указано название города")
	}

	if user.Birthday.IsZero() {
		return fmt.Errorf("должна быть указана дата рождения")
	}

	if user.FullName == "" {
		return fmt.Errorf("должны быть указаны ФИО")
	}

	if len(strings.Split(user.FullName, " ")) != 3 {
		return fmt.Errorf("некорректное количество слов (должны быть фамилия, имя и отчество)")
	}

	ctx := context.Background()

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("создание пользователя: %w", err)
	}

	return nil
}

func (s *Service) GetById(id uuid.UUID) (user *domain.User, err error) {
	ctx := context.Background()

	user, err = s.userRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("получение пользователя по id: %w", err)
	}

	return user, nil
}

// TODO: pagination
func (s *Service) GetAll(filters utils.Filters, page int) (users []*domain.User, err error) {
	err = filters.Validate()
	if err != nil {
		return nil, fmt.Errorf("валидация фильтров: %w", err)
	}

	ctx := context.Background()

	users, err = s.userRepo.GetAll(ctx, filters, page)
	if err != nil {
		return nil, fmt.Errorf("получение списка всех пользователей: %w", err)
	}

	return users, nil
}

func (s *Service) Update(user *domain.User) (err error) {
	ctx := context.Background()

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("обновление информации о пользователе: %w", err)
	}

	return nil
}

func (s *Service) DeleteById(id uuid.UUID) (err error) {
	ctx := context.Background()

	err = s.userRepo.DeleteById(ctx, id)
	if err != nil {
		return fmt.Errorf("удаление пользователя по id: %w", err)
	}

	return nil
}
