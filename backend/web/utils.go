package web

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
)

const (
	errorMsg   = "error"
	successMsg = "success"
)

const eps = 1e-6

type ErrorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type SuccessResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

func errorResponse(w http.ResponseWriter, err string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Status: errorMsg, Error: err})
}

func successResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(SuccessResponse{Status: successMsg, Data: data})
}

func getStringClaimFromJWT(ctx context.Context, claim string) (strVal string, err error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return "", fmt.Errorf("getting claims from JWT: %w", err)
	}

	id, ok := claims[claim]
	if !ok {
		return "", fmt.Errorf("failed getting claim '%s' from JWT token", claim)
	}

	strVal, ok = id.(string)
	if !ok {
		return "", fmt.Errorf("converting interface to string")
	}

	return strVal, nil
}
