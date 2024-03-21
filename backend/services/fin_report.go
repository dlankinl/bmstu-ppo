package services

import (
	"fmt"
	"github.com/google/uuid"
	"ppo/domain"
)

type FinReportService struct {
	finRepo domain.IFinancialReportRepository
}

func (s FinReportService) Create(finReport *domain.FinancialReport) (err error) {
	err = s.finRepo.Create(finReport)
	if err != nil {
		return fmt.Errorf("добавление финансового отчета: %w", err)
	}

	return nil
}

func (s FinReportService) GetById(id uuid.UUID) (finReport *domain.FinancialReport, err error) {
	finReport, err = s.finRepo.GetById(id)
	if err != nil {
		return nil, fmt.Errorf("получение финансового отчета по id: %w", err)
	}

	return finReport, nil
}

func (s FinReportService) GetByPeriod(period domain.Period) (finReport *domain.FinancialReport, err error) {
	return finReport, nil
}

func (s FinReportService) Update(finReport *domain.FinancialReport) (err error) {
	err = s.finRepo.Update(finReport)
	if err != nil {
		return fmt.Errorf("обновление отчета с id=%d: %w", finReport.ID, err)
	}

	return nil
}

func (s FinReportService) DeleteById(id uuid.UUID) (err error) {
	err = s.finRepo.DeleteById(id)
	if err != nil {
		return fmt.Errorf("удаление отчета: %w", err)
	}

	return nil
}
