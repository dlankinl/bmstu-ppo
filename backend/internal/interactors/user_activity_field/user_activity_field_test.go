package user_activity_field

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"ppo/domain"
	"ppo/internal/services/activity_field"
	"ppo/internal/services/company"
	"ppo/internal/services/fin_report"
	"ppo/internal/services/user"
	"ppo/mocks"
	"testing"
)

const eps = 1e-7

func TestInteractor_CalculateUserRating(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockIUserRepository(ctrl)
	finRepo := mocks.NewMockIFinancialReportRepository(ctrl)
	compRepo := mocks.NewMockICompanyRepository(ctrl)
	actFieldRepo := mocks.NewMockIActivityFieldRepository(ctrl)

	userSvc := user.NewService(userRepo, compRepo, actFieldRepo)
	actFieldSvc := activity_field.NewService(actFieldRepo, compRepo)
	compSvc := company.NewService(compRepo)
	finSvc := fin_report.NewService(finRepo)

	interactor := NewInteractor(userSvc, actFieldSvc, compSvc, finSvc)

	testCases := []struct {
		name       string
		userId     uuid.UUID
		beforeTest func(
			userRepo mocks.MockIUserRepository,
			finRepo mocks.MockIFinancialReportRepository,
			compRepo mocks.MockICompanyRepository,
			actFieldRepo mocks.MockIActivityFieldRepository,
		)
		wantErr  bool
		expected float32
		errStr   error
	}{
		{
			name:   "успешное вычисление рейтинга",
			userId: uuid.UUID{1},
			beforeTest: func(userRepo mocks.MockIUserRepository, finRepo mocks.MockIFinancialReportRepository, compRepo mocks.MockICompanyRepository, actFieldRepo mocks.MockIActivityFieldRepository) {
				compRepo.EXPECT().
					GetByOwnerId(context.Background(), uuid.UUID{1}, 0, false).
					Return(
						[]*domain.Company{
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
						}, nil).AnyTimes()

				compRepo.EXPECT().
					GetById(
						context.Background(),
						uuid.UUID{1},
					).
					Return(&domain.Company{
						ID:              uuid.UUID{1},
						OwnerID:         uuid.UUID{1},
						ActivityFieldId: uuid.UUID{1},
						Name:            "a",
						City:            "a",
					}, nil)

				actFieldRepo.EXPECT().
					GetById(
						context.Background(),
						uuid.UUID{1},
					).
					Return(
						&domain.ActivityField{
							ID:   uuid.UUID{1},
							Cost: float32(5.0),
						}, nil)

				actFieldRepo.EXPECT().
					GetMaxCost(context.Background()).
					Return(float32(13.5), nil)

				finRepo.EXPECT().
					GetByCompany(
						context.Background(),
						uuid.UUID{1},
						&domain.Period{
							StartYear:    2023,
							EndYear:      2023,
							StartQuarter: 1,
							EndQuarter:   4,
						},
					).
					Return(&domain.FinancialReportByPeriod{
						Reports: []domain.FinancialReport{
							{
								ID:        uuid.UUID{8},
								Year:      2023,
								Quarter:   1,
								Revenue:   32532513,
								Costs:     5436438,
								CompanyID: uuid.UUID{1},
							},
							{
								ID:        uuid.UUID{9},
								Year:      2023,
								Quarter:   2,
								Revenue:   6743634,
								Costs:     9876967,
								CompanyID: uuid.UUID{1},
							},
							{
								ID:        uuid.UUID{10},
								Year:      2023,
								Quarter:   3,
								Revenue:   4675424,
								Costs:     2436653,
								CompanyID: uuid.UUID{1},
							},
							{
								ID:        uuid.UUID{11},
								Year:      2023,
								Quarter:   4,
								Revenue:   14385253,
								Costs:     7546424,
								CompanyID: uuid.UUID{1},
							},
						},
						Period: &domain.Period{
							StartYear:    2023,
							EndYear:      2023,
							StartQuarter: 1,
							EndQuarter:   4,
						},
					}, nil).AnyTimes()

				finRepo.EXPECT().
					GetByCompany(
						context.Background(),
						uuid.UUID{2},
						&domain.Period{
							StartYear:    2023,
							EndYear:      2023,
							StartQuarter: 1,
							EndQuarter:   4,
						},
					).
					Return(&domain.FinancialReportByPeriod{
						Reports: []domain.FinancialReport{
							{
								ID:        uuid.UUID{8},
								Year:      2023,
								Quarter:   1,
								Revenue:   3253251,
								Costs:     543643,
								CompanyID: uuid.UUID{2},
							},
							{
								ID:        uuid.UUID{9},
								Year:      2023,
								Quarter:   2,
								Revenue:   6743634,
								Costs:     9876967,
								CompanyID: uuid.UUID{2},
							},
							{
								ID:        uuid.UUID{10},
								Year:      2023,
								Quarter:   3,
								Revenue:   4675412,
								Costs:     2436765,
								CompanyID: uuid.UUID{2},
							},
							{
								ID:        uuid.UUID{11},
								Year:      2023,
								Quarter:   4,
								Revenue:   1438525,
								Costs:     754642,
								CompanyID: uuid.UUID{2},
							},
						},
						Period: &domain.Period{
							StartYear:    2023,
							EndYear:      2023,
							StartQuarter: 1,
							EndQuarter:   4,
						},
					}, nil).AnyTimes()
			},
			expected: (5.0/13.5 + float32(32532513+6743634+4675424+14385253+3253251+6743634+4675412+1438525-5436438-9876967-2436653-7546424-543643-9876967-2436765-754642)/float32(32532513+6743634+4675424+14385253+3253251+6743634+4675412+1438525)) / 2.0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*userRepo, *finRepo, *compRepo, *actFieldRepo)
			}

			val, err := interactor.CalculateUserRating(ctx, tc.userId)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.InEpsilon(t, tc.expected, val, eps)
			}
		})
	}
}

