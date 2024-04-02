package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"ppo/domain"
)

type SkillRepository struct {
	db *pgx.Conn
}

func (r *SkillRepository) Create(ctx context.Context, skill *domain.Skill) (err error) {
	query := `insert into ppo.skills(name, description) 
	values ($1, $2)`

	_, err = r.db.ExecEx(
		ctx,
		query,
		nil,
		skill.Name,
		skill.Description,
	)
	if err != nil {
		return fmt.Errorf("создание навыка: %w", err)
	}

	return nil
}

func (r *SkillRepository) GetById(ctx context.Context, id uuid.UUID) (skill *domain.Skill, err error) {
	query := `select name, description from ppo.skills where id = $1`

	skill = new(domain.Skill)
	err = r.db.QueryRowEx(
		ctx,
		query,
		nil,
		id,
	).Scan(
		&skill.Name,
		&skill.Description,
	)
	if err != nil {
		return nil, fmt.Errorf("получение навыка по id: %w", err)
	}

	return skill, nil
}

func (r *SkillRepository) GetAll(ctx context.Context, page int) (skills []*domain.Skill, err error) {
	query := `select name, description from ppo.skills offset $1 limit $2`

	rows, err := r.db.QueryEx(
		ctx,
		query,
		nil,
		(page-1)*pageSize,
		pageSize,
	)

	skills = make([]*domain.Skill, 0)
	for rows.Next() {
		tmp := new(domain.Skill)

		err = rows.Scan(
			&tmp.Name,
			&tmp.Description,
		)

		if err != nil {
			return nil, fmt.Errorf("сканирование полученных строк: %w", err)
		}
	}

	return skills, nil
}

func (r *SkillRepository) Update(ctx context.Context, skill *domain.Skill) (err error) {
	query := `
			update ppo.skills
			set 
			    name = $1, 
			    description = $2 
			where id = $3`

	_, err = r.db.ExecEx(
		ctx,
		query,
		nil,
		skill.Name,
		skill.Description,
		skill.ID,
	)
	if err != nil {
		return fmt.Errorf("обновление информации о навыке: %w", err)
	}

	return nil
}

func (r *SkillRepository) DeleteById(ctx context.Context, id uuid.UUID) (err error) {
	query := `delete from ppo.skills where id = $1`

	_, err = r.db.ExecEx(
		ctx,
		query,
		nil,
		id,
	)
	if err != nil {
		return fmt.Errorf("удаление навыка по id: %w", err)
	}

	return nil
}
