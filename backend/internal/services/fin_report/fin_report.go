package fin_report

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"ppo/domain"
	"ppo/pkg/logger"
	"time"
)

type Service struct {
	finRepo domain.IFinancialReportRepository
	logger  logger.ILogger
}

func NewService(finRepo domain.IFinancialReportRepository, logger logger.ILogger) domain.IFinancialReportService {
	return &Service{
		finRepo: finRepo,
		logger:  logger,
	}
}

func (s *Service) Create(ctx context.Context, finReport *domain.FinancialReport) (err error) {
	prompt := "FinReportCreate"

	if finReport.Revenue < 0 {
		s.logger.Infof("%s: выручка не может быть отрицательной", prompt)
		return fmt.Errorf("выручка не может быть отрицательной")
	}

	if finReport.Costs < 0 {
		s.logger.Infof("%s: расходы не могут быть отрицательными", prompt)
		return fmt.Errorf("расходы не могут быть отрицательными")
	}

	if finReport.Quarter > 4 || finReport.Quarter < 1 {
		s.logger.Infof("%s: значение квартала должно находиться в отрезке от 1 до 4", prompt)
		return fmt.Errorf("значение квартала должно находиться в отрезке от 1 до 4")
	}

	now := time.Now()
	if finReport.Year > now.Year() {
		s.logger.Infof("%s: значение года не может быть больше текущего года", prompt)
		return fmt.Errorf("значение года не может быть больше текущего года")
	}

	if finReport.Year == now.Year() && finReport.Quarter > (int(now.Month()-1)/3) {
		s.logger.Infof("%s: нельзя добавить отчет за квартал, который еще не закончился", prompt)
		return fmt.Errorf("нельзя добавить отчет за квартал, который еще не закончился")
	}

	err = s.finRepo.Create(ctx, finReport)
	if err != nil {
		s.logger.Infof("%s: добавление финансового отчета: %v", prompt, err)
		return fmt.Errorf("добавление финансового отчета: %w", err)
	}

	return nil
}

func (s *Service) CreateByPeriod(ctx context.Context, finReportByPeriod *domain.FinancialReportByPeriod) (err error) {
	prompt := "FinReportCreateByPeriod"

	for _, report := range finReportByPeriod.Reports {
		err = s.Create(ctx, &report)
		if err != nil {
			s.logger.Infof("%s: добавление отчетов за период: %v", prompt, err)
			return fmt.Errorf("добавление отчетов за период: %w", err)
		}
	}

	return nil
}

func (s *Service) GetById(ctx context.Context, id uuid.UUID) (finReport *domain.FinancialReport, err error) {
	prompt := "FinReportGetById"

	finReport, err = s.finRepo.GetById(ctx, id)
	if err != nil {
		s.logger.Infof("%s: получение финансового отчета по id: %v", prompt, err)
		return nil, fmt.Errorf("получение финансового отчета по id: %w", err)
	}

	return finReport, nil
}

func (s *Service) GetByCompany(ctx context.Context, companyId uuid.UUID, period *domain.Period) (
	finReport *domain.FinancialReportByPeriod, err error) {
	prompt := "FinReportGetByCompany"

	if period.StartYear > period.EndYear ||
		(period.StartYear == period.EndYear && period.StartQuarter > period.EndQuarter) {
		s.logger.Infof("%s: дата конца периода должна быть позже даты начала", prompt)
		return nil, fmt.Errorf("дата конца периода должна быть позже даты начала")
	}

	finReport, err = s.finRepo.GetByCompany(ctx, companyId, period)
	if err != nil {
		s.logger.Infof("%s: получение финансового отчета по id компании: %v", prompt, err)
		return nil, fmt.Errorf("получение финансового отчета по id компании: %w", err)
	}

	return finReport, nil
}

func (s *Service) Update(ctx context.Context, finReport *domain.FinancialReport) (err error) {
	prompt := "FinReportUpdate"

	err = s.finRepo.Update(ctx, finReport)
	if err != nil {
		s.logger.Infof("%s: обновление отчета: %v", prompt, err)
		return fmt.Errorf("обновление отчета: %w", err)
	}

	return nil
}

func (s *Service) DeleteById(ctx context.Context, id uuid.UUID) (err error) {
	prompt := "FinReportDeleteById"

	err = s.finRepo.DeleteById(ctx, id)
	if err != nil {
		s.logger.Infof("%s: удаление отчета по id: %v", prompt, err)
		return fmt.Errorf("удаление отчета по id: %w", err)
	}

	return nil
}
