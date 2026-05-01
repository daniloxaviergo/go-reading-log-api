package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/repository"
	"go-reading-log-api-next/internal/service"
)

// MockDashboardRepository is a mock implementation of DashboardRepository
type MockDashboardRepository struct {
	mock.Mock
}

func (m *MockDashboardRepository) GetDailyStats(ctx context.Context, date time.Time) (*dto.DailyStats, error) {
	args := m.Called(ctx, date)
	return args.Get(0).(*dto.DailyStats), args.Error(1)
}

func (m *MockDashboardRepository) GetProjectAggregates(ctx context.Context) ([]*dto.ProjectAggregate, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*dto.ProjectAggregate), args.Error(1)
}

func (m *MockDashboardRepository) GetFaultsByDateRange(ctx context.Context, start, end time.Time) (*dto.FaultStats, error) {
	args := m.Called(ctx, start, end)
	return args.Get(0).(*dto.FaultStats), args.Error(1)
}

func (m *MockDashboardRepository) GetWeekdayFaults(ctx context.Context, start, end time.Time) (*dto.WeekdayFaults, error) {
	args := m.Called(ctx, start, end)
	return args.Get(0).(*dto.WeekdayFaults), args.Error(1)
}

func (m *MockDashboardRepository) GetLogsByDateRange(ctx context.Context, start, end time.Time) ([]*dto.LogEntry, error) {
	args := m.Called(ctx, start, end)
	return args.Get(0).([]*dto.LogEntry), args.Error(1)
}

func (m *MockDashboardRepository) GetProjectWeekdayMean(ctx context.Context, projectID int64, weekday int) (float64, error) {
	args := m.Called(ctx, projectID, weekday)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockDashboardRepository) CalculatePeriodPages(ctx context.Context, start, end time.Time) (int, error) {
	args := m.Called(ctx, start, end)
	return args.Get(0).(int), args.Error(1)
}

func (m *MockDashboardRepository) GetPool() repository.PoolInterface {
	return nil
}

func (m *MockDashboardRepository) GetProjectsWithLogs(ctx context.Context) ([]*dto.ProjectAggregateResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*dto.ProjectAggregateResponse), args.Error(1)
}

func (m *MockDashboardRepository) GetProjectLogs(ctx context.Context, projectID int64, limit int) ([]*dto.LogEntry, error) {
	args := m.Called(ctx, projectID, limit)
	return args.Get(0).([]*dto.LogEntry), args.Error(1)
}

func (m *MockDashboardRepository) GetMaxByWeekday(ctx context.Context, date time.Time) (*float64, error) {
	args := m.Called(ctx, date)
	return args.Get(0).(*float64), args.Error(1)
}

func (m *MockDashboardRepository) GetOverallMean(ctx context.Context, date time.Time) (*float64, error) {
	args := m.Called(ctx, date)
	return args.Get(0).(*float64), args.Error(1)
}

func (m *MockDashboardRepository) GetPreviousPeriodMean(ctx context.Context, date time.Time) (*float64, error) {
	args := m.Called(ctx, date)
	return args.Get(0).(*float64), args.Error(1)
}

func (m *MockDashboardRepository) GetPreviousPeriodSpecMean(ctx context.Context, date time.Time) (*float64, error) {
	args := m.Called(ctx, date)
	return args.Get(0).(*float64), args.Error(1)
}

func (m *MockDashboardRepository) GetMeanByWeekday(ctx context.Context, weekday int) (*float64, error) {
	args := m.Called(ctx, weekday)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*float64), args.Error(1)
}

func (m *MockDashboardRepository) GetRunningProjectsWithLogs(ctx context.Context) ([]*dto.ProjectWithLogs, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*dto.ProjectWithLogs), args.Error(1)
}

