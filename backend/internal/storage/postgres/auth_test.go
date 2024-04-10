package postgres

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"ppo/domain"
	"testing"
)

func TestAuthRepository_Register(t *testing.T) {
	repo := NewAuthRepository(testDbInstance)

	testCases := []struct {
		name     string
		authInfo *domain.UserAuth
		wantErr  bool
		errStr   error
	}{
		{
			name: "успех",
			authInfo: &domain.UserAuth{
				Username:   "test123",
				HashedPass: "test123",
			},
			wantErr: false,
		},
		{
			name: "неуникальное имя пользователь",
			authInfo: &domain.UserAuth{
				Username:   "test123",
				HashedPass: "test123",
			},
			wantErr: true,
			errStr: errors.New("регистрация пользователя: ERROR: duplicate key value violates unique " +
				"constraint \"users_username_key\" (SQLSTATE 23505)"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.Register(context.Background(), tc.authInfo)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestAuthRepository_GetByUsername(t *testing.T) {
	repo := NewAuthRepository(testDbInstance)

	testCases := []struct {
		name     string
		username string
		expected *domain.UserAuth
		wantErr  bool
		errStr   error
	}{
		{
			name:     "пользователь найден",
			username: "user3",
			expected: &domain.UserAuth{
				Username:   "user3",
				HashedPass: "user3hehe",
			},
			wantErr: false,
		},
		{
			name:     "пользователь не найден",
			username: "test1234",
			wantErr:  true,
			errStr:   errors.New("получение пользователя по username: no rows in result set"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := repo.GetByUsername(context.Background(), tc.username)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.Equal(t, tc.expected.HashedPass, res.HashedPass)
			}
		})
	}
}
