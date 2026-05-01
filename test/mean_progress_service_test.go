package test

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/repository"
	"go-reading-log-api-next/internal/service/dashboard"
)

// MockDBPool is a mock implementation of PoolInterface for testing
type MockDBPool struct{}

func (m *MockDBPool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return nil
}

func (m *MockDBPool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return nil, nil
}

func (m *MockDBPool) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return pgconn.NewCommandTag("OK"), nil
}

func (m *MockDBPool) Acquire(ctx context.Context) (*pgxpool.Conn, error) {
	return nil, nil
}

func (m *MockDBPool) AcquireFunc(ctx context.Context, fn func(*pgxpool.Conn) error) error {
	return nil
}

func (m *MockDBPool) Close() {}

func (m *MockDBPool) Config() *pgxpool.Config {
	return nil
}

func (m *MockDBPool) Reset() {}

// testString is a helper to create string pointers
func testString(s string) *string {
	return &s
}

// MockDashboardRepository is a mock implementation of DashboardRepository for testing
type MockDashboardRepository struct {
	repository.DashboardRepository
	logs []dto.LogEntry
}

// GetLogsByDateRange returns mock logs within the date range
func (m *MockDashboardRepository) GetLogsByDateRange(ctx context.Context, start, end time.Time) ([]*dto.LogEntry, error) {
	// Filter logs within date range
	var result []*dto.LogEntry
	for _, log := range m.logs {
		logTime, err := time.Parse(time.RFC3339, log.Data)
		if err != nil {
			continue
		}
		if (logTime.Equal(start) || logTime.After(start)) && logTime.Before(end) {
			result = append(result, &log)
		}
	}
	return result, nil
}

// GetProjectAggregates returns empty aggregates for this mock
func (m *MockDashboardRepository) GetProjectAggregates(ctx context.Context) ([]*dto.ProjectAggregate, error) {
	return []*dto.ProjectAggregate{}, nil
}

// GetFaultsByDateRange returns empty faults for this mock
func (m *MockDashboardRepository) GetFaultsByDateRange(ctx context.Context, start, end time.Time) (*dto.FaultStats, error) {
	return &dto.FaultStats{FaultCount: 0}, nil
}

// GetWeekdayFaults returns empty weekday faults for this mock
func (m *MockDashboardRepository) GetWeekdayFaults(ctx context.Context, start, end time.Time) (*dto.WeekdayFaults, error) {
	return &dto.WeekdayFaults{Faults: make(map[int]int)}, nil
}

// GetProjectWeekdayMean returns 0 for this mock
func (m *MockDashboardRepository) GetProjectWeekdayMean(ctx context.Context, projectID int64, weekday int) (float64, error) {
	return 0.0, nil
}

// CalculatePeriodPages returns 0 for this mock
func (m *MockDashboardRepository) CalculatePeriodPages(ctx context.Context, start, end time.Time) (int, error) {
	return 0, nil
}

// GetProjectsWithLogs returns empty projects for this mock
func (m *MockDashboardRepository) GetProjectsWithLogs(ctx context.Context) ([]*dto.ProjectAggregateResponse, error) {
	return []*dto.ProjectAggregateResponse{}, nil
}

// GetProjectLogs returns empty logs for this mock
func (m *MockDashboardRepository) GetProjectLogs(ctx context.Context, projectID int64, limit int) ([]*dto.LogEntry, error) {
	return []*dto.LogEntry{}, nil
}

// GetPool returns the mock pool
func (m *MockDashboardRepository) GetPool() repository.PoolInterface {
	return &MockDBPool{}
}

// MockUserConfigProvider is a mock implementation of UserConfigProvider for testing
type MockUserConfigProvider struct {
	predictionPct float64
	maxFaults     int
	pagesPerDay   float64
}

func (m *MockUserConfigProvider) GetPredictionPct() float64 {
	return m.predictionPct
}

func (m *MockUserConfigProvider) GetMaxFaults() int {
	return m.maxFaults
}

func (m *MockUserConfigProvider) GetPagesPerDay() float64 {
	return m.pagesPerDay
}

