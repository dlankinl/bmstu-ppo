package user_activity_field

import (
	"context"
	"fmt"
	"math"
	"ppo/domain"
	"ppo/pkg/logger"
	"time"

	"github.com/google/uuid"
)

const (
	quartersInYear = 4
	firstQuarter   = 1
	lastQuarter    = 4
)

type Interactor struct {
	userService     domain.IUserService
	actFieldService domain.IActivityFieldService
	compService     domain.ICompanyService
	finService      domain.IFinancialReportService
	logger          logger.ILogger
}

func NewInteractor(
	userSvc domain.IUserService,
	actFieldSvc domain.IActivityFieldService,
	compSvc domain.ICompanyService,
	finSvc domain.IFinancialReportService,
	logger logger.ILogger,
) *Interactor {
	return &Interactor{
		userService:     userSvc,
		actFieldService: actFieldSvc,
		compService:     compSvc,
		finService:      finSvc,
		logger:          logger,
	}
}

type taxesData struct {
	taxes   float32
	revenue float32
}

func calculateTaxes(reports map[int]*domain.FinancialReportByPeriod) (taxes *taxesData) {
	taxes = new(taxesData)

	for _, v := range reports {
		if len(v.Reports) == quartersInYear {
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

			taxes.taxes += v.Taxes
			taxes.revenue += v.Revenue()
		}
	}

	return taxes
}

func findFullYearReports(rep *domain.FinancialReportByPeriod, period *domain.Period) (fullYearReports map[int]*domain.FinancialReportByPeriod) {
	fullYearReports = make(map[int]*domain.FinancialReportByPeriod)

	var j int
	for year := period.StartYear; year <= period.EndYear; year++ {
		startQtr := firstQuarter
		endQtr := lastQuarter

		if year == period.StartYear {
			startQtr = period.StartQuarter
		}
		if year == period.EndYear {
			endQtr = period.EndQuarter
		}

		var totalFinReport domain.FinancialReportByPeriod

		// цикл нужен для аккумулирования всех отчётов за выбранный период; переменная j нужна для контроля невозможности
		// вылезти за границы слайса, т.к. за год могут быть в наличии отчёты за 1, 3 и 4 квартал и, в таком случае,
		// если итерироваться по quarter, будет печалька ;(
		for quarter := startQtr; quarter <= endQtr; quarter++ {
			if j < len(rep.Reports) {
				totalFinReport.Reports = append(totalFinReport.Reports, rep.Reports[j])
				j++
			}
		}

		if endQtr-startQtr == quartersInYear-1 {
			per := &domain.Period{
				StartYear:    year,
				EndYear:      year,
				StartQuarter: startQtr,
				EndQuarter:   endQtr,
			}

			totalFinReport.Period = per
			fullYearReports[year] = &totalFinReport
		}
	}

	return fullYearReports
}

func calcRating(profit, revenue, cost, maxCost float32) float32 {
	return (cost/maxCost + profit/revenue) / 2.0
}

func (i *Interactor) GetMostProfitableCompany(ctx context.Context, period *domain.Period, companies []*domain.Company) (company *domain.Company, err error) {
	var maxProfit float32

	for _, comp := range companies {
		rep, err := i.finService.GetByCompany(ctx, comp.ID, period)
		if err != nil {
			return nil, fmt.Errorf("получение отчета компании: %w", err)
		}

		if rep.Profit() > maxProfit {
			company = comp
			maxProfit = rep.Profit()
		}
	}

	return company, nil
}

func (i *Interactor) CalculateUserRating(ctx context.Context, id uuid.UUID) (rating float32, err error) {
	prompt := "UserActivityFieldCalculateUserRating"

	companies, _, err := i.compService.GetByOwnerId(ctx, id, 0, false)
	if err != nil {
		i.logger.Infof("%s: получение списка компаний: %v", prompt, err)
		return 0, fmt.Errorf("получение списка компаний: %w", err)
	}

	prevYear := time.Now().AddDate(-1, 0, 0).Year()
	period := &domain.Period{
		StartYear:    prevYear,
		EndYear:      prevYear,
		StartQuarter: firstQuarter,
		EndQuarter:   lastQuarter,
	}

	report, err := i.GetUserFinancialReport(ctx, id, period)
	if err != nil {
		i.logger.Infof("%s: получение финансового отчета пользователя: %v", prompt, err)
		return 0, fmt.Errorf("получение финансового отчета пользователя: %w", err)
	}

	mostProfitableCompany, err := i.GetMostProfitableCompany(ctx, period, companies)
	if err != nil {
		i.logger.Infof("%s: поиск наиболее прибыльной компании: %v", prompt, err)
		return 0, fmt.Errorf("поиск наиболее прибыльной компании: %w", err)
	}
	if mostProfitableCompany == nil {
		return 0.0, nil
	}

	maxCost, err := i.actFieldService.GetMaxCost(ctx)
	if err != nil {
		i.logger.Infof("%s: поиск максимального веса: %v", prompt, err)
		return 0, fmt.Errorf("поиск максимального веса: %w", err)
	}

	cost, err := i.actFieldService.GetCostByCompanyId(ctx, mostProfitableCompany.ID)
	if err != nil {
		i.logger.Infof("%s: получение веса сферы деятельности компании: %v", prompt, err)
		return 0, fmt.Errorf("получение веса сферы деятельности компании: %w", err)
	}

	var totalRevenue, totalProfit float32
	totalRevenue = report.Revenue()
	totalProfit = report.Profit()

	rating = calcRating(totalProfit, totalRevenue, cost, maxCost)

	return rating, nil
}

func (i *Interactor) GetUserFinancialReport(ctx context.Context, id uuid.UUID, period *domain.Period) (report *domain.FinancialReportByPeriod, err error) {
	prompt := "UserActivityFieldGetUserFinancialReport"
	report = new(domain.FinancialReportByPeriod)

	companies, _, err := i.compService.GetByOwnerId(ctx, id, 0, false)
	if err != nil {
		i.logger.Infof("%s: получение списка компаний: %v", prompt, err)
		return nil, fmt.Errorf("получение списка компаний: %w", err)
	}

	var revenueForTaxLoad float32
	report.Reports = make([]domain.FinancialReport, 0)
	for _, comp := range companies {
		rep, err := i.finService.GetByCompany(ctx, comp.ID, period)
		if err != nil {
			i.logger.Infof("%s: получение отчета компании: %v", prompt, err)
			return nil, fmt.Errorf("получение отчета компании: %w", err)
		}

		fullYears := findFullYearReports(rep, period)

		tax := calculateTaxes(fullYears)
		report.Taxes += tax.taxes
		revenueForTaxLoad += tax.revenue

		report.Reports = append(report.Reports, rep.Reports...)
	}

	report.Period = period
	if math.Abs(float64(revenueForTaxLoad)) >= 1e-6 {
		report.TaxLoad = report.Taxes / revenueForTaxLoad * 100
	}

	return report, nil
}
