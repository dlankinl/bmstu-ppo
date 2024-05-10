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
	Create(context.Context, *FinancialReport) error
	GetById(context.Context, uuid.UUID) (*FinancialReport, error)
	GetByCompany(context.Context, uuid.UUID, *Period) (*FinancialReportByPeriod, error)
	Update(context.Context, *FinancialReport) error
	DeleteById(context.Context, uuid.UUID) error
}

type IFinancialReportService interface {
	Create(context.Context, *FinancialReport) error
	CreateByPeriod(context.Context, *FinancialReportByPeriod) error
	GetById(context.Context, uuid.UUID) (*FinancialReport, error)
	GetByCompany(context.Context, uuid.UUID, *Period) (*FinancialReportByPeriod, error)
	Update(context.Context, *FinancialReport) error
	DeleteById(context.Context, uuid.UUID) error
}
