package domain

import (
	"context"
	"github.com/google/uuid"
)

type FinancialReport struct {
	ID        uuid.UUID
	CompanyID uuid.UUID
	Revenue   float32
	Costs     float32
	Year      int
	Quarter   int
}

type FinancialReportByPeriod struct {
	Reports []FinancialReport
	Period  *Period
	Taxes   float32
	TaxLoad float32
}

type Period struct {
	StartYear    int
	StartQuarter int
	EndYear      int
	EndQuarter   int
}

func (r *FinancialReportByPeriod) Revenue() (sum float32) {
	for _, rep := range r.Reports {
		sum += rep.Revenue
	}

	return sum
}

func (r *FinancialReportByPeriod) Costs() (sum float32) {
	for _, rep := range r.Reports {
		sum += rep.Costs
	}

	return sum
}

func (r *FinancialReportByPeriod) Profit() (sum float32) {
	for _, rep := range r.Reports {
		sum += rep.Revenue - rep.Costs
	}

	return sum
}

type IFinancialReportRepository interface {
	Create(ctx context.Context, finRep *FinancialReport) error
	GetById(ctx context.Context, id uuid.UUID) (*FinancialReport, error)
	GetByCompany(ctx context.Context, companyId uuid.UUID, period *Period) (*FinancialReportByPeriod, error)
	Update(ctx context.Context, finRep *FinancialReport) error
	DeleteById(ctx context.Context, id uuid.UUID) error
}

type IFinancialReportService interface {
	Create(ctx context.Context, finRep *FinancialReport) error
	CreateByPeriod(ctx context.Context, finReportByPeriod *FinancialReportByPeriod) error
	GetById(ctx context.Context, id uuid.UUID) (*FinancialReport, error)
	GetByCompany(ctx context.Context, companyId uuid.UUID, period *Period) (*FinancialReportByPeriod, error)
	Update(ctx context.Context, finRep *FinancialReport) error
	DeleteById(ctx context.Context, id uuid.UUID) error
}
