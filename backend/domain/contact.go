package domain

import (
	"context"
	"github.com/google/uuid"
)

type Contact struct {
	ID      uuid.UUID
	OwnerID uuid.UUID
	Name    string
	Value   string
}

type IContactsRepository interface {
	Create(context.Context, *Contact) error
	GetById(context.Context, uuid.UUID) (*Contact, error)
	GetByOwnerId(context.Context, uuid.UUID) ([]*Contact, error)
	Update(context.Context, *Contact) error
	DeleteById(context.Context, uuid.UUID) error
}

type IContactsService interface {
	Create(context.Context, *Contact) error
	GetById(context.Context, uuid.UUID) (*Contact, error)
	GetByOwnerId(context.Context, uuid.UUID) ([]*Contact, error)
	Update(context.Context, *Contact) error
	DeleteById(context.Context, uuid.UUID) error
}
