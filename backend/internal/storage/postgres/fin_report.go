package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"ppo/domain"
)

type FinReportRepository struct {
	db *pgx.ConnPool
}

func NewFinReportRepository(db *pgx.ConnPool) domain.IFinancialReportRepository {
	return &FinReportRepository{
		db: db,
	}
}

func (r *FinReportRepository) Create(ctx context.Context, finReport *domain.FinancialReport) (err error) {
	query := `insert into ppo.fin_reports(company_id, revenue, costs, year, quarter) 
	values ($1, $2, $3, $4, $5)`

	_, err = r.db.ExecEx(
		ctx,
		query,
		nil,
		finReport.CompanyID,
		finReport.Revenue,
		finReport.Costs,
		finReport.Year,
		finReport.Quarter,
	)
	if err != nil {
		return fmt.Errorf("создание финансового отчета: %w", err)
	}

	return nil
}

func (r *FinReportRepository) GetById(ctx context.Context, id uuid.UUID) (report *domain.FinancialReport, err error) {
	query := `select company_id, revenue, costs, year, quarter from ppo.fin_reports where id = $1`

	report = new(domain.FinancialReport)
	err = r.db.QueryRowEx(
		ctx,
		query,
		nil,
		id,
	).Scan(
		&report.CompanyID,
		&report.Revenue,
		&report.Costs,
		&report.Year,
		&report.Quarter,
	)
	if err != nil {
		return nil, fmt.Errorf("получение отчета по id: %w", err)
	}

	return report, nil
}

func (r *FinReportRepository) GetByCompany(ctx context.Context, companyId uuid.UUID, period *domain.Period) (report *domain.FinancialReportByPeriod, err error) {
	query := `select id, company_id, revenue, costs, year, quarter
	from ppo.fin_reports 
	where company_id = $1 and year = $2 and quarter = $3`

	report = new(domain.FinancialReportByPeriod)
	report.Reports = make([]domain.FinancialReport, 0)

	for year := period.StartYear; year <= period.EndYear; year++ {
		startQtr := 1
		endQtr := 4

		if year == period.StartYear {
			startQtr = period.StartQuarter
		}
		if year == period.EndYear {
			endQtr = period.EndQuarter
		}

		for quarter := startQtr; quarter <= endQtr; quarter++ {
			tmp := new(domain.FinancialReport)

			err = r.db.QueryRowEx(
				ctx,
				query,
				nil,
				companyId,
				year,
				quarter,
			).Scan(
				&tmp.ID,
				&tmp.CompanyID,
				&tmp.Revenue,
				&tmp.Costs,
				&tmp.Year,
				&tmp.Quarter,
			)

			if errors.Is(err, sql.ErrNoRows) {
				continue
			} else if err != nil {
				return nil, fmt.Errorf("cканирование записи: %w", err)
			}

			report.Reports = append(report.Reports, *tmp)
		}
	}

	report.Period = period

	return report, nil
}

func (r *FinReportRepository) Update(ctx context.Context, finRep *domain.FinancialReport) (err error) {
	query := `
			update ppo.fin_reports
			set 
			    company_id = $1, 
			    revenue = $2,
			    costs = $3,
			    year = $4,
			    quarter = $5
			where id = $6`

	_, err = r.db.ExecEx(
		ctx,
		query,
		nil,
		finRep.CompanyID,
		finRep.Revenue,
		finRep.Costs,
		finRep.Year,
		finRep.Quarter,
	)
	if err != nil {
		return fmt.Errorf("обновление информации о финансовом отчете: %w", err)
	}

	return nil
}

func (r *FinReportRepository) DeleteById(ctx context.Context, id uuid.UUID) (err error) {
	query := `delete from ppo.fin_reports where id = $1`

	_, err = r.db.ExecEx(
		ctx,
		query,
		nil,
		id,
	)
	if err != nil {
		return fmt.Errorf("удаление отчета по id: %w", err)
	}

	return nil
}
