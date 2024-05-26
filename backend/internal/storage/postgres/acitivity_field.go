package postgres

import (
	"context"
	"fmt"
	"ppo/domain"
	"ppo/internal/config"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (r *ActivityFieldRepository) GetAll(ctx context.Context, page int, isPaginated bool) (fields []*domain.ActivityField, numPages int, err error) {
	query :=
		`select 
    		id, 
    		name,
    		description,
    		cost 
		from ppo.activity_fields`

	var rows pgx.Rows
	if !isPaginated {
		rows, err = r.db.Query(
			ctx,
			query,
		)
	} else {
		rows, err = r.db.Query(
			ctx,
			query+` offset $1 limit $2`,
			(page-1)*config.PageSize,
			config.PageSize,
		)
	}
	if err != nil {
		return nil, 0, fmt.Errorf("получение сфер деятельности: %w", err)
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
			return nil, 0, fmt.Errorf("сканирование полученных строк: %w", err)
		}

		fields = append(fields, tmp)
	}

	var numRecords int
	err = r.db.QueryRow(
		ctx,
		`select count(*) from ppo.activity_fields`,
	).Scan(&numRecords)
	if err != nil {
		return nil, 0, fmt.Errorf("получение числа сфер деятельности: %w", err)
	}

	numPages = numRecords / config.PageSize
	if numRecords%config.PageSize != 0 {
		numPages++
	}

	return fields, numPages, nil
}
