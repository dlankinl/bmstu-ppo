package user_activity_field

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"math"
	"ppo/domain"
	"time"
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

func (i *Interactor) CalculateUserRating(ctx context.Context, id uuid.UUID) (rating float32, err error) {
	companies, err := i.compService.GetByOwnerId(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("получение списка компаний: %w", err)
	}

	var maxProfit float32
	var mostProfitableCompId uuid.UUID
	prevYear := time.Now().AddDate(-1, 0, 0).Year()
	period := &domain.Period{
		StartYear:    prevYear,
		EndYear:      prevYear,
		StartQuarter: 1,
		EndQuarter:   4,
	}

	reports := make([]*domain.FinancialReportByPeriod, 0)
	for _, comp := range companies {
		rep, err := i.finService.GetByCompany(ctx, comp.ID, period)
		if err != nil {
			return 0, fmt.Errorf("получение отчета компании: %w", err)
		}

		if rep.Profit() > maxProfit {
			mostProfitableCompId = comp.ID
			maxProfit = rep.Profit()
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
				totalFinReport.Reports = append(totalFinReport.Reports, rep.Reports[i])
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

				rep.Taxes += v.Taxes
				rep.TaxLoad += v.Taxes / v.Revenue() * 100
			}
		}

		reports = append(reports, rep)
	}

	cost, err := i.actFieldService.GetCostByCompanyId(ctx, mostProfitableCompId)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	var totalRevenue, totalProfit float32
	for _, rep := range reports {
		fmt.Println(rep.Revenue(), rep.Profit(), rep.Period)
		totalRevenue += rep.Revenue()
		totalProfit += rep.Profit()
	}
	fmt.Println(reports)

	fmt.Println(totalProfit, totalRevenue, totalProfit-totalRevenue)
	fmt.Println(32532513-5436438+6743634-9876967+4675424-2436653+14385253-7546424+3253251-543643+6743634-9876967+4675412-2436765+1438525-754642, 32532513+6743634+4675424+14385253+3253251+6743634+4675412+1438525, 5436438+9876967+2436653+7546424+543643+9876967+2436765+754642)
	rating = (1.2*cost*cost + 0.35*totalRevenue + 0.9*float32(math.Pow(float64(totalProfit), 1.5))) * coef

	return rating, nil
}
