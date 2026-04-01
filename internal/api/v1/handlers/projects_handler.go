package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/repository"
)

// ProjectsHandler handles HTTP requests for projects endpoints
type ProjectsHandler struct {
	repo repository.ProjectRepository
}

// NewProjectsHandler creates a new ProjectsHandler with the given repository
func NewProjectsHandler(repo repository.ProjectRepository) *ProjectsHandler {
	return &ProjectsHandler{repo: repo}
}

// Index returns all projects ordered by logs data DESC
func (h *ProjectsHandler) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Use the repository method that eager-loads logs
	projectsWithLogs, err := h.repo.GetAllWithLogs(ctx)
	if err != nil {
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
