package unit

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

	"go-reading-log-api-next/internal/api/v1/handlers"
	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/service"
	"go-reading-log-api-next/internal/service/dashboard"
	"go-reading-log-api-next/test/testutil"
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

// =============================================================================
// PerMeanDay Tests
// =============================================================================

// TestDashboardHandler_Day_PerMeanDay_WithData tests per_mean_day calculation when previous data exists
func TestDashboardHandler_Day_PerMeanDay_WithData(t *testing.T) {
	mockRepo := testutil.NewMockDashboardRepository()
	userConfig := service.NewUserConfigService(service.GetDefaultConfig())
	mockProjectsService := &MockProjectsService{}
	handler := handlers.NewDashboardHandler(mockRepo, userConfig, mockProjectsService)

	testDate := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC) // Monday
	expectedStats := dto.NewDailyStats(140, 7)                 // 140 pages, 7 logs = mean_day = 20

	// Mock GetDailyStats
	mockRepo.On("GetDailyStats", mock.Anything, testDate).Return(expectedStats, nil)
	prevDate := testDate.AddDate(0, 0, -7)
	mockRepo.On("GetDailyStats", mock.Anything, prevDate).Return(dto.NewDailyStats(140, 7), nil)

	// Mock GetProjectAggregates
	mockRepo.On("GetProjectAggregates", mock.Anything).Return([]*dto.ProjectAggregate{}, nil)

	// Mock previous period mean (20.0)
	prevMean := 20.0
	mockRepo.On("GetPreviousPeriodMean", mock.Anything, testDate).Return(&prevMean, nil)

	// Mock other required methods
	maxDay := 50.0
	mockRepo.On("GetMaxByWeekday", mock.Anything, testDate).Return(&maxDay, nil)
	meanGeral := 25.0
	mockRepo.On("GetOverallMean", mock.Anything, testDate).Return(&meanGeral, nil)
	prevSpecMean := 23.0
	mockRepo.On("GetPreviousPeriodSpecMean", mock.Anything, testDate).Return(&prevSpecMean, nil)
	// Mock GetMeanByWeekday - return 20.0 (using V1::MeanLog algorithm)
	meanDay := 20.0
	mockRepo.On("GetMeanByWeekday", mock.Anything, 1).Return(&meanDay, nil) // Monday = weekday 1

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/day.json?date=2024-01-15T10:30:00Z", nil)
	w := httptest.NewRecorder()

	handler.Day(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	statsMap := response["stats"].(map[string]interface{})
	assert.NotNil(t, statsMap["per_mean_day"])
	assert.InDelta(t, 1.0, statsMap["per_mean_day"], 0.001) // 20/20 = 1.0

	mockRepo.AssertExpectations(t)
}

// TestDashboardHandler_Day_PerMeanDay_NoPreviousData tests per_mean_day returns nil when no previous data
func TestDashboardHandler_Day_PerMeanDay_NoPreviousData(t *testing.T) {
	mockRepo := &testutil.MockDashboardRepository{}
	userConfig := service.NewUserConfigService(service.GetDefaultConfig())
	mockProjectsService := &MockProjectsService{}
	handler := handlers.NewDashboardHandler(mockRepo, userConfig, mockProjectsService)

	testDate := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC) // Monday
	expectedStats := dto.NewDailyStats(140, 7)                 // mean_day = 20

	// Mock GetDailyStats
	mockRepo.On("GetDailyStats", mock.Anything, testDate).Return(expectedStats, nil)
	prevDate := testDate.AddDate(0, 0, -7)
	mockRepo.On("GetDailyStats", mock.Anything, prevDate).Return(dto.NewDailyStats(0, 0), nil)

	// Mock GetProjectAggregates
	mockRepo.On("GetProjectAggregates", mock.Anything).Return([]*dto.ProjectAggregate{}, nil)

	// Mock previous period mean = nil (no data)
	mockRepo.On("GetPreviousPeriodMean", mock.Anything, testDate).Return(nil, nil)

	// Mock other required methods
	maxDay := 50.0
	mockRepo.On("GetMaxByWeekday", mock.Anything, testDate).Return(&maxDay, nil)
	meanGeral := 25.0
	mockRepo.On("GetOverallMean", mock.Anything, testDate).Return(&meanGeral, nil)
	mockRepo.On("GetPreviousPeriodSpecMean", mock.Anything, testDate).Return(nil, nil)
	// Mock GetMeanByWeekday - return 20.0 (using V1::MeanLog algorithm)
	meanDay := 20.0
	mockRepo.On("GetMeanByWeekday", mock.Anything, 1).Return(&meanDay, nil) // Monday = weekday 1

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/day.json?date=2024-01-15T10:30:00Z", nil)
	w := httptest.NewRecorder()

	handler.Day(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	statsMap := response["stats"].(map[string]interface{})
	assert.Nil(t, statsMap["per_mean_day"])

	mockRepo.AssertExpectations(t)
}