func TestInteractor_GetMostProfitableCompany(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockIUserRepository(ctrl)
	finRepo := mocks.NewMockIFinancialReportRepository(ctrl)
	compRepo := mocks.NewMockICompanyRepository(ctrl)
	actFieldRepo := mocks.NewMockIActivityFieldRepository(ctrl)

	userSvc := user.NewService(userRepo, compRepo, actFieldRepo)
	actFieldSvc := activity_field.NewService(actFieldRepo, compRepo)
	compSvc := company.NewService(compRepo)
	finSvc := fin_report.NewService(finRepo)

	interactor := NewInteractor(userSvc, actFieldSvc, compSvc, finSvc)

	testCases := []struct {
		name       string
		period     *domain.Period
		companies  []*domain.Company
		beforeTest func(finRepo mocks.MockIFinancialReportRepository)
		expected   *domain.Company
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешный случай",
			period: &domain.Period{
				StartYear:    2023,
				EndYear:      2023,
				StartQuarter: 1,
				EndQuarter:   4,
			},
			companies: []*domain.Company{
				{
					ID: uuid.UUID{1},
				},
				{
					ID: uuid.UUID{2},
				},
			},
			beforeTest: func(finRepo mocks.MockIFinancialReportRepository) {
				finRepo.EXPECT().
					GetByCompany(
						context.Background(),
						uuid.UUID{1},
						&domain.Period{
							StartYear:    2023,
							EndYear:      2023,
							StartQuarter: 1,
							EndQuarter:   4,
						},
					).Return(
					&domain.FinancialReportByPeriod{
						Reports: []domain.FinancialReport{
							{
								ID:        uuid.UUID{1},
								CompanyID: uuid.UUID{1},
								Revenue:   100,
								Costs:     50,
								Year:      2023,
								Quarter:   1,
							},
							{
								ID:        uuid.UUID{2},
								CompanyID: uuid.UUID{1},
								Revenue:   100,
								Costs:     50,
								Year:      2023,
								Quarter:   2,
							},
							{
								ID:        uuid.UUID{3},
								CompanyID: uuid.UUID{1},
								Revenue:   100,
								Costs:     50,
								Year:      2023,
								Quarter:   3,
							},
							{
								ID:        uuid.UUID{4},
								CompanyID: uuid.UUID{1},
								Revenue:   100,
								Costs:     50,
								Year:      2023,
								Quarter:   4,
							},
						},
					}, nil)

				finRepo.EXPECT().
					GetByCompany(
						context.Background(),
						uuid.UUID{2},
						&domain.Period{
							StartYear:    2023,
							EndYear:      2023,
							StartQuarter: 1,
							EndQuarter:   4,
						},
					).Return(
					&domain.FinancialReportByPeriod{
						Reports: []domain.FinancialReport{
							{
								ID:        uuid.UUID{5},
								CompanyID: uuid.UUID{2},
								Revenue:   75,
								Costs:     50,
								Year:      2023,
								Quarter:   1,
							},
							{
								ID:        uuid.UUID{6},
								CompanyID: uuid.UUID{2},
								Revenue:   75,
								Costs:     50,
								Year:      2023,
								Quarter:   2,
							},
							{
								ID:        uuid.UUID{7},
								CompanyID: uuid.UUID{2},
								Revenue:   75,
								Costs:     50,
								Year:      2023,
								Quarter:   3,
							},
							{
								ID:        uuid.UUID{8},
								CompanyID: uuid.UUID{2},
								Revenue:   75,
								Costs:     50,
								Year:      2023,
								Quarter:   4,
							},
						},
					}, nil)
			},
			expected: &domain.Company{
				ID: uuid.UUID{1},
			},
			wantErr: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*finRepo)
			}

			company, err := interactor.GetMostProfitableCompany(ctx, tc.period, tc.companies)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.Equal(t, tc.expected, company)
			}
		})
	}
}

