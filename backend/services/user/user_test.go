package user

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
	"time"
)

//func TestUserService_CalculateRating(t *testing.T) {} // TODO: mainFieldWeight...

func TestUserService_DeleteById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockIUserRepository(ctrl)
	compRepo := mocks.NewMockICompanyRepository(ctrl)
	finRepo := mocks.NewMockIFinancialReportRepository(ctrl)
	actFieldRepo := mocks.NewMockIActivityFieldRepository(ctrl)
	svc := NewService(userRepo, compRepo, finRepo, actFieldRepo)

	curUuid := uuid.New()

	testCases := []struct {
		name       string
		id         uuid.UUID
		beforeTest func(userRepo mocks.MockIUserRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное удаление",
			id:   curUuid,
			beforeTest: func(userRepo mocks.MockIUserRepository) {
				userRepo.EXPECT().
					DeleteById(context.Background(), curUuid).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ошибка выполнения запроса в репозитории",
			id:   curUuid,
			beforeTest: func(userRepo mocks.MockIUserRepository) {
				userRepo.EXPECT().
					DeleteById(context.Background(), curUuid).
					Return(fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("удаление пользователя по id: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.beforeTest != nil {
				tc.beforeTest(*userRepo)
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

func TestUserService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockIUserRepository(ctrl)
	compRepo := mocks.NewMockICompanyRepository(ctrl)
	finRepo := mocks.NewMockIFinancialReportRepository(ctrl)
	actFieldRepo := mocks.NewMockIActivityFieldRepository(ctrl)
	svc := NewService(userRepo, compRepo, finRepo, actFieldRepo)

	testCases := []struct {
		name       string
		beforeTest func(userRepo mocks.MockIUserRepository)
		expected   []*domain.User
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное получение списка всех компаний",
			beforeTest: func(userRepo mocks.MockIUserRepository) {
				userRepo.EXPECT().
					GetAll(context.Background(), nil).
					Return([]*domain.User{
						{
							ID:       uuid.UUID{1},
							Username: "a",
							FullName: "a",
							Gender:   "m",
							Birthday: time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
							City:     "a",
						},
						{
							ID:       uuid.UUID{2},
							Username: "b",
							FullName: "b",
							Gender:   "w",
							Birthday: time.Date(2, 2, 2, 2, 2, 2, 2, time.Local),
							City:     "b",
						},
						{
							ID:       uuid.UUID{3},
							Username: "c",
							FullName: "c",
							Gender:   "m",
							Birthday: time.Date(3, 3, 3, 3, 3, 3, 3, time.Local),
							City:     "c",
						},
					}, nil)
			},
			expected: []*domain.User{
				{
					ID:       uuid.UUID{1},
					Username: "a",
					FullName: "a",
					Gender:   "m",
					Birthday: time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
					City:     "a",
				},
				{
					ID:       uuid.UUID{2},
					Username: "b",
					FullName: "b",
					Gender:   "w",
					Birthday: time.Date(2, 2, 2, 2, 2, 2, 2, time.Local),
					City:     "b",
				},
				{
					ID:       uuid.UUID{3},
					Username: "c",
					FullName: "c",
					Gender:   "m",
					Birthday: time.Date(3, 3, 3, 3, 3, 3, 3, time.Local),
					City:     "c",
				},
			},
			wantErr: false,
		},
		{
			name: "ошибка получения данных в репозитории",
			beforeTest: func(userRepo mocks.MockIUserRepository) {
				userRepo.EXPECT().
					GetAll(context.Background(), nil).
					Return(nil, fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("получение списка всех пользователей: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.beforeTest != nil {
				tc.beforeTest(*userRepo)
			}

			users, err := svc.GetAll(context.Background(), nil)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.Equal(t, users, tc.expected)
			}
		})
	}
}

func TestUserService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockIUserRepository(ctrl)
	compRepo := mocks.NewMockICompanyRepository(ctrl)
	finRepo := mocks.NewMockIFinancialReportRepository(ctrl)
	actFieldRepo := mocks.NewMockIActivityFieldRepository(ctrl)
	svc := NewService(userRepo, compRepo, finRepo, actFieldRepo)

	testCases := []struct {
		name       string
		user       domain.User
		beforeTest func(userRepo mocks.MockIUserRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное добавление",
			user: domain.User{
				ID:       uuid.UUID{1},
				Username: "a",
				FullName: "a b c",
				Gender:   "m",
				Birthday: time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
				City:     "a",
			},
			beforeTest: func(userRepo mocks.MockIUserRepository) {
				userRepo.EXPECT().
					Create(
						context.Background(),
						&domain.User{
							ID:       uuid.UUID{1},
							Username: "a",
							FullName: "a b c",
							Gender:   "m",
							Birthday: time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
							City:     "a",
						},
					).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "некорректное ФИО (не 3 слова)",
			user: domain.User{
				ID:       uuid.UUID{1},
				Username: "a",
				FullName: "a",
				Gender:   "m",
				Birthday: time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
				City:     "a",
			},
			beforeTest: func(userRepo mocks.MockIUserRepository) {
				compRepo.EXPECT().
					Create(
						context.Background(),
						&domain.User{
							ID:       uuid.UUID{1},
							Username: "a",
							FullName: "a",
							Gender:   "m",
							Birthday: time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
							City:     "a",
						},
					).Return(nil).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("некорректное количество слов (должны быть фамилия, имя и отчество)"),
		},
		{
			name: "пустое название города",
			user: domain.User{
				ID:       uuid.UUID{1},
				Username: "a",
				FullName: "a b c",
				Gender:   "m",
				Birthday: time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
				City:     "",
			},
			beforeTest: func(userRepo mocks.MockIUserRepository) {
				compRepo.EXPECT().
					Create(
						context.Background(),
						&domain.User{
							ID:       uuid.UUID{1},
							Username: "a",
							FullName: "a b c",
							Gender:   "m",
							Birthday: time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
							City:     "",
						},
					).Return(nil).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("должно быть указано название города"),
		},
		{
			name: "неизвестный пол",
			user: domain.User{
				ID:       uuid.UUID{1},
				Username: "a",
				FullName: "a b c",
				Gender:   "r",
				Birthday: time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
				City:     "c",
			},
			beforeTest: func(userRepo mocks.MockIUserRepository) {
				compRepo.EXPECT().
					Create(
						context.Background(),
						&domain.User{
							ID:       uuid.UUID{1},
							Username: "a",
							FullName: "a b c",
							Gender:   "r",
							Birthday: time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
							City:     "c",
						},
					).Return(nil).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("неизвестный пол"),
		},
		{
			name: "пустая дата рождения",
			user: domain.User{
				ID:       uuid.UUID{1},
				Username: "a",
				FullName: "a b c",
				Gender:   "w",
				Birthday: time.Time{},
				City:     "c",
			},
			beforeTest: func(userRepo mocks.MockIUserRepository) {
				compRepo.EXPECT().
					Create(
						context.Background(),
						&domain.User{
							ID:       uuid.UUID{1},
							Username: "a",
							FullName: "a b c",
							Gender:   "w",
							Birthday: time.Time{},
							City:     "c",
						},
					).Return(nil).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("должна быть указана дата рождения"),
		},
		{
			name: "ошибка выполнения запроса в репозитории",
			user: domain.User{
				ID:       uuid.UUID{1},
				Username: "a",
				FullName: "a b c",
				Gender:   "w",
				Birthday: time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
				City:     "c",
			},
			beforeTest: func(userRepo mocks.MockIUserRepository) {
				userRepo.EXPECT().
					Create(
						context.Background(),
						&domain.User{
							ID:       uuid.UUID{1},
							Username: "a",
							FullName: "a b c",
							Gender:   "w",
							Birthday: time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
							City:     "c",
						},
					).Return(fmt.Errorf("sql error")).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("создание пользователя: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.beforeTest != nil {
				tc.beforeTest(*userRepo)
			}

			err := svc.Create(context.Background(), &tc.user)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestUserService_GetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockIUserRepository(ctrl)
	compRepo := mocks.NewMockICompanyRepository(ctrl)
	finRepo := mocks.NewMockIFinancialReportRepository(ctrl)
	actFieldRepo := mocks.NewMockIActivityFieldRepository(ctrl)
	svc := NewService(userRepo, compRepo, finRepo, actFieldRepo)

	testCases := []struct {
		name       string
		id         uuid.UUID
		beforeTest func(userRepo mocks.MockIUserRepository)
		expected   *domain.User
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное получение компании по id",
			id:   uuid.UUID{1},
			beforeTest: func(userRepo mocks.MockIUserRepository) {
				userRepo.EXPECT().
					GetById(
						context.Background(),
						uuid.UUID{1},
					).
					Return(&domain.User{
						ID:       uuid.UUID{1},
						Username: "a",
						FullName: "a b c",
						Gender:   "m",
						Birthday: time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
						City:     "a",
					}, nil)
			},
			expected: &domain.User{
				ID:       uuid.UUID{1},
				Username: "a",
				FullName: "a b c",
				Gender:   "m",
				Birthday: time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
				City:     "a",
			},
			wantErr: false,
		},
		{
			name: "ошибка получения данных в репозитории",
			id:   uuid.UUID{1},
			beforeTest: func(userRepo mocks.MockIUserRepository) {
				userRepo.EXPECT().
					GetById(
						context.Background(),
						uuid.UUID{1},
					).
					Return(nil, fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("получение пользователя по id: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.beforeTest != nil {
				tc.beforeTest(*userRepo)
			}

			user, err := svc.GetById(context.Background(), tc.id)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.Equal(t, user, tc.expected)
			}
		})
	}
}

func TestUserService_GetFinancialReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockIUserRepository(ctrl)
	compRepo := mocks.NewMockICompanyRepository(ctrl)
	finRepo := mocks.NewMockIFinancialReportRepository(ctrl)
	actFieldRepo := mocks.NewMockIActivityFieldRepository(ctrl)
	svc := NewService(userRepo, compRepo, finRepo, actFieldRepo)

	testCases := []struct {
		name       string
		id         uuid.UUID
		period     *domain.Period
		beforeTest func(userRepo mocks.MockIUserRepository, compRepo mocks.MockICompanyRepository, finRepo mocks.MockIFinancialReportRepository)
		expected   []*domain.FinancialReportByPeriod
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное получение отчета деятельности пользователя по id",
			id:   uuid.UUID{1},
			period: &domain.Period{
				StartYear:    2021,
				EndYear:      2023,
				StartQuarter: 2,
				EndQuarter:   4,
			},
			beforeTest: func(userRepo mocks.MockIUserRepository, compRepo mocks.MockICompanyRepository, finRepo mocks.MockIFinancialReportRepository) {
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
							},
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

				finRepo.EXPECT().
					GetByCompany(
						context.Background(),
						uuid.UUID{2},
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
			expected: []*domain.FinancialReportByPeriod{
				{
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
					Taxes:   (43635325+50934123+78902453+64352357-12362332-13543623-15326443-23534252)*0.2 + (32532513+6743634+46754124+14385253-5436438-9876967-24367653-7546424)*0.13,
					TaxLoad: (43635325+50934123+78902453+64352357-12362332-13543623-15326443-23534252)*0.2/(43635325+50934123+78902453+64352357)*100 + (32532513+6743634+46754124+14385253-5436438-9876967-24367653-7546424)*0.13/(32532513+6743634+46754124+14385253)*100,
				},
				{
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
					Taxes:   (43635325+50934123+78902453+64352357-12362332-13543623-15326443-23534252)*0.2 + (32532513+6743634+46754124+14385253-5436438-9876967-24367653-7546424)*0.13,
					TaxLoad: (43635325+50934123+78902453+64352357-12362332-13543623-15326443-23534252)*0.2/(43635325+50934123+78902453+64352357)*100 + (32532513+6743634+46754124+14385253-5436438-9876967-24367653-7546424)*0.13/(32532513+6743634+46754124+14385253)*100,
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
			beforeTest: func(userRepo mocks.MockIUserRepository, compRepo mocks.MockICompanyRepository, finRepo mocks.MockIFinancialReportRepository) {
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
						}, nil).
					AnyTimes()

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
					}, nil).
					AnyTimes()

				finRepo.EXPECT().
					GetByCompany(
						context.Background(),
						uuid.UUID{2},
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
					}, nil).
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
			beforeTest: func(userRepo mocks.MockIUserRepository, compRepo mocks.MockICompanyRepository, finRepo mocks.MockIFinancialReportRepository) {
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
						}, nil).
					AnyTimes()

				finRepo.EXPECT().
					GetByCompany(
						context.Background(),
						uuid.UUID{1},
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
					}, nil).
					AnyTimes()

				finRepo.EXPECT().
					GetByCompany(
						context.Background(),
						uuid.UUID{2},
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
					}, nil).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("дата конца периода должна быть позже даты начала"),
		},
		//{
		//	name: "ошибка получения данных в репозитории",
		//	id:   uuid.UUID{1},
		//	period: &domain.Period{
		//		StartYear:    2021,
		//		EndYear:      2023,
		//		StartQuarter: 2,
		//		EndQuarter:   4,
		//	},
		//	beforeTest: func(userRepo mocks.MockIUserRepository, compRepo mocks.MockICompanyRepository, finRepo mocks.MockIFinancialReportRepository) {
		//		compRepo.EXPECT().
		//			GetByOwnerId(context.Background(), uuid.UUID{1}).
		//			Return(nil, fmt.Errorf("sql error")).
		//			AnyTimes()
		//	},
		//	wantErr: true,
		//	errStr:  errors.New("получение списка компаний предпринимателя по id: sql error"),
		//}, // FIXME
		//{
		//	name: "ошибка получения данных в репозитории_2",
		//	id:   uuid.UUID{1},
		//	period: &domain.Period{
		//		StartYear:    2021,
		//		EndYear:      2023,
		//		StartQuarter: 2,
		//		EndQuarter:   4,
		//	},
		//	beforeTest: func(userRepo mocks.MockIUserRepository, compRepo mocks.MockICompanyRepository, finRepo mocks.MockIFinancialReportRepository) {
		//		compRepo.EXPECT().
		//			GetByOwnerId(context.Background(), uuid.UUID{1}).
		//			Return(
		//				[]*domain.Company{
		//					{
		//						ID:      uuid.UUID{1},
		//						OwnerID: uuid.UUID{1},
		//						Name:    "a",
		//						City:    "a",
		//					},
		//					{
		//						ID:      uuid.UUID{2},
		//						OwnerID: uuid.UUID{1},
		//						Name:    "b",
		//						City:    "b",
		//					},
		//				}, nil).
		//			AnyTimes()
		//
		//		finRepo.EXPECT().
		//			GetByCompany(
		//				context.Background(),
		//				uuid.UUID{1},
		//				&domain.Period{
		//					StartYear:    2021,
		//					EndYear:      2023,
		//					StartQuarter: 2,
		//					EndQuarter:   4,
		//				},
		//			).
		//			Return(nil, fmt.Errorf("sql error")).
		//			AnyTimes()
		//
		//		finRepo.EXPECT().
		//			GetByCompany(
		//				context.Background(),
		//				uuid.UUID{2},
		//				&domain.Period{
		//					StartYear:    2021,
		//					EndYear:      2023,
		//					StartQuarter: 2,
		//					EndQuarter:   4,
		//				},
		//			).
		//			Return(&domain.FinancialReportByPeriod{
		//				Reports: []domain.FinancialReport{
		//					{
		//						ID:      uuid.UUID{1},
		//						Year:    2021,
		//						Quarter: 2,
		//						Revenue: 1432523,
		//						Costs:   75423,
		//					},
		//					{
		//						ID:      uuid.UUID{2},
		//						Year:    2021,
		//						Quarter: 3,
		//						Revenue: 7435235,
		//						Costs:   125654,
		//					},
		//					{
		//						ID:      uuid.UUID{3},
		//						Year:    2021,
		//						Quarter: 4,
		//						Revenue: 65742,
		//						Costs:   7845634,
		//					},
		//					{
		//						ID:      uuid.UUID{4},
		//						Year:    2022,
		//						Quarter: 1,
		//						Revenue: 43635325,
		//						Costs:   12362332,
		//					},
		//					{
		//						ID:      uuid.UUID{5},
		//						Year:    2022,
		//						Quarter: 2,
		//						Revenue: 50934123,
		//						Costs:   13543623,
		//					},
		//					{
		//						ID:      uuid.UUID{6},
		//						Year:    2022,
		//						Quarter: 3,
		//						Revenue: 78902453,
		//						Costs:   15326443,
		//					},
		//					{
		//						ID:      uuid.UUID{7},
		//						Year:    2022,
		//						Quarter: 4,
		//						Revenue: 64352357,
		//						Costs:   23534252,
		//					}, // 173 057 608 => 34 611 521.6; 237 824 258 => 14.5534025
		//					{
		//						ID:      uuid.UUID{8},
		//						Year:    2023,
		//						Quarter: 1,
		//						Revenue: 32532513,
		//						Costs:   5436438,
		//					},
		//					{
		//						ID:      uuid.UUID{9},
		//						Year:    2023,
		//						Quarter: 2,
		//						Revenue: 6743634,
		//						Costs:   9876967,
		//					},
		//					{
		//						ID:      uuid.UUID{10},
		//						Year:    2023,
		//						Quarter: 3,
		//						Revenue: 46754124,
		//						Costs:   24367653,
		//					},
		//					{
		//						ID:      uuid.UUID{11},
		//						Year:    2023,
		//						Quarter: 4,
		//						Revenue: 14385253,
		//						Costs:   7546424,
		//					},
		//				},
		//				Period: &domain.Period{
		//					StartYear:    2021,
		//					EndYear:      2023,
		//					StartQuarter: 2,
		//					EndQuarter:   4,
		//				},
		//			}, nil).
		//			AnyTimes()
		//	},
		//	wantErr: true,
		//	errStr:  errors.New("получение финансовой отчетности компании по id: sql error"),
		//}, // FIXME
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.beforeTest != nil {
				tc.beforeTest(*userRepo, *compRepo, *finRepo)
			}

			report, err := svc.GetFinancialReport(context.Background(), tc.id, tc.period)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				for i := 0; i < len(report); i++ {
					require.Equal(t, report[i].Reports, tc.expected[i].Reports)
					require.Equal(t, report[i].Period, tc.expected[i].Period)
					require.InEpsilon(t, report[i].Taxes, tc.expected[i].Taxes, 1e-7)
					require.InEpsilon(t, report[i].TaxLoad, tc.expected[i].TaxLoad, 1e-7)
				}
			}
		})
	}
}

func TestUserService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockIUserRepository(ctrl)
	compRepo := mocks.NewMockICompanyRepository(ctrl)
	finRepo := mocks.NewMockIFinancialReportRepository(ctrl)
	actFieldRepo := mocks.NewMockIActivityFieldRepository(ctrl)
	svc := NewService(userRepo, compRepo, finRepo, actFieldRepo)

	testCases := []struct {
		name       string
		user       domain.User
		beforeTest func(userRepo mocks.MockIUserRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное обновление",
			user: domain.User{
				ID:   uuid.UUID{1},
				City: "a",
			},
			beforeTest: func(userRepo mocks.MockIUserRepository) {
				userRepo.EXPECT().
					Update(
						context.Background(),
						&domain.User{
							ID:   uuid.UUID{1},
							City: "a",
						},
					).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ошибка выполнения запроса в репозитории",
			user: domain.User{
				ID:   uuid.UUID{1},
				City: "a",
			},
			beforeTest: func(userRepo mocks.MockIUserRepository) {
				userRepo.EXPECT().
					Update(
						context.Background(),
						&domain.User{
							ID:   uuid.UUID{1},
							City: "a",
						},
					).Return(fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("обновление информации о пользователе: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.beforeTest != nil {
				tc.beforeTest(*userRepo)
			}

			err := svc.Update(context.Background(), &tc.user)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}
