package unit

import (
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/test/testutil"
)

// =============================================================================
// TestDashboardResponse - Basic functionality tests
// =============================================================================

// TestDashboardResponse_BasicCreation tests basic response creation
func TestDashboardResponse_BasicCreation(t *testing.T) {
	// Arrange
	response := dto.NewDashboardResponse()

	// Act
	response.SetEchart(dto.NewEchartConfig()).
		SetStats(dto.NewStatsData()).
		AddLog(*dto.NewLogEntry(1, "2024-01-15T10:30:00Z", 0, 25, nil, nil))

	// Assert
	assert.NotNil(t, response)
	assert.NotNil(t, response.Echart)
	assert.NotNil(t, response.Stats)
	assert.Len(t, response.Logs, 1)
}

// TestDashboardResponse_JSONSerialization tests JSON marshaling
func TestDashboardResponse_JSONSerialization(t *testing.T) {
	// Arrange
	response := dto.NewDashboardResponse()

	project := dto.NewProject(1, "Test Project", 200, 50)
	log := dto.NewLogEntry(1, "2024-01-15T10:30:00Z", 0, 25, nil, project)
	response.AddLog(*log).
		SetStats(dto.NewStatsData().
			SetPreviousWeekPages(100).
			SetLastWeekPages(150).
			SetPerPages(testutil.FloatPtr(1.5)))

	// Act
	data, err := json.Marshal(response)

	// Assert
	require.NoError(t, err)

	// Verify JSON structure
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, 1.5, result["stats"].(map[string]interface{})["per_pages"])
	assert.Equal(t, float64(100), result["stats"].(map[string]interface{})["previous_week_pages"])
	assert.Equal(t, float64(150), result["stats"].(map[string]interface{})["last_week_pages"])
}

// TestDashboardResponse_Validation tests validation methods
func TestDashboardResponse_Validation(t *testing.T) {
	testCases := []struct {
		name        string
		response    *dto.DashboardResponse
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid response",
			response:    dto.NewDashboardResponse().SetStats(dto.NewStatsData()),
			expectError: false,
		},
		{
			name: "negative previous_week_pages",
			response: dto.NewDashboardResponse().
				SetStats(dto.NewStatsData().SetPreviousWeekPages(-10)),
			expectError: true,
			errorMsg:    "previous_week_pages cannot be negative",
		},
		{
			name: "invalid percentage",
			response: dto.NewDashboardResponse().
				SetStats(dto.NewStatsData().SetPerPages(testutil.FloatPtr(150))),
			expectError: false, // Removed 0-100 constraint
		},
		{
			name: "empty series",
			response: dto.NewDashboardResponse().
				SetEchart(dto.NewEchartConfigWithSeries([]dto.Series{})),
			expectError: true,
			errorMsg:    "series cannot be empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.response.Validate()

			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestDashboardResponse_EmptyValues tests handling of empty/zero values
func TestDashboardResponse_EmptyValues(t *testing.T) {
	// Arrange
	response := dto.NewDashboardResponse()

	// Act - Marshal with zero values
	data, err := json.Marshal(response)

	// Assert
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	// Empty slices should be marshaled as empty arrays or null (both valid)
	logs := result["logs"]
	if logs != nil {
		assert.Equal(t, []interface{}{}, logs)
	}

	// Nil optional fields should be omitted (omitempty)
	assert.NotContains(t, result, "echart")
	assert.NotContains(t, result, "stats")
}

// TestDashboardResponse_ConcurrentAccess tests thread safety with mutex protection
func TestDashboardResponse_ConcurrentAccess(t *testing.T) {
	// Note: DashboardResponse is not inherently thread-safe for concurrent writes.
	// This test verifies that the struct can be used safely when external synchronization is applied.
	response := dto.NewDashboardResponse()

	// Use a mutex to protect concurrent modifications
	var mu sync.Mutex
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(idx int) {
			log := dto.NewLogEntry(int64(idx), "2024-01-15T10:30:00Z", 0, 10, nil, nil)
			mu.Lock()
			response.AddLog(*log)
			mu.Unlock()
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Assert
	assert.Len(t, response.Logs, 10)
}

// TestDashboardResponse_WithAllFields tests response with all fields populated
func TestDashboardResponse_WithAllFields(t *testing.T) {
	// Arrange
	testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)

	// Create sample logs
	logs := []dto.LogEntry{
		*dto.NewLogEntry(1, testDate, 0, 25, nil, dto.NewProject(1, "Project A", 200, 50)),
		*dto.NewLogEntry(2, testDate, 25, 50, nil, dto.NewProject(1, "Project A", 200, 50)),
		*dto.NewLogEntry(3, testDate, 0, 30, nil, dto.NewProject(2, "Project B", 300, 80)),
	}

	// Create stats
	stats := dto.NewStatsData().
		SetPreviousWeekPages(150).
		SetLastWeekPages(200).
		SetPerPages(testutil.FloatPtr(1.333)).
		SetMeanDay(25.5).
		SetSpecMeanDay(29.325)

	// Create ECharts config
	echart := dto.NewEchartConfig().
		SetTitle("Daily Progress").
		AddSeries(dto.Series{
			Name: "Pages",
			Type: "line",
			Data: []interface{}{10, 20, 30, 40, 50},
		})

	// Build response
	response := dto.NewDashboardResponse().
		SetStats(stats).
		SetEchart(echart)

	for _, log := range logs {
		response.AddLog(log)
	}

	// Act - Marshal and unmarshal
	data, err := json.Marshal(response)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, 3, len(result["logs"].([]interface{})))
	assert.Equal(t, 1.333, result["stats"].(map[string]interface{})["per_pages"])
	assert.Equal(t, "Daily Progress", result["echart"].(map[string]interface{})["title"])
}

// =============================================================================
// TestEchartConfig - ECharts configuration tests
// =============================================================================

// TestEchartConfig_JSONSerialization tests ECharts config serialization
func TestEchartConfig_JSONSerialization(t *testing.T) {
	// Arrange
	config := dto.NewEchartConfig().
		SetTitle("Reading Progress").
		SetTooltip(map[string]interface{}{
			"trigger": "axis",
			"axisPointer": map[string]interface{}{
				"type": "shadow",
			},
		}).
		AddSeries(dto.Series{
			Name: "Pages Read",
			Type: "line",
			Data: []interface{}{10, 20, 30, 40, 50},
			ItemStyle: map[string]interface{}{
				"color": "#5470C6",
			},
		})

	// Act
	data, err := json.Marshal(config)

	// Assert
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, "Reading Progress", result["title"])
	assert.Equal(t, "axis", result["tooltip"].(map[string]interface{})["trigger"])
	assert.Len(t, result["series"].([]interface{}), 1)
}

