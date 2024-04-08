package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"os"
	"ppo/domain"
	"testing"
)

var testDbInstance *pgxpool.Pool

func TestMain(m *testing.M) {
	testDB := SetupTestDatabase()
	testDbInstance = testDB.DbInstance
	defer testDB.TearDown()
	os.Exit(m.Run())
}

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
			username: "test123",
			expected: &domain.UserAuth{
				Username:   "test123",
				HashedPass: "test123",
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
