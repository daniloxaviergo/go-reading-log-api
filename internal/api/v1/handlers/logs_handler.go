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

// parseLogDate attempts to parse a date string using multiple formats.
// Supported formats:
//   - RFC3339 (e.g., "2024-01-15T10:30:00Z")
//   - Date only (e.g., "2024-01-15")
//   - Standard datetime (e.g., "2024-01-15 10:30:00")
func parseLogDate(dateStr string) (*time.Time, bool) {
	formats := []string{
		time.RFC3339,          // 2006-01-02T15:04:05Z
		"2006-01-02",          // YYYY-MM-DD
		"2006-01-02 15:04:05", // Standard datetime
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return &t, true
		}
	}

	return nil, false
}

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

	// Build data objects for JSON:API response
	dataObjects := make([]dto.JSONAPIData, limit)

	for i := 0; i < limit; i++ {
		// Parse the data string to time.Time for RFC3339 compliance
		var dataTime *time.Time
		if logs[i].Data != nil && *logs[i].Data != "" {
			dataTime, _ = parseLogDate(*logs[i].Data)
		}

		logResponse := &dto.LogResponse{
			ID:        logs[i].ID,
			Data:      dataTime,
			StartPage: logs[i].StartPage,
			EndPage:   logs[i].EndPage,
			Note:      logs[i].Note,
		}

		// Build JSONAPIData with relationships at the resource level (not in attributes)
		dataObjects[i] = dto.JSONAPIData{
			Type:       "logs",
			ID:         strconv.FormatInt(logs[i].ID, 10), // ID as string per JSON:API spec
			Attributes: logResponse,
			Relationships: map[string]interface{}{
				"project": map[string]string{
					"id":   strconv.FormatInt(project.ID, 10),
					"type": "projects",
				},
			},
		}
	}

	// Build included array with project data
	included := []interface{}{}

	if limit > 0 {
		// Create project response for inclusion
		// Convert time.Time pointers to string pointers for JSON serialization
		startedAtStr := formatTimePtr(project.StartedAt)
		finishedAtStr := formatTimePtr(project.FinishedAt)

		projectResponse := &dto.ProjectResponse{
			ID:         project.ID,
			Name:       project.Name,
			TotalPage:  project.TotalPage,
			Page:       project.Page,
			StartedAt:  startedAtStr,
			Progress:   project.Progress,
			Status:     project.Status,
			LogsCount:  project.LogsCount,
			DaysUnread: project.DaysUnread,
			MedianDay:  project.MedianDay,
			FinishedAt: finishedAtStr,
		}

		// Add project to included array
		included = append(included, dto.NewIncludedProject(projectResponse))
	}

	// Wrap collection in envelope with included resources
	envelope := dto.NewJSONAPIEnvelopeWithIncluded(dataObjects, included)

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(envelope)
}
