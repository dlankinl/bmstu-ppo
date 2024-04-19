package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"ppo/internal/app"
	"ppo/internal/config"
	"ppo/internal/tui"
)

func newConn(ctx context.Context, cfg *config.DBConfig) (pool *pgxpool.Pool, err error) {
	connStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s", cfg.Driver, cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	pool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("подключение к БД: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("пинг БД: %w", err)
	}

	return pool, nil
}

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	pool, err := newConn(context.Background(), &cfg.DBConfig)
	if err != nil {
		log.Fatalln(err)
	}

	app := app.NewApp(pool, cfg)
	termui := tui.NewTUI(app)

	err = termui.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
