package domain

import "github.com/google/uuid"

type FinancialReport struct {
	ID        uuid.UUID
	CompanyID uuid.UUID
	Revenue   float32
	Costs     float32
	Taxes     float32
	TaxLoad   float32
	Year      int
	Quarter   int
}

type FinancialReportByPeriod struct {
	Period  Period
	Reports []FinancialReport
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
	GetByCompany(companyId uuid.UUID, period Period) (*FinancialReportByPeriod, error)
	Update(finRep *FinancialReport) error
	DeleteById(id uuid.UUID) error
}

type IFinancialReportService interface {
	Create(finRep *FinancialReport) error
	CreateByPeriod(finReportByPeriod *FinancialReportByPeriod) error
	GetById(id uuid.UUID) (*FinancialReport, error)
	GetByCompany(companyId uuid.UUID, period Period) (*FinancialReportByPeriod, error)
	Update(finRep *FinancialReport) error
	DeleteById(id uuid.UUID) error
}
