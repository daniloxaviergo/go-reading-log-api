package handlers

import (
	"encoding/json"
	"net/http"

	"go-reading-log-api-next/internal/domain/dto"
)

// HealthHandler handles HTTP requests for health check endpoint
type HealthHandler struct{}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Healthz returns a health check response
func (h *HealthHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	response := dto.NewHealthCheckResponse("healthy")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
