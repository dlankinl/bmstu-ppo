package domain

import "github.com/google/uuid"

type Contact struct {
	ID    uuid.UUID
	Name  string
	Value string
}

type IContactsRepository interface {
	Create(contact *Contact) error
	GetById(id uuid.UUID) (*Contact, error)
	GetAllByUserId(id uuid.UUID) ([]*Contact, error)
	Update(contact *Contact) error
	DeleteById(id uuid.UUID) error
}

type IContactsService interface {
	Create(contact *Contact) error
	GetById(id uuid.UUID) (*Contact, error)
	GetAllByUserId(id uuid.UUID) ([]*Contact, error)
	Update(contact *Contact) error
	DeleteById(id uuid.UUID) error
}
