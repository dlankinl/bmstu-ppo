package postgres

import (
	"context"
	"data-access-layer/config"
	"data-access-layer/domain"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
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

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (user *domain.User, err error) {
	query := `select id, username, full_name, birthday, gender, city, role from ppo.users where username = $1`

	tmp := new(User)
	err = r.db.QueryRow(
		ctx,
		query,
		username,
	).Scan(
		&tmp.ID,
		&tmp.Username,
		&tmp.FullName,
		&tmp.Birthday,
		&tmp.Gender,
		&tmp.City,
		&tmp.Role,
	)
	if err != nil {
		return nil, fmt.Errorf("получение пользователя по username: %w", err)
	}

	return UserDbToUser(tmp), nil
}

func (r *UserRepository) GetById(ctx context.Context, userId uuid.UUID) (user *domain.User, err error) {
	query := `select username, full_name, birthday, gender, city, role from ppo.users where id = $1`

	user = new(domain.User)
	err = r.db.QueryRow(
		ctx,
		query,
		userId,
	).Scan(
		&user.Username,
		&user.FullName,
		&user.Birthday,
		&user.Gender,
		&user.City,
		&user.Role,
	)
	if err != nil {
		return nil, fmt.Errorf("получение пользователя по id: %w", err)
	}

	user.ID = userId
	return user, nil
}

func (r *UserRepository) GetAll(ctx context.Context, page int) (users []*domain.User, err error) {
	query := `select id, username, full_name, birthday, gender, city from ppo.users
	where full_name is not null 
	    and birthday is not null
	    and gender is not null
	    and city is not null
		and role = 'user'
	offset $1 limit $2`

	rows, err := r.db.Query(
		ctx,
		query,
		(page-1)*config.PageSize,
		config.PageSize,
	)
	if err != nil {
		return nil, fmt.Errorf("получение предпринимателей: %w", err)
	}

	users = make([]*domain.User, 0)
	for rows.Next() {
		tmp := new(domain.User)

		err = rows.Scan(
			&tmp.ID,
			&tmp.Username,
			&tmp.FullName,
			&tmp.Birthday,
			&tmp.Gender,
			&tmp.City,
		)

		if err != nil {
			return nil, fmt.Errorf("сканирование полученных строк: %w", err)
		}
		users = append(users, tmp)
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
			    city = $4,
			    role = $5
			where username = $6`

	_, err = r.db.Exec(
		ctx,
		query,
		user.FullName,
		user.Birthday,
		user.Gender,
		user.City,
		user.Role,
		user.Username,
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
