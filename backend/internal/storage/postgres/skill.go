package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"ppo/domain"
	"ppo/internal/config"
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

func (r *SkillRepository) GetAll(ctx context.Context, page int) (skills []*domain.Skill, err error) {
	query := `select id, name, description from ppo.skills offset $1 limit $2`

	rows, err := r.db.Query(
		ctx,
		query,
		(page-1)*config.PageSize,
		config.PageSize,
	)

	skills = make([]*domain.Skill, 0)
	for rows.Next() {
		tmp := new(domain.Skill)

		err = rows.Scan(
			&tmp.ID,
			&tmp.Name,
			&tmp.Description,
		)

		if err != nil {
			return nil, fmt.Errorf("сканирование полученных строк: %w", err)
		}
		skills = append(skills, tmp)
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
	query := `delete from ppo.skills where id = $1`

	_, err = r.db.Exec(
		ctx,
		query,
		id,
	)
	if err != nil {
		return fmt.Errorf("удаление навыка по id: %w", err)
	}

	return nil
}

/*

Ресурсы добываются с помощью космических станций и у них есть норм связь. На них свозят руду и они выступают хабом.
Приходит фирма, которая занимается производством таких станций. Автоматизаций нет. Корабли всё делают руками.
Важные события:
1) Какой-то рудокоп нарыл 10 тонн и привез ко мне на базу.
2) Рудокоп обнаружил астероид с потенциальной жилы.

На станции стабильный интернет.
2 типа кораблей:
1) Кораблики, добывающие руду, без интернета. Подключается проводом к станции и через него
передается.
2) С интернетом. Но интернет есть не везде. Могут отправлять сообщения, но если вылетают из зоны покрытия, то не удается.

Надо:
1) чтобы можно было просматривать данные о событиях с фильтрациями (по кораблю, астероидам и т.д)
2) в системе видеть отчеты, аналитику и т.д.
3) чтоб за большой период (квартал например) в формат типа excel выгрузить отчет, который будет рассчитываться по каким-то
правилам.
4) авторизация нужна.



Задачи:
1) написать требования, usecase
2) определить функциональные требования

В каждом списке не более 20 вопросов. Команды 3-4 человека.

*/
