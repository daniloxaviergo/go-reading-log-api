package validation

import "fmt"

// ValidateStartEndPage validates that startPage and endPage are non-negative
// and that startPage <= endPage.
// Returns a ValidationError if any validation fails.
func ValidateStartEndPage(startPage int, endPage int) *ValidationError {
	// Validate startPage is non-negative
	if startPage < 0 {
		return NewValidationError(
			"start_page_invalid",
			"start_page",
			fmt.Sprintf("start_page (%d) cannot be negative", startPage),
		)
	}

	// Validate endPage is non-negative
	if endPage < 0 {
		return NewValidationError(
			"end_page_invalid",
			"end_page",
			fmt.Sprintf("end_page (%d) cannot be negative", endPage),
		)
	}

	// Validate startPage <= endPage
	if startPage > endPage {
		return NewValidationError(
			"start_page_exceeds_end_page",
			"start_page",
			fmt.Sprintf("start_page (%d) cannot exceed end_page (%d)", startPage, endPage),
		)
	}

	return nil
}

// ValidateLog performs comprehensive validation on log page fields.
// Returns a ValidationErrorList containing all validation errors found.
// Validates: 0 <= startPage <= endPage.
func ValidateLog(startPage int, endPage int) *ValidationErrorList {
	errors := &ValidationErrorList{}

	// Validate startPage and endPage relationship
	if err := ValidateStartEndPage(startPage, endPage); err != nil {
		errors.AddError(err)
	}

	return errors
}