// TestDashboardHandler_Day_PerMeanDay_ZeroPreviousData tests per_mean_day returns nil when previous mean is 0
func TestDashboardHandler_Day_PerMeanDay_ZeroPreviousData(t *testing.T) {
	mockRepo := &testutil.MockDashboardRepository{}
	userConfig := service.NewUserConfigService(service.GetDefaultConfig())
	mockProjectsService := &MockProjectsService{}
	handler := handlers.NewDashboardHandler(mockRepo, userConfig, mockProjectsService)

	testDate := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC) // Monday
	expectedStats := dto.NewDailyStats(140, 7)                 // mean_day = 20

	// Mock GetDailyStats
	mockRepo.On("GetDailyStats", mock.Anything, testDate).Return(expectedStats, nil)
	prevDate := testDate.AddDate(0, 0, -7)
	mockRepo.On("GetDailyStats", mock.Anything, prevDate).Return(dto.NewDailyStats(0, 0), nil)

	// Mock GetProjectAggregates
	mockRepo.On("GetProjectAggregates", mock.Anything).Return([]*dto.ProjectAggregate{}, nil)

	// Mock previous period mean = 0 (avoids division by zero)
	zeroMean := 0.0
	mockRepo.On("GetPreviousPeriodMean", mock.Anything, testDate).Return(&zeroMean, nil)

	// Mock other required methods
	maxDay := 50.0
	mockRepo.On("GetMaxByWeekday", mock.Anything, testDate).Return(&maxDay, nil)
	meanGeral := 25.0
	mockRepo.On("GetOverallMean", mock.Anything, testDate).Return(&meanGeral, nil)
	mockZeroSpecMean := 0.0
	mockRepo.On("GetPreviousPeriodSpecMean", mock.Anything, testDate).Return(&mockZeroSpecMean, nil)
	// Mock GetMeanByWeekday - return 20.0 (using V1::MeanLog algorithm)
	meanDay := 20.0
	mockRepo.On("GetMeanByWeekday", mock.Anything, 1).Return(&meanDay, nil) // Monday = weekday 1

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/day.json?date=2024-01-15T10:30:00Z", nil)
	w := httptest.NewRecorder()

	handler.Day(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	statsMap := response["stats"].(map[string]interface{})
	assert.Nil(t, statsMap["per_mean_day"])

	mockRepo.AssertExpectations(t)
}