func TestInteractor_GetUserFinancialReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockIUserRepository(ctrl)
	finRepo := mocks.NewMockIFinancialReportRepository(ctrl)
	compRepo := mocks.NewMockICompanyRepository(ctrl)
	actFieldRepo := mocks.NewMockIActivityFieldRepository(ctrl)

	userSvc := user.NewService(userRepo, compRepo, actFieldRepo)
	actFieldSvc := activity_field.NewService(actFieldRepo, compRepo)
	compSvc := company.NewService(compRepo)
	finSvc := fin_report.NewService(finRepo)

	interactor := NewInteractor(userSvc, actFieldSvc, compSvc, finSvc)

	testCases := []struct {
		name       string
		userId     uuid.UUID
		beforeTest func(
			userRepo mocks.MockIUserRepository,
			finRepo mocks.MockIFinancialReportRepository,
			compRepo mocks.MockICompanyRepository,
			actFieldRepo mocks.MockIActivityFieldRepository,
		)
		period   *domain.Period
		expected *domain.FinancialReportByPeriod
		wantErr  bool
		errStr   error
	}{
		{
			name:   "успешный тест",
			userId: uuid.UUID{1},
			beforeTest: func(userRepo mocks.MockIUserRepository, finRepo mocks.MockIFinancialReportRepository, compRepo mocks.MockICompanyRepository, actFieldRepo mocks.MockIActivityFieldRepository) {
				compRepo.EXPECT().
					GetByOwnerId(context.Background(), uuid.UUID{1}, 0, false).
					Return(
						[]*domain.Company{
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
						}, nil)

				finRepo.EXPECT().
					GetByCompany(
						context.Background(),
						uuid.UUID{1},
						&domain.Period{
							StartYear:    2023,
							EndYear:      2024,
							StartQuarter: 1,
							EndQuarter:   1,
						},
					).Return(
					&domain.FinancialReportByPeriod{
						Reports: []domain.FinancialReport{
							{
								ID:        uuid.UUID{1},
								CompanyID: uuid.UUID{1},
								Revenue:   100,
								Costs:     50,
								Year:      2023,
								Quarter:   1,
							},
							{
								ID:        uuid.UUID{2},
								CompanyID: uuid.UUID{1},
								Revenue:   100,
								Costs:     50,
								Year:      2023,
								Quarter:   2,
							},
							{
								ID:        uuid.UUID{3},
								CompanyID: uuid.UUID{1},
								Revenue:   100,
								Costs:     50,
								Year:      2023,
								Quarter:   3,
							},
							{
								ID:        uuid.UUID{4},
								CompanyID: uuid.UUID{1},
								Revenue:   100,
								Costs:     50,
								Year:      2023,
								Quarter:   4,
							},
							{
								ID:        uuid.UUID{5},
								CompanyID: uuid.UUID{1},
								Revenue:   100,
								Costs:     50,
								Year:      2024,
								Quarter:   1,
							},
						},
						Period: &domain.Period{
							StartYear:    2023,
							EndYear:      2024,
							StartQuarter: 1,
							EndQuarter:   1,
						},
					}, nil)

				finRepo.EXPECT().
					GetByCompany(
						context.Background(),
						uuid.UUID{2},
						&domain.Period{
							StartYear:    2023,
							EndYear:      2024,
							StartQuarter: 1,
							EndQuarter:   1,
						},
					).Return(
					&domain.FinancialReportByPeriod{
						Reports: []domain.FinancialReport{
							{
								ID:        uuid.UUID{6},
								CompanyID: uuid.UUID{2},
								Revenue:   75,
								Costs:     50,
								Year:      2023,
								Quarter:   1,
							},
							{
								ID:        uuid.UUID{7},
								CompanyID: uuid.UUID{2},
								Revenue:   75,
								Costs:     50,
								Year:      2023,
								Quarter:   2,
							},
							{
								ID:        uuid.UUID{8},
								CompanyID: uuid.UUID{2},
								Revenue:   75,
								Costs:     50,
								Year:      2023,
								Quarter:   3,
							},
							{
								ID:        uuid.UUID{9},
								CompanyID: uuid.UUID{2},
								Revenue:   75,
								Costs:     50,
								Year:      2023,
								Quarter:   4,
							},
							{
								ID:        uuid.UUID{10},
								CompanyID: uuid.UUID{2},
								Revenue:   75,
								Costs:     50,
								Year:      2024,
								Quarter:   1,
							},
						},
						Period: &domain.Period{
							StartYear:    2023,
							EndYear:      2024,
							StartQuarter: 1,
							EndQuarter:   1,
						},
					}, nil)
			},
			period: &domain.Period{
				StartYear:    2023,
				EndYear:      2024,
				StartQuarter: 1,
				EndQuarter:   1,
			},
			expected: &domain.FinancialReportByPeriod{
				Reports: []domain.FinancialReport{
					{
						ID:        uuid.UUID{1},
						CompanyID: uuid.UUID{1},
						Revenue:   100,
						Costs:     50,
						Year:      2023,
						Quarter:   1,
					},
					{
						ID:        uuid.UUID{2},
						CompanyID: uuid.UUID{1},
						Revenue:   100,
						Costs:     50,
						Year:      2023,
						Quarter:   2,
					},
					{
						ID:        uuid.UUID{3},
						CompanyID: uuid.UUID{1},
						Revenue:   100,
						Costs:     50,
						Year:      2023,
						Quarter:   3,
					},
					{
						ID:        uuid.UUID{4},
						CompanyID: uuid.UUID{1},
						Revenue:   100,
						Costs:     50,
						Year:      2023,
						Quarter:   4,
					},
					{
						ID:        uuid.UUID{5},
						CompanyID: uuid.UUID{1},
						Revenue:   100,
						Costs:     50,
						Year:      2024,
						Quarter:   1,
					},
					{
						ID:        uuid.UUID{6},
						CompanyID: uuid.UUID{2},
						Revenue:   75,
						Costs:     50,
						Year:      2023,
						Quarter:   1,
					},
					{
						ID:        uuid.UUID{7},
						CompanyID: uuid.UUID{2},
						Revenue:   75,
						Costs:     50,
						Year:      2023,
						Quarter:   2,
					},
					{
						ID:        uuid.UUID{8},
						CompanyID: uuid.UUID{2},
						Revenue:   75,
						Costs:     50,
						Year:      2023,
						Quarter:   3,
					},
					{
						ID:        uuid.UUID{9},
						CompanyID: uuid.UUID{2},
						Revenue:   75,
						Costs:     50,
						Year:      2023,
						Quarter:   4,
					},
					{
						ID:        uuid.UUID{10},
						CompanyID: uuid.UUID{2},
						Revenue:   75,
						Costs:     50,
						Year:      2024,
						Quarter:   1,
					},
				},
				Period: &domain.Period{
					StartYear:    2023,
					EndYear:      2024,
					StartQuarter: 1,
					EndQuarter:   1,
				},
				Taxes:   float32(((100 - 50) + (75 - 50)) * 4 * 0.04),
				TaxLoad: float32((((100 - 50) + (75 - 50)) * 4 * 0.04) / ((100 + 75) * 4) * 100),
			},
			wantErr: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.beforeTest != nil {
				tc.beforeTest(*userRepo, *finRepo, *compRepo, *actFieldRepo)
			}

			report, err := interactor.GetUserFinancialReport(ctx, tc.userId, tc.period)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.Equal(t, tc.expected.Reports, report.Reports)
				require.Equal(t, tc.expected.Period, report.Period)
				require.InEpsilon(t, tc.expected.Taxes, report.Taxes, eps)
				require.InEpsilon(t, tc.expected.TaxLoad, report.TaxLoad, eps)
			}
		})
	}
}