// TestEchartConfig_Validation tests ECharts config validation
func TestEchartConfig_Validation(t *testing.T) {
	testCases := []struct {
		name        string
		config      *dto.EchartConfig
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid config",
			config:      dto.NewEchartConfig().AddSeries(dto.Series{Name: "Test", Type: "line", Data: []interface{}{1, 2, 3}}),
			expectError: false,
		},
		{
			name:        "empty series",
			config:      dto.NewEchartConfig(),
			expectError: true,
			errorMsg:    "series cannot be empty",
		},
		{
			name: "series without name",
			config: func() *dto.EchartConfig {
				cfg := dto.NewEchartConfig()
				cfg.AddSeries(dto.Series{Type: "line", Data: []interface{}{1, 2, 3}})
				return cfg
			}(),
			expectError: true,
			errorMsg:    "series name is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.Validate()

			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestEchartConfig_SeriesValidation tests series validation
func TestEchartConfig_SeriesValidation(t *testing.T) {
	// Valid config with series
	config := dto.NewEchartConfig().AddSeries(dto.Series{
		Name: "Test Series",
		Type: "line",
		Data: []interface{}{1, 2, 3, 4, 5},
	})
	err := config.Validate()
	assert.NoError(t, err)

	// Invalid series - empty data
	invalidSeries := dto.Series{
		Name: "Empty Data",
		Type: "bar",
		Data: []interface{}{},
	}
	configWithInvalid := dto.NewEchartConfig().AddSeries(invalidSeries)
	err = configWithInvalid.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "series data cannot be empty")
}

// =============================================================================
// TestStatsData - Statistics data tests
// =============================================================================

// TestStatsData_Rounding tests decimal rounding
func TestStatsData_Rounding(t *testing.T) {
	// Arrange
	perPages := 0.333333
	stats := dto.NewStatsData().
		SetPreviousWeekPages(100).
		SetLastWeekPages(33).
		SetPerPages(&perPages).
		SetMeanDay(12.345678).
		SetSpecMeanDay(14.987654)

	// Act
	stats.RoundToThreeDecimals()

	// Assert
	assert.Equal(t, 0.333, *stats.PerPages)
	assert.Equal(t, 12.346, stats.MeanDay)
	assert.Equal(t, 14.988, stats.SpecMeanDay)
}

// TestStatsData_Validation tests statistics validation
func TestStatsData_Validation(t *testing.T) {
	testCases := []struct {
		name        string
		stats       *dto.StatsData
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid stats",
			stats:       dto.NewStatsData().SetPerPages(testutil.FloatPtr(50.0)).SetProgressGeral(75.0),
			expectError: false,
		},
		{
			name: "negative previous_week_pages",
			stats: dto.NewStatsData().
				SetPreviousWeekPages(-10),
			expectError: true,
			errorMsg:    "previous_week_pages cannot be negative",
		},
		{
			name: "percentage over 100 - now allowed",
			stats: dto.NewStatsData().
				SetPerPages(testutil.FloatPtr(150)),
			expectError: false, // Removed 0-100 constraint
		},
		{
			name: "negative per_pages",
			stats: dto.NewStatsData().
				SetPerPages(testutil.FloatPtr(-10)),
			expectError: true,
			errorMsg:    "per_pages cannot be negative",
		},
		{
			name: "negative mean_day",
			stats: dto.NewStatsData().
				SetMeanDay(-5.0),
			expectError: true,
			errorMsg:    "mean_day cannot be negative",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.stats.Validate()

			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestStatsData_JSONSerialization tests StatsData JSON serialization
func TestStatsData_JSONSerialization(t *testing.T) {
	// Arrange
	stats := dto.NewStatsData().
		SetPreviousWeekPages(100).
		SetLastWeekPages(150).
		SetPerPages(testutil.FloatPtr(1.5)).
		SetMeanDay(25.5).
		SetProgressGeral(75.5)

	// Act
	data, err := json.Marshal(stats)

	// Assert
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, float64(100), result["previous_week_pages"])
	assert.Equal(t, float64(150), result["last_week_pages"])
	assert.Equal(t, 1.5, result["per_pages"])
	assert.Equal(t, 25.5, result["mean_day"])
	assert.Equal(t, 75.5, result["progress_geral"])
}

// =============================================================================
// TestLogEntry - Log entry tests
// =============================================================================

// TestLogEntry_CalculatedFields tests calculated fields
func TestLogEntry_CalculatedFields(t *testing.T) {
	// Arrange
	project := dto.NewProject(1, "Test Project", 200, 50)

	// Act
	log := dto.NewLogEntry(1, "2024-01-15T10:30:00Z", 10, 35, nil, project)

	// Assert
	assert.Equal(t, 25, log.ReadPages) // 35 - 10 = 25

	// Test validation with invalid page range
	invalidLog := dto.NewLogEntry(2, "2024-01-15T10:30:00Z", 35, 10, nil, project)
	err := invalidLog.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "read_pages cannot be negative")
}

// TestLogEntry_JSONSerialization tests LogEntry JSON serialization
func TestLogEntry_JSONSerialization(t *testing.T) {
	// Arrange
	project := dto.NewProject(1, "Test Project", 200, 50)
	log := dto.NewLogEntry(1, "2024-01-15T10:30:00Z", 10, 35, nil, project)

	// Act
	data, err := json.Marshal(log)

	// Assert
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, float64(1), result["id"])
	assert.Equal(t, "2024-01-15T10:30:00Z", result["data"])
	assert.Equal(t, float64(10), result["start_page"])
	assert.Equal(t, float64(35), result["end_page"])
	assert.Equal(t, float64(25), result["read_pages"])

	// Verify project is included
	assert.Contains(t, result, "project")
}

