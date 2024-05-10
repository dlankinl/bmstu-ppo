package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"ppo/domain"
	"ppo/mocks"
	"testing"
	"time"
)

func TestUserService_DeleteById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockIUserRepository(ctrl)
	compRepo := mocks.NewMockICompanyRepository(ctrl)
	actFieldRepo := mocks.NewMockIActivityFieldRepository(ctrl)
	svc := NewService(userRepo, compRepo, actFieldRepo)

	curUuid := uuid.New()

	testCases := []struct {
		name       string
		id         uuid.UUID
		beforeTest func(userRepo mocks.MockIUserRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное удаление",
			id:   curUuid,
			beforeTest: func(userRepo mocks.MockIUserRepository) {
				userRepo.EXPECT().
					DeleteById(context.Background(), curUuid).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ошибка выполнения запроса в репозитории",
			id:   curUuid,
			beforeTest: func(userRepo mocks.MockIUserRepository) {
				userRepo.EXPECT().
					DeleteById(context.Background(), curUuid).
					Return(fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("удаление пользователя по id: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*userRepo)
			}

			err := svc.DeleteById(ctx, tc.id)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestUserService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockIUserRepository(ctrl)
	compRepo := mocks.NewMockICompanyRepository(ctrl)
	actFieldRepo := mocks.NewMockIActivityFieldRepository(ctrl)
	svc := NewService(userRepo, compRepo, actFieldRepo)

	testCases := []struct {
		name       string
		beforeTest func(userRepo mocks.MockIUserRepository)
		expected   []*domain.User
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное получение списка всех компаний",
			beforeTest: func(userRepo mocks.MockIUserRepository) {
				userRepo.EXPECT().
					GetAll(context.Background(), 1).
					Return([]*domain.User{
						{
							ID:       uuid.UUID{1},
							Username: "a",
							FullName: "a",
							Gender:   "m",
							Birthday: time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
							City:     "a",
						},
						{
							ID:       uuid.UUID{2},
							Username: "b",
							FullName: "b",
							Gender:   "w",
							Birthday: time.Date(2, 2, 2, 2, 2, 2, 2, time.Local),
							City:     "b",
						},
						{
							ID:       uuid.UUID{3},
							Username: "c",
							FullName: "c",
							Gender:   "m",
							Birthday: time.Date(3, 3, 3, 3, 3, 3, 3, time.Local),
							City:     "c",
						},
					}, nil)
			},
			expected: []*domain.User{
				{
					ID:       uuid.UUID{1},
					Username: "a",
					FullName: "a",
					Gender:   "m",
					Birthday: time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
					City:     "a",
				},
				{
					ID:       uuid.UUID{2},
					Username: "b",
					FullName: "b",
					Gender:   "w",
					Birthday: time.Date(2, 2, 2, 2, 2, 2, 2, time.Local),
					City:     "b",
				},
				{
					ID:       uuid.UUID{3},
					Username: "c",
					FullName: "c",
					Gender:   "m",
					Birthday: time.Date(3, 3, 3, 3, 3, 3, 3, time.Local),
					City:     "c",
				},
			},
			wantErr: false,
		},
		{
			name: "ошибка получения данных в репозитории",
			beforeTest: func(userRepo mocks.MockIUserRepository) {
				userRepo.EXPECT().
					GetAll(context.Background(), 1).
					Return(nil, fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("получение списка всех пользователей: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*userRepo)
			}

			users, err := svc.GetAll(ctx, 1)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.Equal(t, users, tc.expected)
			}
		})
	}
}

func TestUserService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockIUserRepository(ctrl)
	compRepo := mocks.NewMockICompanyRepository(ctrl)
	actFieldRepo := mocks.NewMockIActivityFieldRepository(ctrl)
	svc := NewService(userRepo, compRepo, actFieldRepo)

	testCases := []struct {
		name       string
		user       *domain.User
		beforeTest func(userRepo mocks.MockIUserRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное добавление",
			user: &domain.User{
				ID:       uuid.UUID{1},
				Username: "a",
				FullName: "a b c",
				Gender:   "m",
				Birthday: time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
				City:     "a",
			},
			beforeTest: func(userRepo mocks.MockIUserRepository) {
				userRepo.EXPECT().
					Create(
						context.Background(),
						&domain.User{
							ID:       uuid.UUID{1},
							Username: "a",
							FullName: "a b c",
							Gender:   "m",
							Birthday: time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
							City:     "a",
						},
					).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "некорректное ФИО (не 3 слова)",
			user: &domain.User{
				ID:       uuid.UUID{1},
				Username: "a",
				FullName: "a",
				Gender:   "m",
				Birthday: time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
				City:     "a",
			},
			beforeTest: func(userRepo mocks.MockIUserRepository) {
				compRepo.EXPECT().
					Create(
						context.Background(),
						&domain.User{
							ID:       uuid.UUID{1},
							Username: "a",
							FullName: "a",
							Gender:   "m",
							Birthday: time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
							City:     "a",
						},
					).Return(nil).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("некорректное количество слов (должны быть фамилия, имя и отчество)"),
		},
		{
			name: "пустое название города",
			user: &domain.User{
				ID:       uuid.UUID{1},
				Username: "a",
				FullName: "a b c",
				Gender:   "m",
				Birthday: time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
				City:     "",
			},
			beforeTest: func(userRepo mocks.MockIUserRepository) {
				compRepo.EXPECT().
					Create(
						context.Background(),
						&domain.User{
							ID:       uuid.UUID{1},
							Username: "a",
							FullName: "a b c",
							Gender:   "m",
							Birthday: time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
							City:     "",
						},
					).Return(nil).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("должно быть указано название города"),
		},
		{
			name: "неизвестный пол",
			user: &domain.User{
				ID:       uuid.UUID{1},
				Username: "a",
				FullName: "a b c",
				Gender:   "r",
				Birthday: time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
				City:     "c",
			},
			beforeTest: func(userRepo mocks.MockIUserRepository) {
				compRepo.EXPECT().
					Create(
						context.Background(),
						&domain.User{
							ID:       uuid.UUID{1},
							Username: "a",
							FullName: "a b c",
							Gender:   "r",
							Birthday: time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
							City:     "c",
						},
					).Return(nil).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("неизвестный пол"),
		},
		{
			name: "пустая дата рождения",
			user: &domain.User{
				ID:       uuid.UUID{1},
				Username: "a",
				FullName: "a b c",
				Gender:   "w",
				Birthday: time.Time{},
				City:     "c",
			},
			beforeTest: func(userRepo mocks.MockIUserRepository) {
				compRepo.EXPECT().
					Create(
						context.Background(),
						&domain.User{
							ID:       uuid.UUID{1},
							Username: "a",
							FullName: "a b c",
							Gender:   "w",
							Birthday: time.Time{},
							City:     "c",
						},
					).Return(nil).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("должна быть указана дата рождения"),
		},
		{
			name: "ошибка выполнения запроса в репозитории",
			user: &domain.User{
				ID:       uuid.UUID{1},
				Username: "a",
				FullName: "a b c",
				Gender:   "w",
				Birthday: time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
				City:     "c",
			},
			beforeTest: func(userRepo mocks.MockIUserRepository) {
				userRepo.EXPECT().
					Create(
						context.Background(),
						&domain.User{
							ID:       uuid.UUID{1},
							Username: "a",
							FullName: "a b c",
							Gender:   "w",
							Birthday: time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
							City:     "c",
						},
					).Return(fmt.Errorf("sql error")).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("создание пользователя: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*userRepo)
			}

			err := svc.Create(ctx, tc.user)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestUserService_GetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockIUserRepository(ctrl)
	compRepo := mocks.NewMockICompanyRepository(ctrl)
	actFieldRepo := mocks.NewMockIActivityFieldRepository(ctrl)
	svc := NewService(userRepo, compRepo, actFieldRepo)

	testCases := []struct {
		name       string
		id         uuid.UUID
		beforeTest func(userRepo mocks.MockIUserRepository)
		expected   *domain.User
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное получение компании по id",
			id:   uuid.UUID{1},
			beforeTest: func(userRepo mocks.MockIUserRepository) {
				userRepo.EXPECT().
					GetById(
						context.Background(),
						uuid.UUID{1},
					).
					Return(&domain.User{
						ID:       uuid.UUID{1},
						Username: "a",
						FullName: "a b c",
						Gender:   "m",
						Birthday: time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
						City:     "a",
					}, nil)
			},
			expected: &domain.User{
				ID:       uuid.UUID{1},
				Username: "a",
				FullName: "a b c",
				Gender:   "m",
				Birthday: time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
				City:     "a",
			},
			wantErr: false,
		},
		{
			name: "ошибка получения данных в репозитории",
			id:   uuid.UUID{1},
			beforeTest: func(userRepo mocks.MockIUserRepository) {
				userRepo.EXPECT().
					GetById(
						context.Background(),
						uuid.UUID{1},
					).
					Return(nil, fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("получение пользователя по id: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*userRepo)
			}

			user, err := svc.GetById(ctx, tc.id)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.Equal(t, user, tc.expected)
			}
		})
	}
}

func TestUserService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockIUserRepository(ctrl)
	compRepo := mocks.NewMockICompanyRepository(ctrl)
	actFieldRepo := mocks.NewMockIActivityFieldRepository(ctrl)
	svc := NewService(userRepo, compRepo, actFieldRepo)

	testCases := []struct {
		name       string
		user       *domain.User
		beforeTest func(userRepo mocks.MockIUserRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное обновление",
			user: &domain.User{
				ID:   uuid.UUID{1},
				City: "a",
				Role: "admin",
			},
			beforeTest: func(userRepo mocks.MockIUserRepository) {
				userRepo.EXPECT().
					Update(
						context.Background(),
						&domain.User{
							ID:   uuid.UUID{1},
							City: "a",
							Role: "admin",
						},
					).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ошибка выполнения запроса в репозитории",
			user: &domain.User{
				ID:   uuid.UUID{1},
				City: "a",
				Role: "admin",
			},
			beforeTest: func(userRepo mocks.MockIUserRepository) {
				userRepo.EXPECT().
					Update(
						context.Background(),
						&domain.User{
							ID:   uuid.UUID{1},
							City: "a",
							Role: "admin",
						},
					).Return(fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("обновление информации о пользователе: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*userRepo)
			}

			err := svc.Update(ctx, tc.user)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}
