package utils

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFilters_ParseToSql(t *testing.T) {
	testCases := []struct {
		name     string
		f        Filters
		expected string
		wantErr  bool
		errStr   error
	}{
		{
			name: "успех",
			f: map[string]string{
				"age":    "-23",
				"gender": "m",
				"price":  "100-200",
				"length": "23-",
			},
			expected: "age < 23 and gender = 'm' and price between 100 and 200 and length > 23",
			wantErr:  false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.f.ParseToSql()

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.Equal(t, tc.expected, res)
			}
		})
	}
}

func TestFilters_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		f       Filters
		wantErr bool
		errStr  error
	}{
		{
			name: "успешная валидация",
			f: map[string]string{
				"age":    "-23",
				"gender": "m",
				"price":  "100-200",
				"length": "23-",
			},
			wantErr: false,
		},
		{
			name: "sql-инъекция 1",
			f: map[string]string{
				"age":    "-23; where 1=1",
				"gender": "m",
				"price":  "100-200",
				"length": "23-",
			},
			wantErr: true,
			errStr:  errors.New("sql-инъекция"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.f.Validate()

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}