// TestLogEntry_WithNote tests log entry with note
func TestLogEntry_WithNote(t *testing.T) {
	// Arrange
	note := "Morning reading session"
	project := dto.NewProject(1, "Test Project", 200, 50)
	log := dto.NewLogEntry(1, "2024-01-15T10:30:00Z", 10, 35, &note, project)

	// Act
	data, err := json.Marshal(log)

	// Assert
	require.NoError(t, err)
	assert.Contains(t, string(data), "Morning reading session")
}

// TestLogEntry_Validation tests log entry validation
func TestLogEntry_Validation(t *testing.T) {
	testCases := []struct {
		name        string
		log         *dto.LogEntry
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid log",
			log:         dto.NewLogEntry(1, "2024-01-15T10:30:00Z", 0, 25, nil, dto.NewProject(1, "Test", 200, 50)),
			expectError: false,
		},
		{
			name:        "negative start_page",
			log:         dto.NewLogEntry(1, "2024-01-15T10:30:00Z", -5, 25, nil, nil),
			expectError: true,
			errorMsg:    "start_page cannot be negative",
		},
		{
			name:        "negative end_page",
			log:         dto.NewLogEntry(1, "2024-01-15T10:30:00Z", 10, -5, nil, nil),
			expectError: true,
			errorMsg:    "end_page cannot be negative",
		},
		{
			name:        "zero id",
			log:         dto.NewLogEntry(0, "2024-01-15T10:30:00Z", 10, 25, nil, nil),
			expectError: true,
			errorMsg:    "id must be positive",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.log.Validate()

			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// =============================================================================
// TestProject - Project DTO tests
// =============================================================================

// TestProject_JSONSerialization tests Project JSON serialization
func TestProject_JSONSerialization(t *testing.T) {
	// Arrange
	project := dto.NewProject(1, "Test Project", 200, 50)

	// Act
	data, err := json.Marshal(project)

	// Assert
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, float64(1), result["id"])
	assert.Equal(t, "Test Project", result["name"])
	assert.Equal(t, float64(200), result["total_page"])
	assert.Equal(t, float64(50), result["page"])
}

// TestProject_Validation tests Project validation
func TestProject_Validation(t *testing.T) {
	testCases := []struct {
		name        string
		project     *dto.Project
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid project",
			project:     dto.NewProject(1, "Test Project", 200, 50),
			expectError: false,
		},
		{
			name: "zero id",
			project: &dto.Project{
				ID: 0, Name: "Test", TotalPage: 200, Page: 50,
			},
			expectError: true,
			errorMsg:    "project id must be positive",
		},
		{
			name: "zero total_page",
			project: &dto.Project{
				ID: 1, Name: "Test", TotalPage: 0, Page: 50,
			},
			expectError: true,
			errorMsg:    "total_page must be positive",
		},
		{
			name: "page exceeds total",
			project: &dto.Project{
				ID: 1, Name: "Test", TotalPage: 100, Page: 150,
			},
			expectError: true,
			errorMsg:    "page cannot exceed total_page",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.project.Validate()

			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// =============================================================================
// TestIntegration - Integration tests
// =============================================================================

// TestIntegration_DashboardResponseWithRealData tests with simulated real data
func TestIntegration_DashboardResponseWithRealData(t *testing.T) {
	// Arrange - Simulate data from database
	testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)

	// Create sample logs
	logs := []dto.LogEntry{
		*dto.NewLogEntry(1, testDate, 0, 25, nil, dto.NewProject(1, "Project A", 200, 50)),
		*dto.NewLogEntry(2, testDate, 25, 50, nil, dto.NewProject(1, "Project A", 200, 50)),
		*dto.NewLogEntry(3, testDate, 0, 30, nil, dto.NewProject(2, "Project B", 300, 80)),
	}

	// Create stats
	stats := dto.NewStatsData().
		SetPreviousWeekPages(150).
		SetLastWeekPages(200).
		SetPerPages(testutil.FloatPtr(1.333)).
		SetMeanDay(25.5).
		SetSpecMeanDay(29.325)

	// Create ECharts config
	echart := dto.NewEchartConfig().
		SetTitle("Daily Progress").
		AddSeries(dto.Series{
			Name: "Pages",
			Type: "line",
			Data: []interface{}{10, 20, 30, 40, 50},
		})

	// Build response
	response := dto.NewDashboardResponse().
		SetStats(stats).
		SetEchart(echart)

	for _, log := range logs {
		response.AddLog(log)
	}

	// Act - Marshal and unmarshal
	data, err := json.Marshal(response)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, 3, len(result["logs"].([]interface{})))
	assert.Equal(t, 1.333, result["stats"].(map[string]interface{})["per_pages"])
	assert.Equal(t, "Daily Progress", result["echart"].(map[string]interface{})["title"])
}

// TestIntegration_RoundTrip tests JSON round-trip preservation
func TestIntegration_RoundTrip(t *testing.T) {
	// Arrange
	original := dto.NewDashboardResponse().
		SetStats(dto.NewStatsData().
			SetPreviousWeekPages(100).
			SetLastWeekPages(150).
			SetPerPages(testutil.FloatPtr(1.5))).
		AddLog(*dto.NewLogEntry(1, "2024-01-15T10:30:00Z", 0, 25, nil, dto.NewProject(1, "Test", 200, 50)))

	// Act
	data, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded dto.DashboardResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, original.Stats.PreviousWeekPages, decoded.Stats.PreviousWeekPages)
	assert.Equal(t, original.Stats.LastWeekPages, decoded.Stats.LastWeekPages)
	assert.Equal(t, original.Stats.PerPages, decoded.Stats.PerPages)
	assert.Len(t, decoded.Logs, len(original.Logs))
}