func Test_calcRating(t *testing.T) {
	testCases := []struct {
		name     string
		profit   float32
		revenue  float32
		cost     float32
		maxCost  float32
		expected float32
	}{
		{
			name:     "успешное вычисление",
			profit:   100,
			revenue:  1000,
			cost:     5.0,
			maxCost:  13.5,
			expected: (5.0/13.5 + 100.0/1000.0) / 2.0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rating := calcRating(tc.profit, tc.revenue, tc.cost, tc.maxCost)

			require.InEpsilon(t, tc.expected, rating, eps)
		})
	}
}

func Test_calculateTaxes(t *testing.T) {
	testCases := []struct {
		name     string
		reports  map[int]*domain.FinancialReportByPeriod
		expected *taxesData
		//expected float32
		wantErr bool
		errStr  error
	}{
		{
			name: "успешное вычисление",
			reports: map[int]*domain.FinancialReportByPeriod{
				1: {
					Reports: []domain.FinancialReport{
						{
							ID:        uuid.UUID{1},
							CompanyID: uuid.UUID{1},
							Revenue:   12432532,
							Costs:     3213214,
							Year:      1,
							Quarter:   2,
						},
						{
							ID:        uuid.UUID{2},
							CompanyID: uuid.UUID{1},
							Revenue:   12432532,
							Costs:     3213214,
							Year:      1,
							Quarter:   3,
						},
						{
							ID:        uuid.UUID{3},
							CompanyID: uuid.UUID{1},
							Revenue:   12432532,
							Costs:     3213214,
							Year:      1,
							Quarter:   4,
						},
					},
					Period: &domain.Period{
						StartYear:    1,
						EndYear:      1,
						StartQuarter: 2,
						EndQuarter:   4,
					},
				},
				2: {
					Reports: []domain.FinancialReport{
						{
							ID:        uuid.UUID{4},
							CompanyID: uuid.UUID{1},
							Revenue:   12432532,
							Costs:     3213214,
							Year:      2,
							Quarter:   1,
						},
						{
							ID:        uuid.UUID{5},
							CompanyID: uuid.UUID{1},
							Revenue:   12432532,
							Costs:     3213214,
							Year:      2,
							Quarter:   2,
						},
						{
							ID:        uuid.UUID{6},
							CompanyID: uuid.UUID{1},
							Revenue:   12432532,
							Costs:     3213214,
							Year:      2,
							Quarter:   3,
						},
						{
							ID:        uuid.UUID{7},
							CompanyID: uuid.UUID{1},
							Revenue:   12432532,
							Costs:     3213214,
							Year:      2,
							Quarter:   4,
						},
					},
					Period: &domain.Period{
						StartYear:    2,
						EndYear:      2,
						StartQuarter: 1,
						EndQuarter:   4,
					},
				},
			},
			expected: &taxesData{
				taxes:   (12432532 - 3213214) * 4.0 * 0.07,
				revenue: 12432532 * 4,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tax := calculateTaxes(tc.reports)

			require.InEpsilon(t, tc.expected.taxes, tax.taxes, eps)
			require.InEpsilon(t, tc.expected.revenue, tax.revenue, eps)
		})
	}
}

