package web

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"ppo/domain"
	"ppo/internal/app"
	"ppo/pkg/base"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func LoginHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "LoginHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		type Req struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}
		var req Req

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		ua := &domain.UserAuth{Username: req.Login, Password: req.Password}
		token, err := app.AuthSvc.Login(r.Context(), ua)
		if err != nil {
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusUnauthorized)
			return
		}

		_, err = base.VerifyAuthToken(token, app.Config.Server.JwtKey)
		if err != nil {
			app.Logger.Infof("%s: проверка JWT-токена: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: проверка JWT-токена: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		cookie := http.Cookie{
			Name:    "access_token",
			Value:   token,
			Path:    "/",
			Secure:  true,
			Expires: time.Now().Add(3600 * 24 * time.Second),
		}

		http.SetCookie(w, &cookie)
		successResponse(wrappedWriter, http.StatusOK, map[string]string{"token": token})
	}
}

func RegisterHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "RegisterHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		type Req struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}
		var req Req

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		ua := &domain.UserAuth{Username: req.Login, Password: req.Password}
		err = app.AuthSvc.Register(r.Context(), ua)
		if err != nil {
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		successResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func ListEntrepreneurs(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "ListEntrepreneursHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		page := r.URL.Query().Get("page")
		if page == "" {
			app.Logger.Infof("%s: пустой номер страницы", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("%s: пустой номер страницы", prompt).Error(), http.StatusBadRequest)
			return
		}

		pageInt, err := strconv.Atoi(page)
		if err != nil {
			app.Logger.Infof("%s: преобразование номера страницы к int: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: преобразование номера страницы к int: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		users, numPages, err := app.UserSvc.GetAll(r.Context(), pageInt)
		if err != nil {
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		usersTransport := make([]User, len(users))
		for i, user := range users {
			usersTransport[i] = toUserTransport(user)
		}

		successResponse(wrappedWriter, http.StatusOK, map[string]interface{}{"num_pages": numPages, "users": usersTransport})
	}
}

func UpdateEntrepreneur(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "UpdateEntrepreneurHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		id := chi.URLParam(r, "id")
		if id == "" {
			app.Logger.Infof("%s: пустой id", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("%s: пустой id", prompt).Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			app.Logger.Infof("%s: преобразование id к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id к uuid: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		userDb, err := app.UserSvc.GetById(r.Context(), idUuid)
		if err != nil {
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		var req User

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, err.Error(), http.StatusBadRequest)
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
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		successResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func DeleteEntrepreneur(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "DeleteEntrepreneurHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		id := chi.URLParam(r, "id")
		if id == "" {
			app.Logger.Infof("%s: пустой id", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("%s: пустой id", prompt).Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			app.Logger.Infof("%s: преобразование id к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id к uuid: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		_, err = app.UserSvc.GetById(r.Context(), idUuid)
		if err != nil {
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		err = app.UserSvc.DeleteById(r.Context(), idUuid)
		if err != nil {
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func GetEntrepreneur(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "GetEntrepreneurHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		id := chi.URLParam(r, "id")
		if id == "" {
			app.Logger.Infof("%s: пустой id", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("%s: пустой id", prompt).Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			app.Logger.Infof("%s: преобразование id к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id к uuid: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		user, err := app.UserSvc.GetById(r.Context(), idUuid)
		if err != nil {
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(wrappedWriter, http.StatusOK, map[string]interface{}{"entrepreneur": toUserTransport(user)})
	}
}

func CreateSkill(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "CreateSkillHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		var req Skill

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		skill := toSkillModel(&req)

		err = app.SkillSvc.Create(r.Context(), &skill)
		if err != nil {
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		successResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func DeleteSkill(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "DeleteSkillHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		id := chi.URLParam(r, "id")
		if id == "" {
			app.Logger.Infof("%s: пустой id", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("%s: пустой id", prompt).Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			app.Logger.Infof("%s: преобразование id к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id к uuid: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		_, err = app.SkillSvc.GetById(r.Context(), idUuid)
		if err != nil {
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		err = app.SkillSvc.DeleteById(r.Context(), idUuid)
		if err != nil {
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func UpdateSkill(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "UpdateSkillHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		id := chi.URLParam(r, "id")
		if id == "" {
			app.Logger.Infof("%s: пустой id", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("%s: пустой id", prompt).Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			app.Logger.Infof("%s: преобразование id к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id к uuid: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		skillDb, err := app.SkillSvc.GetById(r.Context(), idUuid)
		if err != nil {
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		var req Skill

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, err.Error(), http.StatusBadRequest)
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
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func ListSkills(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "ListSkillsHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		page := r.URL.Query().Get("page")
		if page == "" {
			app.Logger.Infof("%s: пустой номер страницы", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("%s: пустой номер страницы", prompt).Error(), http.StatusBadRequest)
			return
		}

		pageInt, err := strconv.Atoi(page)
		if err != nil {
			app.Logger.Infof("%s: преобразование номера страницы к int: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: преобразование номера страницы к int: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		skills, numPages, err := app.SkillSvc.GetAll(r.Context(), pageInt)
		if err != nil {
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		skillsTransport := make([]Skill, len(skills))
		for i, skill := range skills {
			skillsTransport[i] = toSkillTransport(skill)
		}

		successResponse(wrappedWriter, http.StatusOK, map[string]interface{}{"num_pages": numPages, "skills": skillsTransport})
	}
}

func GetSkill(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "GetSkillHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		id := chi.URLParam(r, "id")
		if id == "" {
			app.Logger.Infof("%s: пустой id", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("%s: пустой id", prompt).Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			app.Logger.Infof("%s: преобразование id к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id к uuid: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		skill, err := app.SkillSvc.GetById(r.Context(), idUuid)
		if err != nil {
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(wrappedWriter, http.StatusOK, map[string]interface{}{"skill": toSkillTransport(skill)})
	}
}

func CreateContact(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "CreateContactHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		idStr, err := getStringClaimFromJWT(r.Context(), "sub")
		if err != nil {
			app.Logger.Infof("%s: получение записей из JWT: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: получение записей из JWT: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(idStr)
		if err != nil {
			app.Logger.Infof("%s: преобразование id к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id к uuid: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		var req Contact
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, err.Error(), http.StatusBadRequest)
			return
		}

		contact := toContactModel(&req)
		contact.OwnerID = idUuid

		err = app.ConSvc.Create(r.Context(), &contact)
		if err != nil {
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		successResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func DeleteContact(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "DeleteContactHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		ownerIdStr, err := getStringClaimFromJWT(r.Context(), "sub")
		if err != nil {
			app.Logger.Infof("%s: получение записей из JWT: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("получение записей из JWT: %w", err).Error(), http.StatusBadRequest)
			return
		}

		ownerIdUuid, err := uuid.Parse(ownerIdStr)
		if err != nil {
			app.Logger.Infof("%s: преобразование строки к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование строки к uuid: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		id := chi.URLParam(r, "id")
		if id == "" {
			app.Logger.Infof("%s: пустой id", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("пустой id").Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			app.Logger.Infof("%s: преобразование строки к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование строки к uuid: %w", err).Error(), http.StatusBadRequest)
			return
		}

		contact, err := app.ConSvc.GetById(r.Context(), idUuid)
		if err != nil {
			app.Logger.Infof("%s: удаление средства связи по id: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("удаление средства связи по id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		if ownerIdUuid != contact.OwnerID {
			app.Logger.Infof("%s: только владелец может удалить своё средство связи", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("только владелец может удалить своё средство связи").Error(), http.StatusInternalServerError)
			return
		}

		err = app.ConSvc.DeleteById(r.Context(), idUuid)
		if err != nil {
			app.Logger.Infof("%s: удаление средства связи по id: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("удаление средства связи по id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func UpdateContact(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "UpdateContactHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		ownerIdStr, err := getStringClaimFromJWT(r.Context(), "sub")
		if err != nil {
			app.Logger.Infof("%s: получение записей из JWT: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("получение записей из JWT: %w", err).Error(), http.StatusBadRequest)
			return
		}

		ownerIdUuid, err := uuid.Parse(ownerIdStr)
		if err != nil {
			app.Logger.Infof("%s: преобразование строки к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование строки к uuid: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		id := chi.URLParam(r, "id")
		if id == "" {
			app.Logger.Infof("%s: пустой id", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("пустой id").Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			app.Logger.Infof("%s: преобразование id к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование id к uuid: %w", err).Error(), http.StatusBadRequest)
			return
		}

		conDb, err := app.ConSvc.GetById(r.Context(), idUuid)
		if err != nil {
			app.Logger.Infof("%s: получение средства связи по id: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("получение средства связи по id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		if ownerIdUuid != conDb.OwnerID {
			app.Logger.Infof("%s: только владелец может обновлять информацию о своих средствах связи", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("только владелец может обновлять информацию о своих средствах связи").Error(), http.StatusInternalServerError)
			return
		}

		var req Contact

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			app.Logger.Infof("%s: %v", err)
			errorResponse(wrappedWriter, err.Error(), http.StatusBadRequest)
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
			app.Logger.Infof("%s: обновление информации о средстве связи: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("обновление информации о средстве связи: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func GetContact(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "GetContactHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		id := chi.URLParam(r, "id")
		if id == "" {
			app.Logger.Infof("%s: пустой id", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("пустой id").Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			app.Logger.Infof("%s: преобразование id к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование id к uuid: %w", err).Error(), http.StatusBadRequest)
			return
		}

		contact, err := app.ConSvc.GetById(r.Context(), idUuid)
		if err != nil {
			app.Logger.Infof("%s: получение средства связи по id: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("получение средства связи по id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(wrappedWriter, http.StatusOK, map[string]interface{}{"contact": toContactTransport(contact)})
	}
}

func ListEntrepreneurContacts(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "ListEntrepreneursContactsHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		entId := r.URL.Query().Get("entrepreneur-id")
		if entId == "" {
			app.Logger.Infof("%s: пустой id предпринимателя", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("пустой id предпринимателя").Error(), http.StatusBadRequest)
			return
		}

		entUuid, err := uuid.Parse(entId)
		if err != nil {
			app.Logger.Infof("%s: преобразование id к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование id к uuid: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		contacts, err := app.ConSvc.GetByOwnerId(r.Context(), entUuid)
		if err != nil {
			app.Logger.Infof("%s: получение списка контактов: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("получение списка контактов: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		contactsTransport := make([]Contact, len(contacts))
		for i, contact := range contacts {
			contactsTransport[i] = toContactTransport(contact)
		}

		successResponse(wrappedWriter, http.StatusOK, map[string]interface{}{"entrepreneur_id": entId, "contacts": contactsTransport})
	}
}

func CreateActivityField(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "CreateActivityFieldHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		var req ActivityField
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, err.Error(), http.StatusBadRequest)
			return
		}

		actField := toActFieldModel(&req)

		err = app.ActFieldSvc.Create(r.Context(), &actField)
		if err != nil {
			app.Logger.Infof("%s: создание сферы деятельности: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("создание сферы деятельности: %w", err).Error(), http.StatusBadRequest)
			return
		}

		successResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func DeleteActivityField(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "DeleteActivityFieldHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		id := chi.URLParam(r, "id")
		if id == "" {
			app.Logger.Infof("%s: пустой id", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("пустой id").Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			app.Logger.Infof("%s: преобразование id к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование id к uuid: %w", err).Error(), http.StatusBadRequest)
			return
		}

		_, err = app.ActFieldSvc.GetById(r.Context(), idUuid)
		if err != nil {
			app.Logger.Infof("%s: получение сферы деятельности по id: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("получение сферы деятельности по id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		err = app.ActFieldSvc.DeleteById(r.Context(), idUuid)
		if err != nil {
			app.Logger.Infof("%s: удаление сферы деятельности по id: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("удаление сферы деятельности по id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func UpdateActivityField(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "UpdateActivityFieldHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		id := chi.URLParam(r, "id")
		if id == "" {
			app.Logger.Infof("%s: пустой id", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("пустой id").Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			app.Logger.Infof("%s: преобразование id к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование id к uuid: %w", err).Error(), http.StatusBadRequest)
			return
		}

		actFieldDb, err := app.ActFieldSvc.GetById(r.Context(), idUuid)
		if err != nil {
			app.Logger.Infof("%s: получение сферы деятельности по id: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("получение сферы деятельности по id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		var req ActivityField

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, err.Error(), http.StatusBadRequest)
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
			app.Logger.Infof("%s: обновление информации о сфере деятельности: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("обновление информации о сфере деятельности: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func GetActivityField(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "GetActivityFieldHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		id := chi.URLParam(r, "id")
		if id == "" {
			app.Logger.Infof("%s: пустой id", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("пустой id").Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			app.Logger.Infof("%s: преобразование id к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование id к uuid: %w", err).Error(), http.StatusBadRequest)
			return
		}

		actField, err := app.ActFieldSvc.GetById(r.Context(), idUuid)
		if err != nil {
			app.Logger.Infof("%s: получение сферы деятельности по id: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("получение сферы деятельности по id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(wrappedWriter, http.StatusOK, map[string]interface{}{"activity_field": toActFieldTransport(actField)})
	}
}

func ListActivityFields(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "ListActivityFieldsHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		var paginated bool
		var pageInt int
		var err error

		page := r.URL.Query().Get("page")
		if page != "" {
			paginated = true

			pageInt, err = strconv.Atoi(page)
			if err != nil {
				app.Logger.Infof("%s: преобразование страницы к int: %v", prompt, err)
				errorResponse(wrappedWriter, fmt.Errorf("преобразование страницы к int: %w", err).Error(), http.StatusBadRequest)
				return
			}
		}

		actFields, numPages, err := app.ActFieldSvc.GetAll(r.Context(), pageInt, paginated)
		if err != nil {
			app.Logger.Infof("%s: получение списка сфер деятельности: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("получение списка сфер деятельности: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		actFieldsTransport := make([]ActivityField, len(actFields))
		for i, actField := range actFields {
			actFieldsTransport[i] = toActFieldTransport(actField)
		}

		successResponse(wrappedWriter, http.StatusOK, map[string]interface{}{"activity_fields": actFieldsTransport, "num_pages": numPages})
	}
}

func CreateCompany(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "CreateCompanyHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		idStr, err := getStringClaimFromJWT(r.Context(), "sub")
		if err != nil {
			app.Logger.Infof("%s: получение записей из JWT: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("получение записей из JWT: %w", err).Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(idStr)
		if err != nil {
			app.Logger.Infof("%s: преобразование строки к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование строки к uuid: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		var req Company
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, err.Error(), http.StatusBadRequest)
			return
		}

		company := toCompanyModel(&req)
		company.OwnerID = idUuid

		err = app.CompSvc.Create(r.Context(), &company)
		if err != nil {
			app.Logger.Infof("%s: создание компании: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("создание компании: %w", err).Error(), http.StatusBadRequest)
			return
		}

		successResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func DeleteCompany(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "DeleteCompanyHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		ownerIdStr, err := getStringClaimFromJWT(r.Context(), "sub")
		if err != nil {
			app.Logger.Infof("%s: получение записей из JWT: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("получение записей из JWT: %w", err).Error(), http.StatusBadRequest)
			return
		}

		ownerIdUuid, err := uuid.Parse(ownerIdStr)
		if err != nil {
			app.Logger.Infof("%s: преобразование строки к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование строки к uuid: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		id := chi.URLParam(r, "id")
		if id == "" {
			app.Logger.Infof("%s: пустой id", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("пустой id").Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			app.Logger.Infof("%s: преобразование id к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование id к uuid: %w", err).Error(), http.StatusBadRequest)
			return
		}

		company, err := app.CompSvc.GetById(r.Context(), idUuid)
		if err != nil {
			app.Logger.Infof("%s: удаление компании по id: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("удаление компании по id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		if ownerIdUuid != company.OwnerID {
			app.Logger.Infof("%s: только владелец может удалять свои компании", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("только владелец может удалять свои компании").Error(), http.StatusInternalServerError)
			return
		}

		err = app.CompSvc.DeleteById(r.Context(), idUuid)
		if err != nil {
			app.Logger.Infof("%s: удаление компании по id: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("удаление компании по id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func UpdateCompany(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "UpdateCompanyHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		ownerIdStr, err := getStringClaimFromJWT(r.Context(), "sub")
		if err != nil {
			app.Logger.Infof("%s: получение записей из JWT: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("получение записей из JWT: %w", err).Error(), http.StatusBadRequest)
			return
		}

		ownerIdUuid, err := uuid.Parse(ownerIdStr)
		if err != nil {
			app.Logger.Infof("%s: преобразование строки к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование строки к uuid: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		id := chi.URLParam(r, "id")
		if id == "" {
			app.Logger.Infof("%s: пустой id", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("пустой id").Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			app.Logger.Infof("%s: преобразование id к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование id к uuid: %w", err).Error(), http.StatusBadRequest)
			return
		}

		compDb, err := app.CompSvc.GetById(r.Context(), idUuid)
		if err != nil {
			app.Logger.Infof("%s: получение компании по id: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("получение компании по id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		if ownerIdUuid != compDb.OwnerID {
			app.Logger.Infof("%s: только владелец может обновлять информацию о своих компаниях", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("только владелец может обновлять информацию о своих компаниях").Error(), http.StatusInternalServerError)
			return
		}

		var req Company

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, err.Error(), http.StatusBadRequest)
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
			app.Logger.Infof("%s: обновление информации о компании: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("обновление информации о компании: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func GetCompany(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "GetCompanyHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		id := chi.URLParam(r, "id")
		if id == "" {
			app.Logger.Infof("%s: пустой id", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("пустой id").Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			app.Logger.Infof("%s: преобразование id к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование id к uuid: %w", err).Error(), http.StatusBadRequest)
			return
		}

		company, err := app.CompSvc.GetById(r.Context(), idUuid)
		if err != nil {
			app.Logger.Infof("%s: получение компании по id: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("получение компании по id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(wrappedWriter, http.StatusOK, map[string]interface{}{"company": toCompanyTransport(company)})
	}
}

func ListEntrepreneurCompanies(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "ListEntrepreneurCompaniesHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		page := r.URL.Query().Get("page")
		if page == "" {
			app.Logger.Infof("%s: пустой номер страницы", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("пустой номер страницы").Error(), http.StatusBadRequest)
			return
		}

		pageInt, err := strconv.Atoi(page)
		if err != nil {
			app.Logger.Infof("%s: преобразование к int: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование к int: %w", err).Error(), http.StatusBadRequest)
			return
		}

		entId := r.URL.Query().Get("entrepreneur-id")
		if page == "" {
			app.Logger.Infof("%s: пустой id предпринимателя", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("пустой id предпринимателя").Error(), http.StatusBadRequest)
			return
		}

		entUuid, err := uuid.Parse(entId)
		if err != nil {
			app.Logger.Infof("%s: преобразование id предпринимателя к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование id предпринимателя к uuid: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		companies, numPages, err := app.CompSvc.GetByOwnerId(r.Context(), entUuid, pageInt, true)
		if err != nil {
			app.Logger.Infof("%s: получение списка компаний: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("получение списка компаний: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		companiesTransport := make([]Company, len(companies))
		for i, company := range companies {
			companiesTransport[i] = toCompanyTransport(company)
		}

		successResponse(wrappedWriter, http.StatusOK, map[string]interface{}{"entrepreneur_id": entId, "companies": companiesTransport, "num_pages": numPages})
	}
}

func CreateUserSkill(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "CreateUserSkillHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		idStr, err := getStringClaimFromJWT(r.Context(), "sub")
		if err != nil {
			app.Logger.Infof("%s: получение записей из JWT: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("получение записей из JWT: %w", err).Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(idStr)
		if err != nil {
			app.Logger.Infof("%s: преобразование строки к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование строки к uuid: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		var req UserSkill
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, err.Error(), http.StatusBadRequest)
			return
		}

		userSkill := toUserSkillModel(&req)
		userSkill.UserId = idUuid

		err = app.UserSkillSvc.Create(r.Context(), &userSkill)
		if err != nil {
			app.Logger.Infof("%s: создание пары навык-предприниматель: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("создание пары навык-предприниматель: %w", err).Error(), http.StatusBadRequest)
			return
		}

		successResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func DeleteUserSkill(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "DeleteUserSkillHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		ownerIdStr, err := getStringClaimFromJWT(r.Context(), "sub")
		if err != nil {
			app.Logger.Infof("%s: создание пары навык-предприниматель: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("получение записей из JWT: %w", err).Error(), http.StatusBadRequest)
			return
		}

		ownerIdUuid, err := uuid.Parse(ownerIdStr)
		if err != nil {
			app.Logger.Infof("%s: преобразование строки к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование строки к uuid: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		id := chi.URLParam(r, "id")
		if id == "" {
			app.Logger.Infof("%s: пустой id", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("пустой id").Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			app.Logger.Infof("%s: преобразование id к uuid", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование id к uuid: %w", err).Error(), http.StatusBadRequest)
			return
		}

		err = app.UserSkillSvc.Delete(r.Context(), &domain.UserSkill{UserId: ownerIdUuid, SkillId: idUuid})
		if err != nil {
			app.Logger.Infof("%s: удаление пары навык-предприниматель: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("удаление пары навык-предприниматель: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func ListEntrepreneurSkills(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "ListEntrepreneurSkillsHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		page := r.URL.Query().Get("page")
		if page == "" {
			app.Logger.Infof("%s: пустой номер страницы", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("пустой номер страницы").Error(), http.StatusBadRequest)
			return
		}

		pageInt, err := strconv.Atoi(page)
		if err != nil {
			app.Logger.Infof("%s: преобразование страницы к int: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование страницы к int: %w", err).Error(), http.StatusBadRequest)
			return
		}

		entId := r.URL.Query().Get("entrepreneur-id")
		if entId == "" {
			app.Logger.Infof("%s: пустой id предпринимателя", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("пустой id предпринимателя").Error(), http.StatusBadRequest)
			return
		}

		entUuid, err := uuid.Parse(entId)
		if err != nil {
			app.Logger.Infof("%s: преобразование id предпринимателя к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование id предпринимателя к uuid: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		skills, numPages, err := app.UserSkillSvc.GetSkillsForUser(r.Context(), entUuid, pageInt, true)
		if err != nil {
			app.Logger.Infof("%s: получение списка компаний: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("получение списка компаний: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		skillsTransport := make([]Skill, len(skills))
		for i, skill := range skills {
			skillsTransport[i] = toSkillTransport(skill)
		}

		successResponse(wrappedWriter, http.StatusOK, map[string]interface{}{"entrepreneur_id": entId, "skills": skillsTransport, "num_pages": numPages})
	}
}

func CreateReport(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "CreateReportHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		userIdStr, err := getStringClaimFromJWT(r.Context(), "sub")
		if err != nil {
			app.Logger.Infof("%s: получение списка компаний: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("получение записей из JWT: %w", err).Error(), http.StatusBadRequest)
			return
		}

		userIdUuid, err := uuid.Parse(userIdStr)
		if err != nil {
			app.Logger.Infof("%s: преобразование строки к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование строки к uuid: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		compIdStr := chi.URLParam(r, "id")
		if compIdStr == "" {
			app.Logger.Infof("%s: пустой id компании", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("пустой id компании").Error(), http.StatusBadRequest)
			return
		}

		compIdUuid, err := uuid.Parse(compIdStr)
		if err != nil {
			app.Logger.Infof("%s: преобразование строки к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование строки к uuid: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		company, err := app.CompSvc.GetById(r.Context(), compIdUuid)
		if err != nil {
			app.Logger.Infof("%s: создание финансового отчета: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("создание финансового отчета: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		if company.OwnerID != userIdUuid {
			app.Logger.Infof("%s: только владелец компании может добавлять финансовые отчеты", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("только владелец компании может добавлять финансовые отчеты").Error(), http.StatusInternalServerError)
			return
		}

		var req FinancialReport
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			app.Logger.Infof("%s: %м", prompt, err)
			errorResponse(wrappedWriter, err.Error(), http.StatusBadRequest)
			return
		}

		report := toFinReportModel(&req)
		report.CompanyID = compIdUuid

		err = app.FinSvc.Create(r.Context(), &report)
		if err != nil {
			app.Logger.Infof("%s: создание финансового отчета: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("создание финансового отчета: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func DeleteFinReport(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "DeleteFinReportHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		userIdStr, err := getStringClaimFromJWT(r.Context(), "sub")
		if err != nil {
			app.Logger.Infof("%s: получение записей из JWT: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("получение записей из JWT: %w", err).Error(), http.StatusBadRequest)
			return
		}

		userIdUuid, err := uuid.Parse(userIdStr)
		if err != nil {
			app.Logger.Infof("%s: преобразование строки к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование строки к uuid: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		reportIdStr := chi.URLParam(r, "id")
		if reportIdStr == "" {
			app.Logger.Infof("%s: пустой id отчета", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("пустой id отчета").Error(), http.StatusBadRequest)
			return
		}

		reportIdUuid, err := uuid.Parse(reportIdStr)
		if err != nil {
			app.Logger.Infof("%s: преобразование строки к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование строки к uuid: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		report, err := app.FinSvc.GetById(r.Context(), reportIdUuid)
		if err != nil {
			app.Logger.Infof("%s: получение финансового отчета: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("получение финансового отчета: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		company, err := app.CompSvc.GetById(r.Context(), report.CompanyID)
		if err != nil {
			app.Logger.Infof("%s: получение компании: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("получение компании: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		if company.OwnerID != userIdUuid {
			app.Logger.Infof("%s: только владелец компании может удалять финансовые отчеты", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("только владелец компании может удалять финансовые отчеты").Error(), http.StatusInternalServerError)
			return
		}

		err = app.FinSvc.DeleteById(r.Context(), reportIdUuid)
		if err != nil {
			app.Logger.Infof("%s: только владелец компании может удалять финансовые отчеты", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("удаление финансового отчета по id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func UpdateFinReport(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "UpdateFinReportHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		userIdStr, err := getStringClaimFromJWT(r.Context(), "sub")
		if err != nil {
			app.Logger.Infof("%s: получение записей из JWT: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("получение записей из JWT: %w", err).Error(), http.StatusBadRequest)
			return
		}

		userIdUuid, err := uuid.Parse(userIdStr)
		if err != nil {
			app.Logger.Infof("%s: преобразование строки к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование строки к uuid: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		reportIdStr := chi.URLParam(r, "id")
		if reportIdStr == "" {
			app.Logger.Infof("%s: пустой id отчета", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("пустой id отчета").Error(), http.StatusBadRequest)
			return
		}

		reportIdUuid, err := uuid.Parse(reportIdStr)
		if err != nil {
			app.Logger.Infof("%s: преобразование строки к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование строки к uuid: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		reportDb, err := app.FinSvc.GetById(r.Context(), reportIdUuid)
		if err != nil {
			app.Logger.Infof("%s: получение финансового отчета: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("получение финансового отчета: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		company, err := app.CompSvc.GetById(r.Context(), reportDb.CompanyID)
		if err != nil {
			app.Logger.Infof("%s: получение компании: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("получение компании: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		if company.OwnerID != userIdUuid {
			app.Logger.Infof("%s: только владелец компании может изменять финансовый отчет", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("только владелец компании может изменять финансовый отчет").Error(), http.StatusInternalServerError)
			return
		}

		var req FinancialReport

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, err.Error(), http.StatusBadRequest)
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
			app.Logger.Infof("%s: обновление информации о финансовом отчете: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("обновление информации о финансовом отчете: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func GetFinReport(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "GetFinReportHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		id := chi.URLParam(r, "id")
		if id == "" {
			app.Logger.Infof("%s: пустой id", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("пустой id").Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			app.Logger.Infof("%s: преобразование id к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование id к uuid: %w", err).Error(), http.StatusBadRequest)
			return
		}

		report, err := app.FinSvc.GetById(r.Context(), idUuid)
		if err != nil {
			app.Logger.Infof("%s: получение финансового отчета по id: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("получение финансового отчета по id: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(wrappedWriter, http.StatusOK, map[string]interface{}{"financial_report": toFinReportTransport(report)})
	}
}

func ListCompanyReports(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "ListCompanyReportsHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		period, err := parsePeriodFromURL(r)
		if err != nil {
			app.Logger.Infof("%s: парсинг периода из URL: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("парсинг периода из URL: %w", err).Error(), http.StatusBadRequest)
			return
		}

		compIdUuid, err := parseUUIDFromURL(r, "id", "company")
		if err != nil {
			app.Logger.Infof("%s: парсинг id компании из URL: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("парсинг id компании из URL: %w", err).Error(), http.StatusBadRequest)
			return
		}

		reports, err := app.FinSvc.GetByCompany(r.Context(), compIdUuid, period)
		if err != nil {
			app.Logger.Infof("%s: получение отчетов компании: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("получение отчетов компании: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		reportsTransport := make([]FinancialReport, len(reports.Reports))
		for i, rep := range reports.Reports {
			reportsTransport[i] = toFinReportTransport(&rep)
		}

		successResponse(wrappedWriter, http.StatusOK,
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

func CalculateRating(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "CalculateRatingHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		id := chi.URLParam(r, "id")
		if id == "" {
			app.Logger.Infof("%s: пустой id", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("пустой id").Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			app.Logger.Infof("%s: преобразование id к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование id к uuid: %w", err).Error(), http.StatusBadRequest)
			return
		}

		rating, err := app.Interactor.CalculateUserRating(r.Context(), idUuid)
		if err != nil {
			app.Logger.Infof("%s: вычисление рейтинга предпринимателя: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("вычисление рейтинга предпринимателя: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(wrappedWriter, http.StatusOK, map[string]float32{"rating": rating})
	}
}

func GetEntrepreneurFinancials(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "GetEntrepreneurFinancials"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		id := r.URL.Query().Get("entrepreneur-id")
		if id == "" {
			app.Logger.Infof("%s: пустой id предпринимателя", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("пустой id предпринимателя").Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			app.Logger.Infof("%s: преобразование id к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование id к uuid: %w", err).Error(), http.StatusBadRequest)
			return
		}

		prevYear := time.Now().AddDate(-1, 0, 0).Year()
		period := &domain.Period{
			StartYear:    prevYear,
			EndYear:      prevYear,
			StartQuarter: 1,
			EndQuarter:   4,
		}

		rep, err := app.Interactor.GetUserFinancialReport(r.Context(), idUuid, period)
		if err != nil {
			app.Logger.Infof("%s: получение финансового отчета предпринимателя: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("получение финансового отчета предпринимателя: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(wrappedWriter, http.StatusOK, map[string]float32{
			"revenue": rep.Revenue(),
			"costs":   rep.Costs(),
			"profit":  rep.Profit(),
			"taxes":   rep.Taxes,
			"taxLoad": rep.TaxLoad,
		})
	}
}

func GetEntrepreneurReviews(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "GetEntrepreneurReviewsHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		id := r.URL.Query().Get("entrepreneur-id")
		if id == "" {
			app.Logger.Infof("%s: пустой id предпринимателя", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("пустой id предпринимателя").Error(), http.StatusBadRequest)
			return
		}

		entUuid, err := uuid.Parse(id)
		if err != nil {
			app.Logger.Infof("%s: преобразование id к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование id к uuid: %w", err).Error(), http.StatusBadRequest)
			return
		}

		page := r.URL.Query().Get("page")
		if page == "" {
			app.Logger.Infof("%s: пустой номер страницы", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("пустой номер страницы").Error(), http.StatusBadRequest)
			return
		}

		pageInt, err := strconv.Atoi(page)
		if err != nil {
			app.Logger.Infof("%s: преобразование страницы к int: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование страницы к int: %w", err).Error(), http.StatusBadRequest)
			return
		}

		revs, numPages, err := app.RevSvc.GetAllForTarget(r.Context(), entUuid, pageInt)
		if err != nil {
			app.Logger.Infof("%s: получение отзывов: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("получение отзывов: %w", err).Error(), http.StatusBadRequest)
			return
		}

		reviewsTransport := make([]Review, len(revs))
		for i, rev := range revs {
			reviewsTransport[i] = toReviewTransport(rev)
		}

		successResponse(wrappedWriter, http.StatusOK, map[string]interface{}{"entrepreneur_id": entUuid, "reviews": reviewsTransport, "num_pages": numPages})
	}
}

func GetAuthorReviews(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "GetAuthorReviewsHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		idStr, err := getStringClaimFromJWT(r.Context(), "sub")
		if err != nil {
			app.Logger.Infof("%s: получение записей из JWT: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: получение записей из JWT: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		entUuid, err := uuid.Parse(idStr)
		if err != nil {
			app.Logger.Infof("%s: преобразование id к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id к uuid: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		page := r.URL.Query().Get("page")
		if page == "" {
			app.Logger.Infof("%s: пустой номер страницы", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("пустой номер страницы").Error(), http.StatusBadRequest)
			return
		}

		pageInt, err := strconv.Atoi(page)
		if err != nil {
			app.Logger.Infof("%s: преобразование страницы к int: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("преобразование страницы к int: %w", err).Error(), http.StatusBadRequest)
			return
		}

		revs, numPages, err := app.RevSvc.GetAllForReviewer(r.Context(), entUuid, pageInt)

		reviewsTransport := make([]Review, len(revs))
		for i, rev := range revs {
			reviewsTransport[i] = toReviewTransport(rev)
		}

		successResponse(wrappedWriter, http.StatusOK, map[string]interface{}{"entrepreneur_id": entUuid, "reviews": reviewsTransport, "num_pages": numPages})
	}
}

func CreateReview(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "CreateReviewHandler"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		idStr, err := getStringClaimFromJWT(r.Context(), "sub")
		if err != nil {
			app.Logger.Infof("%s: получение записей из JWT: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: получение записей из JWT: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(idStr)
		if err != nil {
			app.Logger.Infof("%s: преобразование id к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id к uuid: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		var req Review
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, err.Error(), http.StatusBadRequest)
			return
		}

		rev := toReviewModel(&req)
		rev.Reviewer = idUuid

		err = app.RevSvc.Create(r.Context(), &rev)
		if err != nil {
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		successResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func DeleteReview(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "удаление отзыва"
		start := time.Now()

		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			observeRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		}()

		id := chi.URLParam(r, "id")
		if id == "" {
			app.Logger.Infof("%s: пустой id", prompt)
			errorResponse(wrappedWriter, fmt.Errorf("%s: пустой id", prompt).Error(), http.StatusBadRequest)
			return
		}

		idUuid, err := uuid.Parse(id)
		if err != nil {
			app.Logger.Infof("%s: преобразование id к uuid: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id к uuid: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		_, err = app.RevSvc.Get(r.Context(), idUuid)
		if err != nil {
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		err = app.RevSvc.Delete(r.Context(), idUuid)
		if err != nil {
			app.Logger.Infof("%s: %v", prompt, err)
			errorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		successResponse(wrappedWriter, http.StatusOK, nil)
	}
}
