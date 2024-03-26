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

			val, err := interactor.CalculateUserRating(context.Background(), tc.userId)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.Equal(t, tc.expected, val)
			}
		})
	}
}
