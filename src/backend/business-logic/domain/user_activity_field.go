package domain

import "github.com/google/uuid"

type IInteractor interface {
	GetMostProfitableCompany(period *Period, companies []*Company) (*Company, error)
	CalculateUserRating(id uuid.UUID) (float32, error)
	GetUserFinancialReport(id uuid.UUID, period *Period) (*FinancialReportByPeriod, error)
}