// =============================================================================
// TestEdgeCases - Edge case tests
// =============================================================================

// TestEdgeCases_NilValues handles nil pointer scenarios
func TestEdgeCases_NilValues(t *testing.T) {
	// Arrange
	response := dto.NewDashboardResponse()

	// Act - Marshal with nil optional fields
	data, err := json.Marshal(response)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	// Assert - Nil fields should be omitted
	assert.NotContains(t, result, "echart")
	assert.NotContains(t, result, "stats")

	// Test validation with nil response
	err = response.Validate()
	assert.NoError(t, err) // Should not error on empty response
}

// TestEdgeCases_MathRounding tests mathematical rounding precision
func TestEdgeCases_MathRounding(t *testing.T) {
	// Arrange
	perPages := 0.3335
	stats := dto.NewStatsData().
		SetPerPages(&perPages).   // Should round to 0.334
		SetMeanDay(12.3456789).   // Should round to 12.346
		SetProgressGeral(99.9999) // Should round to 100.0

	// Act
	stats.RoundToThreeDecimals()

	// Assert
	assert.Equal(t, 0.334, *stats.PerPages)
	assert.Equal(t, 12.346, stats.MeanDay)
	assert.Equal(t, 100.0, stats.ProgressGeral)
}

// TestEdgeCases_LargeNumbers handles large number scenarios
func TestEdgeCases_LargeNumbers(t *testing.T) {
	// Arrange
	stats := dto.NewStatsData().
		SetPreviousWeekPages(math.MaxInt32).
		SetLastWeekPages(math.MaxInt32).
		SetTotalPages(math.MaxInt32)

	// Act & Assert - Should handle large numbers without error
	err := stats.Validate()
	assert.NoError(t, err)

	data, err := json.Marshal(stats)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, float64(math.MaxInt32), result["previous_week_pages"])
}

// TestEdgeCases_SeriesTypes tests different chart series types
func TestEdgeCases_SeriesTypes(t *testing.T) {
	config := dto.NewEchartConfig()

	// Test various series types
	seriesTypes := []string{"line", "bar", "pie", "gauge", "radar"}

	for _, typ := range seriesTypes {
		t.Run(typ, func(t *testing.T) {
			series := dto.Series{
				Name: fmt.Sprintf("Series %s", typ),
				Type: typ,
				Data: []interface{}{1, 2, 3},
			}

			err := series.Validate()
			assert.NoError(t, err)

			config.AddSeries(series)
		})
	}

	err := config.Validate()
	assert.NoError(t, err)
	assert.Len(t, config.Series, len(seriesTypes))
}

