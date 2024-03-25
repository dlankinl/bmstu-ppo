package services

import (
	"context"
	"errors"
	"fmt"
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
	svc := NewUserSkillService(userSkillRepo, userRepo, skillRepo)

	testCases := []struct {
		name       string
		pair       domain.UserSkill
		beforeTest func(userSkillRepo mocks.MockIUserSkillRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное добавление",
			pair: domain.UserSkill{
				UserId:  [16]byte{1},
				SkillId: [16]byte{1},
			},
			beforeTest: func(userSkillRepo mocks.MockIUserSkillRepository) {
				userSkillRepo.EXPECT().
					Create(
						context.Background(),
						&domain.UserSkill{
							UserId:  [16]byte{1},
							SkillId: [16]byte{1},
						},
					).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ошибка выполнения запроса в репозитории",
			pair: domain.UserSkill{
				UserId:  [16]byte{1},
				SkillId: [16]byte{1},
			},
			beforeTest: func(userSkillRepo mocks.MockIUserSkillRepository) {
				userSkillRepo.EXPECT().
					Create(
						context.Background(),
						&domain.UserSkill{
							UserId:  [16]byte{1},
							SkillId: [16]byte{1},
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

			err := svc.Create(context.Background(), &tc.pair)

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
	svc := NewUserSkillService(userSkillRepo, userRepo, skillRepo)

	testCases := []struct {
		name       string
		pair       domain.UserSkill
		beforeTest func(userSkillRepo mocks.MockIUserSkillRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное добавление",
			pair: domain.UserSkill{
				UserId:  [16]byte{1},
				SkillId: [16]byte{1},
			},
			beforeTest: func(userSkillRepo mocks.MockIUserSkillRepository) {
				userSkillRepo.EXPECT().
					Delete(
						context.Background(),
						&domain.UserSkill{
							UserId:  [16]byte{1},
							SkillId: [16]byte{1},
						},
					).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ошибка выполнения запроса в репозитории",
			pair: domain.UserSkill{
				UserId:  [16]byte{1},
				SkillId: [16]byte{1},
			},
			beforeTest: func(userSkillRepo mocks.MockIUserSkillRepository) {
				userSkillRepo.EXPECT().
					Delete(
						context.Background(),
						&domain.UserSkill{
							UserId:  [16]byte{1},
							SkillId: [16]byte{1},
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

			err := svc.Delete(context.Background(), &tc.pair)

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
	svc := NewUserSkillService(userSkillRepo, userRepo, skillRepo)

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
					UserId:  [16]byte{1},
					SkillId: [16]byte{1},
				},
				{
					UserId:  [16]byte{1},
					SkillId: [16]byte{2},
				},
				{
					UserId:  [16]byte{1},
					SkillId: [16]byte{3},
				},
				{
					UserId:  [16]byte{2},
					SkillId: [16]byte{1},
				},
			},
			beforeTest: func(userSkillRepo mocks.MockIUserSkillRepository, userRepo mocks.MockIUserRepository, skillRepo mocks.MockISkillRepository) {
				userSkillRepo.EXPECT().
					GetUserSkillsByUserId(
						context.Background(),
						[16]byte{1},
					).
					Return([]*domain.UserSkill{
						{
							UserId:  [16]byte{1},
							SkillId: [16]byte{1},
						},
						{
							UserId:  [16]byte{1},
							SkillId: [16]byte{2},
						},
						{
							UserId:  [16]byte{1},
							SkillId: [16]byte{3},
						},
						{
							UserId:  [16]byte{2},
							SkillId: [16]byte{1},
						},
					}, nil)

				skillRepo.EXPECT().
					GetById(
						context.Background(),
						[16]byte{1},
					).
					Return(&domain.Skill{ID: [16]byte{1}, Name: "a", Description: "a"}, nil)

				skillRepo.EXPECT().
					GetById(
						context.Background(),
						[16]byte{2},
					).
					Return(&domain.Skill{ID: [16]byte{2}, Name: "b", Description: "b"}, nil)

				skillRepo.EXPECT().
					GetById(
						context.Background(),
						[16]byte{3},
					).
					Return(&domain.Skill{ID: [16]byte{3}, Name: "c", Description: "c"}, nil)
			},
			expected: []*domain.Skill{
				{
					ID:          [16]byte{1},
					Name:        "a",
					Description: "a",
				},
				{
					ID:          [16]byte{2},
					Name:        "b",
					Description: "b",
				},
				{
					ID:          [16]byte{3},
					Name:        "c",
					Description: "c",
				},
			},
			wantErr: false,
		},
		//{
		//	name: "ошибка выполнения запроса в репозитории",
		//	pairs: domain.Skill{
		//		UserId:  [16]byte{1},
		//		SkillId: [16]byte{1},
		//	},
		//	beforeTest: func(userSkillRepo mocks.MockIUserSkillRepository, skillRepo mocks.MockISkillRepository) {
		//		userSkillRepo.EXPECT().
		//			Create(
		//				context.Background(),
		//				&domain.UserSkill{
		//					UserId:  [16]byte{1},
		//					SkillId: [16]byte{1},
		//				},
		//			).Return(fmt.Errorf("sql error")).
		//			AnyTimes()
		//	},
		//	wantErr: true,
		//	errStr:  errors.New("связывание пользователя и навыка: sql error"),
		//},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.beforeTest != nil {
				tc.beforeTest(*userSkillRepo, *userRepo, *skillRepo)
			}

			skills, err := svc.GetSkillsForUser(context.Background(), [16]byte{1})
			fmt.Println(skills)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.Equal(t, tc.expected, skills)
			}
		})
	}
}

//func TestUserSkillService_GetUsersForSkill(t *testing.T) {
//	type fields struct {
//		userSkillRepo domain.IUserSkillRepository
//		userRepo      domain.IUserRepository
//		skillRepo     domain.ISkillRepository
//	}
//	type args struct {
//		ctx     context.Context
//		skillId uuid.UUID
//	}
//	tests := []struct {
//		name      string
//		fields    fields
//		args      args
//		wantUsers []*domain.User
//		wantErr   bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := UserSkillService{
//				userSkillRepo: tt.fields.userSkillRepo,
//				userRepo:      tt.fields.userRepo,
//				skillRepo:     tt.fields.skillRepo,
//			}
//			gotUsers, err := s.GetUsersForSkill(tt.args.ctx, tt.args.skillId)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetUsersForSkill() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(gotUsers, tt.wantUsers) {
//				t.Errorf("GetUsersForSkill() gotUsers = %v, want %v", gotUsers, tt.wantUsers)
//			}
//		})
//	}
//}
