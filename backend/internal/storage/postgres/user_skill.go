package postgres

import (
	"context"
	"fmt"
	"ppo/domain"
	"ppo/internal/config"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserSkillRepository struct {
	db *pgxpool.Pool
}

func NewUserSkillRepository(db *pgxpool.Pool) domain.IUserSkillRepository {
	return &UserSkillRepository{
		db: db,
	}
}

func (r *UserSkillRepository) Create(ctx context.Context, pair *domain.UserSkill) (err error) {
	query := `insert into ppo.user_skills(user_id, skill_id) 
	values ($1, $2)`

	_, err = r.db.Exec(
		ctx,
		query,
		pair.UserId,
		pair.SkillId,
	)
	if err != nil {
		return fmt.Errorf("создание cвязи пользователь-навык: %w", err)
	}

	return nil
}

func (r *UserSkillRepository) Delete(ctx context.Context, pair *domain.UserSkill) (err error) {
	query := `delete from ppo.user_skills where user_id = $1 and skill_id = $2`

	_, err = r.db.Exec(
		ctx,
		query,
		pair.UserId,
		pair.SkillId,
	)
	if err != nil {
		return fmt.Errorf("удаление пары пользователь-навык: %w", err)
	}

	return nil
}

func (r *UserSkillRepository) GetUserSkillsByUserId(ctx context.Context, userId uuid.UUID, page int, isPaginated bool) (pairs []*domain.UserSkill, numPages int, err error) {
	query := `
		select skill_id 
		from ppo.user_skills 
		where user_id = $1`

	var rows pgx.Rows
	if !isPaginated {
		rows, err = r.db.Query(
			ctx,
			query,
			userId,
		)
	} else {
		rows, err = r.db.Query(
			ctx,
			query+` offset $2 limit $3`,
			userId,
			(page-1)*config.PageSize,
			config.PageSize,
		)
	}
	if err != nil {
		return nil, 0, fmt.Errorf("получение навыков пользователя: %w", err)
	}

	pairs = make([]*domain.UserSkill, 0)
	for rows.Next() {
		tmp := new(domain.UserSkill)

		err = rows.Scan(
			&tmp.SkillId,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("сканирование строки: %w", err)
		}

		tmp.UserId = userId
		pairs = append(pairs, tmp)
	}

	var numRecords int
	err = r.db.QueryRow(
		ctx,
		`select count(*) from ppo.user_skills where user_id = $1`,
		userId,
	).Scan(&numRecords)
	if err != nil {
		return nil, 0, fmt.Errorf("получение количества навыков предпринимателя: %w", err)
	}

	numPages = numRecords / config.PageSize
	if numRecords%config.PageSize != 0 {
		numPages++
	}

	return pairs, numPages, nil
}

func (r *UserSkillRepository) GetUserSkillsBySkillId(ctx context.Context, skillId uuid.UUID, page int) (pairs []*domain.UserSkill, err error) {
	query := `
		select user_id 
		from ppo.user_skills 
		where skill_id = $1
		offset $2
		limit $3`

	rows, err := r.db.Query(
		ctx,
		query,
		skillId,
		(page-1)*config.PageSize,
		config.PageSize,
	)
	if err != nil {
		return nil, fmt.Errorf("получение пользователей по навыку: %w", err)
	}

	pairs = make([]*domain.UserSkill, 0)
	for rows.Next() {
		tmp := new(domain.UserSkill)

		err = rows.Scan(
			&tmp.UserId,
		)
		if err != nil {
			return nil, fmt.Errorf("сканирование строки: %w", err)
		}

		tmp.SkillId = skillId
		pairs = append(pairs, tmp)
	}

	return pairs, nil
}
