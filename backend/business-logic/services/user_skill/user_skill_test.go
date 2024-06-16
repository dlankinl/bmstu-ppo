package user_skill

import (
	"business-logic/domain"
	mocks2 "business-logic/domain/mocks"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestUserSkillService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userSkillRepo := mocks2.NewMockIUserSkillRepository(ctrl)
	userRepo := mocks2.NewMockIUserRepository(ctrl)
	skillRepo := mocks2.NewMockISkillRepository(ctrl)
	svc := NewService(userSkillRepo, userRepo, skillRepo)

	testCases := []struct {
		name       string
		pair       *domain.UserSkill
		beforeTest func(userSkillRepo mocks2.MockIUserSkillRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное добавление",
			pair: &domain.UserSkill{
				UserId:  uuid.UUID{1},
				SkillId: uuid.UUID{1},
			},
			beforeTest: func(userSkillRepo mocks2.MockIUserSkillRepository) {
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
			beforeTest: func(userSkillRepo mocks2.MockIUserSkillRepository) {
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
			if tc.beforeTest != nil {
				tc.beforeTest(*userSkillRepo)
			}

			err := svc.Create(tc.pair)

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

	userSkillRepo := mocks2.NewMockIUserSkillRepository(ctrl)
	userRepo := mocks2.NewMockIUserRepository(ctrl)
	skillRepo := mocks2.NewMockISkillRepository(ctrl)
	svc := NewService(userSkillRepo, userRepo, skillRepo)

	testCases := []struct {
		name       string
		pair       *domain.UserSkill
		beforeTest func(userSkillRepo mocks2.MockIUserSkillRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное удаление",
			pair: &domain.UserSkill{
				UserId:  uuid.UUID{1},
				SkillId: uuid.UUID{1},
			},
			beforeTest: func(userSkillRepo mocks2.MockIUserSkillRepository) {
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
			beforeTest: func(userSkillRepo mocks2.MockIUserSkillRepository) {
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
			if tc.beforeTest != nil {
				tc.beforeTest(*userSkillRepo)
			}

			err := svc.Delete(tc.pair)

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

	userSkillRepo := mocks2.NewMockIUserSkillRepository(ctrl)
	userRepo := mocks2.NewMockIUserRepository(ctrl)
	skillRepo := mocks2.NewMockISkillRepository(ctrl)
	svc := NewService(userSkillRepo, userRepo, skillRepo)

	testCases := []struct {
		name       string
		pairs      []*domain.UserSkill
		beforeTest func(userSkillRepo mocks2.MockIUserSkillRepository, userRepo mocks2.MockIUserRepository, skillRepo mocks2.MockISkillRepository)
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
			beforeTest: func(userSkillRepo mocks2.MockIUserSkillRepository, userRepo mocks2.MockIUserRepository, skillRepo mocks2.MockISkillRepository) {
				userSkillRepo.EXPECT().
					GetUserSkillsByUserId(
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
			beforeTest: func(userSkillRepo mocks2.MockIUserSkillRepository, userRepo mocks2.MockIUserRepository, skillRepo mocks2.MockISkillRepository) {
				userSkillRepo.EXPECT().
					GetUserSkillsByUserId(
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
			beforeTest: func(userSkillRepo mocks2.MockIUserSkillRepository, userRepo mocks2.MockIUserRepository, skillRepo mocks2.MockISkillRepository) {
				userSkillRepo.EXPECT().
					GetUserSkillsByUserId(
						context.Background(),
						uuid.UUID{1},
						1,
					).
					Return(nil, fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("получение связок пользователь-навык по userId: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.beforeTest != nil {
				tc.beforeTest(*userSkillRepo, *userRepo, *skillRepo)
			}

			skills, err := svc.GetSkillsForUser(uuid.UUID{1}, 1)

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

	userSkillRepo := mocks2.NewMockIUserSkillRepository(ctrl)
	userRepo := mocks2.NewMockIUserRepository(ctrl)
	skillRepo := mocks2.NewMockISkillRepository(ctrl)
	svc := NewService(userSkillRepo, userRepo, skillRepo)

	testCases := []struct {
		name       string
		pairs      []*domain.UserSkill
		beforeTest func(userSkillRepo mocks2.MockIUserSkillRepository, userRepo mocks2.MockIUserRepository, skillRepo mocks2.MockISkillRepository)
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
			beforeTest: func(userSkillRepo mocks2.MockIUserSkillRepository, userRepo mocks2.MockIUserRepository, skillRepo mocks2.MockISkillRepository) {
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
			beforeTest: func(userSkillRepo mocks2.MockIUserSkillRepository, userRepo mocks2.MockIUserRepository, skillRepo mocks2.MockISkillRepository) {
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
			beforeTest: func(userSkillRepo mocks2.MockIUserSkillRepository, userRepo mocks2.MockIUserRepository, skillRepo mocks2.MockISkillRepository) {
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
			if tc.beforeTest != nil {
				tc.beforeTest(*userSkillRepo, *userRepo, *skillRepo)
			}

			users, err := svc.GetUsersForSkill(uuid.UUID{1}, 1)

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

	userSkillRepo := mocks2.NewMockIUserSkillRepository(ctrl)
	userRepo := mocks2.NewMockIUserRepository(ctrl)
	skillRepo := mocks2.NewMockISkillRepository(ctrl)
	svc := NewService(userSkillRepo, userRepo, skillRepo)

	testCases := []struct {
		name       string
		pairs      []*domain.UserSkill
		beforeTest func(userSkillRepo mocks2.MockIUserSkillRepository, userRepo mocks2.MockIUserRepository, skillRepo mocks2.MockISkillRepository)
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
			beforeTest: func(userSkillRepo mocks2.MockIUserSkillRepository, userRepo mocks2.MockIUserRepository, skillRepo mocks2.MockISkillRepository) {
				userSkillRepo.EXPECT().
					GetUserSkillsByUserId(
						context.Background(),
						uuid.UUID{1},
						0,
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
			beforeTest: func(userSkillRepo mocks2.MockIUserSkillRepository, userRepo mocks2.MockIUserRepository, skillRepo mocks2.MockISkillRepository) {
				userSkillRepo.EXPECT().
					GetUserSkillsByUserId(
						context.Background(),
						uuid.UUID{1},
						0,
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
			beforeTest: func(userSkillRepo mocks2.MockIUserSkillRepository, userRepo mocks2.MockIUserRepository, skillRepo mocks2.MockISkillRepository) {
				userSkillRepo.EXPECT().
					GetUserSkillsByUserId(
						context.Background(),
						uuid.UUID{1},
						0,
					).
					Return(nil, fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("получение связок пользователь-навык по userId: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.beforeTest != nil {
				tc.beforeTest(*userSkillRepo, *userRepo, *skillRepo)
			}

			err := svc.DeleteSkillsForUser(uuid.UUID{1})

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}
