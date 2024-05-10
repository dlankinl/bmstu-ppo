package user_skill

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
)

func TestUserSkillService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userSkillRepo := mocks.NewMockIUserSkillRepository(ctrl)
	userRepo := mocks.NewMockIUserRepository(ctrl)
	skillRepo := mocks.NewMockISkillRepository(ctrl)
	svc := NewService(userSkillRepo, userRepo, skillRepo)

	testCases := []struct {
		name       string
		pair       *domain.UserSkill
		beforeTest func(userSkillRepo mocks.MockIUserSkillRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное добавление",
			pair: &domain.UserSkill{
				UserId:  uuid.UUID{1},
				SkillId: uuid.UUID{1},
			},
			beforeTest: func(userSkillRepo mocks.MockIUserSkillRepository) {
				userSkillRepo.EXPECT().
					Create(
						context.Background(),
						&domain.UserSkill{
							UserId:  uuid.UUID{1},
							SkillId: uuid.UUID{1},
						},
					).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ошибка выполнения запроса в репозитории",
			pair: &domain.UserSkill{
				UserId:  uuid.UUID{1},
				SkillId: uuid.UUID{1},
			},
			beforeTest: func(userSkillRepo mocks.MockIUserSkillRepository) {
				userSkillRepo.EXPECT().
					Create(
						context.Background(),
						&domain.UserSkill{
							UserId:  uuid.UUID{1},
							SkillId: uuid.UUID{1},
						},
					).Return(fmt.Errorf("sql error")).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("связывание пользователя и навыка: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*userSkillRepo)
			}

			err := svc.Create(ctx, tc.pair)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestUserSkillService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userSkillRepo := mocks.NewMockIUserSkillRepository(ctrl)
	userRepo := mocks.NewMockIUserRepository(ctrl)
	skillRepo := mocks.NewMockISkillRepository(ctrl)
	svc := NewService(userSkillRepo, userRepo, skillRepo)

	testCases := []struct {
		name       string
		pair       *domain.UserSkill
		beforeTest func(userSkillRepo mocks.MockIUserSkillRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное удаление",
			pair: &domain.UserSkill{
				UserId:  uuid.UUID{1},
				SkillId: uuid.UUID{1},
			},
			beforeTest: func(userSkillRepo mocks.MockIUserSkillRepository) {
				userSkillRepo.EXPECT().
					Delete(
						context.Background(),
						&domain.UserSkill{
							UserId:  uuid.UUID{1},
							SkillId: uuid.UUID{1},
						},
					).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ошибка выполнения запроса в репозитории",
			pair: &domain.UserSkill{
				UserId:  uuid.UUID{1},
				SkillId: uuid.UUID{1},
			},
			beforeTest: func(userSkillRepo mocks.MockIUserSkillRepository) {
				userSkillRepo.EXPECT().
					Delete(
						context.Background(),
						&domain.UserSkill{
							UserId:  uuid.UUID{1},
							SkillId: uuid.UUID{1},
						},
					).Return(fmt.Errorf("sql error")).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("удаление связи пользователь-навык: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*userSkillRepo)
			}

			err := svc.Delete(ctx, tc.pair)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestUserSkillService_GetSkillsForUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userSkillRepo := mocks.NewMockIUserSkillRepository(ctrl)
	userRepo := mocks.NewMockIUserRepository(ctrl)
	skillRepo := mocks.NewMockISkillRepository(ctrl)
	svc := NewService(userSkillRepo, userRepo, skillRepo)

	testCases := []struct {
		name       string
		pairs      []*domain.UserSkill
		beforeTest func(userSkillRepo mocks.MockIUserSkillRepository, userRepo mocks.MockIUserRepository, skillRepo mocks.MockISkillRepository)
		expected   []*domain.Skill
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное получение набора пользователь-навык",
			pairs: []*domain.UserSkill{
				{
					UserId:  uuid.UUID{1},
					SkillId: uuid.UUID{1},
				},
				{
					UserId:  uuid.UUID{1},
					SkillId: uuid.UUID{2},
				},
				{
					UserId:  uuid.UUID{1},
					SkillId: uuid.UUID{3},
				},
			},
			beforeTest: func(userSkillRepo mocks.MockIUserSkillRepository, userRepo mocks.MockIUserRepository, skillRepo mocks.MockISkillRepository) {
				userSkillRepo.EXPECT().
					GetUserSkillsByUserId(
						context.Background(),
						uuid.UUID{1},
						1,
						true,
					).
					Return([]*domain.UserSkill{
						{
							UserId:  uuid.UUID{1},
							SkillId: uuid.UUID{1},
						},
						{
							UserId:  uuid.UUID{1},
							SkillId: uuid.UUID{2},
						},
						{
							UserId:  uuid.UUID{1},
							SkillId: uuid.UUID{3},
						},
					}, nil)

				skillRepo.EXPECT().
					GetById(
						context.Background(),
						uuid.UUID{1},
					).
					Return(&domain.Skill{ID: uuid.UUID{1}, Name: "a", Description: "a"}, nil)

				skillRepo.EXPECT().
					GetById(
						context.Background(),
						uuid.UUID{2},
					).
					Return(&domain.Skill{ID: uuid.UUID{2}, Name: "b", Description: "b"}, nil)

				skillRepo.EXPECT().
					GetById(
						context.Background(),
						uuid.UUID{3},
					).
					Return(&domain.Skill{ID: uuid.UUID{3}, Name: "c", Description: "c"}, nil)
			},
			expected: []*domain.Skill{
				{
					ID:          uuid.UUID{1},
					Name:        "a",
					Description: "a",
				},
				{
					ID:          uuid.UUID{2},
					Name:        "b",
					Description: "b",
				},
				{
					ID:          uuid.UUID{3},
					Name:        "c",
					Description: "c",
				},
			},
			wantErr: false,
		},
		{
			name: "ошибка при получении данных из репозитория",
			pairs: []*domain.UserSkill{
				{
					UserId:  uuid.UUID{1},
					SkillId: uuid.UUID{1},
				},
			},
			beforeTest: func(userSkillRepo mocks.MockIUserSkillRepository, userRepo mocks.MockIUserRepository, skillRepo mocks.MockISkillRepository) {
				userSkillRepo.EXPECT().
					GetUserSkillsByUserId(
						context.Background(),
						uuid.UUID{1},
						1,
						true,
					).
					Return([]*domain.UserSkill{
						{
							UserId:  uuid.UUID{1},
							SkillId: uuid.UUID{1},
						},
					}, nil)

				skillRepo.EXPECT().
					GetById(
						context.Background(),
						uuid.UUID{1},
					).
					Return(nil, fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("получение скилла по skillId: sql error"),
		},
		{
			name: "ошибка при получении данных из репозитория_2",
			pairs: []*domain.UserSkill{
				{
					UserId:  uuid.UUID{1},
					SkillId: uuid.UUID{1},
				},
			},
			beforeTest: func(userSkillRepo mocks.MockIUserSkillRepository, userRepo mocks.MockIUserRepository, skillRepo mocks.MockISkillRepository) {
				userSkillRepo.EXPECT().
					GetUserSkillsByUserId(
						context.Background(),
						uuid.UUID{1},
						1,
						true,
					).
					Return(nil, fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("получение связок пользователь-навык по userId: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*userSkillRepo, *userRepo, *skillRepo)
			}

			skills, err := svc.GetSkillsForUser(ctx, uuid.UUID{1}, 1, true)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.Equal(t, tc.expected, skills)
			}
		})
	}
}

func TestUserSkillService_GetUsersForSkill(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userSkillRepo := mocks.NewMockIUserSkillRepository(ctrl)
	userRepo := mocks.NewMockIUserRepository(ctrl)
	skillRepo := mocks.NewMockISkillRepository(ctrl)
	svc := NewService(userSkillRepo, userRepo, skillRepo)

	testCases := []struct {
		name       string
		pairs      []*domain.UserSkill
		beforeTest func(userSkillRepo mocks.MockIUserSkillRepository, userRepo mocks.MockIUserRepository, skillRepo mocks.MockISkillRepository)
		expected   []*domain.User
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное получение набора пользователь-навык",
			pairs: []*domain.UserSkill{
				{
					UserId:  uuid.UUID{1},
					SkillId: uuid.UUID{1},
				},
				{
					UserId:  uuid.UUID{1},
					SkillId: uuid.UUID{2},
				},
				{
					UserId:  uuid.UUID{1},
					SkillId: uuid.UUID{3},
				},
			},
			beforeTest: func(userSkillRepo mocks.MockIUserSkillRepository, userRepo mocks.MockIUserRepository, skillRepo mocks.MockISkillRepository) {
				userSkillRepo.EXPECT().
					GetUserSkillsBySkillId(
						context.Background(),
						uuid.UUID{1},
						1,
					).
					Return([]*domain.UserSkill{
						{
							UserId:  uuid.UUID{1},
							SkillId: uuid.UUID{1},
						},
						{
							UserId:  uuid.UUID{2},
							SkillId: uuid.UUID{1},
						},
						{
							UserId:  uuid.UUID{3},
							SkillId: uuid.UUID{1},
						},
					}, nil)

				userRepo.EXPECT().
					GetById(
						context.Background(),
						uuid.UUID{1},
					).
					Return(&domain.User{ID: uuid.UUID{1}, Username: "a", FullName: "a"}, nil)

				userRepo.EXPECT().
					GetById(
						context.Background(),
						uuid.UUID{2},
					).
					Return(&domain.User{ID: uuid.UUID{2}, Username: "b", FullName: "b"}, nil)

				userRepo.EXPECT().
					GetById(
						context.Background(),
						uuid.UUID{3},
					).
					Return(&domain.User{ID: uuid.UUID{3}, Username: "c", FullName: "c"}, nil)
			},
			expected: []*domain.User{
				{
					ID:       uuid.UUID{1},
					Username: "a",
					FullName: "a",
				},
				{
					ID:       uuid.UUID{2},
					Username: "b",
					FullName: "b",
				},
				{
					ID:       uuid.UUID{3},
					Username: "c",
					FullName: "c",
				},
			},
			wantErr: false,
		},
		{
			name: "ошибка при получении данных из репозитория",
			pairs: []*domain.UserSkill{
				{
					UserId:  uuid.UUID{1},
					SkillId: uuid.UUID{1},
				},
			},
			beforeTest: func(userSkillRepo mocks.MockIUserSkillRepository, userRepo mocks.MockIUserRepository, skillRepo mocks.MockISkillRepository) {
				userSkillRepo.EXPECT().
					GetUserSkillsBySkillId(
						context.Background(),
						uuid.UUID{1},
						1,
					).
					Return([]*domain.UserSkill{
						{
							UserId:  uuid.UUID{1},
							SkillId: uuid.UUID{1},
						},
					}, nil)

				userRepo.EXPECT().
					GetById(
						context.Background(),
						uuid.UUID{1},
					).
					Return(nil, fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("получение пользователя по userId: sql error"),
		},
		{
			name: "ошибка при получении данных из репозитория_2",
			pairs: []*domain.UserSkill{
				{
					UserId:  uuid.UUID{1},
					SkillId: uuid.UUID{1},
				},
			},
			beforeTest: func(userSkillRepo mocks.MockIUserSkillRepository, userRepo mocks.MockIUserRepository, skillRepo mocks.MockISkillRepository) {
				userSkillRepo.EXPECT().
					GetUserSkillsBySkillId(
						context.Background(),
						uuid.UUID{1},
						1,
					).
					Return(nil, fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("получение связок пользователь-навык по skillId: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*userSkillRepo, *userRepo, *skillRepo)
			}

			users, err := svc.GetUsersForSkill(ctx, uuid.UUID{1}, 1)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.Equal(t, tc.expected, users)
			}
		})
	}
}

func TestUserSkillService_DeleteSkillsForUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userSkillRepo := mocks.NewMockIUserSkillRepository(ctrl)
	userRepo := mocks.NewMockIUserRepository(ctrl)
	skillRepo := mocks.NewMockISkillRepository(ctrl)
	svc := NewService(userSkillRepo, userRepo, skillRepo)

	testCases := []struct {
		name       string
		pairs      []*domain.UserSkill
		beforeTest func(userSkillRepo mocks.MockIUserSkillRepository, userRepo mocks.MockIUserRepository, skillRepo mocks.MockISkillRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное удаление набора пользователь-навык по userId",
			pairs: []*domain.UserSkill{
				{
					UserId:  uuid.UUID{1},
					SkillId: uuid.UUID{1},
				},
				{
					UserId:  uuid.UUID{1},
					SkillId: uuid.UUID{2},
				},
			},
			beforeTest: func(userSkillRepo mocks.MockIUserSkillRepository, userRepo mocks.MockIUserRepository, skillRepo mocks.MockISkillRepository) {
				userSkillRepo.EXPECT().
					GetUserSkillsByUserId(
						context.Background(),
						uuid.UUID{1},
						0,
						false,
					).
					Return([]*domain.UserSkill{
						{
							UserId:  uuid.UUID{1},
							SkillId: uuid.UUID{1},
						},
						{
							UserId:  uuid.UUID{1},
							SkillId: uuid.UUID{2},
						},
					}, nil)

				userSkillRepo.EXPECT().
					Delete(
						context.Background(),
						&domain.UserSkill{UserId: uuid.UUID{1}, SkillId: uuid.UUID{1}},
					).
					Return(nil)

				userSkillRepo.EXPECT().
					Delete(
						context.Background(),
						&domain.UserSkill{UserId: uuid.UUID{1}, SkillId: uuid.UUID{2}},
					).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ошибка выполнении запроса в репозитории",
			pairs: []*domain.UserSkill{
				{
					UserId:  uuid.UUID{1},
					SkillId: uuid.UUID{1},
				},
			},
			beforeTest: func(userSkillRepo mocks.MockIUserSkillRepository, userRepo mocks.MockIUserRepository, skillRepo mocks.MockISkillRepository) {
				userSkillRepo.EXPECT().
					GetUserSkillsByUserId(
						context.Background(),
						uuid.UUID{1},
						0,
						false,
					).
					Return([]*domain.UserSkill{
						{
							UserId:  uuid.UUID{1},
							SkillId: uuid.UUID{1},
						},
						{
							UserId:  uuid.UUID{1},
							SkillId: uuid.UUID{2},
						},
					}, nil)

				userSkillRepo.EXPECT().
					Delete(
						context.Background(),
						&domain.UserSkill{UserId: uuid.UUID{1}, SkillId: uuid.UUID{1}},
					).
					Return(fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("удаление пары пользователь-навык: sql error"),
		},
		{
			name: "ошибка выполнении запроса в репозитории_2",
			pairs: []*domain.UserSkill{
				{
					UserId:  uuid.UUID{1},
					SkillId: uuid.UUID{1},
				},
			},
			beforeTest: func(userSkillRepo mocks.MockIUserSkillRepository, userRepo mocks.MockIUserRepository, skillRepo mocks.MockISkillRepository) {
				userSkillRepo.EXPECT().
					GetUserSkillsByUserId(
						context.Background(),
						uuid.UUID{1},
						0,
						false,
					).
					Return(nil, fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("получение связок пользователь-навык по userId: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*userSkillRepo, *userRepo, *skillRepo)
			}

			err := svc.DeleteSkillsForUser(ctx, uuid.UUID{1})

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}
