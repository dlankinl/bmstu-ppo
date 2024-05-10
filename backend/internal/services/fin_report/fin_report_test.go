package fin_report

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

func TestFinReportService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	finRepo := mocks.NewMockIFinancialReportRepository(ctrl)
	svc := NewService(finRepo)

	testCases := []struct {
		name       string
		data       *domain.FinancialReport
		beforeTest func(finRepo mocks.MockIFinancialReportRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное добавление",
			data: &domain.FinancialReport{
				CompanyID: uuid.UUID{1},
				Revenue:   1,
				Costs:     1,
				Year:      1,
				Quarter:   1,
			},
			beforeTest: func(finRepo mocks.MockIFinancialReportRepository) {
				finRepo.EXPECT().
					Create(
						context.Background(),
						&domain.FinancialReport{
							CompanyID: uuid.UUID{1},
							Revenue:   1,
							Costs:     1,
							Year:      1,
							Quarter:   1,
						},
					).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "отрицательная выручка",
			data: &domain.FinancialReport{
				CompanyID: uuid.UUID{1},
				Revenue:   -1,
				Costs:     1,
				Year:      1,
				Quarter:   1,
			},
			beforeTest: func(finRepo mocks.MockIFinancialReportRepository) {
				finRepo.EXPECT().
					Create(
						context.Background(),
						&domain.FinancialReport{
							CompanyID: uuid.UUID{1},
							Revenue:   -1,
							Costs:     1,
							Year:      1,
							Quarter:   1,
						},
					).
					Return(nil).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("выручка не может быть отрицательной"),
		},
		{
			name: "отрицательные расходы",
			data: &domain.FinancialReport{
				CompanyID: uuid.UUID{1},
				Revenue:   1,
				Costs:     -1,
				Year:      1,
				Quarter:   1,
			},
			beforeTest: func(finRepo mocks.MockIFinancialReportRepository) {
				finRepo.EXPECT().
					Create(
						context.Background(),
						&domain.FinancialReport{
							CompanyID: uuid.UUID{1},
							Revenue:   1,
							Costs:     -1,
							Year:      1,
							Quarter:   1,
						},
					).
					Return(nil).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("расходы не могут быть отрицательными"),
		},
		{
			name: "некорректное значение квартала",
			data: &domain.FinancialReport{
				CompanyID: uuid.UUID{1},
				Revenue:   1,
				Costs:     1,
				Year:      1,
				Quarter:   5,
			},
			beforeTest: func(finRepo mocks.MockIFinancialReportRepository) {
				finRepo.EXPECT().
					Create(
						context.Background(),
						&domain.FinancialReport{
							CompanyID: uuid.UUID{1},
							Revenue:   1,
							Costs:     1,
							Year:      1,
							Quarter:   5,
						},
					).
					Return(nil).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("значение квартала должно находиться в отрезке от 1 до 4"),
		},
		{
			name: "указан год, больший текущего",
			data: &domain.FinancialReport{
				CompanyID: uuid.UUID{1},
				Revenue:   1,
				Costs:     1,
				Year:      2025,
				Quarter:   1,
			},
			beforeTest: func(finRepo mocks.MockIFinancialReportRepository) {
				finRepo.EXPECT().
					Create(
						context.Background(),
						&domain.FinancialReport{
							CompanyID: uuid.UUID{1},
							Revenue:   1,
							Costs:     1,
							Year:      2025,
							Quarter:   1,
						},
					).
					Return(nil).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("значение года не может быть больше текущего года"),
		},
		{
			name: "указан квартал, который еще не закончен",
			data: &domain.FinancialReport{
				CompanyID: uuid.UUID{1},
				Revenue:   1,
				Costs:     1,
				Year:      2024,
				Quarter:   2,
			},
			beforeTest: func(finRepo mocks.MockIFinancialReportRepository) {
				finRepo.EXPECT().
					Create(
						context.Background(),
						&domain.FinancialReport{
							CompanyID: uuid.UUID{1},
							Revenue:   1,
							Costs:     1,
							Year:      2024,
							Quarter:   2,
						},
					).
					Return(nil).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("нельзя добавить отчет за квартал, который еще не закончился"),
		},
		{
			name: "ошибка выполнения запроса в репозитории",
			data: &domain.FinancialReport{
				CompanyID: uuid.UUID{1},
				Revenue:   1,
				Costs:     1,
				Year:      2023,
				Quarter:   1,
			},
			beforeTest: func(finRepo mocks.MockIFinancialReportRepository) {
				finRepo.EXPECT().
					Create(
						context.Background(),
						&domain.FinancialReport{
							CompanyID: uuid.UUID{1},
							Revenue:   1,
							Costs:     1,
							Year:      2023,
							Quarter:   1,
						},
					).
					Return(fmt.Errorf("sql error")).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("добавление финансового отчета: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*finRepo)
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

func TestFinReportService_DeleteById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	finRepo := mocks.NewMockIFinancialReportRepository(ctrl)
	svc := NewService(finRepo)

	curUuid := uuid.New()

	testCases := []struct {
		name       string
		id         uuid.UUID
		beforeTest func(finRepo mocks.MockIFinancialReportRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное удаление",
			id:   curUuid,
			beforeTest: func(finRepo mocks.MockIFinancialReportRepository) {
				finRepo.EXPECT().
					DeleteById(context.Background(), curUuid).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ошибка выполнения запроса в репозитории",
			id:   curUuid,
			beforeTest: func(finRepo mocks.MockIFinancialReportRepository) {
				finRepo.EXPECT().
					DeleteById(context.Background(), curUuid).
					Return(fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("удаление отчета по id: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*finRepo)
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

func TestFinReportService_GetByCompany(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	finRepo := mocks.NewMockIFinancialReportRepository(ctrl)
	svc := NewService(finRepo)

	testCases := []struct {
		name       string
		id         uuid.UUID
		period     *domain.Period
		beforeTest func(finRepo mocks.MockIFinancialReportRepository)
		expected   *domain.FinancialReportByPeriod
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное получение отчета компании по id",
			id:   uuid.UUID{1},
			period: &domain.Period{
				StartYear:    2021,
				EndYear:      2023,
				StartQuarter: 2,
				EndQuarter:   4,
			},
			beforeTest: func(finRepo mocks.MockIFinancialReportRepository) {
				finRepo.EXPECT().
					GetByCompany(
						context.Background(),
						uuid.UUID{1},
						&domain.Period{
							StartYear:    2021,
							EndYear:      2023,
							StartQuarter: 2,
							EndQuarter:   4,
						},
					).
					Return(&domain.FinancialReportByPeriod{
						Reports: []domain.FinancialReport{
							{
								ID:      uuid.UUID{1},
								Year:    2021,
								Quarter: 2,
								Revenue: 1432523,
								Costs:   75423,
							},
							{
								ID:      uuid.UUID{2},
								Year:    2021,
								Quarter: 3,
								Revenue: 7435235,
								Costs:   125654,
							},
							{
								ID:      uuid.UUID{3},
								Year:    2021,
								Quarter: 4,
								Revenue: 65742,
								Costs:   7845634,
							},
							{
								ID:      uuid.UUID{4},
								Year:    2022,
								Quarter: 1,
								Revenue: 43635325,
								Costs:   12362332,
							},
							{
								ID:      uuid.UUID{5},
								Year:    2022,
								Quarter: 2,
								Revenue: 50934123,
								Costs:   13543623,
							},
							{
								ID:      uuid.UUID{6},
								Year:    2022,
								Quarter: 3,
								Revenue: 78902453,
								Costs:   15326443,
							},
							{
								ID:      uuid.UUID{7},
								Year:    2022,
								Quarter: 4,
								Revenue: 64352357,
								Costs:   23534252,
							}, // 173 057 608 => 34 611 521.6; 237 824 258 => 14.5534025
							{
								ID:      uuid.UUID{8},
								Year:    2023,
								Quarter: 1,
								Revenue: 32532513,
								Costs:   5436438,
							},
							{
								ID:      uuid.UUID{9},
								Year:    2023,
								Quarter: 2,
								Revenue: 6743634,
								Costs:   9876967,
							},
							{
								ID:      uuid.UUID{10},
								Year:    2023,
								Quarter: 3,
								Revenue: 46754124,
								Costs:   24367653,
							},
							{
								ID:      uuid.UUID{11},
								Year:    2023,
								Quarter: 4,
								Revenue: 14385253,
								Costs:   7546424,
							},
						},
						Period: &domain.Period{
							StartYear:    2021,
							EndYear:      2023,
							StartQuarter: 2,
							EndQuarter:   4,
						},
					}, nil)
			},
			expected: &domain.FinancialReportByPeriod{
				Reports: []domain.FinancialReport{
					{
						ID:      uuid.UUID{1},
						Year:    2021,
						Quarter: 2,
						Revenue: 1432523,
						Costs:   75423,
					},
					{
						ID:      uuid.UUID{2},
						Year:    2021,
						Quarter: 3,
						Revenue: 7435235,
						Costs:   125654,
					},
					{
						ID:      uuid.UUID{3},
						Year:    2021,
						Quarter: 4,
						Revenue: 65742,
						Costs:   7845634,
					},
					{
						ID:      uuid.UUID{4},
						Year:    2022,
						Quarter: 1,
						Revenue: 43635325,
						Costs:   12362332,
					},
					{
						ID:      uuid.UUID{5},
						Year:    2022,
						Quarter: 2,
						Revenue: 50934123,
						Costs:   13543623,
					},
					{
						ID:      uuid.UUID{6},
						Year:    2022,
						Quarter: 3,
						Revenue: 78902453,
						Costs:   15326443,
					},
					{
						ID:      uuid.UUID{7},
						Year:    2022,
						Quarter: 4,
						Revenue: 64352357,
						Costs:   23534252,
					}, // 173 057 608 => 34 611 521.6; 237 824 258 => 14.5534025
					{
						ID:      uuid.UUID{8},
						Year:    2023,
						Quarter: 1,
						Revenue: 32532513,
						Costs:   5436438,
					},
					{
						ID:      uuid.UUID{9},
						Year:    2023,
						Quarter: 2,
						Revenue: 6743634,
						Costs:   9876967,
					},
					{
						ID:      uuid.UUID{10},
						Year:    2023,
						Quarter: 3,
						Revenue: 46754124,
						Costs:   24367653,
					},
					{
						ID:      uuid.UUID{11},
						Year:    2023,
						Quarter: 4,
						Revenue: 14385253,
						Costs:   7546424,
					},
				},
				Period: &domain.Period{
					StartYear:    2021,
					EndYear:      2023,
					StartQuarter: 2,
					EndQuarter:   4,
				},
			},
			wantErr: false,
		},
		{
			name: "год начала периода больше года конца периода",
			id:   uuid.UUID{1},
			period: &domain.Period{
				StartYear:    2,
				EndYear:      1,
				StartQuarter: 1,
				EndQuarter:   1,
			},
			beforeTest: func(finRepo mocks.MockIFinancialReportRepository) {
				finRepo.EXPECT().
					GetByCompany(
						context.Background(),
						uuid.UUID{1},
						domain.Period{
							StartYear:    2,
							EndYear:      1,
							StartQuarter: 1,
							EndQuarter:   1,
						},
					).
					Return(nil, nil).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("дата конца периода должна быть позже даты начала"),
		},
		{
			name: "равный год, но квартал начала больше квартала конца",
			id:   uuid.UUID{1},
			period: &domain.Period{
				StartYear:    1,
				EndYear:      1,
				StartQuarter: 3,
				EndQuarter:   1,
			},
			beforeTest: func(finRepo mocks.MockIFinancialReportRepository) {
				finRepo.EXPECT().
					GetByCompany(
						context.Background(),
						uuid.UUID{1},
						domain.Period{
							StartYear:    1,
							EndYear:      1,
							StartQuarter: 3,
							EndQuarter:   1,
						},
					).
					Return(nil, nil).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("дата конца периода должна быть позже даты начала"),
		},
		{
			name: "ошибка получения данных в репозитории",
			id:   uuid.UUID{1},
			period: &domain.Period{
				StartYear:    1,
				EndYear:      1,
				StartQuarter: 1,
				EndQuarter:   2,
			},
			beforeTest: func(finRepo mocks.MockIFinancialReportRepository) {
				finRepo.EXPECT().
					GetByCompany(
						context.Background(),
						uuid.UUID{1},
						&domain.Period{
							StartYear:    1,
							EndYear:      1,
							StartQuarter: 1,
							EndQuarter:   2,
						},
					).
					Return(nil, fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("получение финансового отчета по id компании: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*finRepo)
			}

			report, err := svc.GetByCompany(ctx, tc.id, tc.period)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.Equal(t, report.Reports, tc.expected.Reports)
				require.Equal(t, report.Period, tc.expected.Period)
			}
		})
	}
}

func TestFinReportService_GetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockIFinancialReportRepository(ctrl)
	svc := NewService(repo)

	testCases := []struct {
		name       string
		id         uuid.UUID
		beforeTest func(finRepo mocks.MockIFinancialReportRepository)
		expected   *domain.FinancialReport
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное получение отчета по id",
			id:   uuid.UUID{1},
			beforeTest: func(finRepo mocks.MockIFinancialReportRepository) {
				finRepo.EXPECT().
					GetById(
						context.Background(),
						uuid.UUID{1},
					).
					Return(&domain.FinancialReport{
						ID:        uuid.UUID{1},
						CompanyID: uuid.UUID{1},
						Revenue:   1,
						Costs:     1,
						Year:      1,
						Quarter:   1,
					}, nil)
			},
			expected: &domain.FinancialReport{
				ID:        uuid.UUID{1},
				CompanyID: uuid.UUID{1},
				Revenue:   1,
				Costs:     1,
				Year:      1,
				Quarter:   1,
			},
			wantErr: false,
		},
		{
			name: "ошибка получения данных в репозитории",
			id:   uuid.UUID{1},
			beforeTest: func(finRepo mocks.MockIFinancialReportRepository) {
				finRepo.EXPECT().
					GetById(
						context.Background(),
						uuid.UUID{1},
					).
					Return(nil, fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("получение финансового отчета по id: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*repo)
			}

			report, err := svc.GetById(ctx, tc.id)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.Equal(t, report, tc.expected)
			}
		})
	}
}

func TestFinReportService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockIFinancialReportRepository(ctrl)
	svc := NewService(repo)

	testCases := []struct {
		name       string
		report     *domain.FinancialReport
		beforeTest func(finRepo mocks.MockIFinancialReportRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное обновление",
			report: &domain.FinancialReport{
				ID:      uuid.UUID{1},
				Revenue: 2,
			},
			beforeTest: func(finRepo mocks.MockIFinancialReportRepository) {
				finRepo.EXPECT().
					Update(
						context.Background(),
						&domain.FinancialReport{
							ID:      uuid.UUID{1},
							Revenue: 2,
						},
					).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ошибка выполнения запроса в репозитории",
			report: &domain.FinancialReport{
				ID:      uuid.UUID{1},
				Revenue: 2,
			},
			beforeTest: func(finRepo mocks.MockIFinancialReportRepository) {
				finRepo.EXPECT().
					Update(
						context.Background(),
						&domain.FinancialReport{
							ID:      uuid.UUID{1},
							Revenue: 2,
						},
					).Return(fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("обновление отчета: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*repo)
			}

			err := svc.Update(ctx, tc.report)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}
