package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"ppo/domain"
)

type ActivityFieldRepository struct {
	db *pgx.Conn
}

func (r *ActivityFieldRepository) Create(context.Context, *domain.ActivityField) (err error) {
	return nil
}

func (r *ActivityFieldRepository) DeleteById(context.Context, uuid.UUID) (err error) {
	return nil
}

func (r *ActivityFieldRepository) Update(context.Context, *domain.ActivityField) (err error) {
	return nil
}

func (r *ActivityFieldRepository) GetById(context.Context, uuid.UUID) (field *domain.ActivityField, err error) {
	return nil, nil
}

func (r *ActivityFieldRepository) GetByCompanyId(context.Context, uuid.UUID) (cost float32, err error) {
	return 0, nil
}

func (r *ActivityFieldRepository) GetMaxCost(context.Context) (cost float32, err error) {
	return 0, nil
}
