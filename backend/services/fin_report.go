package services

import (
	"fmt"
	"github.com/google/uuid"
	"ppo/domain"
	"time"
)

type FinReportService struct {
	finRepo domain.IFinancialReportRepository
}

func NewFinReportService(finRepo domain.IFinancialReportRepository) *FinReportService {
	return &FinReportService{
		finRepo: finRepo,
	}
}

func (s FinReportService) Create(finReport *domain.FinancialReport) (err error) {
	if finReport.Revenue < 0 {
		return fmt.Errorf("выручка не может быть отрицательной")
	}

	if finReport.Costs < 0 {
		return fmt.Errorf("расходы не могут быть отрицательными")
	}

	if finReport.Quarter > 4 || finReport.Quarter < 1 {
		return fmt.Errorf("значение квартала должно находиться в отрезке от 1 до 4")
	}

	now := time.Now()
	if finReport.Year > now.Year() {
		return fmt.Errorf("значение года не может быть больше текущего года")
	}

	if finReport.Year == now.Year() && finReport.Quarter > (int(now.Month()-1)/3) { // TODO: проверить юнит-тестами
		return fmt.Errorf("нельзя добавить отчет за квартал, который еще не закончился")
	}

	err = s.finRepo.Create(finReport)
	if err != nil {
		return fmt.Errorf("добавление финансового отчета: %w", err)
	}

	return nil
}

func (s FinReportService) CreateByPeriod(finReportByPeriod *domain.FinancialReportByPeriod) (err error) {
	for _, report := range finReportByPeriod.Reports {
		err = s.Create(&report)
		if err != nil {
			return fmt.Errorf("добавление отчетов за период: %w", err)
		}
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

func (s FinReportService) GetByCompany(companyId uuid.UUID, period domain.Period) (
	finReport *domain.FinancialReportByPeriod, err error) {
	if period.StartYear > period.EndYear ||
		(period.StartYear == period.EndYear && period.StartQuarter > period.EndQuarter) {
		return nil, fmt.Errorf("дата конца периода должна быть позже даты начала")
	}

	finReport, err = s.finRepo.GetByCompany(companyId, period)
	if err != nil {
		return nil, fmt.Errorf("получение финансового отчета по id компании: %w", err)
	}

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
