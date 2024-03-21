package services

import (
	"github.com/google/uuid"
	"ppo/domain"
)

type UserService struct {
	userRepo    domain.IUserRepository
	companyRepo domain.ICompanyRepository
	finRepo     domain.IFinancialReportRepository
}

func (s UserService) Create(user domain.User) (err error) {
	return nil
}

func (s UserService) GetById(id uuid.UUID) (user domain.User, err error) {
	return user, nil
}

func (s UserService) GetAll() (users []domain.User, err error) {
	return users, nil
}

func (s UserService) Update(user domain.User) (err error) {
	return nil
}

func (s UserService) DeleteById(id uuid.UUID) (err error) {
	return nil
}

func (s UserService) GetUserCompanies(id uuid.UUID) (companies []domain.Company, err error) {
	return companies, nil
}

func (s UserService) GetFinancialReport(period domain.Period) (finReport domain.FinancialReport, err error) {
	return finReport, nil
}

func (s UserService) CalculateRating(id uuid.UUID) (rating float32, err error) {
	return 0, nil
}
