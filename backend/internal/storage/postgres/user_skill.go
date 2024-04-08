package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"ppo/domain"
)

type UserSkillRepository struct {
	db *pgx.ConnPool
}

func NewUserSkillRepository(db *pgx.ConnPool) domain.IUserSkillRepository {
	return &UserSkillRepository{
		db: db,
	}
}

func (r *UserSkillRepository) Create(ctx context.Context, pair *domain.UserSkill) (err error) {
	query := `insert into ppo.user_skills(user_id, skill_id) 
	values ($1, $2)`

	_, err = r.db.ExecEx(
		ctx,
		query,
		nil,
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

	_, err = r.db.ExecEx(
		ctx,
		query,
		nil,
		pair.UserId,
		pair.SkillId,
	)
	if err != nil {
		return fmt.Errorf("удаление пары пользователь-навык: %w", err)
	}

	return nil
}

func (r *UserSkillRepository) GetUserSkillsByUserId(ctx context.Context, userId uuid.UUID) (pairs []*domain.UserSkill, err error) {
	query := `select skill_id from ppo.user_skills where user_id = $1`

	rows, err := r.db.QueryEx(
		ctx,
		query,
		nil,
		userId,
	)
	if err != nil {
		return nil, fmt.Errorf("получение навыков пользователя: %w", err)
	}

	pairs = make([]*domain.UserSkill, 0)
	for rows.Next() {
		tmp := new(domain.UserSkill)

		err = rows.Scan(
			&tmp.SkillId,
		)
		if err != nil {
			return nil, fmt.Errorf("сканирование строки: %w", err)
		}

		tmp.UserId = userId
		pairs = append(pairs, tmp)
	}

	return pairs, nil
}

func (r *UserSkillRepository) GetUserSkillsBySkillId(ctx context.Context, skillId uuid.UUID) (pairs []*domain.UserSkill, err error) {
	query := `select user_id from ppo.user_skills where skill_id = $1`

	rows, err := r.db.QueryEx(
		ctx,
		query,
		nil,
		skillId,
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
