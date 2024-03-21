package services

import (
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

func (s UserService) Create(user *domain.User) (err error) {
	return nil
}

func (s UserService) GetById(id uuid.UUID) (user *domain.User, err error) {
	return user, nil
}

func (s UserService) GetAll() (users []*domain.User, err error) {
	return users, nil
}

func (s UserService) Update(user *domain.User) (err error) {
	return nil
}

func (s UserService) DeleteById(id uuid.UUID) (err error) {
	return nil
}

func (s UserService) GetUserCompanies(id uuid.UUID) (companies []*domain.Company, err error) {
	return companies, nil
}

func (s UserService) GetFinancialReport(period domain.Period) (finReport *domain.FinancialReport, err error) {
	return finReport, nil
}

func (s UserService) CalculateRating(id uuid.UUID) (rating float32, err error) {
	var totalProfit, totalRevenue, mainFieldWeight float32 // TODO: mainFieldWeight

	prevYear := time.Now().AddDate(-1, 0, 0).Year()
	period := domain.Period{
		StartYear:    prevYear,
		EndYear:      prevYear,
		StartQuarter: 1,
		EndQuarter:   4,
	}

	companies, err := s.companyRepo.GetByOwnerId(id)
	if err != nil {
		return 0, fmt.Errorf("получение списка компаний предпринимателя с id=%d: %w", id, err)
	}

	finReports := make([]domain.FinancialReport, 0)
	for _, company := range companies {
		report, err := s.finRepo.GetByPeriod(company.ID, period)
		if err != nil {
			return 0, fmt.Errorf("получение финансовой отчетности компании с id=%d: %w", company.ID, err)
		}

		finReports = append(finReports, *report)
	}

	for _, report := range finReports {
		totalRevenue += report.Revenue
		totalProfit += report.Revenue - report.Costs
	}

	rating = 1.2*mainFieldWeight*mainFieldWeight + 0.35*totalRevenue + 0.9*float32(math.Pow(float64(totalProfit), 1.5))

	return rating, nil
}
