package main

import (
	"app/internal/app"
	"app/internal/config"
	"app/internal/tui"
	"context"
	"fmt"
	"github.com/dlankinl/bmstu-ppo/backend/business-logic/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

func newConn(ctx context.Context, cfg *config.Config) (pool *pgxpool.Pool, err error) {
	connStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s", cfg.Database.Driver, cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

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

	pool, err := newConn(context.Background(), cfg)
	if err != nil {
		log.Fatalln(err)
	}

	log, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("cоздание логгера: %v", err)
	}

	app := app.NewApp(pool, cfg, log)
	termui := tui.NewTUI(app)

	err = termui.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
