package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"ppo/internal/app"
	"ppo/internal/config"
	"ppo/web"
)

var tokenAuth *jwtauth.JWTAuth

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

	tokenAuth = jwtauth.New("HS256", []byte(cfg.JwtKey), nil)

	pool, err := newConn(context.Background(), &cfg.DBConfig)
	if err != nil {
		log.Fatalln(err)
	}

	a := app.NewApp(pool, cfg)

	mux := chi.NewMux()

	mux.Use(middleware.Logger)

	mux.Route("/admin", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))
		r.Use(web.ValidateAdminRoleJWT)

		r.Route("/entrepreneurs", func(r chi.Router) {
			r.Patch("/{id}/update", web.UpdateEntrepreneur(a))
			r.Delete("/{id}/delete", web.DeleteEntrepreneur(a))
		})

		r.Route("/skills", func(r chi.Router) {
			r.Post("/create", web.CreateSkill(a))
			r.Delete("/{id}/delete", web.DeleteSkill(a))
			r.Patch("/{id}/update", web.UpdateSkill(a))
		})
	})

	mux.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))
	})

	mux.Post("/login", web.LoginHandler(a))
	mux.Post("/signup", web.RegisterHandler(a))

	mux.Get("/entrepreneurs/{page}", web.ListEntrepreneurs(a))

	fmt.Println("server was started")
	http.ListenAndServe("localhost:8080", mux)
}
