package skill

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

func TestSkillService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	skillRepo := mocks.NewMockISkillRepository(ctrl)
	svc := NewService(skillRepo)

	testCases := []struct {
		name       string
		skill      *domain.Skill
		beforeTest func(skillRepo mocks.MockISkillRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное добавление",
			skill: &domain.Skill{
				Name:        "aaa",
				Description: "bbb",
			},
			beforeTest: func(skillRepo mocks.MockISkillRepository) {
				skillRepo.EXPECT().
					Create(
						context.Background(),
						&domain.Skill{
							Name:        "aaa",
							Description: "bbb",
						},
					).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "пустое название навыка",
			skill: &domain.Skill{
				Name:        "",
				Description: "aaa",
			},
			beforeTest: func(skillRepo mocks.MockISkillRepository) {
				skillRepo.EXPECT().
					Create(
						context.Background(),
						&domain.Skill{
							Name:        "",
							Description: "aaa",
						},
					).Return(fmt.Errorf("должно быть указано название навыка")).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("должно быть указано название навыка"),
		},
		{
			name: "пустое описание навыка",
			skill: &domain.Skill{
				Name:        "aaa",
				Description: "",
			},
			beforeTest: func(skillRepo mocks.MockISkillRepository) {
				skillRepo.EXPECT().
					Create(
						context.Background(),
						&domain.Skill{
							Name:        "aaa",
							Description: "",
						},
					).Return(fmt.Errorf("должно быть указано описание навыка")).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("должно быть указано описание навыка"),
		},
		{
			name: "ошибка выполнения запроса в репозитории",
			skill: &domain.Skill{
				Name:        "aaa",
				Description: "bbb",
			},
			beforeTest: func(skillRepo mocks.MockISkillRepository) {
				skillRepo.EXPECT().
					Create(
						context.Background(),
						&domain.Skill{
							Name:        "aaa",
							Description: "bbb",
						},
					).Return(fmt.Errorf("sql error")).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("добавление навыка: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*skillRepo)
			}

			err := svc.Create(ctx, tc.skill)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestSkillService_DeleteById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	skillRepo := mocks.NewMockISkillRepository(ctrl)
	svc := NewService(skillRepo)

	curUuid := uuid.New()

	testCases := []struct {
		name       string
		id         uuid.UUID
		beforeTest func(skillRepo mocks.MockISkillRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное удаление",
			id:   curUuid,
			beforeTest: func(skillRepo mocks.MockISkillRepository) {
				skillRepo.EXPECT().
					DeleteById(context.Background(), curUuid).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ошибка выполнения запроса в репозитории",
			id:   curUuid,
			beforeTest: func(skillRepo mocks.MockISkillRepository) {
				skillRepo.EXPECT().
					DeleteById(context.Background(), curUuid).
					Return(fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("удаление навыка по id: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*skillRepo)
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

func TestSkillService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	skillRepo := mocks.NewMockISkillRepository(ctrl)
	svc := NewService(skillRepo)

	testCases := []struct {
		name       string
		beforeTest func(skillRepo mocks.MockISkillRepository)
		expected   []*domain.Skill
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное получение списка всех навыков",
			beforeTest: func(skillRepo mocks.MockISkillRepository) {
				skillRepo.EXPECT().
					GetAll(context.Background(), 1).
					Return([]*domain.Skill{
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
					}, nil)
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
			name: "ошибка получения данных в репозитории",
			beforeTest: func(skillRepo mocks.MockISkillRepository) {
				skillRepo.EXPECT().
					GetAll(context.Background(), 1).
					Return(nil, fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("получение списка всех навыков: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*skillRepo)
			}

			skills, err := svc.GetAll(ctx, 1)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.Equal(t, skills, tc.expected)
			}
		})
	}
}

func TestSkillService_GetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	skillRepo := mocks.NewMockISkillRepository(ctrl)
	svc := NewService(skillRepo)

	testCases := []struct {
		name       string
		id         uuid.UUID
		beforeTest func(skillRepo mocks.MockISkillRepository)
		expected   *domain.Skill
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное получение навыка по id",
			id:   uuid.UUID{1},
			beforeTest: func(skillRepo mocks.MockISkillRepository) {
				skillRepo.EXPECT().
					GetById(
						context.Background(),
						uuid.UUID{1},
					).
					Return(&domain.Skill{
						ID:          uuid.UUID{1},
						Name:        "a",
						Description: "a",
					}, nil)
			},
			expected: &domain.Skill{
				ID:          uuid.UUID{1},
				Name:        "a",
				Description: "a",
			},
			wantErr: false,
		},
		{
			name: "ошибка получения данных в репозитории",
			id:   uuid.UUID{1},
			beforeTest: func(skillRepo mocks.MockISkillRepository) {
				skillRepo.EXPECT().
					GetById(
						context.Background(),
						uuid.UUID{1},
					).
					Return(nil, fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("получение навыка по id: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*skillRepo)
			}

			company, err := svc.GetById(ctx, tc.id)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.Equal(t, company, tc.expected)
			}
		})
	}
}

func TestSkillService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	skillRepo := mocks.NewMockISkillRepository(ctrl)
	svc := NewService(skillRepo)

	testCases := []struct {
		name       string
		skill      *domain.Skill
		beforeTest func(skillRepo mocks.MockISkillRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное обновление",
			skill: &domain.Skill{
				ID:   uuid.UUID{1},
				Name: "aaa",
			},
			beforeTest: func(skillRepo mocks.MockISkillRepository) {
				skillRepo.EXPECT().
					Update(
						context.Background(),
						&domain.Skill{
							ID:   uuid.UUID{1},
							Name: "aaa",
						},
					).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ошибка выполнения запроса в репозитории",
			skill: &domain.Skill{
				ID:   uuid.UUID{1},
				Name: "aaa",
			},
			beforeTest: func(skillRepo mocks.MockISkillRepository) {
				skillRepo.EXPECT().
					Update(
						context.Background(),
						&domain.Skill{
							ID:   uuid.UUID{1},
							Name: "aaa",
						},
					).Return(fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("обновление информации о навыке: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*skillRepo)
			}

			err := svc.Update(ctx, tc.skill)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}
