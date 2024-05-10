package domain

import (
	"context"
	"github.com/google/uuid"
)

type IInteractor interface {
	GetMostProfitableCompany(context.Context, *Period, []*Company) (*Company, error)
	CalculateUserRating(context.Context, uuid.UUID) (float32, error)
	GetUserFinancialReport(context.Context, uuid.UUID, *Period) (*FinancialReportByPeriod, error)
}
