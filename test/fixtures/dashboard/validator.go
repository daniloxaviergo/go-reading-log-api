package dashboard

import (
	"fmt"
	"time"
)

// FixtureValidator validates dashboard test fixtures for completeness
type FixtureValidator struct {
	logs           []*LogFixture
	projects       []*ProjectFixture
	validationDate time.Time // For deterministic testing
}

// ValidationOption configures the validator behavior
type ValidationOption func(*FixtureValidator)

// WithValidationDate sets a fixed date for deterministic validation
func WithValidationDate(date time.Time) ValidationOption {
	return func(v *FixtureValidator) {
		v.validationDate = date
	}
}

// NewFixtureValidator creates a new validator instance
func NewFixtureValidator(logs []*LogFixture, projects []*ProjectFixture, options ...ValidationOption) *FixtureValidator {
	validator := &FixtureValidator{
		logs:           logs,
		projects:       projects,
		validationDate: time.Now(),
	}

	for _, opt := range options {
		opt(validator)
	}

	return validator
}

// ValidationResult contains the outcome of a validation
type ValidationResult struct {
	Valid    bool
	Errors   []error
	Warnings []error
}

// IsComplete returns true if validation passed without errors
func (r *ValidationResult) IsComplete() bool {
	return len(r.Errors) == 0
}

// Error returns a formatted error message for all validation failures
func (r *ValidationResult) Error() string {
	if r.Valid {
		return "validation passed"
	}

	var msg string
	for _, err := range r.Errors {
		if msg != "" {
			msg += "; "
		}
		msg += err.Error()
	}
	return msg
}

// Validate runs all validation checks and returns the results
func (v *FixtureValidator) Validate() *ValidationResult {
	result := &ValidationResult{
		Errors:   make([]error, 0),
		Warnings: make([]error, 0),
	}

	// Run all validation checks
	v.validateWeekdayCoverage(result)
	v.validateDataRange(result)
	v.validateProjectLogConsistency(result)
	v.validateDateRange(result)

	result.Valid = len(result.Errors) == 0
	return result
}

// validateWeekdayCoverage ensures all 7 weekdays are covered
func (v *FixtureValidator) validateWeekdayCoverage(result *ValidationResult) {
	if len(v.logs) == 0 {
		result.Errors = append(result.Errors, fmt.Errorf("no logs found - cannot validate weekday coverage"))
		return
	}

	// Collect unique weekdays from logs
	weekdays := make(map[int]bool)
	for _, log := range v.logs {
		// Use the stored WDay field if available, otherwise calculate from Data
		if log.WDay >= 0 && log.WDay <= 6 {
			weekdays[log.WDay] = true
		} else if !log.Data.IsZero() {
			weekdays[int(log.Data.Weekday())] = true
		}
	}

	// Check for complete coverage (0-6 = Sunday-Saturday)
	requiredWeekdays := []int{0, 1, 2, 3, 4, 5, 6}
	missingWeekdays := make([]int, 0)

	for _, wd := range requiredWeekdays {
		if !weekdays[wd] {
			missingWeekdays = append(missingWeekdays, wd)
		}
	}

	if len(missingWeekdays) > 0 {
		result.Errors = append(result.Errors, fmt.Errorf(
			"weekday coverage incomplete: missing weekdays %v (got %d unique days, expected 7)",
			missingWeekdays, len(weekdays),
		))
	}
}

// validateDataRange ensures minimum 30 days of data
func (v *FixtureValidator) validateDataRange(result *ValidationResult) {
	if len(v.logs) == 0 {
		result.Errors = append(result.Errors, fmt.Errorf(
			"insufficient data: no logs found (expected at least 30 log entries)",
		))
		return
	}

	// Collect unique dates from logs
	dates := make(map[string]bool)
	for _, log := range v.logs {
		if !log.Data.IsZero() {
			dateStr := log.Data.Format("2006-01-02")
			dates[dateStr] = true
		}
	}

	minimumDays := 30
	if len(dates) < minimumDays {
		result.Errors = append(result.Errors, fmt.Errorf(
			"insufficient data: only %d unique days found (expected at least %d days)",
			len(dates), minimumDays,
		))
	}
}

// validateProjectLogConsistency ensures all logs have valid project associations
func (v *FixtureValidator) validateProjectLogConsistency(result *ValidationResult) {
	// If no projects provided, skip this validation (for scenarios without projects)
	if len(v.projects) == 0 {
		return
	}

	// Build set of valid project IDs
	validProjectIDs := make(map[int64]bool)
	for _, proj := range v.projects {
		validProjectIDs[proj.ID] = true
	}

	// Check each log has a valid project
	invalidLogs := make([]int64, 0)
	for _, log := range v.logs {
		if !validProjectIDs[log.ProjectID] {
			invalidLogs = append(invalidLogs, log.ID)
		}
	}

	if len(invalidLogs) > 0 {
		result.Errors = append(result.Errors, fmt.Errorf(
			"project consistency error: logs %v reference non-existent projects",
			invalidLogs,
		))
	}
}

// validateDateRange checks for logical date constraints
func (v *FixtureValidator) validateDateRange(result *ValidationResult) {
	if len(v.logs) == 0 {
		return
	}

	// Find min and max dates
	var minDate, maxDate time.Time
	found := false

	for _, log := range v.logs {
		if !log.Data.IsZero() {
			if !found || log.Data.Before(minDate) {
				minDate = log.Data
			}
			if !found || log.Data.After(maxDate) {
				maxDate = log.Data
			}
			found = true
		}
	}

	if !found {
		return
	}

	// Check that we have reasonable date spread
	const minDateSpread = 29 * 24 * time.Hour // At least 29 days of spread
	if maxDate.Sub(minDate) < minDateSpread {
		result.Errors = append(result.Errors, fmt.Errorf(
			"date spread too narrow: only %v covered (expected at least %v)",
			maxDate.Sub(minDate).Round(time.Hour),
			minDateSpread.Round(time.Hour),
		))
	}
}

// ValidateScenario validates a complete scenario
func ValidateScenario(scenario *Scenario, options ...ValidationOption) *ValidationResult {
	validator := NewFixtureValidator(scenario.Logs, scenario.Projects, options...)
	return validator.Validate()
}

// MustValidateScenario validates a scenario and panics if validation fails
func MustValidateScenario(scenario *Scenario, options ...ValidationOption) {
	result := ValidateScenario(scenario, options...)
	if !result.Valid {
		panic(fmt.Sprintf("scenario validation failed: %s", result.Error()))
	}
}
