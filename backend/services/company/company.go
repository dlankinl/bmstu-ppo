package company

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"ppo/domain"
	"ppo/pkg/utils"
)

type Service struct {
	companyRepo domain.ICompanyRepository
	finRepo     domain.IFinancialReportRepository
}

func NewService(companyRepo domain.ICompanyRepository, finRepo domain.IFinancialReportRepository) domain.ICompanyService {
	return &Service{
		companyRepo: companyRepo,
		finRepo:     finRepo,
	}
}

func (s *Service) Create(company *domain.Company) (err error) {
	if company.Name == "" {
		return fmt.Errorf("должно быть указано название компании")
	}

	if company.City == "" {
		return fmt.Errorf("должно быть указано название города")
	}

	ctx := context.Background()

	err = s.companyRepo.Create(ctx, company)
	if err != nil {
		return fmt.Errorf("добавление компании: %w", err)
	}

	return nil
}

func (s *Service) GetById(id uuid.UUID) (company *domain.Company, err error) {
	ctx := context.Background()

	company, err = s.companyRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("получение компании по id: %w", err)
	}

	return company, nil
}

func (s *Service) GetByOwnerId(id uuid.UUID) (companies []*domain.Company, err error) {
	ctx := context.Background()

	companies, err = s.companyRepo.GetByOwnerId(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("получение списка компаний по id владельца: %w", err)
	}

	return companies, nil
}

// TODO: pagination
func (s *Service) GetAll(filters utils.Filters) (companies []*domain.Company, err error) {
	ctx := context.Background()

	companies, err = s.companyRepo.GetAll(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("получение списка всех компаний: %w", err)
	}

	return companies, nil
}

func (s *Service) Update(company *domain.Company) (err error) {
	ctx := context.Background()

	err = s.companyRepo.Update(ctx, company)
	if err != nil {
		return fmt.Errorf("обновление информации о компании: %w", err)
	}

	return nil
}

func (s *Service) DeleteById(id uuid.UUID) (err error) {
	ctx := context.Background()

	err = s.companyRepo.DeleteById(ctx, id)
	if err != nil {
		return fmt.Errorf("удаление компании по id: %w", err)
	}

	return nil
}