// TestDashboardHandler_Day_PerMeanDay_RatioGreaterThan1 tests ratio when current > previous
func TestDashboardHandler_Day_PerMeanDay_RatioGreaterThan1(t *testing.T) {
	mockRepo := &testutil.MockDashboardRepository{}
	userConfig := service.NewUserConfigService(service.GetDefaultConfig())
	mockProjectsService := &MockProjectsService{}
	handler := handlers.NewDashboardHandler(mockRepo, userConfig, mockProjectsService)

	testDate := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC) // Monday
	expectedStats := dto.NewDailyStats(210, 7)                 // mean_day = 30

	// Mock GetDailyStats
	mockRepo.On("GetDailyStats", mock.Anything, testDate).Return(expectedStats, nil)
	prevDate := testDate.AddDate(0, 0, -7)
	mockRepo.On("GetDailyStats", mock.Anything, prevDate).Return(dto.NewDailyStats(140, 7), nil)

	// Mock GetProjectAggregates
	mockRepo.On("GetProjectAggregates", mock.Anything).Return([]*dto.ProjectAggregate{}, nil)

	// Mock previous period mean = 20.0
	prevMean := 20.0
	mockRepo.On("GetPreviousPeriodMean", mock.Anything, testDate).Return(&prevMean, nil)

	// Mock other required methods
	maxDay := 50.0
	mockRepo.On("GetMaxByWeekday", mock.Anything, testDate).Return(&maxDay, nil)
	meanGeral := 25.0
	mockRepo.On("GetOverallMean", mock.Anything, testDate).Return(&meanGeral, nil)
	prevSpecMean := 23.0
	mockRepo.On("GetPreviousPeriodSpecMean", mock.Anything, testDate).Return(&prevSpecMean, nil)
	// Mock GetMeanByWeekday - return 30.0 (using V1::MeanLog algorithm)
	meanDay := 30.0
	mockRepo.On("GetMeanByWeekday", mock.Anything, 1).Return(&meanDay, nil) // Monday = weekday 1

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/day.json?date=2024-01-15T10:30:00Z", nil)
	w := httptest.NewRecorder()

	handler.Day(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	statsMap := response["stats"].(map[string]interface{})
	assert.NotNil(t, statsMap["per_mean_day"])
	assert.InDelta(t, 1.5, statsMap["per_mean_day"], 0.001) // 30/20 = 1.5

	mockRepo.AssertExpectations(t)
}

// TestDashboardHandler_Day_PerMeanDay_RatioLessThan1 tests ratio when current < previous
func TestDashboardHandler_Day_PerMeanDay_RatioLessThan1(t *testing.T) {
	mockRepo := &testutil.MockDashboardRepository{}
	userConfig := service.NewUserConfigService(service.GetDefaultConfig())
	mockProjectsService := &MockProjectsService{}
	handler := handlers.NewDashboardHandler(mockRepo, userConfig, mockProjectsService)

	testDate := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC) // Monday
	expectedStats := dto.NewDailyStats(70, 7)                  // mean_day = 10

	// Mock GetDailyStats
	mockRepo.On("GetDailyStats", mock.Anything, testDate).Return(expectedStats, nil)
	prevDate := testDate.AddDate(0, 0, -7)
	mockRepo.On("GetDailyStats", mock.Anything, prevDate).Return(dto.NewDailyStats(140, 7), nil)

	// Mock GetProjectAggregates
	mockRepo.On("GetProjectAggregates", mock.Anything).Return([]*dto.ProjectAggregate{}, nil)

	// Mock previous period mean = 20.0
	prevMean := 20.0
	mockRepo.On("GetPreviousPeriodMean", mock.Anything, testDate).Return(&prevMean, nil)

	// Mock other required methods
	maxDay := 50.0
	mockRepo.On("GetMaxByWeekday", mock.Anything, testDate).Return(&maxDay, nil)
	meanGeral := 25.0
	mockRepo.On("GetOverallMean", mock.Anything, testDate).Return(&meanGeral, nil)
	prevSpecMean := 23.0
	mockRepo.On("GetPreviousPeriodSpecMean", mock.Anything, testDate).Return(&prevSpecMean, nil)
	// Mock GetMeanByWeekday - return 10.0 (using V1::MeanLog algorithm)
	meanDay := 10.0
	mockRepo.On("GetMeanByWeekday", mock.Anything, 1).Return(&meanDay, nil) // Monday = weekday 1

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/day.json?date=2024-01-15T10:30:00Z", nil)
	w := httptest.NewRecorder()

	handler.Day(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	statsMap := response["stats"].(map[string]interface{})
	assert.NotNil(t, statsMap["per_mean_day"])
	assert.InDelta(t, 0.5, statsMap["per_mean_day"], 0.001) // 10/20 = 0.5

	mockRepo.AssertExpectations(t)
}

