package validation

import "testing"

// =====================
// Error Tests
// =====================

func TestValidationError_Error(t *testing.T) {
	err := NewValidationError("invalid_page", "page", "page must be positive")
	expected := "page: page must be positive"

	if err.Error() != expected {
		t.Errorf("expected error message '%s', got '%s'", expected, err.Error())
	}
}

func TestValidationErrorList_Error_SingleError(t *testing.T) {
	errors := &ValidationErrorList{}
	errors.AddError(NewValidationError("invalid", "field", "is invalid"))

	expected := "field: is invalid"
	if errors.Error() != expected {
		t.Errorf("expected error message '%s', got '%s'", expected, errors.Error())
	}
}

func TestValidationErrorList_Error_MultipleErrors(t *testing.T) {
	errors := &ValidationErrorList{}
	errors.AddError(NewValidationError("invalid1", "field1", "is invalid 1"))
	errors.AddError(NewValidationError("invalid2", "field2", "is invalid 2"))

	expected := "validation error: 2 errors occurred"
	if errors.Error() != expected {
		t.Errorf("expected error message '%s', got '%s'", expected, errors.Error())
	}
}

func TestValidationErrorList_Error_NoErrors(t *testing.T) {
	errors := &ValidationErrorList{}

	expected := "validation error: no errors"
	if errors.Error() != expected {
		t.Errorf("expected error message '%s', got '%s'", expected, errors.Error())
	}
}

func TestValidationErrorList_HasErrors(t *testing.T) {
	tests := []struct {
		name     string
		errors   *ValidationErrorList
		expected bool
	}{
		{
			name:     "no errors",
			errors:   &ValidationErrorList{},
			expected: false,
		},
		{
			name: "with errors",
			errors: func() *ValidationErrorList {
				e := &ValidationErrorList{}
				e.AddError(NewValidationError("test", "field", "message"))
				return e
			}(),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.errors.HasErrors()
			if result != tt.expected {
				t.Errorf("expected HasErrors()=%v, got %v", tt.expected, result)
			}
		})
	}
}

func TestValidationErrorList_ToMap(t *testing.T) {
	tests := []struct {
		name     string
		errors   *ValidationErrorList
		expected map[string]interface{}
	}{
		{
			name:     "no errors",
			errors:   &ValidationErrorList{},
			expected: nil,
		},
		{
			name: "single error",
			errors: func() *ValidationErrorList {
				e := &ValidationErrorList{}
				e.AddError(NewValidationError("invalid", "page", "must be positive"))
				return e
			}(),
			expected: map[string]interface{}{"page": "must be positive"},
		},
		{
			name: "multiple errors",
			errors: func() *ValidationErrorList {
				e := &ValidationErrorList{}
				e.AddError(NewValidationError("invalid1", "page", "must be positive"))
				e.AddError(NewValidationError("invalid2", "status", "must be valid"))
				return e
			}(),
			expected: map[string]interface{}{
				"page":   "must be positive",
				"status": "must be valid",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.errors.ToMap()
			if tt.expected == nil {
				if result != nil {
					t.Errorf("expected nil, got %v", result)
				}
			} else {
				if len(result) != len(tt.expected) {
					t.Errorf("expected %d errors, got %d", len(tt.expected), len(result))
				}
				for k, v := range tt.expected {
					if result[k] != v {
						t.Errorf("expected %s='%v', got '%v'", k, v, result[k])
					}
				}
			}
		})
	}
}

// =====================
// Status Validation Tests
// =====================

func TestValidStatusValues(t *testing.T) {
	expected := []string{
		StatusUnstarted,
		StatusFinished,
		StatusRunning,
		StatusSleeping,
		StatusStopped,
	}

	values := ValidStatusValues()
	if len(values) != len(expected) {
		t.Errorf("expected %d valid status values, got %d", len(expected), len(values))
	}

	for i, expectedStatus := range expected {
		if values[i] != expectedStatus {
			t.Errorf("expected status[%d]='%s', got '%s'", i, expectedStatus, values[i])
		}
	}
}

