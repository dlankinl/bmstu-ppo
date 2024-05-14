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

	mux.Route("/skills", func(r chi.Router) {
		r.Get("/{id}", web.GetSkill(a))
		r.Get("/", web.ListEntrepreneurSkills(a))

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))
			r.Use(web.ValidateAdminRoleJWT)

			r.Post("/create", web.CreateSkill(a))
			r.Delete("/{id}/delete", web.DeleteSkill(a))
			r.Patch("/{id}/update", web.UpdateSkill(a))
		})
	})

	mux.Route("/entrepreneurs", func(r chi.Router) {
		r.Get("/{id}", web.GetEntrepreneur(a))
		r.Get("/", web.ListEntrepreneurs(a))

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))

			r.Group(func(r chi.Router) {
				r.Use(web.ValidateAdminRoleJWT)

				r.Patch("/{id}/update", web.UpdateEntrepreneur(a))
				r.Delete("/{id}/delete", web.DeleteEntrepreneur(a))
			})

			//r.Group(func(r chi.Router) {
			//	r.Use(web.ValidateUserRoleJWT)
			//
			//	r.Post("/{id}/skills/create", web.CreateUserSkill(a))
			//})
		})
	})

	mux.Route("/contacts", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))
			r.Use(web.ValidateUserRoleJWT)

			r.Get("/{id}", web.GetContact(a))
			r.Get("/", web.ListEntrepreneurContacts(a))

			r.Post("/create", web.CreateContact(a))
			r.Patch("/{id}/update", web.UpdateContact(a))
			r.Delete("/{id}/delete", web.DeleteContact(a))
		})
	})

	mux.Route("/activity_fields", func(r chi.Router) {
		r.Get("/{id}", web.GetActivityField(a))
		r.Get("/", web.ListActivityFields(a))

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))
			r.Use(web.ValidateAdminRoleJWT)

			r.Post("/create", web.CreateActivityField(a))
			r.Patch("/{id}/update", web.UpdateActivityField(a))
			r.Delete("/{id}/delete", web.DeleteActivityField(a))
		})
	})

	mux.Route("/companies", func(r chi.Router) {
		r.Get("/{id}", web.GetCompany(a))
		r.Get("/", web.ListEntrepreneurCompanies(a))

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))
			r.Use(web.ValidateUserRoleJWT)

			r.Post("/create", web.CreateCompany(a))
			r.Patch("/{id}/update", web.UpdateCompany(a))
			r.Delete("/{id}/delete", web.DeleteCompany(a))
		})

		r.Route("/{id}/financials", func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))
			r.Use(web.ValidateUserRoleJWT)

			r.Post("/create", web.CreateReport(a))
		})
	})

	mux.Route("/user-skills", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))
			r.Use(web.ValidateUserRoleJWT)

			r.Post("/create", web.CreateUserSkill(a))
			r.Delete("/{id}/delete", web.DeleteUserSkill(a))
		})
	})

	mux.Route("/financials", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))
			r.Use(web.ValidateUserRoleJWT)

			r.Delete("/{id}/delete", web.DeleteFinReport(a))
			r.Patch("/{id}/update", web.UpdateFinReport(a))
		})
	})

	mux.Post("/login", web.LoginHandler(a))
	mux.Post("/signup", web.RegisterHandler(a))

	fmt.Println("server was started")
	http.ListenAndServe("localhost:8080", mux)
}
