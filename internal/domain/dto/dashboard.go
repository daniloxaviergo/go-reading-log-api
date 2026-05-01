// Package dto provides data transfer objects and shared utilities for the dashboard package.
// This file contains the date abstraction layer for deterministic testing.
package dto

import (
	"time"
)

// GetTodayFunc is a variable for dependency injection that allows test-specific date injection.
// Defaults to returning the actual current date truncated to midnight.
// This enables deterministic testing while maintaining production behavior.
//
// ⚠️ WARNING: This global variable is NOT goroutine-safe.
// For parallel tests, each test should use its own isolated context or
// ensure SetTestDate() calls are properly synchronized.
var GetTodayFunc = func() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

// GetToday returns the current date truncated to midnight for consistent date references.
// This ensures all calculations use the same reference point regardless of when they're
// executed within a single day. Uses GetTodayFunc for testability.
func GetToday() time.Time {
	return GetTodayFunc()
}

// SetTestDate sets GetTodayFunc to return a fixed date for deterministic testing.
// Usage: defer SetTestDate(time.Now()) before calling, or call directly in tests.
//
// ⚠️ WARNING: This function modifies a global variable and is NOT goroutine-safe.
// Do not use this function in parallel tests without proper synchronization.
func SetTestDate(date time.Time) {
	GetTodayFunc = func() time.Time {
		now := date
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	}
}

// GetDateRangeLast15Days returns the date range for the last 15 days.
// The end date is today (midnight) and start date is 14 days ago (to include today in the count).
// This ensures exactly 15 days of data coverage including today.
func GetDateRangeLast15Days() (start, end time.Time) {
	end = GetToday()
	start = end.AddDate(0, 0, -14)
	return start, end
}

// GetDateRangeLast30Days returns the date range for the last 30 days.
// The end date is today (midnight) and start date is 30 days ago.
func GetDateRangeLast30Days() (start, end time.Time) {
	end = GetToday()
	start = end.AddDate(0, 0, -30)
	return start, end
}

// GetDateRangeLast6Months returns the date range for the last 6 months.
// The end date is today (midnight) and start date is 6 months ago.
// Uses AddDate(0, -6, 0) for exact month arithmetic that handles varying month lengths.
func GetDateRangeLast6Months() (start, end time.Time) {
	end = GetToday()
	start = end.AddDate(0, -6, 0)
	return start, end
}
