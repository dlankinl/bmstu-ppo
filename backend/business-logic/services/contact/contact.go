package contact

import (
	"business-logic/domain"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type Service struct {
	contactRepo domain.IContactsRepository
}

func NewService(conRepo domain.IContactsRepository) domain.IContactsService {
	return &Service{
		contactRepo: conRepo,
	}
}

func (s *Service) Create(contact *domain.Contact) (err error) {
	if contact.Name == "" {
		return fmt.Errorf("должно быть указано название средства связи")
	}

	if contact.Value == "" {
		return fmt.Errorf("должно быть указано значение средства связи")
	}

	ctx := context.Background()

	err = s.contactRepo.Create(ctx, contact)
	if err != nil {
		return fmt.Errorf("добавление средства связи: %w", err)
	}

	return nil
}

func (s *Service) GetById(id uuid.UUID) (contact *domain.Contact, err error) {
	ctx := context.Background()

	contact, err = s.contactRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("получение средства связи по id: %w", err)
	}

	return contact, nil
}

func (s *Service) GetByOwnerId(id uuid.UUID, page int) (contacts []*domain.Contact, err error) {
	ctx := context.Background()

	contacts, err = s.contactRepo.GetByOwnerId(ctx, id, page)
	if err != nil {
		return nil, fmt.Errorf("получение всех средств связи по id владельца: %w", err)
	}

	return contacts, nil
}

func (s *Service) Update(contact *domain.Contact) (err error) {
	ctx := context.Background()

	err = s.contactRepo.Update(ctx, contact)
	if err != nil {
		return fmt.Errorf("обновление информации о средстве связи: %w", err)
	}

	return nil
}

func (s *Service) DeleteById(id uuid.UUID) (err error) {
	ctx := context.Background()

	err = s.contactRepo.DeleteById(ctx, id)
	if err != nil {
		return fmt.Errorf("удаление средства связи по id: %w", err)
	}

	return nil
}
