package user

import (
	"business-logic/domain"
	"context"
	"fmt"
	"github.com/google/uuid"
	"strings"
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

func (s *Service) GetByUsername(username string) (user *domain.User, err error) {
	ctx := context.Background()

	user, err = s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("получение пользователя по username: %w", err)
	}

	return user, nil
}

func (s *Service) GetById(userId uuid.UUID) (user *domain.User, err error) {
	ctx := context.Background()

	user, err = s.userRepo.GetById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("получение пользователя по id: %w", err)
	}

	return user, nil
}

func (s *Service) GetAll(page int) (users []*domain.User, err error) {
	ctx := context.Background()

	users, err = s.userRepo.GetAll(ctx, page)
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
