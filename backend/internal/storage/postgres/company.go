package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"ppo/domain"
)

type CompanyRepository struct {
	db *pgxpool.Pool
}

func NewCompanyRepository(db *pgxpool.Pool) domain.ICompanyRepository {
	return &CompanyRepository{
		db: db,
	}
}

func (r *CompanyRepository) Create(ctx context.Context, company *domain.Company) (err error) {
	query := `insert into ppo.companies(owner_id, activity_field_id, name, city) 
	values ($1, $2, $3, $4)`

	_, err = r.db.Exec(
		ctx,
		query,
		company.OwnerID,
		company.ActivityFieldId,
		company.Name,
		company.City,
	)
	if err != nil {
		return fmt.Errorf("создание компании: %w", err)
	}

	return nil
}

func (r *CompanyRepository) GetById(ctx context.Context, id uuid.UUID) (company *domain.Company, err error) {
	query := `select owner_id, activity_field_id, name, city from ppo.companies where id = $1`

	company = new(domain.Company)
	err = r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&company.OwnerID,
		&company.ActivityFieldId,
		&company.Name,
		&company.City,
	)
	if err != nil {
		return nil, fmt.Errorf("получение компании по id: %w", err)
	}

	return company, nil
}

func (r *CompanyRepository) GetByOwnerId(ctx context.Context, id uuid.UUID) (companies []*domain.Company, err error) {
	query := `select id, activity_field_id, name, city from ppo.companies where owner_id = $1 `

	rows, err := r.db.Query(
		ctx,
		query,
		id,
	)

	companies = make([]*domain.Company, 0)
	for rows.Next() {
		tmp := new(domain.Company)

		err = rows.Scan(
			&tmp.ID,
			&tmp.ActivityFieldId,
			&tmp.Name,
			&tmp.City,
		)
		tmp.OwnerID = id

		if err != nil {
			return nil, fmt.Errorf("сканирование полученных строк: %w", err)
		}
	}

	return companies, nil
}

func (r *CompanyRepository) Update(ctx context.Context, company *domain.Company) (err error) {
	query := `
			update ppo.companies
			set 
			    owner_id = $1,
			    activity_field_id = $2,
			    name = $3, 
			    city = $4
			where id = $5`

	_, err = r.db.Exec(
		ctx,
		query,
		company.OwnerID,
		company.ActivityFieldId,
		company.Name,
		company.City,
		company.ID,
	)
	if err != nil {
		return fmt.Errorf("обновление информации о компании: %w", err)
	}

	return nil
}

func (r *CompanyRepository) DeleteById(ctx context.Context, id uuid.UUID) (err error) {
	query := `delete from ppo.companies where id = $1`

	_, err = r.db.Exec(
		ctx,
		query,
		id,
	)
	if err != nil {
		return fmt.Errorf("удаление компании по id: %w", err)
	}

	return nil
}

func (r *CompanyRepository) GetAll(ctx context.Context, page int) (companies []*domain.Company, err error) {
	query := `select id, owner_id, activity_field_id, name, city from ppo.companies offset $1 limit $2`

	rows, err := r.db.Query(
		ctx,
		query,
		(page-1)*pageSize,
		pageSize,
	)

	companies = make([]*domain.Company, 0)
	for rows.Next() {
		tmp := new(domain.Company)

		err = rows.Scan(
			&tmp.ID,
			&tmp.OwnerID,
			&tmp.ActivityFieldId,
			&tmp.Name,
			&tmp.City,
		)

		if err != nil {
			return nil, fmt.Errorf("сканирование полученных строк: %w", err)
		}
	}

	return companies, nil
}