// TestDashboardHandler_Day_PerMeanDay_Rounding tests 3 decimal place rounding
func TestDashboardHandler_Day_PerMeanDay_Rounding(t *testing.T) {
	mockRepo := &testutil.MockDashboardRepository{}
	userConfig := service.NewUserConfigService(service.GetDefaultConfig())
	mockProjectsService := &MockProjectsService{}
	handler := handlers.NewDashboardHandler(mockRepo, userConfig, mockProjectsService)

	testDate := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC) // Monday
	// mean_day = 23.333... (163.333 / 7)
	expectedStats := dto.NewDailyStats(163, 7)

	// Mock GetDailyStats
	mockRepo.On("GetDailyStats", mock.Anything, testDate).Return(expectedStats, nil)
	prevDate := testDate.AddDate(0, 0, -7)
	mockRepo.On("GetDailyStats", mock.Anything, prevDate).Return(dto.NewDailyStats(140, 7), nil)

	// Mock GetProjectAggregates
	mockRepo.On("GetProjectAggregates", mock.Anything).Return([]*dto.ProjectAggregate{}, nil)

	// Mock previous period mean = 20.0
	prevMean := 20.0
	mockRepo.On("GetPreviousPeriodMean", mock.Anything, testDate).Return(&prevMean, nil)

	// Mock other required methods
	maxDay := 50.0
	mockRepo.On("GetMaxByWeekday", mock.Anything, testDate).Return(&maxDay, nil)
	meanGeral := 25.0
	mockRepo.On("GetOverallMean", mock.Anything, testDate).Return(&meanGeral, nil)
	prevSpecMean := 23.0
	mockRepo.On("GetPreviousPeriodSpecMean", mock.Anything, testDate).Return(&prevSpecMean, nil)
	// Mock GetMeanByWeekday - return 23.286 (using V1::MeanLog algorithm: 163/7 = 23.286)
	meanDay := 23.286
	mockRepo.On("GetMeanByWeekday", mock.Anything, 1).Return(&meanDay, nil) // Monday = weekday 1

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/day.json?date=2024-01-15T10:30:00Z", nil)
	w := httptest.NewRecorder()

	handler.Day(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	statsMap := response["stats"].(map[string]interface{})
	assert.NotNil(t, statsMap["per_mean_day"])
	// mean_day = 163/7 = 23.285714..., ratio = 23.285714/20 = 1.1642857...
	// Rounded to 3 decimals = 1.164
	assert.InDelta(t, 1.164, statsMap["per_mean_day"], 0.001)

	mockRepo.AssertExpectations(t)
}

// =============================================================================
// PerSpecMeanDay Tests
// =============================================================================

