package auth

import (
	"business-logic/domain"
	"business-logic/pkg/base"
	"context"
	"fmt"
)

type Service struct {
	authRepo domain.IAuthRepository
	crypto   base.IHashCrypto
	jwtKey   string
}

func NewService(repo domain.IAuthRepository, crypto base.IHashCrypto, jwtKey string) domain.IAuthService {
	return &Service{
		authRepo: repo,
		crypto:   crypto,
		jwtKey:   jwtKey,
	}
}

func (s *Service) Register(authInfo *domain.UserAuth) (err error) {
	if authInfo.Username == "" {
		return fmt.Errorf("должно быть указано имя пользователя")
	}

	if authInfo.Password == "" {
		return fmt.Errorf("должен быть указан пароль")
	}

	hashedPass, err := s.crypto.GenerateHashPass(authInfo.Password)
	if err != nil {
		return fmt.Errorf("генерация хэша: %w", err)
	}

	authInfo.HashedPass = hashedPass

	ctx := context.Background()

	err = s.authRepo.Register(ctx, authInfo)
	if err != nil {
		return fmt.Errorf("регистрация пользователя: %w", err)
	}

	return nil
}

func (s *Service) Login(authInfo *domain.UserAuth) (token string, err error) {
	if authInfo.Username == "" {
		return "", fmt.Errorf("должно быть указано имя пользователя")
	}

	if authInfo.Password == "" {
		return "", fmt.Errorf("должен быть указан пароль")
	}

	ctx := context.Background()

	userAuth, err := s.authRepo.GetByUsername(ctx, authInfo.Username)
	if err != nil {
		return "", fmt.Errorf("получение пользователя по username: %w", err)
	}

	if !s.crypto.CheckPasswordHash(authInfo.Password, userAuth.HashedPass) {
		return "", fmt.Errorf("неверный пароль")
	}

	token, err = base.GenerateAuthToken(authInfo.Username, s.jwtKey, userAuth.Role)
	if err != nil {
		return "", fmt.Errorf("генерация токена: %w", err)
	}

	return token, nil
}
