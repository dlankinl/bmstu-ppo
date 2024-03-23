package domain

import "github.com/google/uuid"

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
	Period  Period
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
