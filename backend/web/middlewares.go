package web

import (
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

func ValidateAdminRoleJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			errorResponse(w, fmt.Errorf("getting claims from JWT: %w", err).Error(), http.StatusBadRequest)
			return
		}

		role, ok := claims["role"]
		if !ok {
			errorResponse(w, fmt.Errorf("получение 'role' claim`а из JWT").Error(), http.StatusBadRequest)
			return
		}

		if role != "admin" {
			errorResponse(w, fmt.Errorf("только администраторы могут делать это").Error(), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ValidateUserRoleJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			errorResponse(w, fmt.Errorf("getting claims from JWT: %w", err).Error(), http.StatusBadRequest)
			return
		}

		role, ok := claims["role"]
		if !ok {
			errorResponse(w, fmt.Errorf("получение 'role' claim`а из JWT").Error(), http.StatusBadRequest)
			return
		}

		if role != "user" && role != "admin" {
			errorResponse(w, fmt.Errorf("вам нужно авторизоваться, прежде чем сделать это").Error(), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
