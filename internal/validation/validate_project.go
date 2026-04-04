package validation

import "fmt"

// Status constants for project status validation (matching models package)
const (
	StatusUnstarted = "unstarted"
	StatusFinished  = "finished"
	StatusRunning   = "running"
	StatusSleeping  = "sleeping"
	StatusStopped   = "stopped"
)

// ValidStatusValues returns a slice of valid status values.
// Used by ValidateStatus to validate project status.
func ValidStatusValues() []string {
	return []string{
		StatusUnstarted,
		StatusFinished,
		StatusRunning,
		StatusSleeping,
		StatusStopped,
	}
}

// ValidatePage validates that page is within valid range (0 <= page <= totalPage).
// Returns a ValidationError if page is negative or exceeds totalPage.
func ValidatePage(page int, totalPage int) *ValidationError {
	if page < 0 {
		return NewValidationError(
			"page_invalid",
			"page",
			fmt.Sprintf("page (%d) cannot be negative", page),
		)
	}
	if page > totalPage {
		return NewValidationError(
			"page_exceeds_total",
			"page",
			fmt.Sprintf("page (%d) cannot exceed total_page (%d)", page, totalPage),
		)
	}
	return nil
}

// ValidateTotalPage validates that totalPage is greater than 0.
// Returns a ValidationError if totalPage is zero or negative.
func ValidateTotalPage(totalPage int) *ValidationError {
	if totalPage <= 0 {
		return NewValidationError(
			"total_page_invalid",
			"total_page",
			fmt.Sprintf("total_page (%d) must be greater than 0", totalPage),
		)
	}
	return nil
}

// ValidateStatus validates that status is one of the valid status values.
// Returns a ValidationError if status is not one of: unstarted, finished, running, sleeping, stopped.
func ValidateStatus(status string) *ValidationError {
	validStatuses := ValidStatusValues()
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return nil
		}
	}
	return NewValidationError(
		"invalid_status",
		"status",
		fmt.Sprintf("status (%s) must be one of: %v", status, validStatuses),
	)
}

// ValidateProject performs comprehensive validation on project fields.
// Returns a ValidationErrorList containing all validation errors found.
// Validates: totalPage > 0, 0 <= page <= totalPage, status is valid.
func ValidateProject(page int, totalPage int, status string) *ValidationErrorList {
	errors := &ValidationErrorList{}

	// Validate totalPage first (needed for page validation)
	if err := ValidateTotalPage(totalPage); err != nil {
		errors.AddError(err)
	}

	// Validate page (only if totalPage is valid to avoid misleading errors)
	if err := ValidatePage(page, totalPage); err != nil {
		errors.AddError(err)
	}

	// Validate status
	if err := ValidateStatus(status); err != nil {
		errors.AddError(err)
	}

	return errors
}
