package web

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"ppo/domain"
	"ppo/internal/app"
	"ppo/pkg/base"
	"strconv"
)

func LoginHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Req struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}
		var req Req

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			errorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		ua := &domain.UserAuth{Username: req.Login, Password: req.Password}
		token, err := app.AuthSvc.Login(r.Context(), ua)
		if err != nil {
			errorResponse(w, fmt.Errorf("failed login: %w", err).Error(), http.StatusUnauthorized)
			return
		}

		_, err = base.VerifyAuthToken(token, app.Config.JwtKey)
		if err != nil {
			errorResponse(w, fmt.Errorf("JWT-token verification: %w", err).Error(), http.StatusUnauthorized)
			return
		}

		successResponse(w, http.StatusOK, map[string]string{"token": token})
	}
}

func RegisterHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Req struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}
		var req Req

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			errorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		ua := &domain.UserAuth{Username: req.Login, Password: req.Password}
		err = app.AuthSvc.Register(r.Context(), ua)
		if err != nil {
			errorResponse(w, fmt.Errorf("failed login: %w", err).Error(), http.StatusBadRequest)
			return
		}

		successResponse(w, http.StatusOK, nil)
	}
}

func ListEntrepreneurs(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := chi.URLParam(r, "page")
		if page == "" {
			errorResponse(w, fmt.Errorf("empty page number").Error(), http.StatusBadRequest)
			return
		}

		pageInt, err := strconv.Atoi(page)
		if err != nil {
			errorResponse(w, fmt.Errorf("converting page to int: %w", err).Error(), http.StatusUnauthorized)
			return
		}

		users, err := app.UserSvc.GetAll(r.Context(), pageInt)
		if err != nil {
			errorResponse(w, fmt.Errorf("getting users: %w", err).Error(), http.StatusUnauthorized)
			return
		}

		usersTransport := make([]User, len(users))
		for i, user := range users {
			usersTransport[i] = toUserTransport(user)
		}

		successResponse(w, http.StatusOK, map[string]interface{}{"users": usersTransport})
	}
}

func UpdateEntrepreneur(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			errorResponse(w, fmt.Errorf("empty id").Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			errorResponse(w, fmt.Errorf("converting id to uuid: %w", err).Error(), http.StatusUnauthorized)
			return
		}

		userDb, err := app.UserSvc.GetById(r.Context(), idUuid)
		if err != nil {
			errorResponse(w, fmt.Errorf("getting user from database by id: %w", err).Error(), http.StatusUnauthorized)
			return
		}

		var req User

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			errorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.City != "" {
			userDb.City = req.City
		}
		if req.Role != "" {
			userDb.Role = req.Role
		}
		if req.Gender != "" {
			userDb.Gender = req.Gender
		}
		if !req.Birthday.IsZero() {
			userDb.Birthday = req.Birthday
		}
		if req.FullName != "" {
			userDb.FullName = req.FullName
		}
		if req.Username != "" {
			userDb.Username = req.Username
		}

		err = app.UserSvc.Update(r.Context(), userDb)
		if err != nil {
			errorResponse(w, fmt.Errorf("updating entrepreneur info: %w", err).Error(), http.StatusBadRequest)
			return
		}

		successResponse(w, http.StatusOK, nil)
	}
}

func DeleteEntrepreneur(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			errorResponse(w, fmt.Errorf("empty id").Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			errorResponse(w, fmt.Errorf("converting id to uuid: %w", err).Error(), http.StatusUnauthorized)
			return
		}

		err = app.UserSvc.DeleteById(r.Context(), idUuid)
		if err != nil {
			errorResponse(w, fmt.Errorf("deleting user by id: %w", err).Error(), http.StatusUnauthorized)
			return
		}

		successResponse(w, http.StatusOK, nil)
	}
}

func CreateSkill(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Skill

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			errorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		skill := toSkillModel(&req)

		err = app.SkillSvc.Create(r.Context(), &skill)
		if err != nil {
			errorResponse(w, fmt.Errorf("creating skill: %w", err).Error(), http.StatusBadRequest)
			return
		}

		successResponse(w, http.StatusOK, nil)
	}
}

func DeleteSkill(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			errorResponse(w, fmt.Errorf("empty id").Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			errorResponse(w, fmt.Errorf("converting id to uuid: %w", err).Error(), http.StatusUnauthorized)
			return
		}

		err = app.SkillSvc.DeleteById(r.Context(), idUuid)
		if err != nil {
			errorResponse(w, fmt.Errorf("deleting skill by id: %w", err).Error(), http.StatusUnauthorized)
			return
		}

		successResponse(w, http.StatusOK, nil)
	}
}

func UpdateSkill(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			errorResponse(w, fmt.Errorf("empty id").Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			errorResponse(w, fmt.Errorf("converting id to uuid: %w", err).Error(), http.StatusUnauthorized)
			return
		}

		var req Skill

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			errorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		req.ID = idUuid

		skillModel := toSkillModel(&req)

		err = app.SkillSvc.Update(r.Context(), &skillModel)
		if err != nil {
			errorResponse(w, fmt.Errorf("updating skill info: %w", err).Error(), http.StatusBadRequest)
			return
		}

		successResponse(w, http.StatusOK, nil)
	}
}
