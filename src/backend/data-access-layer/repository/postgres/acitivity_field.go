package postgres

import (
	"context"
	"data-access-layer/config"
	"data-access-layer/domain"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
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

	fmt.Println("HERE", data.ID, data.Name, data.Description, data.Cost)
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

	field.ID = id

	return field, nil
}

func (r *ActivityFieldRepository) GetMaxCost(ctx context.Context) (cost float32, err error) {
	query := `select max(cost)
		from ppo.activity_fields`

	err = r.db.QueryRow(
		ctx,
		query,
	).Scan(&cost)

	if err != nil {
		return 0, fmt.Errorf("получение максимального веса сферы деятельности: %w", err)
	}

	return cost, nil
}

func (r *ActivityFieldRepository) GetAll(ctx context.Context, page int) (fields []*domain.ActivityField, err error) {
	query := `select id, name, description, cost from ppo.activity_fields offset $1 limit $2`

	rows, err := r.db.Query(
		ctx,
		query,
		(page-1)*config.PageSize,
		config.PageSize,
	)
	if err != nil {
		return nil, fmt.Errorf("получение сфер деятельности: %w", err)
	}

	fields = make([]*domain.ActivityField, 0)
	for rows.Next() {
		tmp := new(domain.ActivityField)

		err = rows.Scan(
			&tmp.ID,
			&tmp.Name,
			&tmp.Description,
			&tmp.Cost,
		)

		if err != nil {
			return nil, fmt.Errorf("сканирование полученных строк: %w", err)
		}

		fields = append(fields, tmp)
	}

	return fields, nil
}
