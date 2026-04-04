package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/domain/models"
	"go-reading-log-api-next/internal/repository"
	"go-reading-log-api-next/internal/validation"
)

// ProjectsHandler handles HTTP requests for projects endpoints
type ProjectsHandler struct {
	repo repository.ProjectRepository
}

// NewProjectsHandler creates a new ProjectsHandler with the given repository
func NewProjectsHandler(repo repository.ProjectRepository) *ProjectsHandler {
	return &ProjectsHandler{repo: repo}
}

// Create handles POST /api/v1/projects - creates a new project with page <= total_page validation
func (h *ProjectsHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse request body
	var req dto.ProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Determine status - use default "unstarted" since ProjectRequest doesn't include status
	status := "unstarted"

	// Validate total_page > 0 and page <= total_page constraints
	validationErr := validation.ValidateProject(req.Page, req.TotalPage, status)
	if validationErr != nil && validationErr.HasErrors() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   "validation failed",
			"details": validationErr.ToMap(),
		})
		return
	}

	// Create domain project model
	project := &models.Project{
		Name:      req.Name,
		TotalPage: req.TotalPage,
		Page:      req.Page,
		Reinicia:  req.Reinicia,
	}

	// Handle optional started_at field
	if req.StartedAt != nil {
		parsedTime, err := time.Parse(time.RFC3339, *req.StartedAt)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "invalid date format",
				"details": map[string]string{
					"started_at": "must be in RFC3339 format",
				},
			})
			return
		}
		project.StartedAt = &parsedTime
	}

	// Create project in database
	createdProject, err := h.repo.Create(ctx, project)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Handler error: %v\n", err)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Convert to response DTO
	response := &dto.ProjectResponse{
		ID:        createdProject.ID,
		Name:      createdProject.Name,
		TotalPage: createdProject.TotalPage,
		Page:      createdProject.Page,
	}

	// Convert timestamps to strings for JSON
	if createdProject.StartedAt != nil {
		s := createdProject.StartedAt.Format(time.RFC3339)
		response.StartedAt = &s
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Index returns all projects ordered by logs data DESC
func (h *ProjectsHandler) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Use the repository method that eager-loads logs
	projectsWithLogs, err := h.repo.GetAllWithLogs(ctx)
	if err != nil {
		// Log error to stderr for debugging
		fmt.Fprintf(os.Stderr, "Handler error: %v\n", err)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Convert to response DTOs
	response := make([]*dto.ProjectResponse, len(projectsWithLogs))
	for i, pw := range projectsWithLogs {
		response[i] = pw.Project
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Show returns a single project by ID with eager-loaded logs
func (h *ProjectsHandler) Show(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract project ID from path
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "Invalid project ID"}`, http.StatusBadRequest)
		return
	}

	// Use GetWithLogs to get project with eager-loaded logs
	projectWithLogs, err := h.repo.GetWithLogs(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, `{"error": "project not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projectWithLogs.Project)
}
