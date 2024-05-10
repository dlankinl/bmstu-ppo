package activity_field

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

func TestService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockIActivityFieldRepository(ctrl)
	compRepo := mocks.NewMockICompanyRepository(ctrl)
	svc := NewService(repo, compRepo)

	testCases := []struct {
		name       string
		data       *domain.ActivityField
		beforeTest func(repo mocks.MockIActivityFieldRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное добавление",
			data: &domain.ActivityField{
				Name:        "aaa",
				Description: "aaa",
				Cost:        0.3,
			},
			beforeTest: func(repo mocks.MockIActivityFieldRepository) {
				repo.EXPECT().
					Create(
						context.Background(),
						&domain.ActivityField{
							Name:        "aaa",
							Description: "aaa",
							Cost:        0.3,
						},
					).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "пустое название сферы деятельности",
			data: &domain.ActivityField{
				Name:        "",
				Description: "aaa",
			},
			beforeTest: func(repo mocks.MockIActivityFieldRepository) {
				repo.EXPECT().
					Create(
						context.Background(),
						&domain.ActivityField{
							Name:        "",
							Description: "aaa",
						},
					).Return(nil).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("должно быть указано название сферы деятельности"),
		},
		{
			name: "пустое описание сферы деятельности",
			data: &domain.ActivityField{
				Name:        "aaa",
				Description: "",
				Cost:        0.3,
			},
			beforeTest: func(repo mocks.MockIActivityFieldRepository) {
				repo.EXPECT().
					Create(
						context.Background(),
						&domain.ActivityField{
							Name:        "aaa",
							Description: "",
							Cost:        0.3,
						},
					).Return(nil).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("должно быть указано описание сферы деятельности"),
		},
		{
			name: "нулевой вес сферы деятельности",
			data: &domain.ActivityField{
				Name:        "aaa",
				Description: "aaa",
			},
			beforeTest: func(repo mocks.MockIActivityFieldRepository) {
				repo.EXPECT().
					Create(
						context.Background(),
						&domain.ActivityField{
							Name:        "aaa",
							Description: "aaa",
						},
					).Return(nil).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("вес сферы деятельности не может быть равен 0"),
		},
		{
			name: "ошибка выполнения запроса в репозитории",
			data: &domain.ActivityField{
				Name:        "aaa",
				Description: "aaa",
				Cost:        3,
			},
			beforeTest: func(repo mocks.MockIActivityFieldRepository) {
				repo.EXPECT().
					Create(
						context.Background(),
						&domain.ActivityField{
							Name:        "aaa",
							Description: "aaa",
							Cost:        3,
						},
					).Return(fmt.Errorf("sql error")).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("создание сферы деятельности: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*repo)
			}

			err := svc.Create(ctx, tc.data)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestService_DeleteById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockIActivityFieldRepository(ctrl)
	compRepo := mocks.NewMockICompanyRepository(ctrl)
	svc := NewService(repo, compRepo)

	curUuid := uuid.New()

	testCases := []struct {
		name       string
		id         uuid.UUID
		beforeTest func(repo mocks.MockIActivityFieldRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное удаление",
			id:   curUuid,
			beforeTest: func(repo mocks.MockIActivityFieldRepository) {
				repo.EXPECT().
					DeleteById(context.Background(), curUuid).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ошибка выполнения запроса в репозитории",
			id:   curUuid,
			beforeTest: func(repo mocks.MockIActivityFieldRepository) {
				repo.EXPECT().
					DeleteById(context.Background(), curUuid).
					Return(fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("удаление сферы деятельности по id: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*repo)
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

func TestService_GetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockIActivityFieldRepository(ctrl)
	compRepo := mocks.NewMockICompanyRepository(ctrl)
	svc := NewService(repo, compRepo)

	testCases := []struct {
		name       string
		id         uuid.UUID
		beforeTest func(repo mocks.MockIActivityFieldRepository)
		expected   *domain.ActivityField
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное получение сферы деятельности по id",
			id:   uuid.UUID{1},
			beforeTest: func(repo mocks.MockIActivityFieldRepository) {
				repo.EXPECT().
					GetById(
						context.Background(),
						uuid.UUID{1},
					).
					Return(&domain.ActivityField{
						ID:          uuid.UUID{1},
						Name:        "a",
						Description: "a",
					}, nil)
			},
			expected: &domain.ActivityField{
				ID:          uuid.UUID{1},
				Name:        "a",
				Description: "a",
			},
			wantErr: false,
		},
		{
			name: "ошибка получения данных в репозитории",
			id:   uuid.UUID{1},
			beforeTest: func(repo mocks.MockIActivityFieldRepository) {
				repo.EXPECT().
					GetById(
						context.Background(),
						uuid.UUID{1},
					).
					Return(nil, fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("получение сферы деятельности по id: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*repo)
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

func TestService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockIActivityFieldRepository(ctrl)
	compRepo := mocks.NewMockICompanyRepository(ctrl)
	svc := NewService(repo, compRepo)

	testCases := []struct {
		name       string
		data       *domain.ActivityField
		beforeTest func(repo mocks.MockIActivityFieldRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное обновление",
			data: &domain.ActivityField{
				ID:   uuid.UUID{1},
				Name: "aaa",
			},
			beforeTest: func(repo mocks.MockIActivityFieldRepository) {
				repo.EXPECT().
					Update(
						context.Background(),
						&domain.ActivityField{
							ID:   uuid.UUID{1},
							Name: "aaa",
						},
					).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ошибка выполнения запроса в репозитории",
			data: &domain.ActivityField{
				ID:   uuid.UUID{1},
				Name: "aaa",
			},
			beforeTest: func(repo mocks.MockIActivityFieldRepository) {
				repo.EXPECT().
					Update(
						context.Background(),
						&domain.ActivityField{
							ID:   uuid.UUID{1},
							Name: "aaa",
						},
					).Return(fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("обновление информации о cфере деятельности: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*repo)
			}

			err := svc.Update(ctx, tc.data)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}
