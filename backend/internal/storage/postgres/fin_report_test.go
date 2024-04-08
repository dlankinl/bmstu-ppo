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

func TestFinReportRepository_GetByCompany(t *testing.T) {
	finRepo := NewFinReportRepository(testDbInstance)

	testCases := []struct {
		name      string
		companyId uuid.UUID
		period    *domain.Period
		expected  *domain.FinancialReportByPeriod
		wantErr   bool
		errStr    error
	}{
		{
			name:      "успех",
			companyId: uuid.UUID{2},
			period: &domain.Period{
				StartYear:    1,
				StartQuarter: 1,
				EndYear:      1,
				EndQuarter:   3,
			},
			expected: &domain.FinancialReportByPeriod{
				Reports: []domain.FinancialReport{
					{
						ID:        uuid.UUID{2},
						CompanyID: uuid.UUID{2},
						Revenue:   1.0,
						Costs:     0.5,
						Year:      1,
						Quarter:   1,
					},
					{
						ID:        uuid.UUID{3},
						CompanyID: uuid.UUID{2},
						Revenue:   1.0,
						Costs:     0.5,
						Year:      1,
						Quarter:   2,
					},
					{
						ID:        uuid.UUID{4},
						CompanyID: uuid.UUID{2},
						Revenue:   1.0,
						Costs:     0.5,
						Year:      1,
						Quarter:   3,
					},
				},
				Period: &domain.Period{
					StartYear:    1,
					EndYear:      1,
					StartQuarter: 1,
					EndQuarter:   3,
				},
			},
			wantErr: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			report, err := finRepo.GetByCompany(context.Background(), tc.companyId, tc.period)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.Equal(t, tc.expected, report)
			}
		})
	}
}

func TestFinReportRepository_GetById(t *testing.T) {
	finRepo := NewFinReportRepository(testDbInstance)

	testCases := []struct {
		name     string
		id       uuid.UUID
		expected *domain.FinancialReport
		wantErr  bool
		errStr   error
	}{
		{
			name: "успех",
			id:   uuid.UUID{2},
			expected: &domain.FinancialReport{
				ID:        uuid.UUID{2},
				CompanyID: uuid.UUID{2},
				Revenue:   1.0,
				Costs:     0.5,
				Year:      1,
				Quarter:   1,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			report, err := finRepo.GetById(context.Background(), tc.id)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Equal(t, tc.expected, report)
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
				ID: uuid.UUID{},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

		})
	}
}
