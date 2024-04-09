package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"ppo/domain"
	"reflect"
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

func TestActivityFieldRepository_DeleteById(t *testing.T) {
	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ActivityFieldRepository{
				db: tt.fields.db,
			}
			if err := r.DeleteById(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteById() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestActivityFieldRepository_GetByCompanyId(t *testing.T) {
	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantCost float32
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ActivityFieldRepository{
				db: tt.fields.db,
			}
			gotCost, err := r.GetByCompanyId(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByCompanyId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCost != tt.wantCost {
				t.Errorf("GetByCompanyId() gotCost = %v, want %v", gotCost, tt.wantCost)
			}
		})
	}
}

func TestActivityFieldRepository_GetById(t *testing.T) {
	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantField *domain.ActivityField
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ActivityFieldRepository{
				db: tt.fields.db,
			}
			gotField, err := r.GetById(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotField, tt.wantField) {
				t.Errorf("GetById() gotField = %v, want %v", gotField, tt.wantField)
			}
		})
	}
}

func TestActivityFieldRepository_GetMaxCost(t *testing.T) {
	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantCost float32
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ActivityFieldRepository{
				db: tt.fields.db,
			}
			gotCost, err := r.GetMaxCost(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMaxCost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCost != tt.wantCost {
				t.Errorf("GetMaxCost() gotCost = %v, want %v", gotCost, tt.wantCost)
			}
		})
	}
}

func TestActivityFieldRepository_Update(t *testing.T) {
	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		ctx  context.Context
		data *domain.ActivityField
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ActivityFieldRepository{
				db: tt.fields.db,
			}
			if err := r.Update(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewActivityFieldRepository(t *testing.T) {
	type args struct {
		db *pgxpool.Pool
	}
	tests := []struct {
		name string
		args args
		want domain.IActivityFieldRepository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewActivityFieldRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewActivityFieldRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}
