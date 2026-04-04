package validation

import "fmt"

// ValidationError represents a validation error with a code, field, and message
type ValidationError struct {
	Code    string
	Field   string
	Message string
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrorList represents a list of validation errors
type ValidationErrorList struct {
	Errors []*ValidationError
}

// Error implements the error interface for ValidationErrorList
func (l *ValidationErrorList) Error() string {
	if len(l.Errors) == 0 {
		return "validation error: no errors"
	}
	if len(l.Errors) == 1 {
		return l.Errors[0].Error()
	}
	return fmt.Sprintf("validation error: %d errors occurred", len(l.Errors))
}

// AddError adds a validation error to the list
func (l *ValidationErrorList) AddError(err *ValidationError) {
	l.Errors = append(l.Errors, err)
}

// HasErrors returns true if the error list contains any errors
func (l *ValidationErrorList) HasErrors() bool {
	return len(l.Errors) > 0
}

// ToMap converts the validation error list to a map for JSON response
func (l *ValidationErrorList) ToMap() map[string]interface{} {
	if len(l.Errors) == 0 {
		return nil
	}

	result := make(map[string]interface{})
	for _, err := range l.Errors {
		result[err.Field] = err.Message
	}
	return result
}

// NewValidationError creates a new ValidationError
func NewValidationError(code, field, message string) *ValidationError {
	return &ValidationError{
		Code:    code,
		Field:   field,
		Message: message,
	}
}
