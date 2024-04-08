package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"ppo/domain"
)

type ActivityFieldRepository struct {
	db *pgx.ConnPool
}

func NewActivityField(db *pgx.ConnPool) domain.IActivityFieldRepository {
	return &ActivityFieldRepository{
		db: db,
	}
}

func (r *ActivityFieldRepository) Create(ctx context.Context, data *domain.ActivityField) (err error) {
	query := `insert into ppo.activity_fields(name, description, cost) 
	values ($1, $2, $3)`

	_, err = r.db.ExecEx(
		ctx,
		query,
		nil,
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

	_, err = r.db.ExecEx(
		ctx,
		query,
		nil,
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

	_, err = r.db.ExecEx(
		ctx,
		query,
		nil,
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
	err = r.db.QueryRowEx(
		ctx,
		query,
		nil,
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

// TODO: FIX
func (r *ActivityFieldRepository) GetByCompanyId(ctx context.Context, id uuid.UUID) (cost float32, err error) {
	// query := `select id, name, description, cost from ppo.activity_fields where owner_id = $1 `

	// rows, err := r.db.QueryEx(
	// 	ctx,
	// 	query,
	// 	nil,
	// 	id,
	// )

	// contacts = make([]*domain.Contact, 0)
	// for rows.Next() {
	// 	tmp := new(domain.Contact)

	// 	err = rows.Scan(
	// 		&tmp.ID,
	// 		&tmp.Name,
	// 		&tmp.Value,
	// 	)
	// 	tmp.OwnerID = id

	// 	if err != nil {
	// 		return nil, fmt.Errorf("сканирование полученных строк: %w", err)
	// 	}
	// }

	return 0, nil
}

func (r *ActivityFieldRepository) GetMaxCost(ctx context.Context) (cost float32, err error) {
	query := `select max(cost)
		from ppo.activity_fields`

	var maxVal float32
	err = r.db.QueryRowEx(
		ctx,
		query,
		nil,
	).Scan(&maxVal)

	if err != nil {
		return 0, fmt.Errorf("получение максимального веса сферы деятельности: %w", err)
	}

	return 0, nil
}
