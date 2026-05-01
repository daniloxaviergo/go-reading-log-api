package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/service/dashboard"
)

// MockProjectsService is a mock implementation of ProjectsServiceInterface
type MockProjectsService struct {
	mock.Mock
}

func (m *MockProjectsService) GetRunningProjectsWithLogs(ctx context.Context) ([]*dashboard.ProjectWithLogs, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*dashboard.ProjectWithLogs), args.Error(1)
}

func (m *MockProjectsService) CalculateStats(ctx context.Context) (*dto.StatsData, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.StatsData), args.Error(1)
}

func (m *MockProjectsService) GetDashboardProjects(ctx context.Context) (*dto.DashboardProjectsResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DashboardProjectsResponse), args.Error(1)
}

// TestProjects_Success tests the happy path with valid projects and stats in JSON:API format
func TestProjects_Success(t *testing.T) {
	// Create mock service
	mockService := new(MockProjectsService)

	// Setup mock data - 2 projects in JSON:API format
	attributes1 := dto.NewDashboardProjectAttributes()
	attributes1.Name = "Project 1"
	attributes1.TotalPage = 100
	attributes1.Page = 50
	attributes1.Progress = 50.0
	attributes1.LogsCount = 2
	attributes1.Status = "stopped"
	attributes1.DaysUnreading = 5
	startedAt1 := "2024-01-10T10:30:00Z"
	attributes1.SetStartedAt(startedAt1)

	attributes2 := dto.NewDashboardProjectAttributes()
	attributes2.Name = "Project 2"
	attributes2.TotalPage = 200
	attributes2.Page = 50
	attributes2.Progress = 25.0
	attributes2.LogsCount = 1
	attributes2.Status = "stopped"
	attributes2.DaysUnreading = 10
	startedAt2 := "2024-01-12T14:00:00Z"
	attributes2.SetStartedAt(startedAt2)

	response := dto.NewDashboardProjectsResponse()
	response.AddProject(*dto.NewDashboardProjectItem("1", attributes1))
	response.AddProject(*dto.NewDashboardProjectItem("2", attributes2))

	stats := dto.NewDashboardStats()
	stats.SetTotalPages(300)
	stats.SetPages(100)
	stats.SetProgressGeral(33.333)
	response.SetStats(stats)

	// Setup expectations
	mockService.On("GetDashboardProjects", mock.Anything).Return(response, nil)

	// Create handler with mock service
	handler := &DashboardHandler{
		projectsService: dashboard.ProjectsServiceInterface(mockService),
	}

	// Create test request
	req := httptest.NewRequest("GET", "/v1/dashboard/projects.json", nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.Projects(w, req)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	// Parse response
	var jsonResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &jsonResponse)
	assert.NoError(t, err)

	// Verify response structure - JSON:API format
	assert.Contains(t, jsonResponse, "data")
	assert.Contains(t, jsonResponse, "stats")

	// Verify data array
	dataArray, ok := jsonResponse["data"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, dataArray, 2)

	// Verify first project structure
	project1 := dataArray[0].(map[string]interface{})
	assert.Equal(t, "1", project1["id"])
	assert.Equal(t, "projects", project1["type"])

	attributes, ok := project1["attributes"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "Project 1", attributes["name"])
	assert.Equal(t, "2024-01-10T10:30:00Z", attributes["started-at"])
	assert.Equal(t, 50.0, attributes["progress"])
	assert.Equal(t, 100.0, attributes["total-page"])
	assert.Equal(t, 50.0, attributes["page"])
	assert.Equal(t, "stopped", attributes["status"])
	assert.Equal(t, 2.0, attributes["logs-count"])
	assert.Equal(t, 5.0, attributes["days-unreading"])

	// Verify stats object
	statsObj, ok := jsonResponse["stats"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, 300.0, statsObj["total_pages"])
	assert.Equal(t, 100.0, statsObj["pages"])
	assert.Equal(t, 33.333, statsObj["progress_geral"])

	// Verify mocks were called
	mockService.AssertExpectations(t)
}

