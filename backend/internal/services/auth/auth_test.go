package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"ppo/domain"
	"ppo/mocks"
	"ppo/pkg/base"
	"testing"
)

func TestAuthService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jwtKey := "abcdefgh123"
	repo := mocks.NewMockIAuthRepository(ctrl)
	crypto := mocks.NewMockIHashCrypto(ctrl)
	svc := NewService(repo, crypto, jwtKey)

	testCases := []struct {
		name       string
		authInfo   *domain.UserAuth
		beforeTest func(authRepo mocks.MockIAuthRepository, crypto mocks.MockIHashCrypto)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешная аутентификация",
			authInfo: &domain.UserAuth{
				Username: "test123",
				Password: "pass123",
			},
			beforeTest: func(authRepo mocks.MockIAuthRepository, crypto mocks.MockIHashCrypto) {
				authRepo.EXPECT().
					GetByUsername(
						context.Background(),
						"test123",
					).
					Return(&domain.UserAuth{
						Username:   "test123",
						Password:   "pass123",
						HashedPass: "hashedPass123",
					}, nil)

				crypto.EXPECT().
					CheckPasswordHash("pass123", "hashedPass123").
					Return(true)
			},
			wantErr: false,
		},
		{
			name: "пустое имя пользователя",
			authInfo: &domain.UserAuth{
				Username: "",
				Password: "pass123",
			},
			wantErr: true,
			errStr:  errors.New("должно быть указано имя пользователя"),
		},
		{
			name: "ошибка получения данных из репозитория",
			authInfo: &domain.UserAuth{
				Username: "test123",
				Password: "pass123",
			},
			beforeTest: func(authRepo mocks.MockIAuthRepository, crypto mocks.MockIHashCrypto) {
				authRepo.EXPECT().
					GetByUsername(
						context.Background(),
						"test123",
					).
					Return(nil, fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("получение пользователя по username: sql error"),
		},
		{
			name: "неверный пароль",
			authInfo: &domain.UserAuth{
				Username: "test123",
				Password: "pass123",
			},
			beforeTest: func(authRepo mocks.MockIAuthRepository, crypto mocks.MockIHashCrypto) {
				authRepo.EXPECT().
					GetByUsername(
						context.Background(),
						"test123",
					).
					Return(&domain.UserAuth{
						Username:   "test123",
						Password:   "pass123",
						HashedPass: "hashedPass123",
					}, nil)

				crypto.EXPECT().
					CheckPasswordHash("pass123", "hashedPass123").
					Return(false)
			},
			wantErr: true,
			errStr:  errors.New("неверный пароль"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*repo, *crypto)
			}

			token, err := svc.Login(ctx, tc.authInfo)
			_, verifErr := base.VerifyAuthToken(token, jwtKey)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.Nil(t, verifErr)
			}
		})
	}
}

func TestAuthService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockIAuthRepository(ctrl)
	crypto := mocks.NewMockIHashCrypto(ctrl)
	svc := NewService(repo, crypto, "abcdefgh123")

	testCases := []struct {
		name       string
		authInfo   *domain.UserAuth
		beforeTest func(authRepo mocks.MockIAuthRepository, crypto mocks.MockIHashCrypto)
		expected   *domain.UserAuth
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешная регистрация",
			authInfo: &domain.UserAuth{
				Username: "test123",
				Password: "pass123",
			},
			beforeTest: func(authRepo mocks.MockIAuthRepository, crypto mocks.MockIHashCrypto) {
				crypto.EXPECT().
					GenerateHashPass("pass123").
					Return("hashedPass123", nil)

				authRepo.EXPECT().
					Register(
						context.Background(),
						&domain.UserAuth{
							Username:   "test123",
							Password:   "pass123",
							HashedPass: "hashedPass123",
						},
					).
					Return(nil)
			},
			expected: &domain.UserAuth{
				Username:   "test123",
				Password:   "pass123",
				HashedPass: "hashedPass123",
			},
			wantErr: false,
		},
		{
			name: "пустое имя пользователя",
			authInfo: &domain.UserAuth{
				Username: "",
				Password: "pass123",
			},
			wantErr: true,
			errStr:  errors.New("должно быть указано имя пользователя"),
		},
		{
			name: "пустое имя пользователя",
			authInfo: &domain.UserAuth{
				Username: "test123",
				Password: "",
			},
			wantErr: true,
			errStr:  errors.New("должен быть указан пароль"),
		},
		{
			name: "ошибка выполнения запроса в репозитории",
			authInfo: &domain.UserAuth{
				Username: "test123",
				Password: "pass123",
			},
			beforeTest: func(authRepo mocks.MockIAuthRepository, crypto mocks.MockIHashCrypto) {
				crypto.EXPECT().
					GenerateHashPass("pass123").
					Return("hashedPass123", nil)

				authRepo.EXPECT().
					Register(
						context.Background(),
						&domain.UserAuth{
							Username:   "test123",
							Password:   "pass123",
							HashedPass: "hashedPass123",
						},
					).
					Return(fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("регистрация пользователя: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*repo, *crypto)
			}

			err := svc.Register(ctx, tc.authInfo)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}
