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

func TestFinReportService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	finRepo := mocks.NewMockIFinancialReportRepository(ctrl)
	svc := NewFinReportService(finRepo)

	testCases := []struct {
		name       string
		finReport  domain.FinancialReport
		beforeTest func(finRepo mocks.MockIFinancialReportRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное добавление",
			finReport: domain.FinancialReport{
				CompanyID: [16]byte{1},
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
							CompanyID: [16]byte{1},
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
			finReport: domain.FinancialReport{
				CompanyID: [16]byte{1},
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
							CompanyID: [16]byte{1},
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
			finReport: domain.FinancialReport{
				CompanyID: [16]byte{1},
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
							CompanyID: [16]byte{1},
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
			finReport: domain.FinancialReport{
				CompanyID: [16]byte{1},
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
							CompanyID: [16]byte{1},
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
			finReport: domain.FinancialReport{
				CompanyID: [16]byte{1},
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
							CompanyID: [16]byte{1},
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
			finReport: domain.FinancialReport{
				CompanyID: [16]byte{1},
				Revenue:   1,
				Costs:     1,
				Year:      2024,
				Quarter:   1,
			},
			beforeTest: func(finRepo mocks.MockIFinancialReportRepository) {
				finRepo.EXPECT().
					Create(
						context.Background(),
						&domain.FinancialReport{
							CompanyID: [16]byte{1},
							Revenue:   1,
							Costs:     1,
							Year:      2024,
							Quarter:   1,
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
			finReport: domain.FinancialReport{
				CompanyID: [16]byte{1},
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
							CompanyID: [16]byte{1},
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
			if tc.beforeTest != nil {
				tc.beforeTest(*finRepo)
			}

			err := svc.Create(context.Background(), &tc.finReport)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestFinReportService_CreateByPeriod(t *testing.T) {
	type fields struct {
		finRepo domain.IFinancialReportRepository
	}
	type args struct {
		ctx               context.Context
		finReportByPeriod *domain.FinancialReportByPeriod
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := FinReportService{
				finRepo: tt.fields.finRepo,
			}
			if err := s.CreateByPeriod(tt.args.ctx, tt.args.finReportByPeriod); (err != nil) != tt.wantErr {
				t.Errorf("CreateByPeriod() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
} // TODO

func TestFinReportService_DeleteById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	finRepo := mocks.NewMockIFinancialReportRepository(ctrl)
	svc := NewFinReportService(finRepo)

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
			if tc.beforeTest != nil {
				tc.beforeTest(*finRepo)
			}

			err := svc.DeleteById(context.Background(), tc.id)

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
	svc := NewFinReportService(finRepo)

	testCases := []struct {
		name       string
		id         uuid.UUID
		period     domain.Period
		beforeTest func(finRepo mocks.MockIFinancialReportRepository)
		expected   *domain.FinancialReportByPeriod
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное получение отчета компании по id",
			id:   [16]byte{1},
			period: domain.Period{
				StartYear:    2021,
				EndYear:      2023,
				StartQuarter: 2,
				EndQuarter:   4,
			},
			beforeTest: func(finRepo mocks.MockIFinancialReportRepository) {
				finRepo.EXPECT().
					GetByCompany(
						context.Background(),
						[16]byte{1},
						domain.Period{
							StartYear:    2021,
							EndYear:      2023,
							StartQuarter: 2,
							EndQuarter:   4,
						},
					).
					Return(&domain.FinancialReportByPeriod{
						Reports: []domain.FinancialReport{
							{
								ID:      [16]byte{1},
								Year:    2021,
								Quarter: 2,
								Revenue: 1432523,
								Costs:   75423,
							},
							{
								ID:      [16]byte{2},
								Year:    2021,
								Quarter: 3,
								Revenue: 7435235,
								Costs:   125654,
							},
							{
								ID:      [16]byte{3},
								Year:    2021,
								Quarter: 4,
								Revenue: 65742,
								Costs:   7845634,
							},
							{
								ID:      [16]byte{4},
								Year:    2022,
								Quarter: 1,
								Revenue: 43635325,
								Costs:   12362332,
							},
							{
								ID:      [16]byte{5},
								Year:    2022,
								Quarter: 2,
								Revenue: 50934123,
								Costs:   13543623,
							},
							{
								ID:      [16]byte{6},
								Year:    2022,
								Quarter: 3,
								Revenue: 78902453,
								Costs:   15326443,
							},
							{
								ID:      [16]byte{7},
								Year:    2022,
								Quarter: 4,
								Revenue: 64352357,
								Costs:   23534252,
							}, // 173 057 608 => 34 611 521.6; 237 824 258 => 14.5534025
							{
								ID:      [16]byte{8},
								Year:    2023,
								Quarter: 1,
								Revenue: 32532513,
								Costs:   5436438,
							},
							{
								ID:      [16]byte{9},
								Year:    2023,
								Quarter: 2,
								Revenue: 6743634,
								Costs:   9876967,
							},
							{
								ID:      [16]byte{10},
								Year:    2023,
								Quarter: 3,
								Revenue: 46754124,
								Costs:   24367653,
							},
							{
								ID:      [16]byte{11},
								Year:    2023,
								Quarter: 4,
								Revenue: 14385253,
								Costs:   7546424,
							},
						},
						Period: domain.Period{
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
						ID:      [16]byte{1},
						Year:    2021,
						Quarter: 2,
						Revenue: 1432523,
						Costs:   75423,
					},
					{
						ID:      [16]byte{2},
						Year:    2021,
						Quarter: 3,
						Revenue: 7435235,
						Costs:   125654,
					},
					{
						ID:      [16]byte{3},
						Year:    2021,
						Quarter: 4,
						Revenue: 65742,
						Costs:   7845634,
					},
					{
						ID:      [16]byte{4},
						Year:    2022,
						Quarter: 1,
						Revenue: 43635325,
						Costs:   12362332,
					},
					{
						ID:      [16]byte{5},
						Year:    2022,
						Quarter: 2,
						Revenue: 50934123,
						Costs:   13543623,
					},
					{
						ID:      [16]byte{6},
						Year:    2022,
						Quarter: 3,
						Revenue: 78902453,
						Costs:   15326443,
					},
					{
						ID:      [16]byte{7},
						Year:    2022,
						Quarter: 4,
						Revenue: 64352357,
						Costs:   23534252,
					}, // 173 057 608 => 34 611 521.6; 237 824 258 => 14.5534025
					{
						ID:      [16]byte{8},
						Year:    2023,
						Quarter: 1,
						Revenue: 32532513,
						Costs:   5436438,
					},
					{
						ID:      [16]byte{9},
						Year:    2023,
						Quarter: 2,
						Revenue: 6743634,
						Costs:   9876967,
					},
					{
						ID:      [16]byte{10},
						Year:    2023,
						Quarter: 3,
						Revenue: 46754124,
						Costs:   24367653,
					},
					{
						ID:      [16]byte{11},
						Year:    2023,
						Quarter: 4,
						Revenue: 14385253,
						Costs:   7546424,
					},
				},
				Period: domain.Period{
					StartYear:    2021,
					EndYear:      2023,
					StartQuarter: 2,
					EndQuarter:   4,
				},
				Taxes:   (43635325+50934123+78902453+64352357-12362332-13543623-15326443-23534252)*0.2 + (32532513+6743634+46754124+14385253-5436438-9876967-24367653-7546424)*0.13,
				TaxLoad: (43635325+50934123+78902453+64352357-12362332-13543623-15326443-23534252)*0.2/(43635325+50934123+78902453+64352357)*100 + (32532513+6743634+46754124+14385253-5436438-9876967-24367653-7546424)*0.13/(32532513+6743634+46754124+14385253)*100,
			},
			wantErr: false,
		},
		{
			name: "год начала периода больше года конца периода",
			id:   [16]byte{1},
			period: domain.Period{
				StartYear:    2,
				EndYear:      1,
				StartQuarter: 1,
				EndQuarter:   1,
			},
			beforeTest: func(finRepo mocks.MockIFinancialReportRepository) {
				finRepo.EXPECT().
					GetByCompany(
						context.Background(),
						[16]byte{1},
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
			id:   [16]byte{1},
			period: domain.Period{
				StartYear:    1,
				EndYear:      1,
				StartQuarter: 3,
				EndQuarter:   1,
			},
			beforeTest: func(finRepo mocks.MockIFinancialReportRepository) {
				finRepo.EXPECT().
					GetByCompany(
						context.Background(),
						[16]byte{1},
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
			id:   [16]byte{1},
			period: domain.Period{
				StartYear:    1,
				EndYear:      1,
				StartQuarter: 1,
				EndQuarter:   2,
			},
			beforeTest: func(finRepo mocks.MockIFinancialReportRepository) {
				finRepo.EXPECT().
					GetByCompany(
						context.Background(),
						[16]byte{1},
						domain.Period{
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
			if tc.beforeTest != nil {
				tc.beforeTest(*finRepo)
			}

			report, err := svc.GetByCompany(context.Background(), tc.id, tc.period)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.Equal(t, report.Reports, tc.expected.Reports)
				require.Equal(t, report.Period, tc.expected.Period)
				require.InEpsilon(t, report.Taxes, tc.expected.Taxes, 1e-7)
				require.InEpsilon(t, report.TaxLoad, tc.expected.TaxLoad, 1e-7)
			}
		})
	}
}

func TestFinReportService_GetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	finRepo := mocks.NewMockIFinancialReportRepository(ctrl)
	svc := NewFinReportService(finRepo)

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
			id:   [16]byte{1},
			beforeTest: func(finRepo mocks.MockIFinancialReportRepository) {
				finRepo.EXPECT().
					GetById(
						context.Background(),
						[16]byte{1},
					).
					Return(&domain.FinancialReport{
						ID:        [16]byte{1},
						CompanyID: [16]byte{1},
						Revenue:   1,
						Costs:     1,
						Year:      1,
						Quarter:   1,
					}, nil)
			},
			expected: &domain.FinancialReport{
				ID:        [16]byte{1},
				CompanyID: [16]byte{1},
				Revenue:   1,
				Costs:     1,
				Year:      1,
				Quarter:   1,
			},
			wantErr: false,
		},
		{
			name: "ошибка получения данных в репозитории",
			id:   [16]byte{1},
			beforeTest: func(finRepo mocks.MockIFinancialReportRepository) {
				finRepo.EXPECT().
					GetById(
						context.Background(),
						[16]byte{1},
					).
					Return(nil, fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("получение финансового отчета по id: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.beforeTest != nil {
				tc.beforeTest(*finRepo)
			}

			report, err := svc.GetById(context.Background(), tc.id)

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

	finRepo := mocks.NewMockIFinancialReportRepository(ctrl)
	svc := NewFinReportService(finRepo)

	testCases := []struct {
		name       string
		report     domain.FinancialReport
		beforeTest func(finRepo mocks.MockIFinancialReportRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное обновление",
			report: domain.FinancialReport{
				ID:      [16]byte{1},
				Revenue: 2,
			},
			beforeTest: func(finRepo mocks.MockIFinancialReportRepository) {
				finRepo.EXPECT().
					Update(
						context.Background(),
						&domain.FinancialReport{
							ID:      [16]byte{1},
							Revenue: 2,
						},
					).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ошибка выполнения запроса в репозитории",
			report: domain.FinancialReport{
				ID:      [16]byte{1},
				Revenue: 2,
			},
			beforeTest: func(finRepo mocks.MockIFinancialReportRepository) {
				finRepo.EXPECT().
					Update(
						context.Background(),
						&domain.FinancialReport{
							ID:      [16]byte{1},
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
			if tc.beforeTest != nil {
				tc.beforeTest(*finRepo)
			}

			err := svc.Update(context.Background(), &tc.report)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}
