package services

import (
	"fmt"
	"github.com/google/uuid"
	"math"
	"ppo/domain"
	"strings"
	"time"
)

type UserService struct {
	userRepo    domain.IUserRepository
	companyRepo domain.ICompanyRepository
	finRepo     domain.IFinancialReportRepository
}

func (s UserService) Create(user *domain.User) (err error) {
	if user.Gender != "m" || user.Gender != "w" {
		return fmt.Errorf("неизвестный пол: %w", err)
	}

	if user.Username == "" {
		return fmt.Errorf("должно быть указано имя пользователя: %w", err)
	}

	if user.City == "" {
		return fmt.Errorf("должно быть указано название города: %w", err)
	}

	if user.Birthday.IsZero() {
		return fmt.Errorf("должна быть указана дата рождения: %w", err)
	}

	if user.FullName == "" {
		return fmt.Errorf("должны быть указаны ФИО: %w", err)
	}

	if len(strings.Split(user.FullName, " ")) != 3 {
		return fmt.Errorf("некорректное количество слов (должны быть фамилия, имя и отчество): %w", err)
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return fmt.Errorf("создание пользователя: %w", err)
	}

	return nil
}

func (s UserService) GetById(id uuid.UUID) (user *domain.User, err error) {
	user, err = s.userRepo.GetById(id)
	if err != nil {
		return nil, fmt.Errorf("получение пользователя по id: %w", err)
	}

	return user, nil
}

// TODO: фильтрация
func (s UserService) GetAll() (users []*domain.User, err error) {
	users, err = s.userRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("получение списка всех пользователей: %w", err)
	}

	return users, nil
}

func (s UserService) Update(user *domain.User) (err error) {
	err = s.userRepo.Update(user)
	if err != nil {
		return fmt.Errorf("обновление информации о пользователе: %w", err)
	}

	return nil
}

func (s UserService) DeleteById(id uuid.UUID) (err error) {
	err = s.userRepo.DeleteById(id)
	if err != nil {
		return fmt.Errorf("удаление пользователя по id: %w", err)
	}

	return nil
}

func (s UserService) GetFinancialReport(id uuid.UUID, period domain.Period) (finReport *domain.ExtFinancialReport, err error) {
	companies, err := s.companyRepo.GetByOwnerId(id)
	if err != nil {
		return nil, fmt.Errorf("получение списка компаний предпринимателя с id=%d: %w", id, err)
	}

	finReports := make([]domain.FinancialReport, 0)
	for _, company := range companies {
		report, err := s.finRepo.GetByPeriod(company.ID, period)
		if err != nil {
			return nil, fmt.Errorf("получение финансовой отчетности компании с id=%d: %w", company.ID, err)
		}

		finReports = append(finReports, *report)
	}

	yearReports := make(map[int]domain.ExtFinancialReport)

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

		var totalFinReport domain.FinancialReport
		for quarter := startQtr; quarter <= endQtr; quarter++ {
			totalFinReport.Revenue += finReports[i].Revenue
			totalFinReport.Costs += finReports[i].Costs
			i++
		}

		per := domain.Period{
			StartYear:    year,
			EndYear:      year,
			StartQuarter: startQtr,
			EndQuarter:   endQtr,
		}
		totalFinReport.Period = per
		yearReports[year] = domain.FinReportToExt(totalFinReport)
	}

	for _, v := range yearReports {
		totalProfit := v.FinReport.Revenue - v.FinReport.Costs
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

		finReport.FinReport.Revenue += v.FinReport.Revenue
		finReport.FinReport.Costs += v.FinReport.Costs
		finReport.Taxes += v.Taxes
	}

	return finReport, nil
}

func (s UserService) CalculateRating(id uuid.UUID) (rating float32, err error) {
	var mainFieldWeight float32 // TODO: mainFieldWeight

	prevYear := time.Now().AddDate(-1, 0, 0).Year()
	period := domain.Period{
		StartYear:    prevYear,
		EndYear:      prevYear,
		StartQuarter: 1,
		EndQuarter:   4,
	}

	report, err := s.GetFinancialReport(id, period)
	if err != nil {
		return 0, fmt.Errorf("получение финансового отчета за прошлый год: %w", err)
	}

	rating = 1.2*mainFieldWeight*mainFieldWeight + 0.35*report.FinReport.Revenue + 0.9*float32(math.Pow(float64(report.FinReport.Revenue-report.FinReport.Costs), 1.5))

	return rating, nil
}
