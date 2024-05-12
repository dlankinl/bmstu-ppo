package web

import (
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
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
			errorResponse(w, fmt.Errorf("getting 'role' claim from JWT: %w", err).Error(), http.StatusBadRequest)
			return
		}

		if role != "admin" {
			errorResponse(w, fmt.Errorf("only administrators can use this").Error(), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
