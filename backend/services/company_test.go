package services

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
	finRepo := mocks.NewMockIFinancialReportRepository(ctrl)
	svc := NewCompanyService(compRepo, finRepo)

	testCases := []struct {
		name       string
		company    domain.Company
		beforeTest func(compRepo mocks.MockICompanyRepository, finRepo mocks.MockIFinancialReportRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное добавление",
			company: domain.Company{
				Name:            "aaa",
				FieldOfActivity: "bbb",
				City:            "ccc",
			},
			beforeTest: func(compRepo mocks.MockICompanyRepository, finRepo mocks.MockIFinancialReportRepository) {
				compRepo.EXPECT().
					Create(
						context.Background(),
						&domain.Company{
							Name:            "aaa",
							FieldOfActivity: "bbb",
							City:            "ccc",
						},
					).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "пустое название компании",
			company: domain.Company{
				Name:            "",
				FieldOfActivity: "bbb",
				City:            "ccc",
			},
			beforeTest: func(compRepo mocks.MockICompanyRepository, finRepo mocks.MockIFinancialReportRepository) {
				compRepo.EXPECT().
					Create(
						context.Background(),
						&domain.Company{
							Name:            "",
							FieldOfActivity: "bbb",
							City:            "ccc",
						},
					).Return(nil).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("должно быть указано название компании"),
		},
		{
			name: "пустое название города",
			company: domain.Company{
				Name:            "aaa",
				FieldOfActivity: "bbb",
				City:            "",
			},
			beforeTest: func(compRepo mocks.MockICompanyRepository, finRepo mocks.MockIFinancialReportRepository) {
				compRepo.EXPECT().
					Create(
						context.Background(),
						&domain.Company{
							Name:            "aaa",
							FieldOfActivity: "bbb",
							City:            "",
						},
					).Return(nil).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("должно быть указано название города"),
		},
		{
			name: "пустое название сферы деятельности",
			company: domain.Company{
				Name:            "aaa",
				FieldOfActivity: "",
				City:            "ccc",
			},
			beforeTest: func(compRepo mocks.MockICompanyRepository, finRepo mocks.MockIFinancialReportRepository) {
				compRepo.EXPECT().
					Create(
						context.Background(),
						&domain.Company{
							Name:            "aaa",
							FieldOfActivity: "",
							City:            "ccc",
						},
					).Return(nil).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("должно быть указано название сферы деятельности"),
		},
		{
			name: "ошибка выполнения запроса в репозитории",
			company: domain.Company{
				Name:            "aaa",
				FieldOfActivity: "bbb",
				City:            "ccc",
			},
			beforeTest: func(compRepo mocks.MockICompanyRepository, finRepo mocks.MockIFinancialReportRepository) {
				compRepo.EXPECT().
					Create(
						context.Background(),
						&domain.Company{
							Name:            "aaa",
							FieldOfActivity: "bbb",
							City:            "ccc",
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
			if tc.beforeTest != nil {
				tc.beforeTest(*compRepo, *finRepo)
			}

			err := svc.Create(context.Background(), &tc.company)

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
	finRepo := mocks.NewMockIFinancialReportRepository(ctrl)
	svc := NewCompanyService(compRepo, finRepo)

	curUuid := uuid.New()

	testCases := []struct {
		name       string
		id         uuid.UUID
		beforeTest func(compRepo mocks.MockICompanyRepository, finRepo mocks.MockIFinancialReportRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное удаление",
			id:   curUuid,
			beforeTest: func(compRepo mocks.MockICompanyRepository, finRepo mocks.MockIFinancialReportRepository) {
				compRepo.EXPECT().
					DeleteById(context.Background(), curUuid).
					Return(nil)
			},
			wantErr: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.beforeTest != nil {
				tc.beforeTest(*compRepo, *finRepo)
			}

			err := svc.DeleteById(context.Background(), tc.id)

			if tc.wantErr {
				require.Equal(t, tc.errStr, err)
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
	finRepo := mocks.NewMockIFinancialReportRepository(ctrl)
	svc := NewCompanyService(compRepo, finRepo)

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
					GetAll(context.Background()).
					Return([]*domain.Company{
						{
							ID:              [16]byte{1},
							Name:            "a",
							FieldOfActivity: "a",
							City:            "a",
						},
						{
							ID:              [16]byte{2},
							Name:            "b",
							FieldOfActivity: "b",
							City:            "b",
						},
						{
							ID:              [16]byte{3},
							Name:            "c",
							FieldOfActivity: "c",
							City:            "c",
						},
					}, nil)
			},
			expected: []*domain.Company{
				{
					ID:              [16]byte{1},
					Name:            "a",
					FieldOfActivity: "a",
					City:            "a",
				},
				{
					ID:              [16]byte{2},
					Name:            "b",
					FieldOfActivity: "b",
					City:            "b",
				},
				{
					ID:              [16]byte{3},
					Name:            "c",
					FieldOfActivity: "c",
					City:            "c",
				},
			},
			wantErr: false,
		},
		{
			name: "ошибка получения данных в репозитории",
			beforeTest: func(compRepo mocks.MockICompanyRepository) {
				compRepo.EXPECT().
					GetAll(context.Background()).
					Return(nil, fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("получение списка всех компаний: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.beforeTest != nil {
				tc.beforeTest(*compRepo)
			}

			companies, err := svc.GetAll(context.Background())

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
	finRepo := mocks.NewMockIFinancialReportRepository(ctrl)
	svc := NewCompanyService(compRepo, finRepo)

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
			id:   [16]byte{1},
			beforeTest: func(compRepo mocks.MockICompanyRepository) {
				compRepo.EXPECT().
					GetById(
						context.Background(),
						[16]byte{1},
					).
					Return(&domain.Company{
						ID:              [16]byte{1},
						Name:            "a",
						FieldOfActivity: "a",
						City:            "a",
					}, nil)
			},
			expected: &domain.Company{
				ID:              [16]byte{1},
				Name:            "a",
				FieldOfActivity: "a",
				City:            "a",
			},
			wantErr: false,
		},
		{
			name: "ошибка получения данных в репозитории",
			id:   [16]byte{1},
			beforeTest: func(compRepo mocks.MockICompanyRepository) {
				compRepo.EXPECT().
					GetById(
						context.Background(),
						[16]byte{1},
					).
					Return(nil, fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("получение компании по id: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.beforeTest != nil {
				tc.beforeTest(*compRepo)
			}

			company, err := svc.GetById(context.Background(), tc.id)

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
	finRepo := mocks.NewMockIFinancialReportRepository(ctrl)
	svc := NewCompanyService(compRepo, finRepo)

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
			id:   [16]byte{1},
			beforeTest: func(compRepo mocks.MockICompanyRepository) {
				compRepo.EXPECT().
					GetByOwnerId(
						context.Background(),
						[16]byte{1},
					).
					Return([]*domain.Company{
						{
							ID:              [16]byte{1},
							OwnerId:         [16]byte{1},
							Name:            "a",
							FieldOfActivity: "a",
							City:            "a",
						},
						{
							ID:              [16]byte{2},
							OwnerId:         [16]byte{1},
							Name:            "b",
							FieldOfActivity: "b",
							City:            "b",
						},
						{
							ID:              [16]byte{3},
							OwnerId:         [16]byte{1},
							Name:            "c",
							FieldOfActivity: "c",
							City:            "c",
						},
					}, nil)
			},
			expected: []*domain.Company{
				{
					ID:              [16]byte{1},
					OwnerId:         [16]byte{1},
					Name:            "a",
					FieldOfActivity: "a",
					City:            "a",
				},
				{
					ID:              [16]byte{2},
					OwnerId:         [16]byte{1},
					Name:            "b",
					FieldOfActivity: "b",
					City:            "b",
				},
				{
					ID:              [16]byte{3},
					OwnerId:         [16]byte{1},
					Name:            "c",
					FieldOfActivity: "c",
					City:            "c",
				},
			},
			wantErr: false,
		},
		{
			name: "ошибка получения данных в репозитории",
			id:   [16]byte{1},
			beforeTest: func(compRepo mocks.MockICompanyRepository) {
				compRepo.EXPECT().
					GetByOwnerId(
						context.Background(),
						[16]byte{1},
					).
					Return(nil, fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("получение списка компаний по id владельца: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.beforeTest != nil {
				tc.beforeTest(*compRepo)
			}

			companies, err := svc.GetByOwnerId(context.Background(), tc.id)

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
	finRepo := mocks.NewMockIFinancialReportRepository(ctrl)
	svc := NewCompanyService(compRepo, finRepo)

	testCases := []struct {
		name       string
		company    domain.Company
		beforeTest func(compRepo mocks.MockICompanyRepository, finRepo mocks.MockIFinancialReportRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное обновление",
			company: domain.Company{
				ID:   [16]byte{1},
				Name: "aaa",
			},
			beforeTest: func(compRepo mocks.MockICompanyRepository, finRepo mocks.MockIFinancialReportRepository) {
				compRepo.EXPECT().
					Update(
						context.Background(),
						&domain.Company{
							ID:   [16]byte{1},
							Name: "aaa",
						},
					).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ошибка выполнения запроса в репозитории",
			company: domain.Company{
				ID:   [16]byte{1},
				Name: "aaa",
			},
			beforeTest: func(compRepo mocks.MockICompanyRepository, finRepo mocks.MockIFinancialReportRepository) {
				compRepo.EXPECT().
					Update(
						context.Background(),
						&domain.Company{
							ID:   [16]byte{1},
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
			if tc.beforeTest != nil {
				tc.beforeTest(*compRepo, *finRepo)
			}

			err := svc.Update(context.Background(), &tc.company)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}