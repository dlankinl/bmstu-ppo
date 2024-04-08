package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserActivityFieldRepository struct {
	db *pgxpool.Pool
}

//func NewUserActivityFieldRepository(db *pgxpool.Pool) domain.IInteractor {
//	return &UserActivityFieldRepository{
//		db: db,
//	}
//}
