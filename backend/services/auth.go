package services

import (
	"fmt"
	"ppo/domain"
	"ppo/pkg/base"
)

type AuthService struct {
	authRepo domain.IAuthRepository
	jwtKey   string
}

func (s AuthService) Register(username, password string) (err error) {
	if username == "" {
		return fmt.Errorf("должно быть указано имя пользователя")
	}

	if password == "" {
		return fmt.Errorf("должен быть указан пароль")
	}

	hashedPass, err := base.GenerateHashPass(password)
	if err != nil {
		return fmt.Errorf("генерация хэша: %w", err)
	}

	userAuth := &domain.UserAuth{
		Username:     username,
		PasswordHash: hashedPass,
	}

	err = s.authRepo.Register(userAuth)
	if err != nil {
		return fmt.Errorf("регистрация пользователя: %w", err)
	}

	return nil
}

func (s AuthService) Login(username, password string) (token string, err error) {
	if username == "" {
		return "", fmt.Errorf("должно быть указано имя пользователя")
	}

	if password == "" {
		return "", fmt.Errorf("должен быть указан пароль")
	}

	userAuth, err := s.authRepo.GetByUsername(username)
	if err != nil {
		return "", fmt.Errorf("получение пользователя по username: %w", err) // FIXME: invalid_username
	}

	if !base.CheckPasswordHash(password, userAuth.PasswordHash) {
		return "", fmt.Errorf("неверный пароль") // FIXME: incorrect_credentials
	}

	token, err = base.GenerateAuthToken(username, s.jwtKey)
	if err != nil {
		return "", fmt.Errorf("генерация токена: %w", err)
	}

	return token, nil
}
