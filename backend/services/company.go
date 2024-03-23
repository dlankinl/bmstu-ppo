package services

import (
	"fmt"
	"github.com/google/uuid"
	"ppo/domain"
)

type CompanyService struct {
	companyRepo domain.ICompanyRepository
	finRepo     domain.IFinancialReportRepository
}

func NewCompanyService(companyRepo domain.ICompanyRepository, finRepo domain.IFinancialReportRepository) *CompanyService {
	return &CompanyService{
		companyRepo: companyRepo,
		finRepo:     finRepo,
	}
}

func (s CompanyService) Create(company *domain.Company) (err error) {
	if company.Name == "" {
		return fmt.Errorf("должно быть указано название компании")
	}

	if company.City == "" {
		return fmt.Errorf("должно быть указано название города")
	}

	if company.FieldOfActivity == "" {
		return fmt.Errorf("должно быть указано название сферы деятельности")
	}

	err = s.companyRepo.Create(company)
	if err != nil {
		return fmt.Errorf("добавление компании: %w", err)
	}

	return nil
}

func (s CompanyService) GetById(id uuid.UUID) (company *domain.Company, err error) {
	company, err = s.companyRepo.GetById(id)
	if err != nil {
		return nil, fmt.Errorf("получение компании по id: %w", err)
	}

	return company, nil
}

func (s CompanyService) GetByOwnerId(id uuid.UUID) (companies []*domain.Company, err error) {
	companies, err = s.companyRepo.GetByOwnerId(id)
	if err != nil {
		return nil, fmt.Errorf("получение списка компаний по id владельца: %w", err)
	}

	return companies, nil
}

// TODO: фильтрация
func (s CompanyService) GetAll() (companies []*domain.Company, err error) {
	companies, err = s.companyRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("получение списка всех компаний: %w", err)
	}

	return companies, nil
}

func (s CompanyService) Update(company *domain.Company) (err error) {
	err = s.companyRepo.Update(company)
	if err != nil {
		return fmt.Errorf("обновление компании с id=%d: %w", company.ID, err)
	}

	return nil
}

func (s CompanyService) DeleteById(id uuid.UUID) (err error) {
	err = s.companyRepo.DeleteById(id)
	if err != nil {
		return fmt.Errorf("удаление компании по id: %w", err)
	}

	return nil
}
