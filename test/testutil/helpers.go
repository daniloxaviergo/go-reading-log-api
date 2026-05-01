// helpers.go
// Common test helper functions that can be imported by test files

package testutil

// floatPtr returns a pointer to the given float64 value
func FloatPtr(f float64) *float64 {
	return &f
}

// IntPtr returns a pointer to the given int value
func IntPtr(i int) *int {
	return &i
}

// Int64Ptr returns a pointer to the given int64 value
func Int64Ptr(i int64) *int64 {
	return &i
}

// BoolPtr returns a pointer to the given bool value
func BoolPtr(b bool) *bool {
	return &b
}

// StringPtr returns a pointer to the given string value
func StringPtr(s string) *string {
	return &s
}
