package services

import (
	"github.com/google/uuid"
	"ppo/domain"
)

type FinReportService struct {
	finRepo domain.IFinancialReportRepository
}

func (s FinReportService) Create(contact domain.FinancialReport) (err error) {
	return nil
}

func (s FinReportService) GetById(id uuid.UUID) (finReport domain.FinancialReport, err error) {
	return finReport, nil
}

func (s FinReportService) GetByPeriod(period domain.Period) (finReport domain.FinancialReport, err error) {
	return finReport, nil
}

func (s FinReportService) Update(finReport domain.FinancialReport) (err error) {
	return nil
}

func (s FinReportService) DeleteById(id uuid.UUID) (err error) {
	return nil
}
