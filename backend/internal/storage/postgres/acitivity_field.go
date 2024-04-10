package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"ppo/domain"
)

type ActivityFieldRepository struct {
	db *pgxpool.Pool
}

func NewActivityFieldRepository(db *pgxpool.Pool) domain.IActivityFieldRepository {
	return &ActivityFieldRepository{
		db: db,
	}
}

func (r *ActivityFieldRepository) Create(ctx context.Context, data *domain.ActivityField) (err error) {
	query := `insert into ppo.activity_fields(name, description, cost) 
	values ($1, $2, $3)`

	_, err = r.db.Exec(
		ctx,
		query,
		data.Name,
		data.Description,
		data.Cost,
	)
	if err != nil {
		return fmt.Errorf("создание сферы деятельности: %w", err)
	}

	return nil
}

func (r *ActivityFieldRepository) DeleteById(ctx context.Context, id uuid.UUID) (err error) {
	query := `delete from ppo.activity_fields where id = $1`

	_, err = r.db.Exec(
		ctx,
		query,
		id,
	)
	if err != nil {
		return fmt.Errorf("удаление сферы деятельности по id: %w", err)
	}

	return nil
}

func (r *ActivityFieldRepository) Update(ctx context.Context, data *domain.ActivityField) (err error) {
	query := `
			update ppo.activity_fields
			set 
			    name = $1,
			    description = $2, 
			    cost = $3
			where id = $4`

	_, err = r.db.Exec(
		ctx,
		query,
		data.Name,
		data.Description,
		data.Cost,
		data.ID,
	)
	if err != nil {
		return fmt.Errorf("обновление информации о сфере деятельности: %w", err)
	}

	return nil
}

func (r *ActivityFieldRepository) GetById(ctx context.Context, id uuid.UUID) (field *domain.ActivityField, err error) {
	query := `select name, description, cost from ppo.activity_fields where id = $1`

	field = new(domain.ActivityField)
	err = r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&field.Name,
		&field.Description,
		&field.Cost,
	)
	if err != nil {
		return nil, fmt.Errorf("получение сферы деятельности по id: %w", err)
	}

	return nil, nil
}

func (r *ActivityFieldRepository) GetMaxCost(ctx context.Context) (cost float32, err error) {
	query := `select max(cost)
		from ppo.activity_fields`

	//var maxVal float32
	err = r.db.QueryRow(
		ctx,
		query,
	).Scan(&cost)

	if err != nil {
		return 0, fmt.Errorf("получение максимального веса сферы деятельности: %w", err)
	}

	return cost, nil
}
