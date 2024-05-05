package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"ppo/domain"
	"ppo/internal/config"
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
	company.ID = id

	return company, nil
}

func (r *CompanyRepository) GetByOwnerId(ctx context.Context, id uuid.UUID, page int) (companies []*domain.Company, err error) {
	query :=
		`select 
    		id, 
    		activity_field_id,
    		name,
    		city 
		from ppo.companies 
		where owner_id = $1`

	var rows pgx.Rows
	if page == 0 {
		rows, err = r.db.Query(
			ctx,
			query,
			id,
		)
	} else {
		rows, err = r.db.Query(
			ctx,
			query+` offset $2 limit $3`,
			id,
			(page-1)*config.PageSize,
			config.PageSize,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("получение компаний: %w", err)
	}

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

		companies = append(companies, tmp)
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
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("открытие транзакции: %w", err)
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				err = fmt.Errorf("обработанная ошибка: %w\nоткат транзакции: %v", err, rollbackErr)
			}
		}
	}()

	_, err = tx.Exec(
		ctx,
		`delete from ppo.companies where id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("удаление компании по id: %w", err)
	}

	_, err = tx.Exec(
		ctx,
		`delete from ppo.fin_reports where company_id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("удаление отчетов, связанных с компанией: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("закрытие транзакции: %w", err)
	}

	return nil
}

func (r *CompanyRepository) GetAll(ctx context.Context, page int) (companies []*domain.Company, err error) {
	query := `select id, owner_id, activity_field_id, name, city from ppo.companies offset $1 limit $2`

	rows, err := r.db.Query(
		ctx,
		query,
		(page-1)*config.PageSize,
		config.PageSize,
	)
	if err != nil {
		return nil, fmt.Errorf("получение списка компаний: %w", err)
	}

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
