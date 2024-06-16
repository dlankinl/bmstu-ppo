package postgres

import (
	"context"
	"fmt"
	"ppo/domain"
	"ppo/internal/config"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReviewRepository struct {
	db *pgxpool.Pool
}

func NewReviewRepository(db *pgxpool.Pool) domain.IReviewRepository {
	return &ReviewRepository{
		db: db,
	}
}

func (r *ReviewRepository) Create(ctx context.Context, rev *domain.Review) (err error) {
	query := `insert into ppo.reviews(target_id, reviewer_id, pros, cons, description, rating) 
	values ($1, $2, $3, $4, $5, $6)`

	_, err = r.db.Exec(
		ctx,
		query,
		rev.Target,
		rev.Reviewer,
		rev.Pros,
		rev.Cons,
		rev.Description,
		rev.Rating,
	)
	if err != nil {
		return fmt.Errorf("создание отзыва: %w", err)
	}

	return nil
}

func (r *ReviewRepository) Get(ctx context.Context, id uuid.UUID) (rev *domain.Review, err error) {
	query := `select
		target_id,
		reviewer_id,
		pros,
		cons,
		description,
		rating
	from ppo.reviews 
	where id = $1`

	rev = new(domain.Review)
	err = r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&rev.Target,
		&rev.Reviewer,
		&rev.Pros,
		&rev.Cons,
		&rev.Description,
		&rev.Rating,
	)
	if err != nil {
		return nil, fmt.Errorf("получение отзыва по id: %w", err)
	}

	rev.ID = id
	return rev, nil
}

func (r *ReviewRepository) GetAllForReviewer(ctx context.Context, id uuid.UUID, page int) (revs []*domain.Review, numPages int, err error) {
	query :=
		`select
			id,
    		target_id,
    		pros,
    		cons,
    		description,
			rating
		from ppo.reviews 
		where reviewer_id = $1
		offset $2 limit $3`

	rows, err := r.db.Query(
		ctx,
		query,
		id,
		(page-1)*config.PageSize,
		config.PageSize,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("получение отзывов ревьюера: %w", err)
	}

	revs = make([]*domain.Review, 0)
	for rows.Next() {
		tmp := new(domain.Review)

		err = rows.Scan(
			&tmp.ID,
			&tmp.Target,
			&tmp.Pros,
			&tmp.Cons,
			&tmp.Description,
			&tmp.Rating,
		)
		tmp.Reviewer = id

		if err != nil {
			return nil, 0, fmt.Errorf("сканирование полученных строк: %w", err)
		}

		revs = append(revs, tmp)
	}

	var numRecords int
	err = r.db.QueryRow(
		ctx,
		`select count(*) from ppo.reviews where reviewer_id = $1`,
		id,
	).Scan(&numRecords)
	if err != nil {
		return nil, 0, fmt.Errorf("получение количества отзывов ревьюера: %w", err)
	}

	numPages = numRecords / config.PageSize
	if numRecords%config.PageSize != 0 {
		numPages++
	}

	return revs, numPages, nil
}

func (r *ReviewRepository) GetAllForTarget(ctx context.Context, id uuid.UUID, page int) (revs []*domain.Review, numPages int, err error) {
	query :=
		`select
			id,
    		reviewer_id,
    		pros,
    		cons,
    		description,
			rating
		from ppo.reviews 
		where target_id = $1
		offset $2 limit $3`

	rows, err := r.db.Query(
		ctx,
		query,
		id,
		(page-1)*config.PageSize,
		config.PageSize,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("получение отзывов ревьюера: %w", err)
	}

	revs = make([]*domain.Review, 0)
	for rows.Next() {
		tmp := new(domain.Review)

		err = rows.Scan(
			&tmp.ID,
			&tmp.Reviewer,
			&tmp.Pros,
			&tmp.Cons,
			&tmp.Description,
			&tmp.Rating,
		)
		tmp.Target = id

		if err != nil {
			return nil, 0, fmt.Errorf("сканирование полученных строк: %w", err)
		}

		revs = append(revs, tmp)
	}

	var numRecords int
	err = r.db.QueryRow(
		ctx,
		`select count(*) from ppo.reviews where target_id = $1`,
		id,
	).Scan(&numRecords)
	if err != nil {
		return nil, 0, fmt.Errorf("получение количества отзывов объекта: %w", err)
	}

	numPages = numRecords / config.PageSize
	if numRecords%config.PageSize != 0 {
		numPages++
	}

	return revs, numPages, nil
}

func (r *ReviewRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	query := `delete from ppo.reviews where id = $1`

	_, err = r.db.Exec(
		ctx,
		query,
		id,
	)
	if err != nil {
		return fmt.Errorf("удаление отзыва по id: %w", err)
	}

	return nil
}