// TestEdgeCases_EmptyStringNote handles empty string notes
func TestEdgeCases_EmptyStringNote(t *testing.T) {
	note := ""
	log := dto.NewLogEntry(1, "2024-01-15T10:30:00Z", 0, 25, &note, nil)

	data, err := json.Marshal(log)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, "", result["note"])
}

// =============================================================================
// TestStatsData_NewFields - Tests for new fields added in Phase 1
// =============================================================================

// TestStatsData_NewFields tests the new fields: MaxDay, MeanGeral, PerMeanDay, PerSpecMeanDay
func TestStatsData_NewFields(t *testing.T) {
	// Arrange
	maxDay := 50.5
	meanGeral := 25.75
	perMeanDay := 1.25
	perSpecMeanDay := 1.35

	stats := dto.NewStatsData().
		SetMaxDay(&maxDay).
		SetMeanGeral(&meanGeral).
		SetPerMeanDay(&perMeanDay).
		SetPerSpecMeanDay(&perSpecMeanDay)

	// Assert - Verify fields are set correctly
	assert.NotNil(t, stats.MaxDay)
	assert.Equal(t, 50.5, *stats.MaxDay)
	assert.NotNil(t, stats.MeanGeral)
	assert.Equal(t, 25.75, *stats.MeanGeral)
	assert.NotNil(t, stats.PerMeanDay)
	assert.Equal(t, 1.25, *stats.PerMeanDay)
	assert.NotNil(t, stats.PerSpecMeanDay)
	assert.Equal(t, 1.35, *stats.PerSpecMeanDay)
}

// TestStatsData_NewFields_JSONSerialization tests JSON serialization of new fields
func TestStatsData_NewFields_JSONSerialization(t *testing.T) {
	// Arrange
	maxDay := 50.5
	meanGeral := 25.75
	perMeanDay := 1.25
	perSpecMeanDay := 1.35

	stats := dto.NewStatsData().
		SetMaxDay(&maxDay).
		SetMeanGeral(&meanGeral).
		SetPerMeanDay(&perMeanDay).
		SetPerSpecMeanDay(&perSpecMeanDay)

	// Act
	data, err := json.Marshal(stats)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	// Assert - Verify new fields are serialized correctly
	assert.Equal(t, 50.5, result["max_day"])
	assert.Equal(t, 25.75, result["mean_geral"])
	assert.Equal(t, 1.25, result["per_mean_day"])
	assert.Equal(t, 1.35, result["per_spec_mean_day"])
}

// TestStatsData_NewFields_NullHandling tests null handling for new fields
func TestStatsData_NewFields_NullHandling(t *testing.T) {
	// Arrange - Create StatsData with nil new fields
	stats := dto.NewStatsData()

	// Act
	data, err := json.Marshal(stats)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	// Assert - New fields should be omitted when nil (omitempty)
	assert.NotContains(t, result, "max_day")
	assert.NotContains(t, result, "mean_geral")
	assert.NotContains(t, result, "per_mean_day")
	assert.NotContains(t, result, "per_spec_mean_day")
}

// TestStatsData_NewFields_SetterChaining tests method chaining with new setter methods
func TestStatsData_NewFields_SetterChaining(t *testing.T) {
	// Arrange
	maxDay := 50.5
	meanGeral := 25.75
	perMeanDay := 1.25
	perSpecMeanDay := 1.35

	// Act - Chain setter methods
	stats := dto.NewStatsData().
		SetMaxDay(&maxDay).
		SetMeanGeral(&meanGeral).
		SetPerMeanDay(&perMeanDay).
		SetPerSpecMeanDay(&perSpecMeanDay)

	// Assert - Verify all fields are set
	assert.NotNil(t, stats.MaxDay)
	assert.Equal(t, 50.5, *stats.MaxDay)
	assert.NotNil(t, stats.MeanGeral)
	assert.Equal(t, 25.75, *stats.MeanGeral)
	assert.NotNil(t, stats.PerMeanDay)
	assert.Equal(t, 1.25, *stats.PerMeanDay)
	assert.NotNil(t, stats.PerSpecMeanDay)
	assert.Equal(t, 1.35, *stats.PerSpecMeanDay)
}

// TestStatsData_NewFields_Rounding tests rounding of new pointer fields
func TestStatsData_NewFields_Rounding(t *testing.T) {
	// Arrange
	maxDay := 50.123456
	meanGeral := 25.987654
	perMeanDay := 1.333333
	perSpecMeanDay := 1.444444

	stats := dto.NewStatsData().
		SetMaxDay(&maxDay).
		SetMeanGeral(&meanGeral).
		SetPerMeanDay(&perMeanDay).
		SetPerSpecMeanDay(&perSpecMeanDay)

	// Act
	stats.RoundToThreeDecimals()

	// Assert
	assert.Equal(t, 50.123, *stats.MaxDay)
	assert.Equal(t, 25.988, *stats.MeanGeral)
	assert.Equal(t, 1.333, *stats.PerMeanDay)
	assert.Equal(t, 1.444, *stats.PerSpecMeanDay)
}