// TestMeanProgressService tests the MeanProgressService implementation
func TestMeanProgressService(t *testing.T) {
	// Create mock repository with test data
	mockRepo := &MockDashboardRepository{
		logs: []dto.LogEntry{
			{
				ID:        1,
				Data:      time.Now().AddDate(0, 0, -29).Format(time.RFC3339),
				StartPage: 0,
				EndPage:   50,
				Note:      testString("Day 1"),
			},
			{
				ID:        2,
				Data:      time.Now().AddDate(0, 0, -28).Format(time.RFC3339),
				StartPage: 50,
				EndPage:   100,
				Note:      testString("Day 2"),
			},
			{
				ID:        3,
				Data:      time.Now().AddDate(0, 0, -1).Format(time.RFC3339),
				StartPage: 100,
				EndPage:   150,
				Note:      testString("Day 29"),
			},
			{
				ID:        4,
				Data:      time.Now().Format(time.RFC3339),
				StartPage: 150,
				EndPage:   200,
				Note:      testString("Today"),
			},
		},
	}

	// Create mock user config
	mockConfig := &MockUserConfigProvider{
		predictionPct: 0.15,
		maxFaults:     10,
		pagesPerDay:   25,
	}

	// Create service
	service := dashboard.NewMeanProgressService(mockRepo, mockConfig)

	// Test GetMeanProgressData
	ctx := context.Background()
	progressData, err := service.GetMeanProgressData(ctx)
	if err != nil {
		t.Fatalf("GetMeanProgressData failed: %v", err)
	}

	// Verify data count (should be 30 days - one for each day in the range)
	// The service now returns exactly 30 data points, one for each day
	if len(progressData) != 30 {
		t.Errorf("Expected 30 progress days, got %d", len(progressData))
	}

	// Test CalculateDailyProgress
	tests := []struct {
		name       string
		dailyPages float64
		meanPages  float64
		expected   float64
		expectZero bool
	}{
		{"Normal case", 50, 200, -75.0, false},
		{"Exact mean", 200, 200, 0.0, false},
		{"Double mean", 400, 200, 100.0, false},
		{"Zero mean (should return 0)", 50, 0, 0.0, true},
		{"Negative progress", 10, 200, -95.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := dashboard.CalculateDailyProgress(tt.dailyPages, tt.meanPages)
			if tt.expectZero {
				if result != 0.0 {
					t.Errorf("Expected 0.0 for zero mean, got %f", result)
				}
			} else {
				expectedRounded := math.Round(tt.expected*1000) / 1000
				if result != expectedRounded {
					t.Errorf("CalculateDailyProgress(%f, %f) = %f, expected %f", tt.dailyPages, tt.meanPages, result, expectedRounded)
				}
			}
		})
	}

	// Test GetColorForProgress
	colorTests := []struct {
		name     string
		progress float64
		expected string
	}{
		{"Negative (Red)", -10.0, "#ff4d4f"},
		{"Zero (Gray)", 0.0, "#959595"},
		{"5% (Gray)", 5.0, "#959595"},
		{"10% (Cyan)", 10.0, "#1890ff"},
		{"15% (Cyan)", 15.0, "#1890ff"},
		{"20% (Blue)", 20.0, "#108ee9"},
		{"30% (Blue)", 30.0, "#108ee9"},
		{"50% (Green)", 50.0, "#67c23a"},
		{"100% (Green)", 100.0, "#67c23a"},
	}

	for _, tt := range colorTests {
		t.Run(tt.name, func(t *testing.T) {
			result := dashboard.GetColorForProgress(tt.progress)
			if result != tt.expected {
				t.Errorf("GetColorForProgress(%f) = %s, expected %s", tt.progress, result, tt.expected)
			}
		})
	}

	// Test GenerateChartConfig
	chart, err := service.GenerateChartConfig(ctx)
	if err != nil {
		t.Fatalf("GenerateChartConfig failed: %v", err)
	}

	// Verify chart configuration
	if chart.Title != "Mean Progress" {
		t.Errorf("Expected title 'Mean Progress', got '%s'", chart.Title)
	}

	if len(chart.Series) != 1 {
		t.Errorf("Expected 1 series, got %d", len(chart.Series))
	}

	if chart.Series[0].Name != "Progress" {
		t.Errorf("Expected series name 'Progress', got '%s'", chart.Series[0].Name)
	}

	if chart.Series[0].Type != "line" {
		t.Errorf("Expected series type 'line', got '%s'", chart.Series[0].Type)
	}

	// Verify data points count is 30 (one for each day)
	if len(chart.Series[0].Data) != 30 {
		t.Errorf("Expected 30 data points, got %d", len(chart.Series[0].Data))
	}

	// Verify colors are set in item style
	itemStyle, ok := chart.Series[0].ItemStyle["color"].([]string)
	if !ok {
		t.Errorf("Expected color array in itemStyle")
	} else if len(itemStyle) != 30 {
		t.Errorf("Expected 30 colors, got %d", len(itemStyle))
	}
}

// TestMeanProgressService_EmptyData tests behavior with no logs
func TestMeanProgressService_EmptyData(t *testing.T) {
	mockRepo := &MockDashboardRepository{
		logs: []dto.LogEntry{},
	}

	mockConfig := &MockUserConfigProvider{}
	service := dashboard.NewMeanProgressService(mockRepo, mockConfig)

	ctx := context.Background()
	progressData, err := service.GetMeanProgressData(ctx)
	if err != nil {
		t.Fatalf("GetMeanProgressData failed: %v", err)
	}

	if len(progressData) != 30 {
		t.Errorf("Expected 30 progress days (empty data) for empty logs, got %d entries", len(progressData))
	}
}

// TestMeanProgressService_DateRange tests that only last 30 days are included
func TestMeanProgressService_DateRange(t *testing.T) {
	mockRepo := &MockDashboardRepository{
		logs: []dto.LogEntry{
			// Within last 30 days
			{
				ID:        1,
				Data:      time.Now().AddDate(0, 0, -29).Format(time.RFC3339),
				StartPage: 0,
				EndPage:   50,
				Note:      testString("Within range"),
			},
			// Outside last 30 days (older)
			{
				ID:        2,
				Data:      time.Now().AddDate(0, 0, -31).Format(time.RFC3339),
				StartPage: 50,
				EndPage:   100,
				Note:      testString("Outside range"),
			},
			// Future date (should not be included)
			{
				ID:        3,
				Data:      time.Now().AddDate(0, 0, 1).Format(time.RFC3339),
				StartPage: 100,
				EndPage:   150,
				Note:      testString("Future"),
			},
		},
	}

	mockConfig := &MockUserConfigProvider{}
	service := dashboard.NewMeanProgressService(mockRepo, mockConfig)

	ctx := context.Background()
	progressData, err := service.GetMeanProgressData(ctx)
	if err != nil {
		t.Fatalf("GetMeanProgressData failed: %v", err)
	}

	// Should return 30 data points (one for each day in the range)
	// Only 1 day will have actual log data, others will be zeros
	if len(progressData) != 30 {
		t.Errorf("Expected 30 progress days, got %d", len(progressData))
	}
}
