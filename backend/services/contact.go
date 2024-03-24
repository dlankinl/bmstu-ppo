package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"ppo/domain"
)

type ContactService struct {
	contactRepo domain.IContactsRepository
}

func NewContactService(conRepo domain.IContactsRepository) domain.IContactsService {
	return &ContactService{
		contactRepo: conRepo,
	}
}

func (s ContactService) Create(ctx context.Context, contact *domain.Contact) (err error) {
	if contact.Name == "" {
		return fmt.Errorf("должно быть указано название средства связи")
	}

	if contact.Value == "" {
		return fmt.Errorf("должно быть указано значение средства связи")
	}

	err = s.contactRepo.Create(ctx, contact)
	if err != nil {
		return fmt.Errorf("добавление средства связи: %w", err)
	}

	return nil
}

func (s ContactService) GetById(ctx context.Context, id uuid.UUID) (contact *domain.Contact, err error) {
	contact, err = s.contactRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("получение средства связи по id: %w", err)
	}

	return contact, nil
}

// TODO: фильтрация
func (s ContactService) GetByOwnerId(ctx context.Context, id uuid.UUID) (contacts []*domain.Contact, err error) {
	contacts, err = s.contactRepo.GetByOwnerId(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("получение всех средств связи по id владельца: %w", err)
	}

	return contacts, nil
}

func (s ContactService) Update(ctx context.Context, contact *domain.Contact) (err error) {
	err = s.contactRepo.Update(ctx, contact)
	if err != nil {
		return fmt.Errorf("обновление информации о средстве связи: %w", err)
	}

	return nil
}

func (s ContactService) DeleteById(ctx context.Context, id uuid.UUID) (err error) {
	err = s.contactRepo.DeleteById(ctx, id)
	if err != nil {
		return fmt.Errorf("удаление средства связи по id: %w", err)
	}

	return nil
}
