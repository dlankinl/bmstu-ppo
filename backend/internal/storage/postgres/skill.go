package postgres

import (
	"context"
	"fmt"
	"ppo/domain"
	"ppo/internal/config"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SkillRepository struct {
	db *pgxpool.Pool
}

func NewSkillRepository(db *pgxpool.Pool) domain.ISkillRepository {
	return &SkillRepository{
		db: db,
	}
}

func (r *SkillRepository) Create(ctx context.Context, skill *domain.Skill) (err error) {
	query := `insert into ppo.skills(name, description) 
	values ($1, $2)`

	_, err = r.db.Exec(
		ctx,
		query,
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
	err = r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&skill.Name,
		&skill.Description,
	)
	if err != nil {
		return nil, fmt.Errorf("получение навыка по id: %w", err)
	}

	skill.ID = id
	return skill, nil
}

func (r *SkillRepository) GetAll(ctx context.Context, page int) (skills []*domain.Skill, numPages int, err error) {
	query := `select id, name, description from ppo.skills offset $1 limit $2`

	rows, err := r.db.Query(
		ctx,
		query,
		(page-1)*config.PageSize,
		config.PageSize,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("получение навыков: %w", err)
	}

	skills = make([]*domain.Skill, 0)
	for rows.Next() {
		tmp := new(domain.Skill)

		err = rows.Scan(
			&tmp.ID,
			&tmp.Name,
			&tmp.Description,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("сканирование полученных строк: %w", err)
		}
		skills = append(skills, tmp)
	}

	var numRecords int
	err = r.db.QueryRow(
		ctx,
		`select count(*) from ppo.skills`,
	).Scan(&numRecords)
	if err != nil {
		return nil, 0, fmt.Errorf("получение количества навыков предпринимателя: %w", err)
	}

	numPages = numRecords / config.PageSize
	if numRecords%config.PageSize != 0 {
		numPages++
	}

	return skills, numPages, nil
}

func (r *SkillRepository) Update(ctx context.Context, skill *domain.Skill) (err error) {
	query := `
			update ppo.skills
			set 
			    name = $1, 
			    description = $2 
			where id = $3`

	_, err = r.db.Exec(
		ctx,
		query,
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
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("открытие транзакции: %w", err)
	}

	_, err = tx.Exec(
		ctx,
		`delete from ppo.skills where id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("удаление навыка по id: %w", err)
	}

	_, err = tx.Exec(
		ctx,
		`delete from ppo.user_skills where skill_id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("удаление связанных с навыком записей: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("закрытие транзакции: %w", err)
	}

	return nil
}
