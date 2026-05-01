package integration

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

// T is a placeholder for testing.T - should be passed in from test functions
type T interface {
	Errorf(format string, args ...interface{})
	FailNow()
	Skip(args ...interface{})
	Run(name string, f func(t *testing.T)) bool
}

// ValidationResult holds comparison results
type ValidationResult struct {
	Passed   bool
	Errors   []string
	Warnings []string
}

// Validator provides response validation utilities
type Validator struct {
	tolerance float64 // Floating point comparison tolerance
}

// NewValidator creates a new validator with specified tolerance
func NewValidator(tolerance float64) *Validator {
	return &Validator{tolerance: tolerance}
}

// ValidateDashboardResponse compares Go and Rails responses
func (v *Validator) ValidateDashboardResponse(
	t T,
	goResponse interface{},
	railsResponse interface{},
	endpoint string,
) ValidationResult {
	result := ValidationResult{Passed: true, Errors: []string{}}

	// Convert to map for comparison
	goMap := v.interfaceToMap(goResponse)
	railsMap := v.interfaceToMap(railsResponse)

	// Compare common fields
	commonFields := []string{"status", "data", "meta"}
	for _, field := range commonFields {
		goVal, goOk := goMap[field]
		railsVal, railsOk := railsMap[field]

		if goOk != railsOk {
			result.Errors = append(result.Errors,
				fmt.Sprintf("%s: field '%s' present in Go but not Rails (or vice versa)",
					endpoint, field))
			result.Passed = false
			continue
		}

		if !v.valuesEqual(goVal, railsVal) {
			result.Errors = append(result.Errors,
				fmt.Sprintf("%s: field '%s' differs - Go: %v, Rails: %v",
					endpoint, field, goVal, railsVal))
			result.Passed = false
		}
	}

	// Special handling for calculated fields
	v.validateCalculatedFields(t, goMap, railsMap, endpoint, &result)

	return result
}

// validateCalculatedFields handles special comparison for calculated values
func (v *Validator) validateCalculatedFields(
	t T,
	goMap, railsMap map[string]interface{},
	endpoint string,
	result *ValidationResult,
) {
	// Handle float comparisons with tolerance
	floatFields := []string{"progress", "per_pages", "mean_day", "spec_mean_day"}

	for _, field := range floatFields {
		goVal, goOk := goMap[field]
		railsVal, railsOk := railsMap[field]

		if !goOk || !railsOk {
			continue
		}

		goFloat := v.toFloat64(goVal)
		railsFloat := v.toFloat64(railsVal)

		if math.Abs(goFloat-railsFloat) > v.tolerance {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("%s: %s differs slightly - Go: %.6f, Rails: %.6f (tolerance: %.6f)",
					endpoint, field, goFloat, railsFloat, v.tolerance))
		}
	}
}

func (v *Validator) valuesEqual(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

func (v *Validator) toFloat64(val interface{}) float64 {
	switch v := val.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	default:
		return 0
	}
}

func (v *Validator) interfaceToMap(val interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	switch v := val.(type) {
	case map[string]interface{}:
		return v
		// Add more type conversions as needed
	}

	return result
}
