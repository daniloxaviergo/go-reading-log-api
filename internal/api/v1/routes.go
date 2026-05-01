package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"go-reading-log-api-next/internal/api/v1/handlers"
	"go-reading-log-api-next/internal/repository"
	"go-reading-log-api-next/internal/service"
	"go-reading-log-api-next/internal/service/dashboard"
)

// SetupRoutes registers all API routes with the given router
func SetupRoutes(repo repository.ProjectRepository, logRepo repository.LogRepository, dashboardRepo repository.DashboardRepository, userConfig *service.UserConfigService, projectsService dashboard.ProjectsServiceInterface) http.Handler {
	r := mux.NewRouter()

	// Create handlers
	projectsHandler := handlers.NewProjectsHandler(repo)
	logsHandler := handlers.NewLogsHandler(logRepo, repo)
	healthHandler := handlers.NewHealthHandler()
	dashboardHandler := handlers.NewDashboardHandler(dashboardRepo, userConfig, projectsService)

	// Health check
	r.HandleFunc("/healthz", healthHandler.Healthz).Methods("GET")

	// Projects endpoints
	r.HandleFunc("/v1/projects.json", projectsHandler.Index).Methods("GET")
	r.HandleFunc("/v1/projects.json", projectsHandler.Create).Methods("POST")
	r.HandleFunc("/v1/projects/{id}.json", projectsHandler.Show).Methods("GET")

	// Logs endpoints
	r.HandleFunc("/v1/projects/{project_id}/logs.json", logsHandler.Index).Methods("GET")
	// Phase 1: Read-only - logs creation will be added in Phase 2

	// Dashboard endpoints
	r.HandleFunc("/v1/dashboard/day.json", dashboardHandler.Day).Methods("GET")
	r.HandleFunc("/v1/dashboard/projects.json", dashboardHandler.Projects).Methods("GET")

	return r
}