// TestDashboardHandler_Day_PerSpecMeanDay_WithData tests per_spec_mean_day calculation when previous data exists
func TestDashboardHandler_Day_PerSpecMeanDay_WithData(t *testing.T) {
	mockRepo := &testutil.MockDashboardRepository{}
	userConfig := service.NewUserConfigService(service.GetDefaultConfig())
	mockProjectsService := &MockProjectsService{}
	handler := handlers.NewDashboardHandler(mockRepo, userConfig, mockProjectsService)

	testDate := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC) // Monday
	expectedStats := dto.NewDailyStats(140, 7)                 // mean_day = 20, spec_mean_day = 23

	// Mock GetDailyStats
	mockRepo.On("GetDailyStats", mock.Anything, testDate).Return(expectedStats, nil)
	prevDate := testDate.AddDate(0, 0, -7)
	mockRepo.On("GetDailyStats", mock.Anything, prevDate).Return(dto.NewDailyStats(140, 7), nil)

	// Mock GetProjectAggregates
	mockRepo.On("GetProjectAggregates", mock.Anything).Return([]*dto.ProjectAggregate{}, nil)

	// Mock previous period mean
	prevMean := 20.0
	mockRepo.On("GetPreviousPeriodMean", mock.Anything, testDate).Return(&prevMean, nil)

	// Mock other required methods
	maxDay := 50.0
	mockRepo.On("GetMaxByWeekday", mock.Anything, testDate).Return(&maxDay, nil)
	meanGeral := 25.0
	mockRepo.On("GetOverallMean", mock.Anything, testDate).Return(&meanGeral, nil)

	// Mock previous period spec mean (23.0 = 20 * 1.15)
	prevSpecMean := 23.0
	mockRepo.On("GetPreviousPeriodSpecMean", mock.Anything, testDate).Return(&prevSpecMean, nil)
	// Mock GetMeanByWeekday - return 20.0 (using V1::MeanLog algorithm)
	meanDay := 20.0
	mockRepo.On("GetMeanByWeekday", mock.Anything, 1).Return(&meanDay, nil) // Monday = weekday 1

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/day.json?date=2024-01-15T10:30:00Z", nil)
	w := httptest.NewRecorder()

	handler.Day(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	statsMap := response["stats"].(map[string]interface{})
	assert.NotNil(t, statsMap["per_spec_mean_day"])
	assert.InDelta(t, 1.0, statsMap["per_spec_mean_day"], 0.001) // 23/23 = 1.0

	mockRepo.AssertExpectations(t)
}

// TestDashboardHandler_Day_PerSpecMeanDay_NoPreviousData tests per_spec_mean_day returns nil when no previous data
func TestDashboardHandler_Day_PerSpecMeanDay_NoPreviousData(t *testing.T) {
	mockRepo := &testutil.MockDashboardRepository{}
	userConfig := service.NewUserConfigService(service.GetDefaultConfig())
	mockProjectsService := &MockProjectsService{}
	handler := handlers.NewDashboardHandler(mockRepo, userConfig, mockProjectsService)

	testDate := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC) // Monday
	expectedStats := dto.NewDailyStats(140, 7)                 // mean_day = 20, spec_mean_day = 23

	// Mock GetDailyStats
	mockRepo.On("GetDailyStats", mock.Anything, testDate).Return(expectedStats, nil)
	prevDate := testDate.AddDate(0, 0, -7)
	mockRepo.On("GetDailyStats", mock.Anything, prevDate).Return(dto.NewDailyStats(0, 0), nil)

	// Mock GetProjectAggregates
	mockRepo.On("GetProjectAggregates", mock.Anything).Return([]*dto.ProjectAggregate{}, nil)

	// Mock previous period mean
	prevMean := 20.0
	mockRepo.On("GetPreviousPeriodMean", mock.Anything, testDate).Return(&prevMean, nil)

	// Mock other required methods
	maxDay := 50.0
	mockRepo.On("GetMaxByWeekday", mock.Anything, testDate).Return(&maxDay, nil)
	meanGeral := 25.0
	mockRepo.On("GetOverallMean", mock.Anything, testDate).Return(&meanGeral, nil)

	// Mock previous period spec mean = nil (no data)
	mockRepo.On("GetPreviousPeriodSpecMean", mock.Anything, testDate).Return(nil, nil)
	// Mock GetMeanByWeekday - return 20.0 (using V1::MeanLog algorithm)
	meanDay := 20.0
	mockRepo.On("GetMeanByWeekday", mock.Anything, 1).Return(&meanDay, nil) // Monday = weekday 1

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/day.json?date=2024-01-15T10:30:00Z", nil)
	w := httptest.NewRecorder()

	handler.Day(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	statsMap := response["stats"].(map[string]interface{})
	assert.Nil(t, statsMap["per_spec_mean_day"])

	mockRepo.AssertExpectations(t)
}

