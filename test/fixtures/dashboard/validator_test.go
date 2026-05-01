package dashboard

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestFixtureValidator_WeekdayCoverage tests the 7 weekday coverage requirement
func TestFixtureValidator_WeekdayCoverage(t *testing.T) {
	baseDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	// Create logs covering all 7 weekdays with proper date spread
	// Need at least 30 days of data and 7 unique weekdays
	logs := make([]*LogFixture, 0)

	// Add multiple entries per weekday to ensure we have 30+ unique days
	for day := 0; day < 7; day++ {
		// Add 5 entries for each weekday, spread across different weeks
		for week := 0; week < 5; week++ {
			logs = append(logs, &LogFixture{
				ID:        int64(day*5 + week + 1),
				ProjectID: 1,
				Data:      baseDate.AddDate(0, 0, -(day + week*7)),
				StartPage: 0,
				EndPage:   10,
				WDay:      day,
			})
		}
	}

	validator := NewFixtureValidator(logs, nil)
	result := validator.Validate()

	assert.True(t, result.Valid, "Validation should pass with all 7 weekdays and 30+ days")
}

// TestFixtureValidator_WeekdayCoverage_Missing tests missing weekday detection
func TestFixtureValidator_WeekdayCoverage_Missing(t *testing.T) {
	baseDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	// Create logs covering only 5 weekdays (missing Sunday=0 and Saturday=6)
	// Spread across 30+ days to pass data range check
	logs := make([]*LogFixture, 0)
	for day := 1; day < 6; day++ {
		for i := 0; i < 7; i++ { // Add multiple entries per weekday to spread dates
			logs = append(logs, &LogFixture{
				ID:        int64(day*10 + i),
				ProjectID: 1,
				Data:      baseDate.AddDate(0, 0, -(day*7 + i)),
				StartPage: 0,
				EndPage:   10,
				WDay:      day,
			})
		}
	}

	validator := NewFixtureValidator(logs, nil)
	result := validator.Validate()

	assert.False(t, result.Valid, "Validation should fail with missing weekdays")
	// Check that at least one error mentions missing weekdays
	hasWeekdayError := false
	for _, err := range result.Errors {
		if containsString(err.Error(), "missing weekdays") {
			hasWeekdayError = true
			break
		}
	}
	assert.True(t, hasWeekdayError, "Should have error mentioning missing weekdays")
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestFixtureValidator_WeekdayCoverage_NoLogs tests empty logs handling
func TestFixtureValidator_WeekdayCoverage_NoLogs(t *testing.T) {
	validator := NewFixtureValidator(nil, nil)
	result := validator.Validate()

	assert.False(t, result.Valid, "Validation should fail with no logs")
	assert.Contains(t, result.Errors[0].Error(), "no logs found", "Error should mention no logs")
}

// TestFixtureValidator_DataRange tests the 30-day minimum requirement
func TestFixtureValidator_DataRange(t *testing.T) {
	baseDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	// Create logs for exactly 30 days
	logs := make([]*LogFixture, 30)
	for i := 0; i < 30; i++ {
		logs[i] = &LogFixture{
			ID:        int64(i + 1),
			ProjectID: 1,
			Data:      baseDate.AddDate(0, 0, -i),
			StartPage: i * 5,
			EndPage:   (i + 1) * 5,
			WDay:      int(baseDate.AddDate(0, 0, -i).Weekday()),
		}
	}

	validator := NewFixtureValidator(logs, nil)
	result := validator.Validate()

	assert.True(t, result.Valid, "Validation should pass with exactly 30 days")
}

// TestFixtureValidator_DataRange_Insufficient tests detection of insufficient data
func TestFixtureValidator_DataRange_Insufficient(t *testing.T) {
	baseDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	// Create logs for only 15 days (won't pass date spread check either)
	logs := make([]*LogFixture, 15)
	for i := 0; i < 15; i++ {
		logs[i] = &LogFixture{
			ID:        int64(i + 1),
			ProjectID: 1,
			Data:      baseDate.AddDate(0, 0, -i),
			StartPage: i * 5,
			EndPage:   (i + 1) * 5,
			WDay:      int(baseDate.AddDate(0, 0, -i).Weekday()),
		}
	}

	validator := NewFixtureValidator(logs, nil)
	result := validator.Validate()

	assert.False(t, result.Valid, "Validation should fail with insufficient data")
	// Check that at least one error mentions insufficient data
	hasInsufficientError := false
	for _, err := range result.Errors {
		if containsString(err.Error(), "insufficient data") {
			hasInsufficientError = true
			break
		}
	}
	assert.True(t, hasInsufficientError, "Should have error mentioning insufficient data")
}

// TestFixtureValidator_DataRange_DuplicateDates tests that duplicate dates count as one day
func TestFixtureValidator_DataRange_DuplicateDates(t *testing.T) {
	baseDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	// Create multiple logs on the same date (should count as 1 day)
	logs := make([]*LogFixture, 30)
	for i := 0; i < 30; i++ {
		logs[i] = &LogFixture{
			ID:        int64(i + 1),
			ProjectID: 1,
			Data:      baseDate, // All on same date!
			StartPage: i * 5,
			EndPage:   (i + 1) * 5,
			WDay:      int(baseDate.Weekday()),
		}
	}

	validator := NewFixtureValidator(logs, nil)
	result := validator.Validate()

	assert.False(t, result.Valid, "Validation should fail with duplicate dates")
	// Check that at least one error mentions insufficient data or date spread
	hasError := false
	for _, err := range result.Errors {
		if containsString(err.Error(), "insufficient data") ||
			containsString(err.Error(), "date spread") {
			hasError = true
			break
		}
	}
	assert.True(t, hasError, "Should have error mentioning insufficient data or date spread")
}

// TestFixtureValidator_Combined tests multiple validation failures
func TestFixtureValidator_Combined(t *testing.T) {
	baseDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	// Create logs with multiple issues:
	// - Only 5 days of data (need 30)
	// - Missing several weekdays

	logs := make([]*LogFixture, 0)
	for day := 1; day < 6; day++ {
		logs = append(logs, &LogFixture{
			ID:        int64(day),
			ProjectID: 1,
			Data:      baseDate.AddDate(0, 0, -day),
			StartPage: 0,
			EndPage:   10,
			WDay:      day,
		})
	}

	validator := NewFixtureValidator(logs, nil)
	result := validator.Validate()

	assert.False(t, result.Valid, "Validation should fail with multiple issues")
	// Check that we have errors for both insufficient data and missing weekdays
	hasInsufficientError := false
	hasWeekdayError := false
	for _, err := range result.Errors {
		if containsString(err.Error(), "insufficient data") {
			hasInsufficientError = true
		}
		if containsString(err.Error(), "missing weekdays") {
			hasWeekdayError = true
		}
	}
	assert.True(t, hasInsufficientError, "Should have error for insufficient data")
	assert.True(t, hasWeekdayError, "Should have error for missing weekdays")
}

// TestFixtureValidator_ProjectConsistency tests project-log consistency validation
func TestFixtureValidator_ProjectConsistency(t *testing.T) {
	baseDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	// Create logs with invalid project references
	logs := []*LogFixture{
		{ID: 1, ProjectID: 999, Data: baseDate, StartPage: 0, EndPage: 10}, // Invalid project
	}

	projects := []*ProjectFixture{
		{ID: 1, Name: "Valid Project", TotalPage: 100, Page: 50},
	}

	validator := NewFixtureValidator(logs, projects)
	result := validator.Validate()

	assert.False(t, result.Valid, "Validation should fail with invalid project references")
	// Check that at least one error mentions non-existent projects
	hasProjectError := false
	for _, err := range result.Errors {
		if containsString(err.Error(), "non-existent projects") {
			hasProjectError = true
			break
		}
	}
	assert.True(t, hasProjectError, "Should have error mentioning non-existent projects")
}

// TestFixtureValidator_DateRange tests date range validation
func TestFixtureValidator_DateRange(t *testing.T) {
	baseDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	// Create logs with adequate date spread
	logs := make([]*LogFixture, 30)
	for i := 0; i < 30; i++ {
		logs[i] = &LogFixture{
			ID:        int64(i + 1),
			ProjectID: 1,
			Data:      baseDate.AddDate(0, 0, -i),
			StartPage: i * 5,
			EndPage:   (i + 1) * 5,
			WDay:      int(baseDate.AddDate(0, 0, -i).Weekday()),
		}
	}

	validator := NewFixtureValidator(logs, nil)
	result := validator.Validate()

	assert.True(t, result.Valid, "Validation should pass with adequate date spread")
}

// TestFixtureValidator_DateRange_Narrow tests narrow date range detection
func TestFixtureValidator_DateRange_Narrow(t *testing.T) {
	baseDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	// Create logs spread over only 10 days (need 29+)
	logs := make([]*LogFixture, 30)
	for i := 0; i < 30; i++ {
		logs[i] = &LogFixture{
			ID:        int64(i + 1),
			ProjectID: 1,
			Data:      baseDate.AddDate(0, 0, -(i % 10)), // Only 10 unique days
			StartPage: i * 5,
			EndPage:   (i + 1) * 5,
			WDay:      int(baseDate.AddDate(0, 0, -(i % 10)).Weekday()),
		}
	}

	validator := NewFixtureValidator(logs, nil)
	result := validator.Validate()

	assert.False(t, result.Valid, "Validation should fail with narrow date range")
	// Check that at least one error mentions date spread
	hasSpreadError := false
	for _, err := range result.Errors {
		if containsString(err.Error(), "date spread") {
			hasSpreadError = true
			break
		}
	}
	assert.True(t, hasSpreadError, "Should have error mentioning date spread")
}

// TestValidateScenario tests the ValidateScenario convenience function
func TestValidateScenario(t *testing.T) {
	scenario := ScenarioMeanProgress()
	result := ValidateScenario(scenario)

	assert.True(t, result.Valid, "ScenarioMeanProgress should validate successfully")
}

// TestMustValidateScenario tests the MustValidateScenario panic behavior
func TestMustValidateScenario(t *testing.T) {
	// This should not panic - valid scenario
	assert.NotPanics(t, func() {
		MustValidateScenario(ScenarioMeanProgress())
	})
}

// TestMustValidateScenario_Panic tests that invalid scenarios cause panic
func TestMustValidateScenario_Panic(t *testing.T) {
	// Create an invalid scenario (no logs)
	invalidScenario := &Scenario{
		Name:        "Invalid",
		Description: "No logs",
		Projects:    []*ProjectFixture{},
		Logs:        []*LogFixture{},
	}

	assert.Panics(t, func() {
		MustValidateScenario(invalidScenario)
	})
}

// TestFixtureValidator_Warnings tests warning generation for narrow date range
func TestFixtureValidator_Warnings(t *testing.T) {
	baseDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	// Create logs with adequate data but narrow date spread (should trigger error)
	logs := make([]*LogFixture, 30)
	for i := 0; i < 30; i++ {
		logs[i] = &LogFixture{
			ID:        int64(i + 1),
			ProjectID: 1,
			Data:      baseDate.AddDate(0, 0, -(i % 5)), // Only 5 unique days
			StartPage: i * 5,
			EndPage:   (i + 1) * 5,
			WDay:      int(baseDate.AddDate(0, 0, -(i % 5)).Weekday()),
		}
	}

	validator := NewFixtureValidator(logs, nil)
	result := validator.Validate()

	// Narrow date range is now an error
	assert.False(t, result.Valid, "Validation should fail with narrow date range")

	// Check that we have the expected error
	hasSpreadError := false
	for _, err := range result.Errors {
		if containsString(err.Error(), "date spread") {
			hasSpreadError = true
			break
		}
	}
	assert.True(t, hasSpreadError, "Should have error mentioning date spread")
}
