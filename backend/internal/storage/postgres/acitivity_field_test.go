package postgres

import (
	"context"
	"github.com/google/uuid"
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

//func TestActivityFieldRepository_DeleteById(t *testing.T) {
//	type fields struct {
//		db *pgxpool.Pool
//	}
//	type args struct {
//		ctx context.Context
//		id  uuid.UUID
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			r := &ActivityFieldRepository{
//				db: tt.fields.db,
//			}
//			if err := r.DeleteById(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
//				t.Errorf("DeleteById() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}

//func TestActivityFieldRepository_GetByCompanyId(t *testing.T) {
//	repo := NewActivityFieldRepository(testDbInstance)
//
//	testCases := []struct {
//		name     string
//		id       uuid.UUID
//		expected float32
//		wantErr  bool
//		errStr   error
//	}{
//		{
//			name: "успех",
//			id:   uuid.UUID{1},
//			expected: 0.1,
//			wantErr: false,
//		},
//	}
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			res, err := repo.GetByCompanyId(context.Background(), tc.id)
//
//			if tc.wantErr {
//				require.Equal(t, tc.errStr.Error(), err.Error())
//			} else {
//				require.Nil(t, err)
//				require.InEpsilon(t, tc.expected, res, 1e-7)
//			}
//		})
//	}
//}

func TestActivityFieldRepository_GetById(t *testing.T) {
	repo := NewActivityFieldRepository(testDbInstance)

	testCases := []struct {
		name     string
		uuidStr  string
		expected *domain.ActivityField
		wantErr  bool
		errStr   error
	}{
		{
			name:    "успех",
			uuidStr: "f80426b8-27e7-4bfa-8721-23075f125165",
			expected: &domain.ActivityField{
				Name:        "field1",
				Description: "field1_descr",
				Cost:        0.1,
			},
			wantErr: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uuidV, _ := uuid.Parse(tc.uuidStr)

			res, err := repo.GetById(context.Background(), uuidV)

			if tc.wantErr {
				require.Equal(t, tc.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				//require.InEpsilon(t, tc.expected, res, 1e-7)
				require.Equal(t, tc.expected.Name, res.Name)
				require.Equal(t, tc.expected.Description, res.Description)
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

//func TestActivityFieldRepository_Update(t *testing.T) {
//	type fields struct {
//		db *pgxpool.Pool
//	}
//	type args struct {
//		ctx  context.Context
//		data *domain.ActivityField
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			r := &ActivityFieldRepository{
//				db: tt.fields.db,
//			}
//			if err := r.Update(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
//				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestNewActivityFieldRepository(t *testing.T) {
//	type args struct {
//		db *pgxpool.Pool
//	}
//	tests := []struct {
//		name string
//		args args
//		want domain.IActivityFieldRepository
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := NewActivityFieldRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("NewActivityFieldRepository() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
