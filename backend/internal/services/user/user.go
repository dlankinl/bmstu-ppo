package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"ppo/domain"
	"ppo/internal/cache"
	"ppo/internal/config"
	"ppo/pkg/logger"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	userRepo     domain.IUserRepository
	companyRepo  domain.ICompanyRepository
	actFieldRepo domain.IActivityFieldRepository
	cache        cache.Cache
	logger       logger.ILogger
}

func NewService(
	userRepo domain.IUserRepository,
	companyRepo domain.ICompanyRepository,
	actFieldRepo domain.IActivityFieldRepository,
	cache cache.Cache,
	logger logger.ILogger,
) domain.IUserService {
	return &Service{
		userRepo:     userRepo,
		companyRepo:  companyRepo,
		actFieldRepo: actFieldRepo,
		cache:        cache,
		logger:       logger,
	}
}

func (s *Service) Create(ctx context.Context, user *domain.User) (err error) {
	prompt := "UserCreate"

	if user.Gender != "m" && user.Gender != "w" {
		s.logger.Infof("%s: неизвестный пол", prompt)
		return fmt.Errorf("неизвестный пол")
	}

	if user.City == "" {
		s.logger.Infof("%s: должно быть указано название города", prompt)
		return fmt.Errorf("должно быть указано название города")
	}

	if user.Birthday.IsZero() {
		s.logger.Infof("%s: должна быть указана дата рождения", prompt)
		return fmt.Errorf("должна быть указана дата рождения")
	}

	if user.FullName == "" {
		s.logger.Infof("%s: должны быть указаны ФИО", prompt)
		return fmt.Errorf("должны быть указаны ФИО")
	}

	if len(strings.Split(user.FullName, " ")) != 3 {
		s.logger.Infof("%s: некорректное количество слов (должны быть фамилия, имя и отчество)", prompt)
		return fmt.Errorf("некорректное количество слов (должны быть фамилия, имя и отчество)")
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		s.logger.Infof("%s: создание пользователя: %v", prompt, err)
		return fmt.Errorf("создание пользователя: %w", err)
	}

	return nil
}

func (s *Service) GetByUsername(ctx context.Context, username string) (user *domain.User, err error) {
	prompt := "UserGetByUsername"

	err = s.cache.Get(ctx, username, user)
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		s.logger.Infof("%s: %v", prompt, err)
		return nil, fmt.Errorf("%s: %w", prompt, err)
	}

	if user == nil {
		user, err = s.userRepo.GetByUsername(ctx, username)
		if err != nil {
			s.logger.Infof("%s: получение пользователя по username: %v", prompt, err)
			return nil, fmt.Errorf("получение пользователя по username: %w", err)
		}

		err = s.cache.Set(ctx, username, user, config.ExpirationTime)
		if err != nil {
			s.logger.Infof("%s: %v", prompt, err)
			return nil, fmt.Errorf("%s: %w", prompt, err)
		}
	}

	return user, nil
}

type User struct {
	ID       uuid.UUID
	Username string
	FullName string
	Gender   string
	Birthday time.Time
	City     string
	Role     string
}

func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

func (s *Service) GetById(ctx context.Context, userId uuid.UUID) (user *domain.User, err error) {
	prompt := "UserGetById"

	var tmp User
	err = s.cache.Get(ctx, userId.String(), &tmp)
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		s.logger.Infof("%s: %v", prompt, err)
		return nil, fmt.Errorf("%s: %w", prompt, err)
	}
	fmt.Println("THIS MODEL: ", tmp)

	if user == nil {
		user, err = s.userRepo.GetById(ctx, userId)
		if err != nil {
			s.logger.Infof("%s: получение пользователя по id: %v", prompt, err)
			return nil, fmt.Errorf("получение пользователя по id: %w", err)
		}

		err = s.cache.Set(ctx, userId.String(), user, config.ExpirationTime)
		if err != nil {
			s.logger.Infof("%s: %v", prompt, err)
			return nil, fmt.Errorf("%s: %w", prompt, err)
		}
	}
	fmt.Println("AFTER REPO: ", user)

	return user, nil
}

func (s *Service) GetAll(ctx context.Context, page int) (users []*domain.User, numPages int, err error) {
	prompt := "UserGetAll"

	users, numPages, err = s.userRepo.GetAll(ctx, page)
	if err != nil {
		s.logger.Infof("%s: получение списка всех пользователей: %v", prompt, err)
		return nil, 0, fmt.Errorf("получение списка всех пользователей: %w", err)
	}

	return users, numPages, nil
}

func (s *Service) Update(ctx context.Context, user *domain.User) (err error) {
	prompt := "UserUpdate"

	if user.Gender != "m" && user.Gender != "w" {
		s.logger.Infof("%s: неизвестный пол", prompt)
		return fmt.Errorf("неизвестный пол")
	}

	if user.City == "" {
		s.logger.Infof("%s: должно быть указано название города", prompt)
		return fmt.Errorf("должно быть указано название города")
	}

	if user.Birthday.IsZero() {
		s.logger.Infof("%s: должна быть указана дата рождения", prompt)
		return fmt.Errorf("должна быть указана дата рождения")
	}

	if user.FullName == "" {
		s.logger.Infof("%s: должны быть указаны ФИО", prompt)
		return fmt.Errorf("должны быть указаны ФИО")
	}

	if len(strings.Split(user.FullName, " ")) != 3 {
		s.logger.Infof("%s: некорректное количество слов (должны быть фамилия, имя и отчество)", prompt)
		return fmt.Errorf("некорректное количество слов (должны быть фамилия, имя и отчество)")
	}

	if user.Role != "admin" && user.Role != "user" {
		s.logger.Infof("%s: невалидная роль", prompt)
		return fmt.Errorf("невалидная роль")
	}

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		s.logger.Infof("%s: создание пользователя: %v", prompt, err)
		return fmt.Errorf("обновление информации о пользователе: %w", err)
	}

	return nil
}

func (s *Service) DeleteById(ctx context.Context, id uuid.UUID) (err error) {
	prompt := "UserDeleteById"

	err = s.userRepo.DeleteById(ctx, id)
	if err != nil {
		s.logger.Infof("%s: удаление пользователя по id: %v", prompt, err)
		return fmt.Errorf("удаление пользователя по id: %w", err)
	}

	return nil
}
