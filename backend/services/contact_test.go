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

func TestContactService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	conRepo := mocks.NewMockIContactsRepository(ctrl)
	svc := NewContactService(conRepo)

	testCases := []struct {
		name       string
		contact    domain.Contact
		beforeTest func(conRepo mocks.MockIContactsRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное добавление",
			contact: domain.Contact{
				Name:  "aaa",
				Value: "bbb",
			},
			beforeTest: func(conRepo mocks.MockIContactsRepository) {
				conRepo.EXPECT().
					Create(
						context.Background(),
						&domain.Contact{
							Name:  "aaa",
							Value: "bbb",
						},
					).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "пустое название средства связи",
			contact: domain.Contact{
				Name:  "",
				Value: "bbb",
			},
			beforeTest: func(conRepo mocks.MockIContactsRepository) {
				conRepo.EXPECT().
					Create(
						context.Background(),
						&domain.Contact{
							Name:  "",
							Value: "bbb",
						},
					).Return(nil).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("должно быть указано название средства связи"),
		},
		{
			name: "пустое значение",
			contact: domain.Contact{
				Name:  "aaa",
				Value: "",
			},
			beforeTest: func(conRepo mocks.MockIContactsRepository) {
				conRepo.EXPECT().
					Create(
						context.Background(),
						&domain.Contact{
							Name:  "aaa",
							Value: "",
						},
					).Return(nil).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("должно быть указано значение средства связи"),
		},
		{
			name: "ошибка выполнения запроса в репозитории",
			contact: domain.Contact{
				Name:  "aaa",
				Value: "bbb",
			},
			beforeTest: func(conRepo mocks.MockIContactsRepository) {
				conRepo.EXPECT().
					Create(
						context.Background(),
						&domain.Contact{
							Name:  "aaa",
							Value: "bbb",
						},
					).Return(fmt.Errorf("sql error")).
					AnyTimes()
			},
			wantErr: true,
			errStr:  errors.New("добавление средства связи: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.beforeTest != nil {
				tc.beforeTest(*conRepo)
			}

			err := svc.Create(context.Background(), &tc.contact)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestContactService_DeleteById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	conRepo := mocks.NewMockIContactsRepository(ctrl)
	svc := NewContactService(conRepo)

	curUuid := uuid.New()

	testCases := []struct {
		name       string
		id         uuid.UUID
		beforeTest func(conRepo mocks.MockIContactsRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное удаление",
			id:   curUuid,
			beforeTest: func(conRepo mocks.MockIContactsRepository) {
				conRepo.EXPECT().
					DeleteById(context.Background(), curUuid).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ошибка выполнения запроса в репозитории",
			id:   curUuid,
			beforeTest: func(conRepo mocks.MockIContactsRepository) {
				conRepo.EXPECT().
					DeleteById(context.Background(), curUuid).
					Return(fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("удаление средства связи по id: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.beforeTest != nil {
				tc.beforeTest(*conRepo)
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

func TestContactService_GetByOwnerId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	conRepo := mocks.NewMockIContactsRepository(ctrl)
	svc := NewContactService(conRepo)

	testCases := []struct {
		name       string
		id         uuid.UUID
		beforeTest func(conRepo mocks.MockIContactsRepository)
		expected   []*domain.Contact
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное получение компании по id",
			id:   [16]byte{1},
			beforeTest: func(conRepo mocks.MockIContactsRepository) {
				conRepo.EXPECT().
					GetByOwnerId(
						context.Background(),
						[16]byte{1},
					).
					Return([]*domain.Contact{
						{
							ID:      [16]byte{1},
							OwnerID: [16]byte{1},
							Name:    "a",
							Value:   "a",
						},
						{
							ID:      [16]byte{2},
							OwnerID: [16]byte{1},
							Name:    "b",
							Value:   "b",
						},
						{
							ID:      [16]byte{3},
							OwnerID: [16]byte{1},
							Name:    "c",
							Value:   "c",
						},
					}, nil)
			},
			expected: []*domain.Contact{
				{
					ID:      [16]byte{1},
					OwnerID: [16]byte{1},
					Name:    "a",
					Value:   "a",
				},
				{
					ID:      [16]byte{2},
					OwnerID: [16]byte{1},
					Name:    "b",
					Value:   "b",
				},
				{
					ID:      [16]byte{3},
					OwnerID: [16]byte{1},
					Name:    "c",
					Value:   "c",
				},
			},
			wantErr: false,
		},
		{
			name: "ошибка получения данных в репозитории",
			id:   [16]byte{1},
			beforeTest: func(conRepo mocks.MockIContactsRepository) {
				conRepo.EXPECT().
					GetByOwnerId(
						context.Background(),
						[16]byte{1},
					).
					Return(nil, fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("получение всех средств связи по id владельца: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.beforeTest != nil {
				tc.beforeTest(*conRepo)
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

func TestContactService_GetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	conRepo := mocks.NewMockIContactsRepository(ctrl)
	svc := NewContactService(conRepo)

	testCases := []struct {
		name       string
		id         uuid.UUID
		beforeTest func(conRepo mocks.MockIContactsRepository)
		expected   *domain.Contact
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное получение компании по id",
			id:   [16]byte{1},
			beforeTest: func(conRepo mocks.MockIContactsRepository) {
				conRepo.EXPECT().
					GetById(
						context.Background(),
						[16]byte{1},
					).
					Return(&domain.Contact{
						ID:    [16]byte{1},
						Name:  "a",
						Value: "a",
					}, nil)
			},
			expected: &domain.Contact{
				ID:    [16]byte{1},
				Name:  "a",
				Value: "a",
			},
			wantErr: false,
		},
		{
			name: "ошибка получения данных в репозитории",
			id:   [16]byte{1},
			beforeTest: func(conRepo mocks.MockIContactsRepository) {
				conRepo.EXPECT().
					GetById(
						context.Background(),
						[16]byte{1},
					).
					Return(nil, fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("получение средства связи по id: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.beforeTest != nil {
				tc.beforeTest(*conRepo)
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

func TestContactService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	conRepo := mocks.NewMockIContactsRepository(ctrl)
	svc := NewContactService(conRepo)

	testCases := []struct {
		name       string
		company    domain.Contact
		beforeTest func(conRepo mocks.MockIContactsRepository)
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное обновление",
			company: domain.Contact{
				ID:   [16]byte{1},
				Name: "aaa",
			},
			beforeTest: func(conRepo mocks.MockIContactsRepository) {
				conRepo.EXPECT().
					Update(
						context.Background(),
						&domain.Contact{
							ID:   [16]byte{1},
							Name: "aaa",
						},
					).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ошибка выполнения запроса в репозитории",
			company: domain.Contact{
				ID:   [16]byte{1},
				Name: "aaa",
			},
			beforeTest: func(conRepo mocks.MockIContactsRepository) {
				conRepo.EXPECT().
					Update(
						context.Background(),
						&domain.Contact{
							ID:   [16]byte{1},
							Name: "aaa",
						},
					).Return(fmt.Errorf("sql error"))
			},
			wantErr: true,
			errStr:  errors.New("обновление информации о средстве связи: sql error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.beforeTest != nil {
				tc.beforeTest(*conRepo)
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
