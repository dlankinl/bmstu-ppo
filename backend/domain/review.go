package domain

import (
	"context"

	"github.com/google/uuid"
)

type Review struct {
	ID          uuid.UUID
	Reviewer    uuid.UUID
	Target      uuid.UUID
	Pros        string
	Cons        string
	Description string
	Rating      int
}

type IReviewRepository interface {
	Create(context.Context, *Review) error
	Get(context.Context, uuid.UUID) (*Review, error)
	GetAllForReviewer(context.Context, uuid.UUID, int) ([]*Review, int, error)
	GetAllForTarget(context.Context, uuid.UUID, int) ([]*Review, int, error)
	Delete(context.Context, uuid.UUID) error
}

type IReviewService interface {
	Create(context.Context, *Review) error
	Get(context.Context, uuid.UUID) (*Review, error)
	GetAllForReviewer(context.Context, uuid.UUID, int) ([]*Review, int, error)
	GetAllForTarget(context.Context, uuid.UUID, int) ([]*Review, int, error)
	Delete(context.Context, uuid.UUID) error
}
