package postgres

import "github.com/jackc/pgx/v5/pgxpool"

const pageSize = 20

var testDbInstance *pgxpool.Pool
