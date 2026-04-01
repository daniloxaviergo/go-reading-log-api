package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5"
	"go-reading-log-api-next/internal/api/v1/middleware"
	"go-reading-log-api-next/internal/domain/models"
	"go-reading-log-api-next/internal/repository"
)

// TestSetupRoutes tests the SetupRoutes function
func TestSetupRoutes(t *testing.T) {
	// Create mock repositories
	projectRepo := &MockProjectRepository{}
	logRepo := &MockLogRepository{}

	// Setup routes
	handler := SetupRoutes(projectRepo, logRepo)

	if handler == nil {
		t.Fatal("Expected non-nil handler, got nil")
	}

	// Test that the handler responds to requests
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d for /healthz, got %d", http.StatusOK, w.Code)
	}
}

// TestSetupRoutes_Routes tests that all expected routes are registered
func TestSetupRoutes_Routes(t *testing.T) {
	projectRepo := &MockProjectRepository{}
	logRepo := &MockLogRepository{}

	handler := SetupRoutes(projectRepo, logRepo)

	// Test healthz endpoint
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 for /healthz, got %d", w.Code)
	}

	// Test projects endpoint
	req = httptest.NewRequest(http.MethodGet, "/api/v1/projects", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 for /api/v1/projects, got %d", w.Code)
	}

	// Test projects/{id} endpoint (will return 404 since no project exists)
	req = httptest.NewRequest(http.MethodGet, "/api/v1/projects/1", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	// Expected 404 because mock repo returns not found

	// Test logs endpoint
	req = httptest.NewRequest(http.MethodGet, "/api/v1/projects/1/logs", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 for /api/v1/projects/1/logs, got %d", w.Code)
	}
}

// TestSetupRoutes_MiddlewareChain tests that routes are properly wrapped with middleware
func TestSetupRoutes_MiddlewareChain(t *testing.T) {
	projectRepo := &MockProjectRepository{}
	logRepo := &MockLogRepository{}

	handler := SetupRoutes(projectRepo, logRepo)

	// Wrap with middleware chain
	middlewareChain := middleware.Chain(handler,
		middleware.RecoveryMiddleware,
		middleware.CORSMiddleware,
		middleware.RequestIDMiddleware,
		middleware.LoggingMiddleware,
	)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	middlewareChain.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

// MockProjectRepository is a mock implementation of repository.ProjectRepository
type MockProjectRepository struct {
	Projects map[int64]interface{}
}

func (m *MockProjectRepository) GetByID(ctx context.Context, id int64) (*models.Project, error) {
	if id == 1 {
		return &models.Project{
			ID:        1,
			Name:      "Test Project",
			TotalPage: 100,
			Page:      50,
		}, nil
	}
	return nil, fmt.Errorf("project with ID %d not found: %w", id, pgx.ErrNoRows)
}

func (m *MockProjectRepository) GetAll(ctx context.Context) ([]*models.Project, error) {
	return nil, nil
}

func (m *MockProjectRepository) GetWithLogs(ctx context.Context, id int64) (*repository.ProjectWithLogs, error) {
	return nil, fmt.Errorf("project with ID %d not found: %w", id, pgx.ErrNoRows)
}

func (m *MockProjectRepository) GetAllWithLogs(ctx context.Context) ([]*repository.ProjectWithLogs, error) {
	return nil, nil
}

// MockLogRepository is a mock implementation of repository.LogRepository
type MockLogRepository struct {
	Logs map[int64]interface{}
}

func (m *MockLogRepository) GetByID(ctx context.Context, id int64) (*models.Log, error) {
	return nil, nil
}

func (m *MockLogRepository) GetByProjectID(ctx context.Context, projectID int64) ([]*models.Log, error) {
	return nil, nil
}

func (m *MockLogRepository) GetAll(ctx context.Context) ([]*models.Log, error) {
	return nil, nil
}

func (m *MockLogRepository) GetByProjectIDOrdered(ctx context.Context, projectID int64) ([]*models.Log, error) {
	return nil, nil
}
