package postgres

import (
	"context"
	"github.com/stretchr/testify/require"
	"ppo/domain"
	"testing"
)

func TestActivityFieldRepository_Create(t *testing.T) {
	repo := NewActivityFieldRepository(testDbInstance)

	testCases := []struct {
		name    string
		field   *domain.ActivityField
		wantErr bool
		errStr  error
	}{
		{
			name: "успех",
			field: &domain.ActivityField{
				Name:        "a",
				Description: "a",
				Cost:        1.0,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.Create(context.Background(), tc.field)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestActivityFieldRepository_GetMaxCost(t *testing.T) {
	repo := NewActivityFieldRepository(testDbInstance)

	testCases := []struct {
		name     string
		expected float32
		wantErr  bool
		errStr   error
	}{
		{
			name:     "успех",
			expected: 1.3,
			wantErr:  false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := repo.GetMaxCost(context.Background())

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.InEpsilon(t, tc.expected, res, 1e-4)
			}
		})
	}
}