// TestProjects_EmptyData tests when there are no projects
func TestProjects_EmptyData(t *testing.T) {
	// Create mock service
	mockService := new(MockProjectsService)

	// Setup mock data - empty response
	response := dto.NewDashboardProjectsResponse()

	stats := dto.NewDashboardStats()
	stats.SetTotalPages(0)
	stats.SetPages(0)
	stats.SetProgressGeral(0.0)
	response.SetStats(stats)

	// Setup expectations
	mockService.On("GetDashboardProjects", mock.Anything).Return(response, nil)

	// Create handler with mock service
	handler := &DashboardHandler{
		projectsService: dashboard.ProjectsServiceInterface(mockService),
	}

	// Create test request
	req := httptest.NewRequest("GET", "/v1/dashboard/projects.json", nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.Projects(w, req)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response
	var jsonResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &jsonResponse)
	assert.NoError(t, err)

	// Verify empty data array
	dataArray, ok := jsonResponse["data"].([]interface{})
	assert.True(t, ok)
	assert.Empty(t, dataArray)

	// Verify zero stats
	statsObj, ok := jsonResponse["stats"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, 0.0, statsObj["total_pages"])
	assert.Equal(t, 0.0, statsObj["pages"])
	assert.Equal(t, 0.0, statsObj["progress_geral"])

	// Verify mocks were called
	mockService.AssertExpectations(t)
}

// TestProjects_ServiceError tests when service returns an error
func TestProjects_ServiceError(t *testing.T) {
	// Create mock service
	mockService := new(MockProjectsService)

	// Setup mock error
	testError := assert.AnError

	// Setup expectations
	mockService.On("GetDashboardProjects", mock.Anything).Return((*dto.DashboardProjectsResponse)(nil), testError)

	// Create handler with mock service
	handler := &DashboardHandler{
		projectsService: dashboard.ProjectsServiceInterface(mockService),
	}

	// Create test request
	req := httptest.NewRequest("GET", "/v1/dashboard/projects.json", nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.Projects(w, req)

	// Verify response
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Internal server error")

	// Verify mocks were called
	mockService.AssertExpectations(t)
}

// TestProjects_WithStartedAtNull tests that started-at can be null when no logs exist
func TestProjects_WithStartedAtNull(t *testing.T) {
	// Create mock service
	mockService := new(MockProjectsService)

	// Setup mock data - project without started-at
	attributes := dto.NewDashboardProjectAttributes()
	attributes.Name = "Project Without Logs"
	attributes.TotalPage = 100
	attributes.Page = 0
	attributes.Progress = 0.0
	attributes.LogsCount = 0
	attributes.Status = "stopped"
	attributes.DaysUnreading = 0
	// Note: started-at is not set (will be null/omitted)

	response := dto.NewDashboardProjectsResponse()
	response.AddProject(*dto.NewDashboardProjectItem("1", attributes))

	stats := dto.NewDashboardStats()
	stats.SetTotalPages(100)
	stats.SetPages(0)
	stats.SetProgressGeral(0.0)
	response.SetStats(stats)

	// Setup expectations
	mockService.On("GetDashboardProjects", mock.Anything).Return(response, nil)

	// Create handler with mock service
	handler := &DashboardHandler{
		projectsService: dashboard.ProjectsServiceInterface(mockService),
	}

	// Create test request
	req := httptest.NewRequest("GET", "/v1/dashboard/projects.json", nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.Projects(w, req)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response
	var jsonResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &jsonResponse)
	assert.NoError(t, err)

	// Verify project structure
	dataArray := jsonResponse["data"].([]interface{})
	project1 := dataArray[0].(map[string]interface{})
	attributesMap := project1["attributes"].(map[string]interface{})

	// started-at should not be present when nil (omitted due to omitempty)
	_, hasStartedAt := attributesMap["started-at"]
	assert.False(t, hasStartedAt, "started-at should be omitted when nil")

	// Verify mocks were called
	mockService.AssertExpectations(t)
}
