package postgres

import "github.com/jackc/pgx/v5/pgxpool"

const PageSize = 20

var testDbInstance *pgxpool.Pool
