package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/repository"
)

// LogsHandler handles HTTP requests for logs endpoints
type LogsHandler struct {
	repo        repository.LogRepository
	projectRepo repository.ProjectRepository
}

// NewLogsHandler creates a new LogsHandler with the given repositories
func NewLogsHandler(repo repository.LogRepository, projectRepo repository.ProjectRepository) *LogsHandler {
	return &LogsHandler{repo: repo, projectRepo: projectRepo}
}

// formatTimePtr converts a time.Time pointer to a string pointer for JSON serialization
func formatTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(time.RFC3339)
	return &s
}

// Index returns the first 4 logs for a project with project eager-loaded
func (h *LogsHandler) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract project ID from path
	idStr := mux.Vars(r)["project_id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "Invalid project ID"}`, http.StatusBadRequest)
		return
	}

	// Get project first to verify it exists and eager-load it
	project, err := h.projectRepo.GetByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, `{"error": "project not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf(`{"error": "Internal server error: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	if project == nil {
		http.Error(w, `{"error": "project not found"}`, http.StatusNotFound)
		return
	}

	// Get logs for the project ordered by data DESC
	logs, err := h.repo.GetByProjectIDOrdered(ctx, id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Internal server error: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	// Limit to first 4 logs
	limit := 4
	if len(logs) < limit {
		limit = len(logs)
	}

	// Convert to response DTOs (with nested project object for Rails API compatibility)
	response := make([]*dto.LogResponse, limit)
	for i := 0; i < limit; i++ {
		// Format StartedAt as string
		var startedAtStr *string
		if project.StartedAt != nil {
			formatted := project.StartedAt.Format(time.RFC3339)
			startedAtStr = &formatted
		}

		response[i] = &dto.LogResponse{
			ID:        logs[i].ID,
			Data:      logs[i].Data,
			StartPage: logs[i].StartPage,
			EndPage:   logs[i].EndPage,
			Note:      logs[i].Note,
			Project: &dto.ProjectResponse{
				ID:        project.ID,
				Name:      project.Name,
				TotalPage: project.TotalPage,
				Page:      project.Page,
				StartedAt: startedAtStr,
				Status:    project.Status,
				Progress:  project.Progress,
			},
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
