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
			errorResponse(w, fmt.Errorf("getting contacts: %w", err).Error(), http.StatusInternalServerError)
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

func CreateCompany(app *app.App) http.HandlerFunc {
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

		var req Company
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			errorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		company := toCompanyModel(&req)
		company.OwnerID = idUuid

		err = app.CompSvc.Create(r.Context(), &company)
		if err != nil {
			errorResponse(w, fmt.Errorf("creating company: %w", err).Error(), http.StatusBadRequest)
			return
		}

		successResponse(w, http.StatusOK, nil)
	}
}

func DeleteCompany(app *app.App) http.HandlerFunc {
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

		company, err := app.CompSvc.GetById(r.Context(), idUuid)
		if err != nil {
			errorResponse(w, fmt.Errorf("deleting company by id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		if ownerIdUuid != company.OwnerID {
			errorResponse(w, fmt.Errorf("only owner can delete his company").Error(), http.StatusInternalServerError)
			return
		}

		err = app.CompSvc.DeleteById(r.Context(), idUuid)
		if err != nil {
			errorResponse(w, fmt.Errorf("deleting company by id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(w, http.StatusOK, nil)
	}
}

func UpdateCompany(app *app.App) http.HandlerFunc {
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

		compDb, err := app.CompSvc.GetById(r.Context(), idUuid)
		if err != nil {
			errorResponse(w, fmt.Errorf("getting company from database by id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		if ownerIdUuid != compDb.OwnerID {
			errorResponse(w, fmt.Errorf("only owner can update his company").Error(), http.StatusInternalServerError)
			return
		}

		var req Company

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			errorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.ActivityFieldId.ID() != 0 {
			compDb.ActivityFieldId = req.ActivityFieldId
		}
		if req.Name != "" {
			compDb.Name = req.Name
		}
		if req.City != "" {
			compDb.City = req.City
		}

		err = app.CompSvc.Update(r.Context(), compDb)
		if err != nil {
			errorResponse(w, fmt.Errorf("updating company info: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(w, http.StatusOK, nil)
	}
}

func GetCompany(app *app.App) http.HandlerFunc {
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

		company, err := app.CompSvc.GetById(r.Context(), idUuid)
		if err != nil {
			errorResponse(w, fmt.Errorf("getting company by id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(w, http.StatusOK, map[string]interface{}{"company": toCompanyTransport(company)})
	}
}

func ListEntrepreneurCompanies(app *app.App) http.HandlerFunc {
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

		companies, err := app.CompSvc.GetByOwnerId(r.Context(), entUuid, pageInt, true)
		if err != nil {
			errorResponse(w, fmt.Errorf("getting companies: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		companiesTransport := make([]Company, len(companies))
		for i, company := range companies {
			companiesTransport[i] = toCompanyTransport(company)
		}

		successResponse(w, http.StatusOK, map[string]interface{}{"entrepreneur_id": entId, "companies": companiesTransport})
	}
}

func CreateUserSkill(app *app.App) http.HandlerFunc {
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

		var req UserSkill
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			errorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		userSkill := toUserSkillModel(&req)
		userSkill.UserId = idUuid

		err = app.UserSkillSvc.Create(r.Context(), &userSkill)
		if err != nil {
			errorResponse(w, fmt.Errorf("creating user-skill pair: %w", err).Error(), http.StatusBadRequest)
			return
		}

		successResponse(w, http.StatusOK, nil)
	}
}

func DeleteUserSkill(app *app.App) http.HandlerFunc {
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

		err = app.UserSkillSvc.Delete(r.Context(), &domain.UserSkill{UserId: ownerIdUuid, SkillId: idUuid})
		if err != nil {
			errorResponse(w, fmt.Errorf("deleting user-skill pair: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(w, http.StatusOK, nil)
	}
}

func ListEntrepreneurSkills(app *app.App) http.HandlerFunc {
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
		if entId == "" {
			errorResponse(w, fmt.Errorf("empty entrepreneur id").Error(), http.StatusBadRequest)
			return
		}

		entUuid, err := uuid.Parse(entId)
		if err != nil {
			errorResponse(w, fmt.Errorf("converting entrepreneur id to uuid: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		skills, err := app.UserSkillSvc.GetSkillsForUser(r.Context(), entUuid, pageInt, true)
		if err != nil {
			errorResponse(w, fmt.Errorf("getting companies: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		skillsTransport := make([]Skill, len(skills))
		for i, skill := range skills {
			skillsTransport[i] = toSkillTransport(skill)
		}

		successResponse(w, http.StatusOK, map[string]interface{}{"entrepreneur_id": entId, "skills": skillsTransport})
	}
}

func CreateReport(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIdStr, err := getStringClaimFromJWT(r.Context(), "sub")
		if err != nil {
			errorResponse(w, fmt.Errorf("getting claim from JWT: %w", err).Error(), http.StatusBadRequest)
			return
		}

		userIdUuid, err := uuid.Parse(userIdStr)
		if err != nil {
			errorResponse(w, fmt.Errorf("converting string to uuid: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		compIdStr := chi.URLParam(r, "id")
		if compIdStr == "" {
			errorResponse(w, fmt.Errorf("empty company id").Error(), http.StatusBadRequest)
			return
		}

		compIdUuid, err := uuid.Parse(compIdStr)
		if err != nil {
			errorResponse(w, fmt.Errorf("converting string to uuid: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		company, err := app.CompSvc.GetById(r.Context(), compIdUuid)
		if err != nil {
			errorResponse(w, fmt.Errorf("creating fin report: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		if company.OwnerID != userIdUuid {
			errorResponse(w, fmt.Errorf("only company`s owner can create financial report").Error(), http.StatusInternalServerError)
			return
		}

		var req FinancialReport
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			errorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		report := toFinReportModel(&req)
		report.CompanyID = compIdUuid

		err = app.FinSvc.Create(r.Context(), &report)
		if err != nil {
			errorResponse(w, fmt.Errorf("creating financial report: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(w, http.StatusOK, nil)
	}
}

func DeleteFinReport(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIdStr, err := getStringClaimFromJWT(r.Context(), "sub")
		if err != nil {
			errorResponse(w, fmt.Errorf("getting claim from JWT: %w", err).Error(), http.StatusBadRequest)
			return
		}

		userIdUuid, err := uuid.Parse(userIdStr)
		if err != nil {
			errorResponse(w, fmt.Errorf("converting string to uuid: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		reportIdStr := chi.URLParam(r, "id")
		if reportIdStr == "" {
			errorResponse(w, fmt.Errorf("empty report id").Error(), http.StatusBadRequest)
			return
		}

		reportIdUuid, err := uuid.Parse(reportIdStr)
		if err != nil {
			errorResponse(w, fmt.Errorf("converting string to uuid: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		report, err := app.FinSvc.GetById(r.Context(), reportIdUuid)
		if err != nil {
			errorResponse(w, fmt.Errorf("getting financial report: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		company, err := app.CompSvc.GetById(r.Context(), report.CompanyID)
		if err != nil {
			errorResponse(w, fmt.Errorf("getting financial report: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		if company.OwnerID != userIdUuid {
			errorResponse(w, fmt.Errorf("only company`s owner can delete financial report").Error(), http.StatusInternalServerError)
			return
		}

		err = app.FinSvc.DeleteById(r.Context(), reportIdUuid)
		if err != nil {
			errorResponse(w, fmt.Errorf("deleting company by id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(w, http.StatusOK, nil)
	}
}

func UpdateFinReport(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIdStr, err := getStringClaimFromJWT(r.Context(), "sub")
		if err != nil {
			errorResponse(w, fmt.Errorf("getting claim from JWT: %w", err).Error(), http.StatusBadRequest)
			return
		}

		userIdUuid, err := uuid.Parse(userIdStr)
		if err != nil {
			errorResponse(w, fmt.Errorf("converting string to uuid: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		reportIdStr := chi.URLParam(r, "id")
		if reportIdStr == "" {
			errorResponse(w, fmt.Errorf("empty report id").Error(), http.StatusBadRequest)
			return
		}

		reportIdUuid, err := uuid.Parse(reportIdStr)
		if err != nil {
			errorResponse(w, fmt.Errorf("converting string to uuid: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		reportDb, err := app.FinSvc.GetById(r.Context(), reportIdUuid)
		if err != nil {
			errorResponse(w, fmt.Errorf("getting financial report: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		company, err := app.CompSvc.GetById(r.Context(), reportDb.CompanyID)
		if err != nil {
			errorResponse(w, fmt.Errorf("getting financial report: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		if company.OwnerID != userIdUuid {
			errorResponse(w, fmt.Errorf("only company`s owner can update financial report").Error(), http.StatusInternalServerError)
			return
		}

		var req FinancialReport

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			errorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.Year != 0 {
			reportDb.Year = req.Year
		}
		if req.Quarter != 0 {
			reportDb.Quarter = req.Quarter
		}
		if !(math.Abs(float64(req.Revenue)) < eps) {
			reportDb.Revenue = req.Revenue
		}
		if !(math.Abs(float64(req.Costs)) < eps) {
			reportDb.Costs = req.Costs
		}

		err = app.FinSvc.Update(r.Context(), reportDb)
		if err != nil {
			errorResponse(w, fmt.Errorf("updating financial report info: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(w, http.StatusOK, nil)
	}
}

func GetFinReport(app *app.App) http.HandlerFunc {
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

		report, err := app.FinSvc.GetById(r.Context(), idUuid)
		if err != nil {
			errorResponse(w, fmt.Errorf("getting financial report by id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(w, http.StatusOK, map[string]interface{}{"financial_report": toFinReportTransport(report)})
	}
}

func ListCompanyReports(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//page := r.URL.Query().Get("page")
		//if page == "" {
		//	errorResponse(w, fmt.Errorf("empty page number").Error(), http.StatusBadRequest)
		//	return
		//}
		//
		//pageInt, err := strconv.Atoi(page)
		//if err != nil {
		//	errorResponse(w, fmt.Errorf("converting page to int: %w", err).Error(), http.StatusBadRequest)
		//	return
		//}

		period, err := parsePeriodFromURL(r)
		if err != nil {
			errorResponse(w, fmt.Errorf("parsing period from URL: %w", err).Error(), http.StatusBadRequest)
			return
		}

		//compIdStr := chi.URLParam(r, "id")
		//if compIdStr == "" {
		//	errorResponse(w, fmt.Errorf("empty company id").Error(), http.StatusBadRequest)
		//	return
		//}
		//
		//compIdUuid, err := uuid.Parse(compIdStr)
		//if err != nil {
		//	errorResponse(w, fmt.Errorf("converting company id to uuid: %w", err).Error(), http.StatusInternalServerError)
		//	return
		//}
		compIdUuid, err := parseUUIDFromURL(r, "id", "company")
		if err != nil {
			errorResponse(w, fmt.Errorf("parsing company id from url: %w", err).Error(), http.StatusBadRequest)
			return
		}

		reports, err := app.FinSvc.GetByCompany(r.Context(), compIdUuid, period)
		if err != nil {
			errorResponse(w, fmt.Errorf("getting companies: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		reportsTransport := make([]FinancialReport, len(reports.Reports))
		for i, rep := range reports.Reports {
			reportsTransport[i] = toFinReportTransport(&rep)
		}

		successResponse(w, http.StatusOK,
			map[string]interface{}{
				"company_id": compIdUuid,
				"period":     toPeriodTransport(period),
				"revenue":    reports.Revenue(),
				"costs":      reports.Costs(),
				"profit":     reports.Profit(),
				"reports":    reportsTransport},
		)
	}
}
