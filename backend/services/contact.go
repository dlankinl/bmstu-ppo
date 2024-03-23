package services

import (
	"fmt"
	"github.com/google/uuid"
	"ppo/domain"
)

type ContactService struct {
	contactRepo domain.IContactsRepository
}

func (s ContactService) Create(contact *domain.Contact) (err error) {
	if contact.Name == "" {
		return fmt.Errorf("должно быть указано название средства связи")
	}

	if contact.Value == "" {
		return fmt.Errorf("должно быть указано значение средства связи")
	}

	err = s.contactRepo.Create(contact)
	if err != nil {
		return fmt.Errorf("добавление контакта для связи: %w", err)
	}

	return nil
}

func (s ContactService) GetById(id uuid.UUID) (contact *domain.Contact, err error) {
	contact, err = s.contactRepo.GetById(id)
	if err != nil {
		return nil, fmt.Errorf("получение контакта для связи по id: %w", err)
	}

	return contact, nil
}

// TODO: фильтрация
func (s ContactService) GetAllByUserId(id uuid.UUID) (contacts []*domain.Contact, err error) {
	contacts, err = s.contactRepo.GetAllByUserId(id)
	if err != nil {
		return nil, fmt.Errorf("получение всех средств связи по id владельца: %w", err)
	}

	return contacts, nil
}

func (s ContactService) Update(contact *domain.Contact) (err error) {
	err = s.contactRepo.Update(contact)
	if err != nil {
		return fmt.Errorf("обновление контакта связи: %w", err)
	}

	return nil
}

func (s ContactService) DeleteById(id uuid.UUID) (err error) {
	err = s.DeleteById(id)
	if err != nil {
		return fmt.Errorf("удаление средства связи по id: %w", err)
	}

	return nil
}
