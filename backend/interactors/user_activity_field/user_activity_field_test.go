package user_activity_field

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"math"
	"ppo/domain"
	"ppo/mocks"
	"ppo/services/activity_field"
	"ppo/services/company"
	"ppo/services/fin_report"
	"ppo/services/user"
	"reflect"
	"testing"
)

func TestInteractor_CalculateUserRating(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockIUserRepository(ctrl)
	finRepo := mocks.NewMockIFinancialReportRepository(ctrl)
	compRepo := mocks.NewMockICompanyRepository(ctrl)
	actFieldRepo := mocks.NewMockIActivityFieldRepository(ctrl)

	userSvc := user.NewService(userRepo, compRepo, finRepo, actFieldRepo)
	actFieldSvc := activity_field.NewService(actFieldRepo)
	compSvc := company.NewService(compRepo, finRepo)
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
					GetByOwnerId(context.Background(), uuid.UUID{1}).
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

				actFieldRepo.EXPECT().
					GetByCompanyId(
						context.Background(),
						uuid.UUID{1},
					).
					Return(float32(5.0), nil)

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
								Revenue: 4675424,
								Costs:   2436653,
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
							StartYear:    2023,
							EndYear:      2023,
							StartQuarter: 1,
							EndQuarter:   4,
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
					).
					Return(&domain.FinancialReportByPeriod{
						Reports: []domain.FinancialReport{
							{
								ID:      uuid.UUID{8},
								Year:    2023,
								Quarter: 1,
								Revenue: 3253251,
								Costs:   543643,
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
								Revenue: 4675412,
								Costs:   2436765,
							},
							{
								ID:      uuid.UUID{11},
								Year:    2023,
								Quarter: 4,
								Revenue: 1438525,
								Costs:   754642,
							},
						},
						Period: &domain.Period{
							StartYear:    2023,
							EndYear:      2023,
							StartQuarter: 1,
							EndQuarter:   4,
						},
					}, nil)
			},
			//expected: 1e-9 * (float32(math.Pow(float64(32532513-5436438+6743634-9876967+4675424-2436653+14385253-7546424), 0.5))*0.9 + 1.2*5*5 + 0.35*(32532513+6743634+4675424+14385253)),
			expected: 1e-9 * (float32(math.Pow(float64(32532513-5436438+6743634-9876967+4675424-2436653+14385253-7546424+3253251-543643+6743634-9876967+4675412-2436765+1438525-754642), 0.5))*0.9 + 1.2*5*5 + 0.35*(32532513+6743634+4675424+14385253+3253251+6743634+4675412+1438525)),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.beforeTest != nil {
				tc.beforeTest(*userRepo, *finRepo, *compRepo, *actFieldRepo)
			}

			val, err := interactor.CalculateUserRating(tc.userId)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.Equal(t, tc.expected, val)
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

	userSvc := user.NewService(userRepo, compRepo, finRepo, actFieldRepo)
	actFieldSvc := activity_field.NewService(actFieldRepo)
	compSvc := company.NewService(compRepo, finRepo)
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
			if tc.beforeTest != nil {
				tc.beforeTest(*finRepo)
			}

			company, err := interactor.GetMostProfitableCompany(tc.period, tc.companies)

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
	type fields struct {
		userService     domain.IUserService
		actFieldService domain.IActivityFieldService
		compService     domain.ICompanyService
		finService      domain.IFinancialReportService
	}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantReport *domain.FinancialReportByPeriod
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Interactor{
				userService:     tt.fields.userService,
				actFieldService: tt.fields.actFieldService,
				compService:     tt.fields.compService,
				finService:      tt.fields.finService,
			}
			gotReport, err := i.GetUserFinancialReport(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserFinancialReport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotReport, tt.wantReport) {
				t.Errorf("GetUserFinancialReport() gotReport = %v, want %v", gotReport, tt.wantReport)
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

			require.InEpsilon(t, tc.expected, rating, 1e-7)
		})
	}
}

func Test_calculateTaxes(t *testing.T) {
	testCases := []struct {
		name     string
		reports  map[int]domain.FinancialReportByPeriod
		expected *taxes
		wantErr  bool
		errStr   error
	}{
		{
			name: "успешное вычисление",
			reports: map[int]domain.FinancialReportByPeriod{
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
			expected: &taxes{
				Sum:  (12432532 - 3213214) * 4.0 * 0.07,
				Load: ((12432532 - 3213214) * 4.0 * 0.07) / (12432532 * 4) * 100,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tax := calculateTaxes(tc.reports)

			require.InEpsilon(t, tc.expected.Sum, tax.Sum, 1e-7)
			require.InEpsilon(t, tc.expected.Load, tax.Load, 1e-7)
		})
	}
}

func Test_findFullYearReports(t *testing.T) {
	testCases := []struct {
		name     string
		reports  *domain.FinancialReportByPeriod
		period   *domain.Period
		expected map[int]domain.FinancialReportByPeriod
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
			expected: map[int]domain.FinancialReportByPeriod{
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
