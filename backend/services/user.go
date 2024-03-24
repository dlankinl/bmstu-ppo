package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"math"
	"ppo/domain"
	"time"
)

type UserService struct {
	userRepo    domain.IUserRepository
	companyRepo domain.ICompanyRepository
	finRepo     domain.IFinancialReportRepository
}

//func (s UserService) Create(user *domain.User) (err error) {
//	if user.Gender != "m" && user.Gender != "w" {
//		return fmt.Errorf("неизвестный пол")
//	}
//
//	if user.Username == "" {
//		return fmt.Errorf("должно быть указано имя пользователя")
//	}
//
//	if user.City == "" {
//		return fmt.Errorf("должно быть указано название города")
//	}
//
//	if user.Birthday.IsZero() {
//		return fmt.Errorf("должна быть указана дата рождения")
//	}
//
//	if user.FullName == "" {
//		return fmt.Errorf("должны быть указаны ФИО")
//	}
//
//	if len(strings.Split(user.FullName, " ")) != 3 {
//		return fmt.Errorf("некорректное количество слов (должны быть фамилия, имя и отчество)")
//	}
//
//	err = s.userRepo.Create(user)
//	if err != nil {
//		return fmt.Errorf("создание пользователя: %w", err)
//	}
//
//	return nil
//}

func (s UserService) GetById(ctx context.Context, id uuid.UUID) (user *domain.User, err error) {
	user, err = s.userRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("получение пользователя по id: %w", err)
	}

	return user, nil
}

// TODO: фильтрация
func (s UserService) GetAll(ctx context.Context) (users []*domain.User, err error) {
	users, err = s.userRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("получение списка всех пользователей: %w", err)
	}

	return users, nil
}

func (s UserService) Update(ctx context.Context, user *domain.User) (err error) {
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("обновление информации о пользователе: %w", err)
	}

	return nil
}

func (s UserService) DeleteById(ctx context.Context, id uuid.UUID) (err error) {
	err = s.userRepo.DeleteById(ctx, id)
	if err != nil {
		return fmt.Errorf("удаление пользователя по id: %w", err)
	}

	return nil
}

func (s UserService) GetFinancialReport(ctx context.Context, id uuid.UUID, period domain.Period) (finReports []*domain.FinancialReportByPeriod, err error) {
	if period.StartYear > period.EndYear ||
		(period.StartYear == period.EndYear && period.StartQuarter > period.EndQuarter) {
		return nil, fmt.Errorf("дата конца периода должна быть позже даты начала")
	}

	companies, err := s.companyRepo.GetByOwnerId(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("получение списка компаний предпринимателя с id=%d: %w", id, err)
	}

	finReports = make([]*domain.FinancialReportByPeriod, 0)
	for _, company := range companies {
		report, err := s.finRepo.GetByCompany(ctx, company.ID, period)
		if err != nil {
			return nil, fmt.Errorf("получение финансовой отчетности компании с id=%d: %w", company.ID, err)
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
				totalFinReport.Reports = append(totalFinReport.Reports, report.Reports[i])
				i++
			}

			per := domain.Period{
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

				report.Taxes += v.Taxes
				report.TaxLoad += v.Taxes / v.Revenue() * 100
			}
		}

		finReports = append(finReports, report)
	}

	return finReports, nil
}

func (s UserService) CalculateRating(ctx context.Context, id uuid.UUID) (rating float32, err error) {
	var mainFieldWeight float32 // TODO: mainFieldWeight

	prevYear := time.Now().AddDate(-1, 0, 0).Year()
	period := domain.Period{
		StartYear:    prevYear,
		EndYear:      prevYear,
		StartQuarter: 1,
		EndQuarter:   4,
	}

	reports, err := s.GetFinancialReport(ctx, id, period)
	if err != nil {
		return 0, fmt.Errorf("получение финансового отчета за прошлый год: %w", err)
	}

	var totalRevenue, totalProfit float32
	for _, rep := range reports {
		totalRevenue += rep.Revenue()
		totalProfit += rep.Profit()
	}

	rating = 1.2*mainFieldWeight*mainFieldWeight + 0.35*totalRevenue + 0.9*float32(math.Pow(float64(totalProfit), 1.5))

	return rating, nil
}
