package auth

import (
	"context"
	"fmt"
	"ppo/domain"
	"ppo/pkg/base"
	"ppo/pkg/logger"
)

type Service struct {
	authRepo domain.IAuthRepository
	crypto   base.IHashCrypto
	jwtKey   string
	logger   logger.ILogger
}

func NewService(
	repo domain.IAuthRepository,
	crypto base.IHashCrypto,
	jwtKey string,
	logger logger.ILogger,
) domain.IAuthService {
	return &Service{
		authRepo: repo,
		crypto:   crypto,
		jwtKey:   jwtKey,
		logger:   logger,
	}
}

func (s *Service) Register(ctx context.Context, authInfo *domain.UserAuth) (err error) {
	prompt := "AuthRegister"
	if authInfo.Username == "" {
		s.logger.Infof("%s: должно быть указано имя пользователя", prompt)
		return fmt.Errorf("должно быть указано имя пользователя")
	}

	if authInfo.Password == "" {
		s.logger.Infof("%s: должен быть указан пароль", prompt)
		return fmt.Errorf("должен быть указан пароль")
	}

	hashedPass, err := s.crypto.GenerateHashPass(authInfo.Password)
	if err != nil {
		s.logger.Infof("%s: генерация хэша: %v", prompt, err)
		return fmt.Errorf("генерация хэша: %w", err)
	}

	authInfo.HashedPass = hashedPass

	err = s.authRepo.Register(ctx, authInfo)
	if err != nil {
		s.logger.Infof("%s: регистрация пользователя: %v", prompt, err)
		return fmt.Errorf("регистрация пользователя: %w", err)
	}

	return nil
}

func (s *Service) Login(ctx context.Context, authInfo *domain.UserAuth) (token string, err error) {
	prompt := "AuthLogin"

	if authInfo.Username == "" {
		s.logger.Infof("%s: должно быть указано имя пользователя", prompt)
		return "", fmt.Errorf("должно быть указано имя пользователя")
	}

	if authInfo.Password == "" {
		s.logger.Infof("%s: должен быть указан пароль", prompt)
		return "", fmt.Errorf("должен быть указан пароль")
	}

	userAuth, err := s.authRepo.GetByUsername(ctx, authInfo.Username)
	if err != nil {
		s.logger.Infof("%s: получение пользователя по username: %v", prompt, err)
		return "", fmt.Errorf("получение пользователя по username: %w", err)
	}

	if !s.crypto.CheckPasswordHash(authInfo.Password, userAuth.HashedPass) {
		s.logger.Infof("%s: неверный пароль", prompt)
		return "", fmt.Errorf("неверный пароль")
	}

	token, err = base.GenerateAuthToken(userAuth.ID.String(), s.jwtKey, userAuth.Role)
	if err != nil {
		s.logger.Infof("%s: генерация токена: %v", prompt, err)
		return "", fmt.Errorf("генерация токена: %w", err)
	}

	return token, nil
}
