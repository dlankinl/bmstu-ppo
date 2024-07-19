package web

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"net/http"
	"ppo/domain"
	"strconv"
)

const (
	errorMsg   = "error"
	successMsg = "success"
)

const eps = 1e-6

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusResponseWriter) StatusCode() int {
	return w.statusCode
}

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

func parsePeriodFromURL(r *http.Request) (period *domain.Period, err error) {
	yearStartStr := chi.URLParam(r, "year-start")
	if yearStartStr == "" {
		return nil, fmt.Errorf("empty year start")
	}

	yearStart, err := strconv.Atoi(yearStartStr)
	if err != nil {
		return nil, fmt.Errorf("converting start year to int: %w", err)
	}

	yearEndStr := chi.URLParam(r, "year-end")
	if yearEndStr == "" {
		return nil, fmt.Errorf("empty year end")
	}

	yearEnd, err := strconv.Atoi(yearEndStr)
	if err != nil {
		return nil, fmt.Errorf("converting end year to int: %w", err)
	}

	quarterStartStr := chi.URLParam(r, "quarter-start")
	if quarterStartStr == "" {
		return nil, fmt.Errorf("empty quarter start")
	}

	quarterStart, err := strconv.Atoi(quarterStartStr)
	if err != nil {
		return nil, fmt.Errorf("converting start quarter to int: %w, err")
	}

	quarterEndStr := chi.URLParam(r, "quarter-end")
	if quarterEndStr == "" {
		return nil, fmt.Errorf("empty quarter end")
	}

	quarterEnd, err := strconv.Atoi(quarterEndStr)
	if err != nil {
		return nil, fmt.Errorf("converting end quarter to int: %w", err)
	}

	period = &domain.Period{
		StartYear:    yearStart,
		StartQuarter: quarterStart,
		EndYear:      yearEnd,
		EndQuarter:   quarterEnd,
	}

	return period, nil
}

func parseUUIDFromURL(r *http.Request, key, entityName string) (val uuid.UUID, err error) {
	compIdStr := chi.URLParam(r, key)
	if compIdStr == "" {
		return uuid.UUID{}, fmt.Errorf("empty %s %s", entityName, key)
	}

	val, err = uuid.Parse(compIdStr)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("converting %s %s to uuid: %w", entityName, key, err)
	}

	return val, nil
}
