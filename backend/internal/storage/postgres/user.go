package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"ppo/domain"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) domain.IUserRepository {
	return &UserRepository{
		db: db,
	}
}

// FIXME: не нужен получается
func (r *UserRepository) Create(ctx context.Context, user *domain.User) (err error) {
	query := `update ppo.users
		set 
		    full_name = $1,
		    birthday = $2,
		    gender = $3,
		    city = $4
		where id = $5`

	_, err = r.db.Exec(
		ctx,
		query,
		user.FullName,
		user.Birthday,
		user.Gender,
		user.City,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("создание пользователя: %w", err)
	}

	return nil
}

func (r *UserRepository) GetById(ctx context.Context, id uuid.UUID) (user *domain.User, err error) {
	query := `select username, full_name, birthday, gender, city from ppo.users where id = $1`

	user = new(domain.User)
	err = r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&user.Username,
		&user.FullName,
		&user.Birthday,
		&user.Gender,
		&user.City,
	)
	if err != nil {
		return nil, fmt.Errorf("получение пользователя по id: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetAll(ctx context.Context, page int) (users []*domain.User, err error) {
	query := `select username, full_name, birthday, gender, city from ppo.users offset $1 limit $2`

	rows, err := r.db.Query(
		ctx,
		query,
		(page-1)*pageSize,
		pageSize,
	)

	users = make([]*domain.User, 0)
	for rows.Next() {
		tmp := new(domain.User)

		err = rows.Scan(
			&tmp.Username,
			&tmp.FullName,
			&tmp.Birthday,
			&tmp.Gender,
			&tmp.City,
		)

		if err != nil {
			return nil, fmt.Errorf("сканирование полученных строк: %w", err)
		}
	}

	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) (err error) {
	query := `
			update ppo.users
			set 
			    full_name = $1, 
			    birthday = $2, 
			    gender = $3, 
			    city = $4
			where id = $5`

	_, err = r.db.Exec(
		ctx,
		query,
		user.FullName,
		user.Birthday,
		user.Gender,
		user.City,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("обновление информации о пользователе: %w", err)
	}

	return nil
}

func (r *UserRepository) DeleteById(ctx context.Context, id uuid.UUID) (err error) {
	query := `delete from ppo.users where id = $1`

	_, err = r.db.Exec(
		ctx,
		query,
		id,
	)
	if err != nil {
		return fmt.Errorf("удаление пользователя по id: %w", err)
	}

	return nil
}