func Test_findFullYearReports(t *testing.T) {
	testCases := []struct {
		name     string
		reports  *domain.FinancialReportByPeriod
		period   *domain.Period
		expected map[int]*domain.FinancialReportByPeriod
	}{
		{
			name: "успешное получение отчетов за полные годы",
			reports: &domain.FinancialReportByPeriod{
				Reports: []domain.FinancialReport{
					{
						ID:        uuid.UUID{1},
						CompanyID: uuid.UUID{1},
						Revenue:   12432532,
						Costs:     3213214,
						Year:      1,
						Quarter:   2,
					},
					{
						ID:        uuid.UUID{2},
						CompanyID: uuid.UUID{1},
						Revenue:   12432532,
						Costs:     3213214,
						Year:      1,
						Quarter:   3,
					},
					{
						ID:        uuid.UUID{3},
						CompanyID: uuid.UUID{1},
						Revenue:   12432532,
						Costs:     3213214,
						Year:      1,
						Quarter:   4,
					},
					{
						ID:        uuid.UUID{4},
						CompanyID: uuid.UUID{1},
						Revenue:   12432532,
						Costs:     3213214,
						Year:      2,
						Quarter:   1,
					},
					{
						ID:        uuid.UUID{5},
						CompanyID: uuid.UUID{1},
						Revenue:   12432532,
						Costs:     3213214,
						Year:      2,
						Quarter:   2,
					},
					{
						ID:        uuid.UUID{6},
						CompanyID: uuid.UUID{1},
						Revenue:   12432532,
						Costs:     3213214,
						Year:      2,
						Quarter:   3,
					},
					{
						ID:        uuid.UUID{7},
						CompanyID: uuid.UUID{1},
						Revenue:   12432532,
						Costs:     3213214,
						Year:      2,
						Quarter:   4,
					},
				},
				Period: &domain.Period{
					StartYear:    1,
					EndYear:      2,
					StartQuarter: 2,
					EndQuarter:   4,
				},
			},
			period: &domain.Period{
				StartYear:    1,
				EndYear:      2,
				StartQuarter: 2,
				EndQuarter:   4,
			},
			expected: map[int]*domain.FinancialReportByPeriod{
				2: {
					Reports: []domain.FinancialReport{
						{
							ID:        uuid.UUID{4},
							CompanyID: uuid.UUID{1},
							Revenue:   12432532,
							Costs:     3213214,
							Year:      2,
							Quarter:   1,
						},
						{
							ID:        uuid.UUID{5},
							CompanyID: uuid.UUID{1},
							Revenue:   12432532,
							Costs:     3213214,
							Year:      2,
							Quarter:   2,
						},
						{
							ID:        uuid.UUID{6},
							CompanyID: uuid.UUID{1},
							Revenue:   12432532,
							Costs:     3213214,
							Year:      2,
							Quarter:   3,
						},
						{
							ID:        uuid.UUID{7},
							CompanyID: uuid.UUID{1},
							Revenue:   12432532,
							Costs:     3213214,
							Year:      2,
							Quarter:   4,
						},
					},
					Period: &domain.Period{
						StartYear:    2,
						EndYear:      2,
						StartQuarter: 1,
						EndQuarter:   4,
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reports := findFullYearReports(tc.reports, tc.period)

			require.Equal(t, tc.expected, reports)
		})
	}
}
