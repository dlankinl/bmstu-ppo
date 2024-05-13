package company

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"ppo/domain"
)

type Service struct {
	actFieldRepo domain.IActivityFieldRepository
	companyRepo  domain.ICompanyRepository
}

func NewService(
	companyRepo domain.ICompanyRepository,
	actFieldRepo domain.IActivityFieldRepository,
) domain.ICompanyService {
	return &Service{
		companyRepo:  companyRepo,
		actFieldRepo: actFieldRepo,
	}
}

func (s *Service) Create(ctx context.Context, company *domain.Company) (err error) {
	if company.Name == "" {
		return fmt.Errorf("должно быть указано название компании")
	}

	if company.City == "" {
		return fmt.Errorf("должно быть указано название города")
	}

	_, err = s.actFieldRepo.GetById(ctx, company.ActivityFieldId)
	if err != nil {
		return fmt.Errorf("добавление компании (поиск сферы деятельности): %w", err)
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

func (s *Service) GetByOwnerId(ctx context.Context, id uuid.UUID, page int, isPaginated bool) (companies []*domain.Company, err error) {
	companies, err = s.companyRepo.GetByOwnerId(ctx, id, page, isPaginated)
	if err != nil {
		return nil, fmt.Errorf("получение списка компаний по id владельца: %w", err)
	}

	return companies, nil
}

func (s *Service) GetAll(ctx context.Context, page int) (companies []*domain.Company, err error) {
	companies, err = s.companyRepo.GetAll(ctx, page)
	if err != nil {
		return nil, fmt.Errorf("получение списка всех компаний: %w", err)
	}

	return companies, nil
}

func (s *Service) Update(ctx context.Context, company *domain.Company) (err error) {
	_, err = s.actFieldRepo.GetById(ctx, company.ActivityFieldId)
	if err != nil {
		return fmt.Errorf("обновление информации о компании (поиск сферы деятельности): %w", err)
	}

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
