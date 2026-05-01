package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"go-reading-log-api-next/internal/api/v1/middleware"
	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/domain/models"
	"go-reading-log-api-next/internal/repository"
	"go-reading-log-api-next/internal/service"
	"go-reading-log-api-next/internal/service/dashboard"
)

// MockProjectsService is a mock implementation of ProjectsServiceInterface
type MockProjectsService struct{}

func (m *MockProjectsService) GetRunningProjectsWithLogs(ctx context.Context) ([]*dashboard.ProjectWithLogs, error) {
	return []*dashboard.ProjectWithLogs{}, nil
}

func (m *MockProjectsService) CalculateStats(ctx context.Context) (*dto.StatsData, error) {
	return dto.NewStatsData(), nil
}

func (m *MockProjectsService) GetDashboardProjects(ctx context.Context) (*dto.DashboardProjectsResponse, error) {
	return dto.NewDashboardProjectsResponse(), nil
}

var _ dashboard.ProjectsServiceInterface = (*MockProjectsService)(nil)

// TestSetupRoutes tests the SetupRoutes function
func TestSetupRoutes(t *testing.T) {
	// Create mock repositories
	projectRepo := &MockProjectRepository{}
	logRepo := &MockLogRepository{}
	dashboardRepo := &MockDashboardRepository{}
	userConfig := service.NewUserConfigService(service.GetDefaultConfig())
	projectsService := &MockProjectsService{}

	// Setup routes
	handler := SetupRoutes(projectRepo, logRepo, dashboardRepo, userConfig, dashboard.ProjectsServiceInterface(projectsService))

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
	dashboardRepo := &MockDashboardRepository{}
	userConfig := service.NewUserConfigService(service.GetDefaultConfig())
	projectsService := &MockProjectsService{}

	handler := SetupRoutes(projectRepo, logRepo, dashboardRepo, userConfig, dashboard.ProjectsServiceInterface(projectsService))

	// Test healthz endpoint
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 for /healthz, got %d", w.Code)
	}

	// Test projects endpoint
	req = httptest.NewRequest(http.MethodGet, "/v1/projects.json", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 for /v1/projects, got %d", w.Code)
	}

	// Test projects/{id} endpoint (will return 404 since no project exists)
	req = httptest.NewRequest(http.MethodGet, "/v1/projects/1.json", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	// Expected 404 because mock repo returns not found

	// Test logs endpoint
	req = httptest.NewRequest(http.MethodGet, "/v1/projects/1/logs.json", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 for /v1/projects/1/logs, got %d", w.Code)
	}
}

// TestSetupRoutes_MiddlewareChain tests that routes are properly wrapped with middleware
func TestSetupRoutes_MiddlewareChain(t *testing.T) {
	projectRepo := &MockProjectRepository{}
	logRepo := &MockLogRepository{}
	dashboardRepo := &MockDashboardRepository{}
	userConfig := service.NewUserConfigService(service.GetDefaultConfig())
	projectsService := &MockProjectsService{}

	handler := SetupRoutes(projectRepo, logRepo, dashboardRepo, userConfig, dashboard.ProjectsServiceInterface(projectsService))

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

func (m *MockProjectRepository) Create(ctx context.Context, project *models.Project) (*models.Project, error) {
	if m.Projects == nil {
		m.Projects = make(map[int64]interface{})
	}

	// Generate a new ID
	var maxID int64
	for id := range m.Projects {
		if id > maxID {
			maxID = id
		}
	}
	project.ID = maxID + 1
	m.Projects[project.ID] = project

	return project, nil
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

func (m *MockLogRepository) Create(ctx context.Context, log *models.Log) (*models.Log, error) {
	// Generate a new ID
	var maxID int64
	for id := range m.Logs {
		if id > maxID {
			maxID = id
		}
	}
	log.ID = maxID + 1
	m.Logs[log.ID] = log

	return log, nil
}

// MockDashboardRepository is a mock implementation of repository.DashboardRepository
type MockDashboardRepository struct{}

func (m *MockDashboardRepository) GetPool() repository.PoolInterface {
	return nil
}

func (m *MockDashboardRepository) GetDailyStats(ctx context.Context, date time.Time) (*dto.DailyStats, error) {
	return dto.NewDailyStats(0, 0), nil
}

func (m *MockDashboardRepository) GetProjectAggregates(ctx context.Context) ([]*dto.ProjectAggregate, error) {
	return []*dto.ProjectAggregate{}, nil
}

func (m *MockDashboardRepository) GetFaultsByDateRange(ctx context.Context, start, end time.Time) (*dto.FaultStats, error) {
	return &dto.FaultStats{FaultCount: 0}, nil
}

func (m *MockDashboardRepository) GetWeekdayFaults(ctx context.Context, start, end time.Time) (*dto.WeekdayFaults, error) {
	return dto.NewWeekdayFaults(make(map[int]int)), nil
}

func (m *MockDashboardRepository) GetLogsByDateRange(ctx context.Context, start, end time.Time) ([]*dto.LogEntry, error) {
	return []*dto.LogEntry{}, nil
}

func (m *MockDashboardRepository) GetProjectWeekdayMean(ctx context.Context, projectID int64, weekday int) (float64, error) {
	return 0.0, nil
}

func (m *MockDashboardRepository) CalculatePeriodPages(ctx context.Context, start, end time.Time) (int, error) {
	return 0, nil
}

func (m *MockDashboardRepository) GetProjectsWithLogs(ctx context.Context) ([]*dto.ProjectAggregateResponse, error) {
	return []*dto.ProjectAggregateResponse{}, nil
}

func (m *MockDashboardRepository) GetProjectLogs(ctx context.Context, projectID int64, limit int) ([]*dto.LogEntry, error) {
	return []*dto.LogEntry{}, nil
}

func (m *MockDashboardRepository) GetMaxByWeekday(ctx context.Context, date time.Time) (*float64, error) {
	return nil, nil
}

func (m *MockDashboardRepository) GetOverallMean(ctx context.Context, date time.Time) (*float64, error) {
	return nil, nil
}

func (m *MockDashboardRepository) GetPreviousPeriodMean(ctx context.Context, date time.Time) (*float64, error) {
	return nil, nil
}

func (m *MockDashboardRepository) GetPreviousPeriodSpecMean(ctx context.Context, date time.Time) (*float64, error) {
	return nil, nil
}

func (m *MockDashboardRepository) GetMeanByWeekday(ctx context.Context, weekday int) (*float64, error) {
	return nil, nil
}

func (m *MockDashboardRepository) GetRunningProjectsWithLogs(ctx context.Context) ([]*dto.ProjectWithLogs, error) {
	return nil, nil
}
