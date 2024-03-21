package domain

import "github.com/google/uuid"

type FinancialReport struct {
	ID      uuid.UUID
	Revenue float32
	Costs   float32
	period  Period
}

type Period struct {
	StartYear    int
	StartQuarter int
	EndYear      int
	EndQuarter   int
}

type IFinancialReportRepository interface {
	Create(finRep *FinancialReport) error
	GetById(id uuid.UUID) (*FinancialReport, error)
	GetByPeriod(companyId uuid.UUID, period Period) (*FinancialReport, error) // FIXME: бред
	Update(finRep *FinancialReport) error
	DeleteById(id uuid.UUID) error
}

type IFinancialReportService interface {
	Create(finRep *FinancialReport) error
	GetById(id uuid.UUID) (*FinancialReport, error)
	GetByPeriod(companyId uuid.UUID, period Period) (*FinancialReport, error) // FIXME: бред
	Update(finRep *FinancialReport) error
	DeleteById(id uuid.UUID) error
}