// TestStatsData_NewFields_Validation tests validation of new fields
func TestStatsData_NewFields_Validation(t *testing.T) {
	testCases := []struct {
		name        string
		stats       *dto.StatsData
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid new fields",
			stats: func() *dto.StatsData {
				maxDay := 50.5
				meanGeral := 25.75
				perMeanDay := 1.25
				perSpecMeanDay := 1.35
				return dto.NewStatsData().
					SetMaxDay(&maxDay).
					SetMeanGeral(&meanGeral).
					SetPerMeanDay(&perMeanDay).
					SetPerSpecMeanDay(&perSpecMeanDay)
			}(),
			expectError: false,
		},
		{
			name:        "nil new fields - should pass",
			stats:       dto.NewStatsData(),
			expectError: false,
		},
		{
			name: "negative max_day",
			stats: func() *dto.StatsData {
				maxDay := -10.0
				return dto.NewStatsData().SetMaxDay(&maxDay)
			}(),
			expectError: true,
			errorMsg:    "max_day cannot be negative",
		},
		{
			name: "negative mean_geral",
			stats: func() *dto.StatsData {
				meanGeral := -5.0
				return dto.NewStatsData().SetMeanGeral(&meanGeral)
			}(),
			expectError: true,
			errorMsg:    "mean_geral cannot be negative",
		},
		{
			name: "negative per_mean_day",
			stats: func() *dto.StatsData {
				perMeanDay := -1.0
				return dto.NewStatsData().SetPerMeanDay(&perMeanDay)
			}(),
			expectError: true,
			errorMsg:    "per_mean_day cannot be negative",
		},
		{
			name: "negative per_spec_mean_day",
			stats: func() *dto.StatsData {
				perSpecMeanDay := -2.0
				return dto.NewStatsData().SetPerSpecMeanDay(&perSpecMeanDay)
			}(),
			expectError: true,
			errorMsg:    "per_spec_mean_day cannot be negative",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.stats.Validate()

			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestStatsData_PerPages_NullHandling tests PerPages null handling
func TestStatsData_PerPages_NullHandling(t *testing.T) {
	// Arrange - Create StatsData with nil PerPages
	stats := dto.NewStatsData()

	// Act
	data, err := json.Marshal(stats)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	// Assert - PerPages should be omitted when nil
	assert.NotContains(t, result, "per_pages")

	// Test validation with nil PerPages
	err = stats.Validate()
	assert.NoError(t, err)
}

// TestStatsData_PerPagesWithValue tests PerPages with value
func TestStatsData_PerPagesWithValue(t *testing.T) {
	// Arrange
	perPages := 150.0 // Value > 100, which should now be allowed
	stats := dto.NewStatsData().SetPerPages(&perPages)

	// Act
	data, err := json.Marshal(stats)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, 150.0, result["per_pages"])

	// Validation should pass (removed 0-100 constraint)
	err = stats.Validate()
	assert.NoError(t, err)
}

// =============================================================================
// TestValidationScenarios - Comprehensive validation scenario tests
// =============================================================================

// TestValidationScenarios_Comprehensive validates all validation paths
func TestValidationScenarios_Comprehensive(t *testing.T) {
	testCases := []struct {
		name        string
		validator   func() error
		expectError bool
		errorMsg    string
	}{
		{
			name: "nil DashboardResponse",
			validator: func() error {
				var d *dto.DashboardResponse
				return d.Validate()
			},
			expectError: true,
			errorMsg:    "dashboard response is nil",
		},
		{
			name: "nil StatsData",
			validator: func() error {
				var s *dto.StatsData
				return s.Validate()
			},
			expectError: true,
			errorMsg:    "stats data is nil",
		},
		{
			name: "nil EchartConfig",
			validator: func() error {
				var e *dto.EchartConfig
				return e.Validate()
			},
			expectError: true,
			errorMsg:    "echart config is nil",
		},
		{
			name: "negative pages in StatsData",
			validator: func() error {
				return dto.NewStatsData().SetPages(-10).Validate()
			},
			expectError: true,
			errorMsg:    "pages cannot be negative",
		},
		{
			name: "negative total_pages in StatsData",
			validator: func() error {
				return dto.NewStatsData().SetTotalPages(-10).Validate()
			},
			expectError: true,
			errorMsg:    "total_pages cannot be negative",
		},
		{
			name: "nil LogEntry",
			validator: func() error {
				var l *dto.LogEntry
				return l.Validate()
			},
			expectError: true,
			errorMsg:    "log entry is nil",
		},
		{
			name: "nil Project",
			validator: func() error {
				var p *dto.Project
				return p.Validate()
			},
			expectError: true,
			errorMsg:    "project is nil",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.validator()

			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestValidationScenarios_ValidScenarios tests all valid scenarios
func TestValidationScenarios_ValidScenarios(t *testing.T) {
	// All valid configuration
	response := dto.NewDashboardResponse()
	response.SetStats(dto.NewStatsData().
		SetPreviousWeekPages(100).
		SetLastWeekPages(150).
		SetPerPages(testutil.FloatPtr(50.0)).
		SetMeanDay(25.5).
		SetProgressGeral(75.0))

	config := dto.NewEchartConfig().
		SetTitle("Valid Chart").
		AddSeries(dto.Series{
			Name: "Valid Series",
			Type: "line",
			Data: []interface{}{1, 2, 3, 4, 5},
		})
	response.SetEchart(config)

	log := dto.NewLogEntry(1, "2024-01-15T10:30:00Z", 0, 25, nil,
		dto.NewProject(1, "Valid Project", 200, 50))
	response.AddLog(*log)

	err := response.Validate()
	assert.NoError(t, err)
}

// TestValidationScenarios_InvalidScenarios tests all invalid scenarios
func TestValidationScenarios_InvalidScenarios(t *testing.T) {
	testCases := []struct {
		name     string
		modifier func(*dto.DashboardResponse)
		errorMsg string
	}{
		{
			name: "negative previous_week_pages",
			modifier: func(r *dto.DashboardResponse) {
				r.SetStats(dto.NewStatsData().SetPreviousWeekPages(-10))
			},
			errorMsg: "previous_week_pages cannot be negative",
		},
		{
			name: "negative per_pages",
			modifier: func(r *dto.DashboardResponse) {
				r.SetStats(dto.NewStatsData().SetPerPages(testutil.FloatPtr(-10)))
			},
			errorMsg: "per_pages cannot be negative",
		},
		{
			name: "negative mean_day",
			modifier: func(r *dto.DashboardResponse) {
				r.SetStats(dto.NewStatsData().SetMeanDay(-5.0))
			},
			errorMsg: "mean_day cannot be negative",
		},
		{
			name: "negative last_week_pages",
			modifier: func(r *dto.DashboardResponse) {
				r.SetStats(dto.NewStatsData().SetLastWeekPages(-20))
			},
			errorMsg: "last_week_pages cannot be negative",
		},
		{
			name: "negative total_pages",
			modifier: func(r *dto.DashboardResponse) {
				r.SetStats(dto.NewStatsData().SetTotalPages(-100))
			},
			errorMsg: "total_pages cannot be negative",
		},
		{
			name: "negative pages",
			modifier: func(r *dto.DashboardResponse) {
				r.SetStats(dto.NewStatsData().SetPages(-50))
			},
			errorMsg: "pages cannot be negative",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response := dto.NewDashboardResponse()
			tc.modifier(response)

			err := response.Validate()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.errorMsg)
		})
	}
}

// =============================================================================
// Helper functions
// =============================================================================

// formatTimePtr converts a time.Time pointer to a string pointer for JSON serialization
func formatTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(time.RFC3339)
	return &s
}

// =============================================================================
// TestStatsData_RatioFields_NullValidation - Comprehensive null validation tests
// =============================================================================

// TestStatsData_RatioFields_NullValidation tests all null scenarios for ratio fields
func TestStatsData_RatioFields_NullValidation(t *testing.T) {
	testCases := []struct {
		name        string
		stats       *dto.StatsData
		expectError bool
		errorMsg    string
	}{
		{
			name:        "all ratio fields nil - should pass",
			stats:       dto.NewStatsData(),
			expectError: false,
		},
		{
			name: "PerPages nil, PerMeanDay and PerSpecMeanDay with values - should pass",
			stats: func() *dto.StatsData {
				perMeanDay := 1.5
				perSpecMeanDay := 2.0
				return dto.NewStatsData().
					SetPerMeanDay(&perMeanDay).
					SetPerSpecMeanDay(&perSpecMeanDay)
			}(),
			expectError: false,
		},
		{
			name: "PerMeanDay nil, PerPages and PerSpecMeanDay with values - should pass",
			stats: func() *dto.StatsData {
				perPages := 1.25
				perSpecMeanDay := 2.0
				return dto.NewStatsData().
					SetPerPages(&perPages).
					SetPerSpecMeanDay(&perSpecMeanDay)
			}(),
			expectError: false,
		},
		{
			name: "PerSpecMeanDay nil, PerPages and PerMeanDay with values - should pass",
			stats: func() *dto.StatsData {
				perPages := 1.25
				perMeanDay := 1.5
				return dto.NewStatsData().
					SetPerPages(&perPages).
					SetPerMeanDay(&perMeanDay)
			}(),
			expectError: false,
		},
		{
			name: "all ratio fields with valid non-negative values - should pass",
			stats: func() *dto.StatsData {
				perPages := 1.25
				perMeanDay := 1.5
				perSpecMeanDay := 2.0
				return dto.NewStatsData().
					SetPerPages(&perPages).
					SetPerMeanDay(&perMeanDay).
					SetPerSpecMeanDay(&perSpecMeanDay)
			}(),
			expectError: false,
		},
		{
			name: "PerPages with negative value - should fail",
			stats: func() *dto.StatsData {
				perPages := -1.0
				return dto.NewStatsData().SetPerPages(&perPages)
			}(),
			expectError: true,
			errorMsg:    "per_pages cannot be negative",
		},
		{
			name: "PerMeanDay with negative value - should fail",
			stats: func() *dto.StatsData {
				perMeanDay := -1.5
				return dto.NewStatsData().SetPerMeanDay(&perMeanDay)
			}(),
			expectError: true,
			errorMsg:    "per_mean_day cannot be negative",
		},
		{
			name: "PerSpecMeanDay with negative value - should fail",
			stats: func() *dto.StatsData {
				perSpecMeanDay := -2.0
				return dto.NewStatsData().SetPerSpecMeanDay(&perSpecMeanDay)
			}(),
			expectError: true,
			errorMsg:    "per_spec_mean_day cannot be negative",
		},
		{
			name: "multiple ratio fields with negative values - should fail",
			stats: func() *dto.StatsData {
				perPages := -1.0
				perMeanDay := -1.5
				return dto.NewStatsData().
					SetPerPages(&perPages).
					SetPerMeanDay(&perMeanDay)
			}(),
			expectError: true,
			errorMsg:    "per_pages cannot be negative",
		},
		{
			name: "PerPages zero value (nil) with valid others - should pass",
			stats: func() *dto.StatsData {
				perMeanDay := 0.0
				perSpecMeanDay := 0.0
				return dto.NewStatsData().
					SetPerMeanDay(&perMeanDay).
					SetPerSpecMeanDay(&perSpecMeanDay)
			}(),
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.stats.Validate()

			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestStatsData_AllNullRatioFields_Validation specifically tests acceptance criteria
func TestStatsData_AllNullRatioFields_Validation(t *testing.T) {
	// Arrange - Create StatsData with all ratio fields as nil
	stats := dto.NewStatsData()

	// Act - Validate
	err := stats.Validate()

	// Assert - AC1: Validate() accepts null for ratio fields
	assert.NoError(t, err)

	// Act - Marshal to JSON
	data, err := json.Marshal(stats)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	// Assert - AC3: No validation errors for valid null values
	// Ratio fields should be omitted when nil (omitempty)
	assert.NotContains(t, result, "per_pages")
	assert.NotContains(t, result, "per_mean_day")
	assert.NotContains(t, result, "per_spec_mean_day")
}

// TestStatsData_MixedNullAndValue_RatioFields tests mixed null/value scenarios
func TestStatsData_MixedNullAndValue_RatioFields(t *testing.T) {
	testCases := []struct {
		name        string
		stats       *dto.StatsData
		expectError bool
	}{
		{
			name: "null PerPages with valid PerMeanDay - should pass",
			stats: func() *dto.StatsData {
				perMeanDay := 1.5
				return dto.NewStatsData().SetPerMeanDay(&perMeanDay)
			}(),
			expectError: false,
		},
		{
			name: "null PerMeanDay with valid PerPages - should pass",
			stats: func() *dto.StatsData {
				perPages := 1.25
				return dto.NewStatsData().SetPerPages(&perPages)
			}(),
			expectError: false,
		},
		{
			name: "null PerSpecMeanDay with valid PerPages and PerMeanDay - should pass",
			stats: func() *dto.StatsData {
				perPages := 1.25
				perMeanDay := 1.5
				return dto.NewStatsData().
					SetPerPages(&perPages).
					SetPerMeanDay(&perMeanDay)
			}(),
			expectError: false,
		},
		{
			name: "PerPages and PerMeanDay nil, PerSpecMeanDay with value - should pass",
			stats: func() *dto.StatsData {
				perSpecMeanDay := 2.0
				return dto.NewStatsData().SetPerSpecMeanDay(&perSpecMeanDay)
			}(),
			expectError: false,
		},
		{
			name: "PerPages nil, PerMeanDay nil, PerSpecMeanDay negative - should fail",
			stats: func() *dto.StatsData {
				perSpecMeanDay := -2.0
				return dto.NewStatsData().SetPerSpecMeanDay(&perSpecMeanDay)
			}(),
			expectError: true,
		},
		{
			name: "PerPages negative, PerMeanDay nil, PerSpecMeanDay nil - should fail",
			stats: func() *dto.StatsData {
				perPages := -1.0
				return dto.NewStatsData().SetPerPages(&perPages)
			}(),
			expectError: true,
		},
		{
			name: "PerPages nil, PerMeanDay negative, PerSpecMeanDay nil - should fail",
			stats: func() *dto.StatsData {
				perMeanDay := -1.5
				return dto.NewStatsData().SetPerMeanDay(&perMeanDay)
			}(),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.stats.Validate()

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestStatsData_RatioFields_JSONSerialization tests JSON serialization of null ratio fields
func TestStatsData_RatioFields_JSONSerialization(t *testing.T) {
	// Arrange - Create StatsData with nil ratio fields
	stats := dto.NewStatsData().
		SetPreviousWeekPages(100).
		SetLastWeekPages(150)

	// Act
	data, err := json.Marshal(stats)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	// Assert - Nil ratio fields should be omitted
	assert.NotContains(t, result, "per_pages")
	assert.NotContains(t, result, "per_mean_day")
	assert.NotContains(t, result, "per_spec_mean_day")

	// Arrange - Create StatsData with non-nil ratio fields
	perPages := 1.5
	perMeanDay := 1.25
	perSpecMeanDay := 2.0
	statsWithValues := dto.NewStatsData().
		SetPreviousWeekPages(100).
		SetLastWeekPages(150).
		SetPerPages(&perPages).
		SetPerMeanDay(&perMeanDay).
		SetPerSpecMeanDay(&perSpecMeanDay)

	// Act
	data, err = json.Marshal(statsWithValues)
	require.NoError(t, err)

	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	// Assert - Non-nil ratio fields should be present
	assert.Equal(t, 1.5, result["per_pages"])
	assert.Equal(t, 1.25, result["per_mean_day"])
	assert.Equal(t, 2.0, result["per_spec_mean_day"])
}
