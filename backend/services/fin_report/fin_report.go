package fin_report

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"ppo/domain"
	"time"
)

type Service struct {
	finRepo domain.IFinancialReportRepository
}

func NewService(finRepo domain.IFinancialReportRepository) domain.IFinancialReportService {
	return &Service{
		finRepo: finRepo,
	}
}

func (s *Service) Create(finReport *domain.FinancialReport) (err error) {
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

	ctx := context.Background()

	err = s.finRepo.Create(ctx, finReport)
	if err != nil {
		return fmt.Errorf("добавление финансового отчета: %w", err)
	}

	return nil
}

func (s *Service) CreateByPeriod(finReportByPeriod *domain.FinancialReportByPeriod) (err error) {
	for _, report := range finReportByPeriod.Reports {
		err = s.Create(&report)
		if err != nil {
			return fmt.Errorf("добавление отчетов за период: %w", err)
		}
	}

	return nil
}

func (s *Service) GetById(id uuid.UUID) (finReport *domain.FinancialReport, err error) {
	ctx := context.Background()

	finReport, err = s.finRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("получение финансового отчета по id: %w", err)
	}

	return finReport, nil
}

func (s *Service) GetByCompany(companyId uuid.UUID, period *domain.Period) (
	finReport *domain.FinancialReportByPeriod, err error) {
	if period.StartYear > period.EndYear ||
		(period.StartYear == period.EndYear && period.StartQuarter > period.EndQuarter) {
		return nil, fmt.Errorf("дата конца периода должна быть позже даты начала")
	}

	ctx := context.Background()

	finReport, err = s.finRepo.GetByCompany(ctx, companyId, period)
	if err != nil {
		return nil, fmt.Errorf("получение финансового отчета по id компании: %w", err)
	}

	yearReports := make(map[int]domain.FinancialReportByPeriod)

	var i int
	for year := period.StartYear; year <= period.EndYear; year++ {
		startQtr := 1
		endQtr := 4

		if year == period.StartYear {
			startQtr = period.StartQuarter
		}
		if year == period.EndYear {
			endQtr = period.EndQuarter
		}

		var totalFinReport domain.FinancialReportByPeriod
		for quarter := startQtr; quarter <= endQtr; quarter++ {
			totalFinReport.Reports = append(totalFinReport.Reports, finReport.Reports[i])
			i++
		}

		per := &domain.Period{
			StartYear:    year,
			EndYear:      year,
			StartQuarter: startQtr,
			EndQuarter:   endQtr,
		}
		totalFinReport.Period = per
		yearReports[year] = totalFinReport
	}

	for _, v := range yearReports {
		if len(v.Reports) == 4 {
			totalProfit := v.Profit()
			var taxFare int
			switch true {
			case totalProfit < 10000000:
				taxFare = 4
			case totalProfit < 50000000:
				taxFare = 7
			case totalProfit < 150000000:
				taxFare = 13
			case totalProfit < 500000000:
				taxFare = 20
			default:
				taxFare = 30
			}

			v.Taxes = totalProfit * (float32(taxFare) / 100)

			finReport.Taxes += v.Taxes
			finReport.TaxLoad += v.Taxes / v.Revenue() * 100
		}
	}

	return finReport, nil
}

func (s *Service) Update(finReport *domain.FinancialReport) (err error) {
	ctx := context.Background()

	err = s.finRepo.Update(ctx, finReport)
	if err != nil {
		return fmt.Errorf("обновление отчета: %w", err)
	}

	return nil
}

func (s *Service) DeleteById(id uuid.UUID) (err error) {
	ctx := context.Background()

	err = s.finRepo.DeleteById(ctx, id)
	if err != nil {
		return fmt.Errorf("удаление отчета по id: %w", err)
	}

	return nil
}