// TestDashboardHandler_Day_PerSpecMeanDay_ZeroPreviousData tests per_spec_mean_day returns nil when previous spec mean is 0
func TestDashboardHandler_Day_PerSpecMeanDay_ZeroPreviousData(t *testing.T) {
	mockRepo := &testutil.MockDashboardRepository{}
	userConfig := service.NewUserConfigService(service.GetDefaultConfig())
	mockProjectsService := &MockProjectsService{}
	handler := handlers.NewDashboardHandler(mockRepo, userConfig, mockProjectsService)

	testDate := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC) // Monday
	expectedStats := dto.NewDailyStats(140, 7)                 // mean_day = 20, spec_mean_day = 23

	// Mock GetDailyStats
	mockRepo.On("GetDailyStats", mock.Anything, testDate).Return(expectedStats, nil)
	prevDate := testDate.AddDate(0, 0, -7)
	mockRepo.On("GetDailyStats", mock.Anything, prevDate).Return(dto.NewDailyStats(0, 0), nil)

	// Mock GetProjectAggregates
	mockRepo.On("GetProjectAggregates", mock.Anything).Return([]*dto.ProjectAggregate{}, nil)

	// Mock previous period mean
	prevMean := 20.0
	mockRepo.On("GetPreviousPeriodMean", mock.Anything, testDate).Return(&prevMean, nil)

	// Mock other required methods
	maxDay := 50.0
	mockRepo.On("GetMaxByWeekday", mock.Anything, testDate).Return(&maxDay, nil)
	meanGeral := 25.0
	mockRepo.On("GetOverallMean", mock.Anything, testDate).Return(&meanGeral, nil)

	// Mock previous period spec mean = 0 (avoids division by zero)
	zeroSpecMean := 0.0
	mockRepo.On("GetPreviousPeriodSpecMean", mock.Anything, testDate).Return(&zeroSpecMean, nil)
	// Mock GetMeanByWeekday - return 20.0 (using V1::MeanLog algorithm)
	meanDay := 20.0
	mockRepo.On("GetMeanByWeekday", mock.Anything, 1).Return(&meanDay, nil) // Monday = weekday 1

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/day.json?date=2024-01-15T10:30:00Z", nil)
	w := httptest.NewRecorder()

	handler.Day(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	statsMap := response["stats"].(map[string]interface{})
	assert.Nil(t, statsMap["per_spec_mean_day"])

	mockRepo.AssertExpectations(t)
}

// TestDashboardHandler_Day_PerSpecMeanDay_Rounding tests 3 decimal place rounding for per_spec_mean_day
func TestDashboardHandler_Day_PerSpecMeanDay_Rounding(t *testing.T) {
	mockRepo := &testutil.MockDashboardRepository{}
	userConfig := service.NewUserConfigService(service.GetDefaultConfig())
	mockProjectsService := &MockProjectsService{}
	handler := handlers.NewDashboardHandler(mockRepo, userConfig, mockProjectsService)

	testDate := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC) // Monday
	// mean_day = 23.285714..., spec_mean_day = 23.285714 * 1.15 = 26.778571...
	expectedStats := dto.NewDailyStats(163, 7)

	// Mock GetDailyStats
	mockRepo.On("GetDailyStats", mock.Anything, testDate).Return(expectedStats, nil)
	prevDate := testDate.AddDate(0, 0, -7)
	mockRepo.On("GetDailyStats", mock.Anything, prevDate).Return(dto.NewDailyStats(140, 7), nil)

	// Mock GetProjectAggregates
	mockRepo.On("GetProjectAggregates", mock.Anything).Return([]*dto.ProjectAggregate{}, nil)

	// Mock previous period mean
	prevMean := 20.0
	mockRepo.On("GetPreviousPeriodMean", mock.Anything, testDate).Return(&prevMean, nil)

	// Mock other required methods
	maxDay := 50.0
	mockRepo.On("GetMaxByWeekday", mock.Anything, testDate).Return(&maxDay, nil)
	meanGeral := 25.0
	mockRepo.On("GetOverallMean", mock.Anything, testDate).Return(&meanGeral, nil)

	// Mock previous period spec mean (23.0 = 20 * 1.15)
	prevSpecMean := 23.0
	mockRepo.On("GetPreviousPeriodSpecMean", mock.Anything, testDate).Return(&prevSpecMean, nil)
	// Mock GetMeanByWeekday - return 23.286 (using V1::MeanLog algorithm: 163/7 = 23.286)
	meanDay := 23.286
	mockRepo.On("GetMeanByWeekday", mock.Anything, 1).Return(&meanDay, nil) // Monday = weekday 1

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/day.json?date=2024-01-15T10:30:00Z", nil)
	w := httptest.NewRecorder()

	handler.Day(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	statsMap := response["stats"].(map[string]interface{})
	assert.NotNil(t, statsMap["per_spec_mean_day"])
	// spec_mean_day = 23.285714 * 1.15 = 26.778571
	// prev_spec_mean = 23.0
	// ratio = 26.778571 / 23.0 = 1.1642857...
	// Rounded to 3 decimals = 1.164
	assert.InDelta(t, 1.164, statsMap["per_spec_mean_day"], 0.001)

	mockRepo.AssertExpectations(t)
}