// TestDashboardHandler_Day tests the Day handler
func TestDashboardHandler_Day(t *testing.T) {
	mockRepo := &MockDashboardRepository{}
	userConfig := service.NewUserConfigService(service.GetDefaultConfig())
	mockProjectsService := &MockProjectsService{}
	handler := NewDashboardHandler(mockRepo, userConfig, mockProjectsService)

	testDate := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	expectedStats := dto.NewDailyStats(100, 5)

	// Mock GetDailyStats for both current date and previous period (7 days before)
	mockRepo.On("GetDailyStats", mock.Anything, testDate).Return(expectedStats, nil)
	prevDate := testDate.AddDate(0, 0, -7)
	mockRepo.On("GetDailyStats", mock.Anything, prevDate).Return(dto.NewDailyStats(75, 4), nil)

	// Mock GetProjectAggregates to avoid panic - return empty slice
	mockRepo.On("GetProjectAggregates", mock.Anything).Return([]*dto.ProjectAggregate{}, nil)

	// Mock new repository methods for additional calculations
	maxDay := 50.0
	mockRepo.On("GetMaxByWeekday", mock.Anything, testDate).Return(&maxDay, nil)
	meanGeral := 25.0
	mockRepo.On("GetOverallMean", mock.Anything, testDate).Return(&meanGeral, nil)
	prevMean := 20.0
	mockRepo.On("GetPreviousPeriodMean", mock.Anything, testDate).Return(&prevMean, nil)
	prevSpecMean := 23.0
	mockRepo.On("GetPreviousPeriodSpecMean", mock.Anything, testDate).Return(&prevSpecMean, nil)
	// Mock GetMeanByWeekday - return 40.0 (using V1::MeanLog algorithm: total_pages / count_reads)
	meanDay := 40.0
	mockRepo.On("GetMeanByWeekday", mock.Anything, 1).Return(&meanDay, nil) // Monday = weekday 1

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/day.json?date=2024-01-15T10:30:00Z", nil)
	w := httptest.NewRecorder()

	handler.Day(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	// Verify flat JSON structure (no JSON:API envelope)
	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	// Verify no JSON:API envelope keys present
	_, hasData := response["data"]
	_, hasType := response["type"]
	_, hasId := response["id"]
	_, hasAttributes := response["attributes"]
	assert.False(t, hasData, "Response should not have 'data' key")
	assert.False(t, hasType, "Response should not have 'type' key")
	assert.False(t, hasId, "Response should not have 'id' key")
	assert.False(t, hasAttributes, "Response should not have 'attributes' key")

	// Verify stats key at root level
	statsMap, ok := response["stats"].(map[string]interface{})
	require.True(t, ok, "Response should have 'stats' key at root level")

	assert.Equal(t, float64(100), statsMap["total_pages"])
	assert.Equal(t, float64(40), statsMap["mean_day"]) // Using V1::MeanLog algorithm
	assert.NotNil(t, statsMap["per_pages"])            // Should have per_pages since prevStats > 0
	assert.NotNil(t, statsMap["max_day"])              // Should have max_day
	assert.NotNil(t, statsMap["mean_geral"])           // Should have mean_geral
	assert.NotNil(t, statsMap["per_mean_day"])         // Should have per_mean_day

	mockRepo.AssertExpectations(t)
}

// TestDashboardHandler_Day_EmptyData tests Day handler with no data
func TestDashboardHandler_Day_EmptyData(t *testing.T) {
	mockRepo := &MockDashboardRepository{}
	userConfig := service.NewUserConfigService(service.GetDefaultConfig())
	mockProjectsService := &MockProjectsService{}
	handler := NewDashboardHandler(mockRepo, userConfig, mockProjectsService)

	testDate := time.Date(2024, 1, 20, 10, 30, 0, 0, time.UTC)
	emptyStats := dto.NewDailyStats(0, 0)

	// Mock GetDailyStats for both current date and previous period (7 days before)
	mockRepo.On("GetDailyStats", mock.Anything, testDate).Return(emptyStats, nil)
	prevDate := testDate.AddDate(0, 0, -7)
	mockRepo.On("GetDailyStats", mock.Anything, prevDate).Return(emptyStats, nil)

	// Mock GetProjectAggregates to avoid panic - return empty slice
	mockRepo.On("GetProjectAggregates", mock.Anything).Return([]*dto.ProjectAggregate{}, nil)

	// Mock new repository methods for additional calculations (return nil for empty data)
	mockRepo.On("GetMaxByWeekday", mock.Anything, testDate).Return((*float64)(nil), nil)
	mockRepo.On("GetOverallMean", mock.Anything, testDate).Return((*float64)(nil), nil)
	mockRepo.On("GetPreviousPeriodMean", mock.Anything, testDate).Return((*float64)(nil), nil)
	mockRepo.On("GetPreviousPeriodSpecMean", mock.Anything, testDate).Return((*float64)(nil), nil)
	// Mock GetMeanByWeekday - return nil for empty data
	mockRepo.On("GetMeanByWeekday", mock.Anything, 6).Return((*float64)(nil), nil) // Saturday = weekday 6

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/day.json?date=2024-01-20T10:30:00Z", nil)
	w := httptest.NewRecorder()

	handler.Day(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	// Verify flat JSON structure (no JSON:API envelope)
	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	// Verify no JSON:API envelope keys present
	_, hasData := response["data"]
	_, hasType := response["type"]
	_, hasId := response["id"]
	_, hasAttributes := response["attributes"]
	assert.False(t, hasData, "Response should not have 'data' key")
	assert.False(t, hasType, "Response should not have 'type' key")
	assert.False(t, hasId, "Response should not have 'id' key")
	assert.False(t, hasAttributes, "Response should not have 'attributes' key")

	// Verify stats key at root level
	statsMap, ok := response["stats"].(map[string]interface{})
	require.True(t, ok, "Response should have 'stats' key at root level")

	assert.Equal(t, float64(0), statsMap["total_pages"])
	assert.Equal(t, float64(0), statsMap["mean_day"])
	assert.Nil(t, statsMap["per_pages"]) // Should be null when no previous data

	mockRepo.AssertExpectations(t)
}

// TestDashboardHandler_Day_InvalidDate tests Day handler with invalid date
func TestDashboardHandler_Day_InvalidDate(t *testing.T) {
	mockRepo := &MockDashboardRepository{}
	userConfig := service.NewUserConfigService(service.GetDefaultConfig())
	mockProjectsService := &MockProjectsService{}
	handler := NewDashboardHandler(mockRepo, userConfig, mockProjectsService)

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/day.json?date=invalid-date", nil)
	w := httptest.NewRecorder()

	handler.Day(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errorResp map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&errorResp)
	require.NoError(t, err)

	assert.Equal(t, "invalid date format", errorResp["error"])
	mockRepo.AssertNotCalled(t, "GetDailyStats")
}

// TestDashboardHandler_Day_NullPerPages tests that per_pages is null when previous period has 0 pages
func TestDashboardHandler_Day_NullPerPages(t *testing.T) {
	mockRepo := &MockDashboardRepository{}
	userConfig := service.NewUserConfigService(service.GetDefaultConfig())
	mockProjectsService := &MockProjectsService{}
	handler := NewDashboardHandler(mockRepo, userConfig, mockProjectsService)

	testDate := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	expectedStats := dto.NewDailyStats(100, 5)

	// Mock GetDailyStats for current date
	mockRepo.On("GetDailyStats", mock.Anything, testDate).Return(expectedStats, nil)
	prevDate := testDate.AddDate(0, 0, -7)
	// Mock previous period with 0 pages (should result in null per_pages)
	mockRepo.On("GetDailyStats", mock.Anything, prevDate).Return(dto.NewDailyStats(0, 0), nil)

	// Mock GetProjectAggregates to avoid panic - return empty slice
	mockRepo.On("GetProjectAggregates", mock.Anything).Return([]*dto.ProjectAggregate{}, nil)

	// Mock new repository methods for additional calculations
	maxDay := 50.0
	mockRepo.On("GetMaxByWeekday", mock.Anything, testDate).Return(&maxDay, nil)
	meanGeral := 25.0
	mockRepo.On("GetOverallMean", mock.Anything, testDate).Return(&meanGeral, nil)
	mockRepo.On("GetPreviousPeriodMean", mock.Anything, testDate).Return((*float64)(nil), nil)
	mockRepo.On("GetPreviousPeriodSpecMean", mock.Anything, testDate).Return((*float64)(nil), nil)
	// Mock GetMeanByWeekday - return 40.0 (using V1::MeanLog algorithm)
	meanDay := 40.0
	mockRepo.On("GetMeanByWeekday", mock.Anything, 1).Return(&meanDay, nil) // Monday = weekday 1

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/day.json?date=2024-01-15T10:30:00Z", nil)
	w := httptest.NewRecorder()

	handler.Day(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	// Verify flat JSON structure
	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	// Verify stats key at root level
	statsMap, ok := response["stats"].(map[string]interface{})
	require.True(t, ok)

	// Verify per_pages is null when previous period has 0 pages
	assert.Nil(t, statsMap["per_pages"], "per_pages should be null when previous period has 0 pages")

	// Verify other fields are present
	assert.Equal(t, float64(100), statsMap["total_pages"])
	assert.NotNil(t, statsMap["max_day"])
	assert.NotNil(t, statsMap["mean_geral"])

	mockRepo.AssertExpectations(t)
}

// TestDashboardHandler_Projects tests the Projects handler with JSON:API format
func TestDashboardHandler_Projects(t *testing.T) {
	mockRepo := &MockDashboardRepository{}
	userConfig := service.NewUserConfigService(service.GetDefaultConfig())
	mockProjectsService := &MockProjectsService{}
	handler := NewDashboardHandler(mockRepo, userConfig, mockProjectsService)

	// Mock service calls for the new implementation - JSON:API format
	response := dto.NewDashboardProjectsResponse()

	stats := dto.NewDashboardStats()
	stats.SetTotalPages(0)
	stats.SetPages(0)
	stats.SetProgressGeral(0.0)
	response.SetStats(stats)

	mockProjectsService.On("GetDashboardProjects", mock.Anything).Return(response, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/projects.json", nil)
	w := httptest.NewRecorder()

	handler.Projects(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Parse JSON:API response format
	var jsonResponse map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&jsonResponse)
	require.NoError(t, err)

	// Verify response structure - data array + stats object
	assert.Contains(t, jsonResponse, "data")
	assert.Contains(t, jsonResponse, "stats")

	dataArr, ok := jsonResponse["data"].([]interface{})
	require.True(t, ok)
	assert.Empty(t, dataArr)

	statsObj, ok := jsonResponse["stats"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, 0.0, statsObj["total_pages"])
	assert.Equal(t, 0.0, statsObj["pages"])
	assert.Equal(t, 0.0, statsObj["progress_geral"])

	mockProjectsService.AssertExpectations(t)
}

// TestDashboardHandler_Projects_Empty tests Projects handler with no data
func TestDashboardHandler_Projects_Empty(t *testing.T) {
	mockRepo := &MockDashboardRepository{}
	userConfig := service.NewUserConfigService(service.GetDefaultConfig())
	mockProjectsService := &MockProjectsService{}
	handler := NewDashboardHandler(mockRepo, userConfig, mockProjectsService)

	// Mock service calls for the new implementation - JSON:API format
	response := dto.NewDashboardProjectsResponse()

	stats := dto.NewDashboardStats()
	stats.SetTotalPages(0)
	stats.SetPages(0)
	stats.SetProgressGeral(0.0)
	response.SetStats(stats)

	mockProjectsService.On("GetDashboardProjects", mock.Anything).Return(response, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/projects.json", nil)
	w := httptest.NewRecorder()

	handler.Projects(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Parse JSON:API response format
	var jsonResponse map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&jsonResponse)
	require.NoError(t, err)

	// Verify response structure - data array + stats object
	assert.Contains(t, jsonResponse, "data")
	assert.Contains(t, jsonResponse, "stats")

	dataArr, ok := jsonResponse["data"].([]interface{})
	require.True(t, ok)
	assert.Empty(t, dataArr)

	statsObj, ok := jsonResponse["stats"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, 0.0, statsObj["total_pages"])
	assert.Equal(t, 0.0, statsObj["pages"])
	assert.Equal(t, 0.0, statsObj["progress_geral"])

	mockProjectsService.AssertExpectations(t)
}

// TestDashboardHandler_Faults tests the Faults ECharts endpoint
func TestDashboardHandler_Faults(t *testing.T) {
	mockRepo := &MockDashboardRepository{}
	userConfig := service.NewUserConfigService(service.GetDefaultConfig())
	mockProjectsService := &MockProjectsService{}
	handler := NewDashboardHandler(mockRepo, userConfig, mockProjectsService)

	// Mock GetFaultsByDateRange to return 8 faults
	mockRepo.On("GetFaultsByDateRange", mock.Anything, mock.Anything, mock.Anything).
		Return(dto.NewFaultStats(8), nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/echart/faults.json", nil)
	w := httptest.NewRecorder()

	handler.Faults(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var envelope dto.JSONAPIEnvelope
	err := json.NewDecoder(w.Body).Decode(&envelope)
	require.NoError(t, err)

	echartData, ok := envelope.Data.(map[string]interface{})
	require.True(t, ok)

	attrsMap, ok := echartData["attributes"].(map[string]interface{})
	require.True(t, ok)

	title, ok := attrsMap["title"].(string)
	require.True(t, ok)
	assert.Equal(t, "Fault Percentage by Weekday", title)

	seriesArr, ok := attrsMap["series"].([]interface{})
	require.True(t, ok)
	assert.Len(t, seriesArr, 1)

	seriesMap, ok := seriesArr[0].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "gauge", seriesMap["type"])

	dataArr, ok := seriesMap["data"].([]interface{})
	require.True(t, ok)
	assert.Len(t, dataArr, 1)

	dataVal, ok := dataArr[0].(float64)
	require.True(t, ok)
	assert.Equal(t, 80.0, dataVal)

	mockRepo.AssertExpectations(t)
}

// TestDashboardHandler_SpeculateActual tests the SpeculateActual ECharts endpoint
func TestDashboardHandler_SpeculateActual(t *testing.T) {
	mockRepo := &MockDashboardRepository{}
	userConfig := service.NewUserConfigService(service.GetDefaultConfig())
	mockProjectsService := &MockProjectsService{}
	handler := NewDashboardHandler(mockRepo, userConfig, mockProjectsService)

	mockRepo.On("GetFaultsByDateRange", mock.Anything, mock.Anything, mock.Anything).
		Return(dto.NewFaultStats(50), nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/echart/speculate_actual.json", nil)
	w := httptest.NewRecorder()

	handler.SpeculateActual(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var envelope dto.JSONAPIEnvelope
	err := json.NewDecoder(w.Body).Decode(&envelope)
	require.NoError(t, err)

	echartData, ok := envelope.Data.(map[string]interface{})
	require.True(t, ok)

	attrsMap, ok := echartData["attributes"].(map[string]interface{})
	require.True(t, ok)

	title, ok := attrsMap["title"].(string)
	require.True(t, ok)
	assert.Equal(t, "Speculated vs Actual Faults", title)

	seriesArr, ok := attrsMap["series"].([]interface{})
	require.True(t, ok)
	assert.Len(t, seriesArr, 2)

	name0, ok := seriesArr[0].(map[string]interface{})["name"].(string)
	require.True(t, ok)
	assert.Equal(t, "Actual", name0)

	name1, ok := seriesArr[1].(map[string]interface{})["name"].(string)
	require.True(t, ok)
	assert.Equal(t, "Speculated", name1)

	data0, ok := seriesArr[0].(map[string]interface{})["data"].([]interface{})
	require.True(t, ok)
	actualVal, ok := data0[0].(float64)
	require.True(t, ok)
	assert.Equal(t, 50.0, actualVal)

	data1, ok := seriesArr[1].(map[string]interface{})["data"].([]interface{})
	require.True(t, ok)
	predictedVal, ok := data1[0].(float64)
	require.True(t, ok)
	// Predicted = 50 * (1 + 0.15) = 57.5
	assert.InDelta(t, 57.5, predictedVal, 1)

	mockRepo.AssertExpectations(t)
}

// TestDashboardHandler_WeekdayFaults tests the WeekdayFaults ECharts endpoint
func TestDashboardHandler_WeekdayFaults(t *testing.T) {
	mockRepo := &MockDashboardRepository{}
	userConfig := service.NewUserConfigService(service.GetDefaultConfig())
	mockProjectsService := &MockProjectsService{}
	handler := NewDashboardHandler(mockRepo, userConfig, mockProjectsService)

	mockRepo.On("GetWeekdayFaults", mock.Anything, mock.Anything, mock.Anything).
		Return(dto.NewWeekdayFaults(map[int]int{
			0: 5,
			1: 8,
			2: 0,
			3: 3,
			4: 0,
			5: 0,
			6: 0,
		}), nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/echart/faults_week_day.json", nil)
	w := httptest.NewRecorder()

	handler.WeekdayFaults(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var envelope dto.JSONAPIEnvelope
	err := json.NewDecoder(w.Body).Decode(&envelope)
	require.NoError(t, err)

	echartData, ok := envelope.Data.(map[string]interface{})
	require.True(t, ok)

	attrsMap, ok := echartData["attributes"].(map[string]interface{})
	require.True(t, ok)

	title, ok := attrsMap["title"].(string)
	require.True(t, ok)
	assert.Equal(t, "Faults by Weekday", title)

	seriesArr, ok := attrsMap["series"].([]interface{})
	require.True(t, ok)
	assert.Len(t, seriesArr, 1)

	seriesMap, ok := seriesArr[0].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "radar", seriesMap["type"])

	dataArr, ok := seriesMap["data"].([]interface{})
	require.True(t, ok)
	assert.Len(t, dataArr, 7)

	assert.Equal(t, float64(5), dataArr[0])
	assert.Equal(t, float64(8), dataArr[1])
	assert.Equal(t, float64(0), dataArr[2])
	assert.Equal(t, float64(3), dataArr[3])

	mockRepo.AssertExpectations(t)
}

// TestDashboardHandler_MeanProgress tests the MeanProgress ECharts endpoint
func TestDashboardHandler_MeanProgress(t *testing.T) {
	mockRepo := &MockDashboardRepository{}
	userConfig := service.NewUserConfigService(service.GetDefaultConfig())
	mockProjectsService := &MockProjectsService{}
	handler := NewDashboardHandler(mockRepo, userConfig, mockProjectsService)

	// Mock GetLogsByDateRange to return sample logs within the last 30 days
	note := "Test log"
	sampleLogs := []*dto.LogEntry{
		dto.NewLogEntry(1, time.Now().AddDate(0, 0, -10).Format(time.RFC3339), 0, 100, &note, nil),
	}
	mockRepo.On("GetLogsByDateRange", mock.Anything, mock.Anything, mock.Anything).
		Return(sampleLogs, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/echart/mean_progress.json", nil)
	w := httptest.NewRecorder()

	handler.MeanProgress(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var envelope dto.JSONAPIEnvelope
	err := json.NewDecoder(w.Body).Decode(&envelope)
	require.NoError(t, err)

	echartData, ok := envelope.Data.(map[string]interface{})
	require.True(t, ok)

	attrsMap, ok := echartData["attributes"].(map[string]interface{})
	require.True(t, ok)

	title, ok := attrsMap["title"].(string)
	require.True(t, ok)
	assert.Equal(t, "Mean Progress", title)

	seriesArr, ok := attrsMap["series"].([]interface{})
	require.True(t, ok)
	assert.Len(t, seriesArr, 1)

	seriesMap, ok := seriesArr[0].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "line", seriesMap["type"])

	dataArr, ok := seriesMap["data"].([]interface{})
	require.True(t, ok)
	// Should return 30 data points (one for each day in the last 30 days)
	assert.Len(t, dataArr, 30)

	mockRepo.AssertExpectations(t)
}

// TestDashboardHandler_MeanProgress_Empty tests MeanProgress with no projects
func TestDashboardHandler_MeanProgress_Empty(t *testing.T) {
	mockRepo := &MockDashboardRepository{}
	userConfig := service.NewUserConfigService(service.GetDefaultConfig())
	mockProjectsService := &MockProjectsService{}
	handler := NewDashboardHandler(mockRepo, userConfig, mockProjectsService)

	// Return empty logs slice to test empty data handling
	mockRepo.On("GetLogsByDateRange", mock.Anything, mock.Anything, mock.Anything).
		Return([]*dto.LogEntry{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/echart/mean_progress.json", nil)
	w := httptest.NewRecorder()

	handler.MeanProgress(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var envelope dto.JSONAPIEnvelope
	err := json.NewDecoder(w.Body).Decode(&envelope)
	require.NoError(t, err)

	echartData, ok := envelope.Data.(map[string]interface{})
	require.True(t, ok)

	attrsMap, ok := echartData["attributes"].(map[string]interface{})
	require.True(t, ok)

	seriesArr, ok := attrsMap["series"].([]interface{})
	require.True(t, ok)

	dataArr, ok := seriesArr[0].(map[string]interface{})["data"].([]interface{})
	require.True(t, ok)

	// When no logs are returned, the data array is empty
	// The service returns an empty slice for progressData when logs are empty
	if len(dataArr) > 0 {
		dataVal, ok := dataArr[0].(float64)
		require.True(t, ok)
		assert.Equal(t, 0.0, dataVal)
	} else {
		// Empty data array is expected when no logs exist
		assert.Empty(t, dataArr)
	}

	mockRepo.AssertExpectations(t)
}

// TestDashboardHandler_YearlyTotal tests the YearlyTotal ECharts endpoint
func TestDashboardHandler_YearlyTotal(t *testing.T) {
	mockRepo := &MockDashboardRepository{}
	userConfig := service.NewUserConfigService(service.GetDefaultConfig())
	mockProjectsService := &MockProjectsService{}
	handler := NewDashboardHandler(mockRepo, userConfig, mockProjectsService)

	// Mock GetLogsByDateRange to return sample logs
	note := "Test log"
	sampleLogs := []*dto.LogEntry{
		dto.NewLogEntry(1, time.Now().AddDate(0, 0, -10).Format(time.RFC3339), 0, 10, &note, nil),
		dto.NewLogEntry(2, time.Now().AddDate(0, 0, -20).Format(time.RFC3339), 10, 20, &note, nil),
		dto.NewLogEntry(3, time.Now().AddDate(0, 0, -30).Format(time.RFC3339), 20, 30, &note, nil),
	}

	mockRepo.On("GetLogsByDateRange", mock.Anything, mock.Anything, mock.Anything).
		Return(sampleLogs, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/echart/last_year_total.json", nil)
	w := httptest.NewRecorder()

	handler.YearlyTotal(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var envelope dto.JSONAPIEnvelope
	err := json.NewDecoder(w.Body).Decode(&envelope)
	require.NoError(t, err)

	echartData, ok := envelope.Data.(map[string]interface{})
	require.True(t, ok)

	attrsMap, ok := echartData["attributes"].(map[string]interface{})
	require.True(t, ok)

	title, ok := attrsMap["title"].(string)
	require.True(t, ok)
	assert.Equal(t, "Yearly Total Faults", title)

	seriesArr, ok := attrsMap["series"].([]interface{})
	require.True(t, ok)
	assert.Len(t, seriesArr, 1)

	// Check that it's a line chart
	series0, ok := seriesArr[0].(map[string]interface{})
	require.True(t, ok)
	seriesType, ok := series0["type"].(string)
	require.True(t, ok)
	assert.Equal(t, "line", seriesType)

	// Check that we have 52 data points (weekly aggregates)
	data, ok := series0["data"].([]interface{})
	require.True(t, ok)
	assert.Len(t, data, 52)

	mockRepo.AssertExpectations(t)
}
