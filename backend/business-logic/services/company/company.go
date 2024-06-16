package company

import (
	"business-logic/domain"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type Service struct {
	companyRepo domain.ICompanyRepository
}

func NewService(
	companyRepo domain.ICompanyRepository,
) domain.ICompanyService {
	return &Service{
		companyRepo: companyRepo,
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

func (s *Service) GetByOwnerId(id uuid.UUID, page int) (companies []*domain.Company, err error) {
	ctx := context.Background()

	companies, err = s.companyRepo.GetByOwnerId(ctx, id, page)
	if err != nil {
		return nil, fmt.Errorf("получение списка компаний по id владельца: %w", err)
	}

	return companies, nil
}

func (s *Service) GetAll(page int) (companies []*domain.Company, err error) {
	ctx := context.Background()

	companies, err = s.companyRepo.GetAll(ctx, page)
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
