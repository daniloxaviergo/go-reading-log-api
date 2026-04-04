package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"go-reading-log-api-next/internal/api/v1/handlers"
	"go-reading-log-api-next/internal/repository"
)

// SetupRoutes registers all API routes with the given router
func SetupRoutes(repo repository.ProjectRepository, logRepo repository.LogRepository) http.Handler {
	r := mux.NewRouter()

	// Create handlers
	projectsHandler := handlers.NewProjectsHandler(repo)
	logsHandler := handlers.NewLogsHandler(logRepo, repo)
	healthHandler := handlers.NewHealthHandler()

	// Health check
	r.HandleFunc("/healthz", healthHandler.Healthz).Methods("GET")

	// Projects endpoints
	r.HandleFunc("/api/v1/projects", projectsHandler.Index).Methods("GET")
	r.HandleFunc("/api/v1/projects", projectsHandler.Create).Methods("POST")
	r.HandleFunc("/api/v1/projects/{id}", projectsHandler.Show).Methods("GET")

	// Logs endpoints
	r.HandleFunc("/api/v1/projects/{project_id}/logs", logsHandler.Index).Methods("GET")
	// Phase 1: Read-only - logs creation will be added in Phase 2

	return r
}