func TestValidateStatus_Valid(t *testing.T) {
	validStatuses := []string{
		StatusUnstarted,
		StatusFinished,
		StatusRunning,
		StatusSleeping,
		StatusStopped,
	}

	for _, status := range validStatuses {
		t.Run(status, func(t *testing.T) {
			err := ValidateStatus(status)
			if err != nil {
				t.Errorf("expected nil error for valid status '%s', got %v", status, err)
			}
		})
	}
}

func TestValidateStatus_Invalid(t *testing.T) {
	invalidStatuses := []string{
		"invalid",
		"",
		"unknown",
		"running ", // trailing space
		" Running", // leading space
	}

	for _, status := range invalidStatuses {
		t.Run(status, func(t *testing.T) {
			err := ValidateStatus(status)
			if err == nil {
				t.Errorf("expected error for invalid status '%s', got nil", status)
			}
			if err != nil && err.Code != "invalid_status" {
				t.Errorf("expected error code 'invalid_status', got '%s'", err.Code)
			}
		})
	}
}

// =====================
// Page Validation Tests
// =====================

func TestValidatePage_Valid(t *testing.T) {
	tests := []struct {
		name      string
		page      int
		totalPage int
	}{
		{"page equals total_page", 100, 100},
		{"page less than total_page", 50, 100},
		{"page is zero", 0, 100},
		{"total_page is zero (edge case)", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePage(tt.page, tt.totalPage)
			if err != nil {
				t.Errorf("expected nil error for page=%d, totalPage=%d, got %v", tt.page, tt.totalPage, err)
			}
		})
	}
}

func TestValidatePage_Negative(t *testing.T) {
	err := ValidatePage(-1, 100)
	if err == nil {
		t.Errorf("expected error for negative page, got nil")
	}
	if err.Code != "page_invalid" {
		t.Errorf("expected error code 'page_invalid', got '%s'", err.Code)
	}
}

func TestValidatePage_ExceedsTotal(t *testing.T) {
	err := ValidatePage(101, 100)
	if err == nil {
		t.Errorf("expected error for page > total_page, got nil")
	}
	if err.Code != "page_exceeds_total" {
		t.Errorf("expected error code 'page_exceeds_total', got '%s'", err.Code)
	}
}

func TestValidateTotalPage_Valid(t *testing.T) {
	err := ValidateTotalPage(100)
	if err != nil {
		t.Errorf("expected nil error for valid total_page, got %v", err)
	}
}

func TestValidateTotalPage_Invalid(t *testing.T) {
	tests := []struct {
		name      string
		totalPage int
	}{
		{"zero", 0},
		{"negative", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTotalPage(tt.totalPage)
			if err == nil {
				t.Errorf("expected error for total_page=%d, got nil", tt.totalPage)
			}
			if err.Code != "total_page_invalid" {
				t.Errorf("expected error code 'total_page_invalid', got '%s'", err.Code)
			}
		})
	}
}

// =====================
// Log Validation Tests
// =====================

func TestValidateStartEndPage_Valid(t *testing.T) {
	tests := []struct {
		name      string
		startPage int
		endPage   int
	}{
		{"start equals end", 10, 10},
		{"start less than end", 10, 20},
		{"both zero", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStartEndPage(tt.startPage, tt.endPage)
			if err != nil {
				t.Errorf("expected nil error for startPage=%d, endPage=%d, got %v", tt.startPage, tt.endPage, err)
			}
		})
	}
}

func TestValidateStartEndPage_StartNegative(t *testing.T) {
	err := ValidateStartEndPage(-1, 10)
	if err == nil {
		t.Errorf("expected error for negative start_page, got nil")
	}
	if err.Code != "start_page_invalid" {
		t.Errorf("expected error code 'start_page_invalid', got '%s'", err.Code)
	}
}

func TestValidateStartEndPage_EndNegative(t *testing.T) {
	err := ValidateStartEndPage(10, -1)
	if err == nil {
		t.Errorf("expected error for negative end_page, got nil")
	}
	if err.Code != "end_page_invalid" {
		t.Errorf("expected error code 'end_page_invalid', got '%s'", err.Code)
	}
}

