package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"ppo/domain"
)

type UserRepository struct {
	db *pgx.Conn
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) (err error) {
	query := `insert into ppo.`

	r.db.ExecEx(ctx)

	return nil
}

func (r *UserRepository) GetById(ctx context.Context, id uuid.UUID) (user *domain.User, err error) {
	return user, nil
}

func (r *UserRepository) GetAll(ctx context.Context, page int) (users []*domain.User, err error) {
	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) (err error) {
	return nil
}

func (r *UserRepository) DeleteById(ctx context.Context, id uuid.UUID) (err error) {
	return nil
}
