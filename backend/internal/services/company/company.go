package company

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"ppo/domain"
	"ppo/pkg/logger"
)

type Service struct {
	actFieldRepo domain.IActivityFieldRepository
	companyRepo  domain.ICompanyRepository
	logger       logger.ILogger
}

func NewService(
	companyRepo domain.ICompanyRepository,
	actFieldRepo domain.IActivityFieldRepository,
	logger logger.ILogger,
) domain.ICompanyService {
	return &Service{
		companyRepo:  companyRepo,
		actFieldRepo: actFieldRepo,
		logger:       logger,
	}
}

func (s *Service) Create(ctx context.Context, company *domain.Company) (err error) {
	prompt := "CompanyCreate"

	if company.Name == "" {
		s.logger.Infof("%s: должно быть указано название компании", prompt)
		return fmt.Errorf("должно быть указано название компании")
	}

	if company.City == "" {
		s.logger.Infof("%s: должно быть указано название города", prompt)
		return fmt.Errorf("должно быть указано название города")
	}

	_, err = s.actFieldRepo.GetById(ctx, company.ActivityFieldId)
	if err != nil {
		s.logger.Infof("%s: поиск сферы деятельности: %v", prompt, err)
		return fmt.Errorf("добавление компании (поиск сферы деятельности): %w", err)
	}

	err = s.companyRepo.Create(ctx, company)
	if err != nil {
		s.logger.Infof("%s: добавление компании: %v", prompt, err)
		return fmt.Errorf("добавление компании: %w", err)
	}

	return nil
}

func (s *Service) GetById(ctx context.Context, id uuid.UUID) (company *domain.Company, err error) {
	prompt := "CompanyGetById"

	company, err = s.companyRepo.GetById(ctx, id)
	if err != nil {
		s.logger.Infof("%s: получение компании по id: %v", prompt, err)
		return nil, fmt.Errorf("получение компании по id: %w", err)
	}

	return company, nil
}

func (s *Service) GetByOwnerId(ctx context.Context, id uuid.UUID, page int, isPaginated bool) (companies []*domain.Company, numPages int, err error) {
	prompt := "CompanyGetByOwnerId"

	companies, numPages, err = s.companyRepo.GetByOwnerId(ctx, id, page, isPaginated)
	if err != nil {
		s.logger.Infof("%s: получение списка компаний по id владельца: %v", prompt, err)
		return nil, 0, fmt.Errorf("получение списка компаний по id владельца: %w", err)
	}

	return companies, numPages, nil
}

func (s *Service) GetAll(ctx context.Context, page int) (companies []*domain.Company, err error) {
	prompt := "CompanyGetAll"

	companies, err = s.companyRepo.GetAll(ctx, page)
	if err != nil {
		s.logger.Infof("%s: получение списка всех компаний: %v", prompt, err)
		return nil, fmt.Errorf("получение списка всех компаний: %w", err)
	}

	return companies, nil
}

func (s *Service) Update(ctx context.Context, company *domain.Company) (err error) {
	prompt := "CompanyUpdate"

	_, err = s.actFieldRepo.GetById(ctx, company.ActivityFieldId)
	if err != nil {
		s.logger.Infof("%s: поиск сферы деятельности: %v", prompt, err)
		return fmt.Errorf("обновление информации о компании (поиск сферы деятельности): %w", err)
	}

	err = s.companyRepo.Update(ctx, company)
	if err != nil {
		s.logger.Infof("%s: обновление информации о компании: %v", prompt, err)
		return fmt.Errorf("обновление информации о компании: %w", err)
	}

	return nil
}

func (s *Service) DeleteById(ctx context.Context, id uuid.UUID) (err error) {
	prompt := "CompanyDeleteById"

	err = s.companyRepo.DeleteById(ctx, id)
	if err != nil {
		s.logger.Infof("%s: удаление компании по id: %v", prompt, err)
		return fmt.Errorf("удаление компании по id: %w", err)
	}

	return nil
}
