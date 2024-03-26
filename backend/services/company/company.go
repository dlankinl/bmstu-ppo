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

func (s *Service) Create(ctx context.Context, company *domain.Company) (err error) {
	if company.Name == "" {
		return fmt.Errorf("должно быть указано название компании")
	}

	if company.City == "" {
		return fmt.Errorf("должно быть указано название города")
	}

	if company.FieldOfActivity == "" {
		return fmt.Errorf("должно быть указано название сферы деятельности")
	}

	err = s.companyRepo.Create(ctx, company)
	if err != nil {
		return fmt.Errorf("добавление компании: %w", err)
	}

	return nil
}

func (s *Service) GetById(ctx context.Context, id uuid.UUID) (company *domain.Company, err error) {
	company, err = s.companyRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("получение компании по id: %w", err)
	}

	return company, nil
}

func (s *Service) GetByOwnerId(ctx context.Context, id uuid.UUID) (companies []*domain.Company, err error) {
	companies, err = s.companyRepo.GetByOwnerId(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("получение списка компаний по id владельца: %w", err)
	}

	return companies, nil
}

func (s *Service) GetAll(ctx context.Context, filters utils.Filters) (companies []*domain.Company, err error) {
	companies, err = s.companyRepo.GetAll(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("получение списка всех компаний: %w", err)
	}

	return companies, nil
}

func (s *Service) Update(ctx context.Context, company *domain.Company) (err error) {
	err = s.companyRepo.Update(ctx, company)
	if err != nil {
		return fmt.Errorf("обновление информации о компании: %w", err)
	}

	return nil
}

func (s *Service) DeleteById(ctx context.Context, id uuid.UUID) (err error) {
	err = s.companyRepo.DeleteById(ctx, id)
	if err != nil {
		return fmt.Errorf("удаление компании по id: %w", err)
	}

	return nil
}