func TestValidateStartEndPage_StartExceedsEnd(t *testing.T) {
	err := ValidateStartEndPage(20, 10)
	if err == nil {
		t.Errorf("expected error for start_page > end_page, got nil")
	}
	if err.Code != "start_page_exceeds_end_page" {
		t.Errorf("expected error code 'start_page_exceeds_end_page', got '%s'", err.Code)
	}
}

// =====================
// Project Validation Tests
// =====================

func TestValidateProject_Valid(t *testing.T) {
	errors := ValidateProject(50, 100, StatusRunning)
	if errors.HasErrors() {
		t.Errorf("expected no errors, got %d", len(errors.Errors))
	}
}

func TestValidateProject_PageExceedsTotal(t *testing.T) {
	errors := ValidateProject(150, 100, StatusRunning)
	if !errors.HasErrors() {
		t.Errorf("expected validation errors, got none")
	}
	if len(errors.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(errors.Errors))
	}
	if errors.Errors[0].Code != "page_exceeds_total" {
		t.Errorf("expected error code 'page_exceeds_total', got '%s'", errors.Errors[0].Code)
	}
}

func TestValidateProject_InvalidStatus(t *testing.T) {
	errors := ValidateProject(50, 100, "invalid_status")
	if !errors.HasErrors() {
		t.Errorf("expected validation errors, got none")
	}
	if len(errors.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(errors.Errors))
	}
	if errors.Errors[0].Code != "invalid_status" {
		t.Errorf("expected error code 'invalid_status', got '%s'", errors.Errors[0].Code)
	}
}

func TestValidateProject_MultipleErrors(t *testing.T) {
	errors := ValidateProject(150, 100, "invalid_status")
	if !errors.HasErrors() {
		t.Errorf("expected validation errors, got none")
	}
	// Should have 2 errors: page_exceeds_total and invalid_status
	if len(errors.Errors) != 2 {
		t.Errorf("expected 2 errors, got %d: %v", len(errors.Errors), errors.Errors)
	}
}

func TestValidateProject_ZeroTotalPage(t *testing.T) {
	errors := ValidateProject(0, 0, StatusRunning)
	if !errors.HasErrors() {
		t.Errorf("expected validation errors, got none")
	}
	if len(errors.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(errors.Errors))
	}
	if errors.Errors[0].Code != "total_page_invalid" {
		t.Errorf("expected error code 'total_page_invalid', got '%s'", errors.Errors[0].Code)
	}
}

// =====================
// Log Validation Tests
// =====================

func TestValidateLog_Valid(t *testing.T) {
	errors := ValidateLog(10, 20)
	if errors.HasErrors() {
		t.Errorf("expected no errors, got %d", len(errors.Errors))
	}
}

func TestValidateLog_StartExceedsEnd(t *testing.T) {
	errors := ValidateLog(20, 10)
	if !errors.HasErrors() {
		t.Errorf("expected validation errors, got none")
	}
	if len(errors.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(errors.Errors))
	}
	if errors.Errors[0].Code != "start_page_exceeds_end_page" {
		t.Errorf("expected error code 'start_page_exceeds_end_page', got '%s'", errors.Errors[0].Code)
	}
}

func TestValidateLog_NegativeStartPage(t *testing.T) {
	errors := ValidateLog(-1, 10)
	if !errors.HasErrors() {
		t.Errorf("expected validation errors, got none")
	}
	if len(errors.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(errors.Errors))
	}
	if errors.Errors[0].Code != "start_page_invalid" {
		t.Errorf("expected error code 'start_page_invalid', got '%s'", errors.Errors[0].Code)
	}
}

func TestValidateLog_NegativeEndPage(t *testing.T) {
	errors := ValidateLog(10, -1)
	if !errors.HasErrors() {
		t.Errorf("expected validation errors, got none")
	}
	if len(errors.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(errors.Errors))
	}
	if errors.Errors[0].Code != "end_page_invalid" {
		t.Errorf("expected error code 'end_page_invalid', got '%s'", errors.Errors[0].Code)
	}
}
