package company

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

func TestCompanyService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	compRepo := mocks.NewMockICompanyRepository(ctrl)
	svc := NewService(compRepo)

	testCases := []struct {
		name       string
		company    *domain.Company
		beforeTest func(compRepo mocks.MockICompanyRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное добавление",
			company: &domain.Company{
				Name: "aaa",
				City: "ccc",
			},
			beforeTest: func(compRepo mocks.MockICompanyRepository) {
				compRepo.EXPECT().
					Create(
						context.Background(),
						&domain.Company{
							Name: "aaa",
							City: "ccc",
						},
					).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "пустое название компании",
			company: &domain.Company{
				Name: "",
				City: "ccc",
			},
			beforeTest: func(compRepo mocks.MockICompanyRepository) {
				compRepo.EXPECT().
					Create(
						context.Background(),
						&domain.Company{
							Name: "",
							City: "ccc",
						},
					).Return(nil).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("должно быть указано название компании"),
		},
		{
			name: "пустое название города",
			company: &domain.Company{
				Name: "aaa",
				City: "",
			},
			beforeTest: func(compRepo mocks.MockICompanyRepository) {
				compRepo.EXPECT().
					Create(
						context.Background(),
						&domain.Company{
							Name: "aaa",
							City: "",
						},
					).Return(nil).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("должно быть указано название города"),
		},
		{
			name: "ошибка выполнения запроса в репозитории",
			company: &domain.Company{
				Name: "aaa",
				City: "ccc",
			},
			beforeTest: func(compRepo mocks.MockICompanyRepository) {
				compRepo.EXPECT().
					Create(
						context.Background(),
						&domain.Company{
							Name: "aaa",
							City: "ccc",
						},
					).Return(fmt.Errorf("sql error")).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("добавление компании: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*compRepo)
			}

			err := svc.Create(ctx, tc.company)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestCompanyService_DeleteById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	compRepo := mocks.NewMockICompanyRepository(ctrl)
	svc := NewService(compRepo)

	curUuid := uuid.New()

	testCases := []struct {
		name       string
		id         uuid.UUID
		beforeTest func(compRepo mocks.MockICompanyRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное удаление",
			id:   curUuid,
			beforeTest: func(compRepo mocks.MockICompanyRepository) {
				compRepo.EXPECT().
					DeleteById(context.Background(), curUuid).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ошибка выполнения запроса в репозитории",
			id:   curUuid,
			beforeTest: func(compRepo mocks.MockICompanyRepository) {
				compRepo.EXPECT().
					DeleteById(context.Background(), curUuid).
					Return(fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("удаление компании по id: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*compRepo)
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

func TestCompanyService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	compRepo := mocks.NewMockICompanyRepository(ctrl)
	svc := NewService(compRepo)

	testCases := []struct {
		name       string
		beforeTest func(compRepo mocks.MockICompanyRepository)
		expected   []*domain.Company
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное получение списка всех компаний",
			beforeTest: func(compRepo mocks.MockICompanyRepository) {
				compRepo.EXPECT().
					GetAll(context.Background(), 1).
					Return([]*domain.Company{
						{
							ID:   uuid.UUID{1},
							Name: "a",
							City: "a",
						},
						{
							ID:   uuid.UUID{2},
							Name: "b",
							City: "b",
						},
						{
							ID:   uuid.UUID{3},
							Name: "c",
							City: "c",
						},
					}, nil)
			},
			expected: []*domain.Company{
				{
					ID:   uuid.UUID{1},
					Name: "a",
					City: "a",
				},
				{
					ID:   uuid.UUID{2},
					Name: "b",
					City: "b",
				},
				{
					ID:   uuid.UUID{3},
					Name: "c",
					City: "c",
				},
			},
			wantErr: false,
		},
		{
			name: "ошибка получения данных в репозитории",
			beforeTest: func(compRepo mocks.MockICompanyRepository) {
				compRepo.EXPECT().
					GetAll(context.Background(), 1).
					Return(nil, fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("получение списка всех компаний: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*compRepo)
			}

			companies, err := svc.GetAll(ctx, 1)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.Equal(t, companies, tc.expected)
			}
		})
	}
}

func TestCompanyService_GetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	compRepo := mocks.NewMockICompanyRepository(ctrl)
	svc := NewService(compRepo)

	testCases := []struct {
		name       string
		id         uuid.UUID
		beforeTest func(compRepo mocks.MockICompanyRepository)
		expected   *domain.Company
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное получение компании по id",
			id:   uuid.UUID{1},
			beforeTest: func(compRepo mocks.MockICompanyRepository) {
				compRepo.EXPECT().
					GetById(
						context.Background(),
						uuid.UUID{1},
					).
					Return(&domain.Company{
						ID:   uuid.UUID{1},
						Name: "a",
						City: "a",
					}, nil)
			},
			expected: &domain.Company{
				ID:   uuid.UUID{1},
				Name: "a",
				City: "a",
			},
			wantErr: false,
		},
		{
			name: "ошибка получения данных в репозитории",
			id:   uuid.UUID{1},
			beforeTest: func(compRepo mocks.MockICompanyRepository) {
				compRepo.EXPECT().
					GetById(
						context.Background(),
						uuid.UUID{1},
					).
					Return(nil, fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("получение компании по id: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*compRepo)
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

func TestCompanyService_GetByOwnerId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	compRepo := mocks.NewMockICompanyRepository(ctrl)
	svc := NewService(compRepo)

	testCases := []struct {
		name       string
		id         uuid.UUID
		beforeTest func(compRepo mocks.MockICompanyRepository)
		expected   []*domain.Company
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное получение компании по id",
			id:   uuid.UUID{1},
			beforeTest: func(compRepo mocks.MockICompanyRepository) {
				compRepo.EXPECT().
					GetByOwnerId(
						context.Background(),
						uuid.UUID{1},
						1,
						true,
					).
					Return([]*domain.Company{
						{
							ID:      uuid.UUID{1},
							OwnerID: uuid.UUID{1},
							Name:    "a",
							City:    "a",
						},
						{
							ID:      uuid.UUID{2},
							OwnerID: uuid.UUID{1},
							Name:    "b",
							City:    "b",
						},
						{
							ID:      uuid.UUID{3},
							OwnerID: uuid.UUID{1},
							Name:    "c",
							City:    "c",
						},
					}, nil)
			},
			expected: []*domain.Company{
				{
					ID:      uuid.UUID{1},
					OwnerID: uuid.UUID{1},
					Name:    "a",
					City:    "a",
				},
				{
					ID:      uuid.UUID{2},
					OwnerID: uuid.UUID{1},
					Name:    "b",
					City:    "b",
				},
				{
					ID:      uuid.UUID{3},
					OwnerID: uuid.UUID{1},
					Name:    "c",
					City:    "c",
				},
			},
			wantErr: false,
		},
		{
			name: "ошибка получения данных в репозитории",
			id:   uuid.UUID{1},
			beforeTest: func(compRepo mocks.MockICompanyRepository) {
				compRepo.EXPECT().
					GetByOwnerId(
						context.Background(),
						uuid.UUID{1},
						1,
						true,
					).
					Return(nil, fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("получение списка компаний по id владельца: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*compRepo)
			}

			companies, err := svc.GetByOwnerId(ctx, tc.id, 1, true)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.Equal(t, companies, tc.expected)
			}
		})
	}
}

func TestCompanyService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	compRepo := mocks.NewMockICompanyRepository(ctrl)
	svc := NewService(compRepo)

	testCases := []struct {
		name       string
		company    *domain.Company
		beforeTest func(compRepo mocks.MockICompanyRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное обновление",
			company: &domain.Company{
				ID:   uuid.UUID{1},
				Name: "aaa",
			},
			beforeTest: func(compRepo mocks.MockICompanyRepository) {
				compRepo.EXPECT().
					Update(
						context.Background(),
						&domain.Company{
							ID:   uuid.UUID{1},
							Name: "aaa",
						},
					).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ошибка выполнения запроса в репозитории",
			company: &domain.Company{
				ID:   uuid.UUID{1},
				Name: "aaa",
			},
			beforeTest: func(compRepo mocks.MockICompanyRepository) {
				compRepo.EXPECT().
					Update(
						context.Background(),
						&domain.Company{
							ID:   uuid.UUID{1},
							Name: "aaa",
						},
					).Return(fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("обновление информации о компании: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*compRepo)
			}

			err := svc.Update(ctx, tc.company)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}
