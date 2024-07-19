package contact

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"ppo/domain"
	"ppo/internal/config"
	"ppo/pkg/logger"
)

type Service struct {
	contactRepo domain.IContactsRepository
	logger      logger.ILogger
}

func NewService(conRepo domain.IContactsRepository, logger logger.ILogger) domain.IContactsService {
	return &Service{
		contactRepo: conRepo,
		logger:      logger,
	}
}

func (s *Service) Create(ctx context.Context, contact *domain.Contact) (err error) {
	prompt := "ContactCreate"

	if contact.Name == "" {
		s.logger.Infof("%s: должно быть указано название средства связи", prompt)
		return fmt.Errorf("должно быть указано название средства связи")
	}

	if contact.Value == "" {
		s.logger.Infof("%s: должно быть указано значение средства связи", prompt)
		return fmt.Errorf("должно быть указано значение средства связи")
	}

	contacts, err := s.contactRepo.GetByOwnerId(ctx, contact.OwnerID)
	if err != nil {
		s.logger.Infof("%s: добавление средства связи: %v", prompt, err)
		return fmt.Errorf("добавление средства связи: %w", err)
	}

	if len(contacts) >= config.MaxContacts {
		s.logger.Infof("%s: количество не должно быть более %d", config.MaxContacts)
		return fmt.Errorf("добавление средства связи: количество не должно быть более %d", config.MaxContacts)
	}

	err = s.contactRepo.Create(ctx, contact)
	if err != nil {
		s.logger.Infof("%s: добавление средства связи: %v", prompt, err)
		return fmt.Errorf("добавление средства связи: %w", err)
	}

	return nil
}

func (s *Service) GetById(ctx context.Context, id uuid.UUID) (contact *domain.Contact, err error) {
	prompt := "ContactGetById"

	contact, err = s.contactRepo.GetById(ctx, id)
	if err != nil {
		s.logger.Infof("%s: получение средства связи по id: %v", prompt, err)
		return nil, fmt.Errorf("получение средства связи по id: %w", err)
	}

	return contact, nil
}

func (s *Service) GetByOwnerId(ctx context.Context, id uuid.UUID) (contacts []*domain.Contact, err error) {
	prompt := "ContactGetByOwnerId"

	contacts, err = s.contactRepo.GetByOwnerId(ctx, id)
	if err != nil {
		s.logger.Infof("%s: получение всех средств связи по id владельца: %v", prompt, err)
		return nil, fmt.Errorf("получение всех средств связи по id владельца: %w", err)
	}

	return contacts, nil
}

func (s *Service) Update(ctx context.Context, contact *domain.Contact) (err error) {
	prompt := "ContactUpdate"

	err = s.contactRepo.Update(ctx, contact)
	if err != nil {
		s.logger.Infof("%s: обновление информации о средстве связи: %v", prompt, err)
		return fmt.Errorf("обновление информации о средстве связи: %w", err)
	}

	return nil
}

func (s *Service) DeleteById(ctx context.Context, id uuid.UUID) (err error) {
	prompt := "ContactDeleteById"

	err = s.contactRepo.DeleteById(ctx, id)
	if err != nil {
		s.logger.Infof("%s: удаление средства связи по id: %v", prompt, err)
		return fmt.Errorf("удаление средства связи по id: %w", err)
	}

	return nil
}
