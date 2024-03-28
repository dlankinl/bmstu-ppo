package user_activity_field

import (
	"fmt"
	"github.com/google/uuid"
	"ppo/domain"
	"time"
)

const (
	quartersInYear = 4
	firstQuarter   = 1
	lastQuarter    = 4
)
const coef = 1e-9

type Interactor struct {
	userService     domain.IUserService
	actFieldService domain.IActivityFieldService
	compService     domain.ICompanyService
	finService      domain.IFinancialReportService
}

func NewInteractor(
	userSvc domain.IUserService,
	actFieldSvc domain.IActivityFieldService,
	compSvc domain.ICompanyService,
	finSvc domain.IFinancialReportService,
) *Interactor {
	return &Interactor{
		userService:     userSvc,
		actFieldService: actFieldSvc,
		compService:     compSvc,
		finService:      finSvc,
	}
}

type taxes struct {
	Sum  float32
	Load float32
}

func calculateTaxes(reports map[int]domain.FinancialReportByPeriod) (tax *taxes) {
	tax = new(taxes)

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

			tax.Sum += v.Taxes
			tax.Load += v.Taxes / v.Revenue() * 100
		}
	}

	return tax
}

func findFullYearReports(rep *domain.FinancialReportByPeriod, period *domain.Period) (fullYearReports map[int]domain.FinancialReportByPeriod) {
	fullYearReports = make(map[int]domain.FinancialReportByPeriod)

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
		for quarter := startQtr; quarter <= endQtr; quarter++ {
			totalFinReport.Reports = append(totalFinReport.Reports, rep.Reports[j])
			j++
		}

		if endQtr-startQtr == quartersInYear-1 {
			per := &domain.Period{
				StartYear:    year,
				EndYear:      year,
				StartQuarter: startQtr,
				EndQuarter:   endQtr,
			}

			totalFinReport.Period = per
			fullYearReports[year] = totalFinReport
		}
	}

	return fullYearReports
}

func calcRating(profit, revenue, cost, maxCost float32) float32 {
	return (cost/maxCost + profit/revenue) / 2.0
}

func (i *Interactor) GetMostProfitableCompany(period *domain.Period, companies []*domain.Company) (company *domain.Company, err error) {
	var maxProfit float32

	for _, comp := range companies {
		rep, err := i.finService.GetByCompany(comp.ID, period)
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

func (i *Interactor) CalculateUserRating(id uuid.UUID) (rating float32, err error) {
	companies, err := i.compService.GetByOwnerId(id)
	if err != nil {
		return 0, fmt.Errorf("получение списка компаний: %w", err)
	}

	prevYear := time.Now().AddDate(-1, 0, 0).Year()
	period := &domain.Period{
		StartYear:    prevYear,
		EndYear:      prevYear,
		StartQuarter: firstQuarter,
		EndQuarter:   lastQuarter,
	}

	report, err := i.GetUserFinancialReport(id, period)
	if err != nil {
		return 0, fmt.Errorf("получение финансового отчета пользователя: %w", err)
	}

	mostProfitableCompany, err := i.GetMostProfitableCompany(period, companies)
	if err != nil {
		return 0, fmt.Errorf("поиск наиболее прибыльной компании: %w", err)
	}

	maxCost, err := i.actFieldService.GetMaxCost()

	cost, err := i.actFieldService.GetCostByCompanyId(mostProfitableCompany.ID)
	if err != nil {
		return 0, fmt.Errorf("получение веса сферы деятельности компании: %w", err)
	}

	var totalRevenue, totalProfit float32
	totalRevenue = report.Revenue()
	totalProfit = report.Profit()

	rating = calcRating(totalProfit, totalRevenue, cost, maxCost)

	return rating, nil
}

func (i *Interactor) GetUserFinancialReport(id uuid.UUID, period *domain.Period) (report *domain.FinancialReportByPeriod, err error) {
	report = new(domain.FinancialReportByPeriod)

	companies, err := i.compService.GetByOwnerId(id)
	if err != nil {
		return nil, fmt.Errorf("получение списка компаний: %w", err)
	}

	var profit, revenue float32
	report.Reports = make([]domain.FinancialReport, 0)
	for _, comp := range companies {
		rep, err := i.finService.GetByCompany(comp.ID, period)
		if err != nil {
			return nil, fmt.Errorf("получение отчета компании: %w", err)
		}

		profit += rep.Profit()
		revenue += rep.Revenue()

		report.Reports = append(report.Reports, rep.Reports...)
	}

	report.Period = period
	return report, nil
}
