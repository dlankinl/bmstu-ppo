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
	Create(ctx context.Context, contact *Contact) error
	GetById(ctx context.Context, id uuid.UUID) (*Contact, error)
	GetByOwnerId(ctx context.Context, id uuid.UUID, page int) ([]*Contact, error)
	Update(ctx context.Context, contact *Contact) error
	DeleteById(ctx context.Context, id uuid.UUID) error
}
