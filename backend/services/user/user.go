package user

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"ppo/domain"
	"ppo/pkg/utils"
	"strings"
)

type Service struct {
	userRepo     domain.IUserRepository
	companyRepo  domain.ICompanyRepository
	finRepo      domain.IFinancialReportRepository
	actFieldRepo domain.IActivityFieldRepository
}

func NewService(
	userRepo domain.IUserRepository,
	companyRepo domain.ICompanyRepository,
	finRepo domain.IFinancialReportRepository,
	actFieldRepo domain.IActivityFieldRepository,
) domain.IUserService {
	return &Service{
		userRepo:     userRepo,
		companyRepo:  companyRepo,
		finRepo:      finRepo,
		actFieldRepo: actFieldRepo,
	}
}

func (s *Service) Create(user *domain.User) (err error) {
	if user.Gender != "m" && user.Gender != "w" {
		return fmt.Errorf("неизвестный пол")
	}

	if user.City == "" {
		return fmt.Errorf("должно быть указано название города")
	}

	if user.Birthday.IsZero() {
		return fmt.Errorf("должна быть указана дата рождения")
	}

	if user.FullName == "" {
		return fmt.Errorf("должны быть указаны ФИО")
	}

	if len(strings.Split(user.FullName, " ")) != 3 {
		return fmt.Errorf("некорректное количество слов (должны быть фамилия, имя и отчество)")
	}

	var ctx context.Context

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("создание пользователя: %w", err)
	}

	return nil
}

func (s *Service) GetById(id uuid.UUID) (user *domain.User, err error) {
	var ctx context.Context

	user, err = s.userRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("получение пользователя по id: %w", err)
	}

	return user, nil
}

// TODO: pagination
func (s *Service) GetAll(filters utils.Filters) (users []*domain.User, err error) {
	var ctx context.Context

	users, err = s.userRepo.GetAll(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("получение списка всех пользователей: %w", err)
	}

	return users, nil
}

func (s *Service) Update(user *domain.User) (err error) {
	var ctx context.Context

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("обновление информации о пользователе: %w", err)
	}

	return nil
}

func (s *Service) DeleteById(id uuid.UUID) (err error) {
	var ctx context.Context

	err = s.userRepo.DeleteById(ctx, id)
	if err != nil {
		return fmt.Errorf("удаление пользователя по id: %w", err)
	}

	return nil
}

func (s *Service) GetFinancialReport(companies []*domain.Company, period *domain.Period) (finReports []*domain.FinancialReportByPeriod, err error) {
	var ctx context.Context

	if period.StartYear > period.EndYear ||
		(period.StartYear == period.EndYear && period.StartQuarter > period.EndQuarter) {
		return nil, fmt.Errorf("дата конца периода должна быть позже даты начала")
	}

	finReports = make([]*domain.FinancialReportByPeriod, 0)
	for _, company := range companies {
		report, err := s.finRepo.GetByCompany(ctx, company.ID, period)
		if err != nil {
			return nil, fmt.Errorf("получение финансовой отчетности компании по id: %w", err)
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

				report.Taxes += v.Taxes
				report.TaxLoad += v.Taxes / v.Revenue() * 100
			}
		}

		finReports = append(finReports, report)
	}

	return finReports, nil
}
