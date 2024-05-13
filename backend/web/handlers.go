package web

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"math"
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
			errorResponse(w, fmt.Errorf("JWT-token verification: %w", err).Error(), http.StatusInternalServerError)
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
		page := r.URL.Query().Get("page")
		if page == "" {
			errorResponse(w, fmt.Errorf("empty page number").Error(), http.StatusBadRequest)
			return
		}

		pageInt, err := strconv.Atoi(page)
		if err != nil {
			errorResponse(w, fmt.Errorf("converting page to int: %w", err).Error(), http.StatusBadRequest)
			return
		}

		users, err := app.UserSvc.GetAll(r.Context(), pageInt)
		if err != nil {
			errorResponse(w, fmt.Errorf("getting users: %w", err).Error(), http.StatusInternalServerError)
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
			errorResponse(w, fmt.Errorf("converting id to uuid: %w", err).Error(), http.StatusBadRequest)
			return
		}

		userDb, err := app.UserSvc.GetById(r.Context(), idUuid)
		if err != nil {
			errorResponse(w, fmt.Errorf("getting user from database by id: %w", err).Error(), http.StatusInternalServerError)
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
			errorResponse(w, fmt.Errorf("converting id to uuid: %w", err).Error(), http.StatusBadRequest)
			return
		}

		_, err = app.UserSvc.GetById(r.Context(), idUuid)
		if err != nil {
			errorResponse(w, fmt.Errorf("deleting user by id: %w", err).Error(), http.StatusBadRequest)
			return
		}

		err = app.UserSvc.DeleteById(r.Context(), idUuid)
		if err != nil {
			errorResponse(w, fmt.Errorf("deleting user by id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(w, http.StatusOK, nil)
	}
}

func GetEntrepreneur(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			errorResponse(w, fmt.Errorf("empty id").Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			errorResponse(w, fmt.Errorf("converting id to uuid: %w", err).Error(), http.StatusBadRequest)
			return
		}

		user, err := app.UserSvc.GetById(r.Context(), idUuid)
		if err != nil {
			errorResponse(w, fmt.Errorf("getting user by id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(w, http.StatusOK, map[string]interface{}{"entrepreneur": toUserTransport(user)})
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
			errorResponse(w, fmt.Errorf("converting id to uuid: %w", err).Error(), http.StatusBadRequest)
			return
		}

		_, err = app.SkillSvc.GetById(r.Context(), idUuid)
		if err != nil {
			errorResponse(w, fmt.Errorf("deleting skill by id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		err = app.SkillSvc.DeleteById(r.Context(), idUuid)
		if err != nil {
			errorResponse(w, fmt.Errorf("deleting skill by id: %w", err).Error(), http.StatusInternalServerError)
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
			errorResponse(w, fmt.Errorf("converting id to uuid: %w", err).Error(), http.StatusBadRequest)
			return
		}

		skillDb, err := app.SkillSvc.GetById(r.Context(), idUuid)
		if err != nil {
			errorResponse(w, fmt.Errorf("getting skill from database by id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		var req Skill

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			errorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.Name != "" {
			skillDb.Name = req.Name
		}
		if req.Description != "" {
			skillDb.Description = req.Description
		}

		err = app.SkillSvc.Update(r.Context(), skillDb)
		if err != nil {
			errorResponse(w, fmt.Errorf("updating skill info: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(w, http.StatusOK, nil)
	}
}

func GetSkill(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			errorResponse(w, fmt.Errorf("empty id").Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			errorResponse(w, fmt.Errorf("converting id to uuid: %w", err).Error(), http.StatusBadRequest)
			return
		}

		skill, err := app.SkillSvc.GetById(r.Context(), idUuid)
		if err != nil {
			errorResponse(w, fmt.Errorf("getting skill by id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(w, http.StatusOK, map[string]interface{}{"skill": toSkillTransport(skill)})
	}
}

func CreateContact(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr, err := getStringClaimFromJWT(r.Context(), "sub")
		if err != nil {
			errorResponse(w, fmt.Errorf("getting claim from JWT: %w", err).Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(idStr)
		if err != nil {
			errorResponse(w, fmt.Errorf("converting string to uuid: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		var req Contact
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			errorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		contact := toContactModel(&req)
		contact.OwnerID = idUuid

		err = app.ConSvc.Create(r.Context(), &contact)
		if err != nil {
			errorResponse(w, fmt.Errorf("creating contact: %w", err).Error(), http.StatusBadRequest)
			return
		}

		successResponse(w, http.StatusOK, nil)
	}
}

func DeleteContact(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ownerIdStr, err := getStringClaimFromJWT(r.Context(), "sub")
		if err != nil {
			errorResponse(w, fmt.Errorf("getting claim from JWT: %w", err).Error(), http.StatusBadRequest)
			return
		}

		ownerIdUuid, err := uuid.Parse(ownerIdStr)
		if err != nil {
			errorResponse(w, fmt.Errorf("converting string to uuid: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		id := chi.URLParam(r, "id")
		if id == "" {
			errorResponse(w, fmt.Errorf("empty id").Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			errorResponse(w, fmt.Errorf("converting id to uuid: %w", err).Error(), http.StatusBadRequest)
			return
		}

		contact, err := app.ConSvc.GetById(r.Context(), idUuid)
		if err != nil {
			errorResponse(w, fmt.Errorf("deleting contact by id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		if ownerIdUuid != contact.OwnerID {
			errorResponse(w, fmt.Errorf("only owner can delete his contact").Error(), http.StatusInternalServerError)
			return
		}

		err = app.ConSvc.DeleteById(r.Context(), idUuid)
		if err != nil {
			errorResponse(w, fmt.Errorf("deleting contact by id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(w, http.StatusOK, nil)
	}
}

func UpdateContact(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ownerIdStr, err := getStringClaimFromJWT(r.Context(), "sub")
		if err != nil {
			errorResponse(w, fmt.Errorf("getting claim from JWT: %w", err).Error(), http.StatusBadRequest)
			return
		}

		ownerIdUuid, err := uuid.Parse(ownerIdStr)
		if err != nil {
			errorResponse(w, fmt.Errorf("converting string to uuid: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		id := chi.URLParam(r, "id")
		if id == "" {
			errorResponse(w, fmt.Errorf("empty id").Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			errorResponse(w, fmt.Errorf("converting id to uuid: %w", err).Error(), http.StatusBadRequest)
			return
		}

		conDb, err := app.ConSvc.GetById(r.Context(), idUuid)
		if err != nil {
			errorResponse(w, fmt.Errorf("getting contact from database by id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		if ownerIdUuid != conDb.OwnerID {
			errorResponse(w, fmt.Errorf("only owner can update his contact").Error(), http.StatusInternalServerError)
			return
		}

		var req Contact

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			errorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.Name != "" {
			conDb.Name = req.Name
		}
		if req.Value != "" {
			conDb.Value = req.Value
		}

		err = app.ConSvc.Update(r.Context(), conDb)
		if err != nil {
			errorResponse(w, fmt.Errorf("updating contact info: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(w, http.StatusOK, nil)
	}
}

func GetContact(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			errorResponse(w, fmt.Errorf("empty id").Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			errorResponse(w, fmt.Errorf("converting id to uuid: %w", err).Error(), http.StatusBadRequest)
			return
		}

		contact, err := app.ConSvc.GetById(r.Context(), idUuid)
		if err != nil {
			errorResponse(w, fmt.Errorf("getting contact by id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(w, http.StatusOK, map[string]interface{}{"contact": toContactTransport(contact)})
	}
}

func ListEntrepreneurContacts(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		if page == "" {
			errorResponse(w, fmt.Errorf("empty page number").Error(), http.StatusBadRequest)
			return
		}

		pageInt, err := strconv.Atoi(page)
		if err != nil {
			errorResponse(w, fmt.Errorf("converting page to int: %w", err).Error(), http.StatusBadRequest)
			return
		}

		entId := r.URL.Query().Get("entrepreneur-id")
		if page == "" {
			errorResponse(w, fmt.Errorf("empty entrepreneur id").Error(), http.StatusBadRequest)
			return
		}

		entUuid, err := uuid.Parse(entId)
		if err != nil {
			errorResponse(w, fmt.Errorf("converting entrepreneur id to uuid: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		contacts, err := app.ConSvc.GetByOwnerId(r.Context(), entUuid, pageInt, true)
		if err != nil {
			errorResponse(w, fmt.Errorf("getting users: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		contactsTransport := make([]Contact, len(contacts))
		for i, contact := range contacts {
			contactsTransport[i] = toContactTransport(contact)
		}

		successResponse(w, http.StatusOK, map[string]interface{}{"entrepreneur_id": entId, "contacts": contactsTransport})
	}
}

func CreateActivityField(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ActivityField
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			errorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		actField := toActFieldModel(&req)

		err = app.ActFieldSvc.Create(r.Context(), &actField)
		if err != nil {
			errorResponse(w, fmt.Errorf("creating activity field: %w", err).Error(), http.StatusBadRequest)
			return
		}

		successResponse(w, http.StatusOK, nil)
	}
}

func DeleteActivityField(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			errorResponse(w, fmt.Errorf("empty id").Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			errorResponse(w, fmt.Errorf("converting id to uuid: %w", err).Error(), http.StatusBadRequest)
			return
		}

		_, err = app.ActFieldSvc.GetById(r.Context(), idUuid)
		if err != nil {
			errorResponse(w, fmt.Errorf("deleting activity field by id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		err = app.ActFieldSvc.DeleteById(r.Context(), idUuid)
		if err != nil {
			errorResponse(w, fmt.Errorf("deleting activity field by id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(w, http.StatusOK, nil)
	}
}

func UpdateActivityField(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			errorResponse(w, fmt.Errorf("empty id").Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			errorResponse(w, fmt.Errorf("converting id to uuid: %w", err).Error(), http.StatusBadRequest)
			return
		}

		actFieldDb, err := app.ActFieldSvc.GetById(r.Context(), idUuid)
		if err != nil {
			errorResponse(w, fmt.Errorf("getting activity field from database by id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		var req ActivityField

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			errorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.Name != "" {
			actFieldDb.Name = req.Name
		}
		if req.Description != "" {
			actFieldDb.Description = req.Description
		}
		if !(math.Abs(float64(req.Cost)) < eps) {
			actFieldDb.Cost = req.Cost
		}

		err = app.ActFieldSvc.Update(r.Context(), actFieldDb)
		if err != nil {
			errorResponse(w, fmt.Errorf("updating activity field info: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(w, http.StatusOK, nil)
	}
}

func GetActivityField(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			errorResponse(w, fmt.Errorf("empty id").Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			errorResponse(w, fmt.Errorf("converting id to uuid: %w", err).Error(), http.StatusBadRequest)
			return
		}

		actField, err := app.ActFieldSvc.GetById(r.Context(), idUuid)
		if err != nil {
			errorResponse(w, fmt.Errorf("getting activity field by id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(w, http.StatusOK, map[string]interface{}{"activity_field": toActFieldTransport(actField)})
	}
}

func ListActivityFields(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		if page == "" {
			errorResponse(w, fmt.Errorf("empty page number").Error(), http.StatusBadRequest)
			return
		}

		pageInt, err := strconv.Atoi(page)
		if err != nil {
			errorResponse(w, fmt.Errorf("converting page to int: %w", err).Error(), http.StatusBadRequest)
			return
		}

		actFields, err := app.ActFieldSvc.GetAll(r.Context(), pageInt)
		if err != nil {
			errorResponse(w, fmt.Errorf("getting users: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		actFieldsTransport := make([]ActivityField, len(actFields))
		for i, actField := range actFields {
			actFieldsTransport[i] = toActFieldTransport(actField)
		}

		successResponse(w, http.StatusOK, map[string]interface{}{"activity_fields": actFieldsTransport})
	}
}
