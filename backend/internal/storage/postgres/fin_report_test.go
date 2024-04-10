package postgres

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"ppo/domain"
	"testing"
)

func TestFinReportRepository_Create(t *testing.T) {
	finRepo := NewFinReportRepository(testDbInstance)

	testCases := []struct {
		name    string
		report  *domain.FinancialReport
		wantErr bool
		errStr  error
	}{
		{
			name: "успех",
			report: &domain.FinancialReport{
				CompanyID: uuid.UUID{1},
				Revenue:   1.32,
				Costs:     1.23,
				Year:      2024,
				Quarter:   1,
			},
			wantErr: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := finRepo.Create(context.Background(), tc.report)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestFinReportRepository_DeleteById(t *testing.T) {
	finRepo := NewFinReportRepository(testDbInstance)

	testCases := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
		errStr  error
	}{
		{
			name:    "успех",
			id:      uuid.UUID{1},
			wantErr: false,
		},
		{
			name:    "несуществующий id",
			id:      uuid.UUID{25},
			wantErr: false,
			errStr:  errors.New(""),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := finRepo.DeleteById(context.Background(), tc.id)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestFinReportRepository_Update(t *testing.T) {
	finRepo := NewFinReportRepository(testDbInstance)

	testCases := []struct {
		name    string
		report  *domain.FinancialReport
		wantErr bool
		errStr  error
	}{
		{
			name: "успех",
			report: &domain.FinancialReport{
				ID:      uuid.UUID{1},
				Revenue: 2.0,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := finRepo.Update(context.Background(), tc.report)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}
