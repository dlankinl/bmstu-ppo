package user

import (
	"context"
	"fmt"
	"ppo/domain"
	"strings"

	"github.com/google/uuid"
)

type Service struct {
	userRepo     domain.IUserRepository
	companyRepo  domain.ICompanyRepository
	actFieldRepo domain.IActivityFieldRepository
}

func NewService(
	userRepo domain.IUserRepository,
	companyRepo domain.ICompanyRepository,
	actFieldRepo domain.IActivityFieldRepository,
) domain.IUserService {
	return &Service{
		userRepo:     userRepo,
		companyRepo:  companyRepo,
		actFieldRepo: actFieldRepo,
	}
}

func (s *Service) Create(ctx context.Context, user *domain.User) (err error) {
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

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("создание пользователя: %w", err)
	}

	return nil
}

func (s *Service) GetByUsername(ctx context.Context, username string) (user *domain.User, err error) {
	user, err = s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("получение пользователя по username: %w", err)
	}

	return user, nil
}

func (s *Service) GetById(ctx context.Context, userId uuid.UUID) (user *domain.User, err error) {
	user, err = s.userRepo.GetById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("получение пользователя по id: %w", err)
	}

	return user, nil
}

func (s *Service) GetAll(ctx context.Context, page int) (users []*domain.User, numPages int, err error) {
	users, numPages, err = s.userRepo.GetAll(ctx, page)
	if err != nil {
		return nil, 0, fmt.Errorf("получение списка всех пользователей: %w", err)
	}

	return users, numPages, nil
}

func (s *Service) Update(ctx context.Context, user *domain.User) (err error) {
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

	if user.Role != "admin" && user.Role != "user" {
		return fmt.Errorf("невалидная роль")
	}

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("обновление информации о пользователе: %w", err)
	}

	return nil
}

func (s *Service) DeleteById(ctx context.Context, id uuid.UUID) (err error) {
	err = s.userRepo.DeleteById(ctx, id)
	if err != nil {
		return fmt.Errorf("удаление пользователя по id: %w", err)
	}

	return nil
}