// =============================================================================
// Combined Tests
// =============================================================================

// TestDashboardHandler_Day_PerMeanDayAndPerSpecMeanDay_Together tests both ratios are calculated correctly together
func TestDashboardHandler_Day_PerMeanDayAndPerSpecMeanDay_Together(t *testing.T) {
	mockRepo := &testutil.MockDashboardRepository{}
	userConfig := service.NewUserConfigService(service.GetDefaultConfig())
	mockProjectsService := &MockProjectsService{}
	handler := handlers.NewDashboardHandler(mockRepo, userConfig, mockProjectsService)

	testDate := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC) // Monday
	expectedStats := dto.NewDailyStats(210, 7)                 // mean_day = 30, spec_mean_day = 34.5

	// Mock GetDailyStats
	mockRepo.On("GetDailyStats", mock.Anything, testDate).Return(expectedStats, nil)
	prevDate := testDate.AddDate(0, 0, -7)
	mockRepo.On("GetDailyStats", mock.Anything, prevDate).Return(dto.NewDailyStats(140, 7), nil)

	// Mock GetProjectAggregates
	mockRepo.On("GetProjectAggregates", mock.Anything).Return([]*dto.ProjectAggregate{}, nil)

	// Mock previous period mean = 20.0
	prevMean := 20.0
	mockRepo.On("GetPreviousPeriodMean", mock.Anything, testDate).Return(&prevMean, nil)

	// Mock other required methods
	maxDay := 50.0
	mockRepo.On("GetMaxByWeekday", mock.Anything, testDate).Return(&maxDay, nil)
	meanGeral := 25.0
	mockRepo.On("GetOverallMean", mock.Anything, testDate).Return(&meanGeral, nil)

	// Mock previous period spec mean = 23.0 (20 * 1.15)
	prevSpecMean := 23.0
	mockRepo.On("GetPreviousPeriodSpecMean", mock.Anything, testDate).Return(&prevSpecMean, nil)
	// Mock GetMeanByWeekday - return 30.0 (using V1::MeanLog algorithm)
	meanDay := 30.0
	mockRepo.On("GetMeanByWeekday", mock.Anything, 1).Return(&meanDay, nil) // Monday = weekday 1

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/day.json?date=2024-01-15T10:30:00Z", nil)
	w := httptest.NewRecorder()

	handler.Day(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	statsMap := response["stats"].(map[string]interface{})

	// Check per_mean_day: 30/20 = 1.5
	assert.NotNil(t, statsMap["per_mean_day"])
	assert.InDelta(t, 1.5, statsMap["per_mean_day"], 0.001)

	// Check per_spec_mean_day: 34.5/23 = 1.5
	assert.NotNil(t, statsMap["per_spec_mean_day"])
	assert.InDelta(t, 1.5, statsMap["per_spec_mean_day"], 0.001)

	mockRepo.AssertExpectations(t)
}
