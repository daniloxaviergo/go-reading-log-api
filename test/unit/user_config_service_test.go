package unit

import (
	"os"
	"path/filepath"
	"testing"

	. "go-reading-log-api-next/internal/service"
)

// TestUserConfigService_LoadFromFile tests loading configuration from a valid YAML file.
func TestUserConfigService_LoadFromFile(t *testing.T) {
	// Create a temporary directory for the test config file
	tempDir, err := os.MkdirTemp("", "user_config_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create test YAML content
	testConfig := `max_faults: 15
prediction_pct: 0.20
pages_per_day: 30.5
`

	// Write to temp file
	configPath := filepath.Join(tempDir, "dashboard.yaml")
	if err := os.WriteFile(configPath, []byte(testConfig), 0644); err != nil {
		t.Fatal(err)
	}

	// Load configuration
	svc, err := LoadDashboardConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify values
	if svc.GetMaxFaults() != 15 {
		t.Errorf("Expected max_faults=15, got %d", svc.GetMaxFaults())
	}
	if svc.GetPredictionPct() != 0.20 {
		t.Errorf("Expected prediction_pct=0.20, got %f", svc.GetPredictionPct())
	}
	if svc.GetPagesPerDay() != 30.5 {
		t.Errorf("Expected pages_per_day=30.5, got %f", svc.GetPagesPerDay())
	}
}

// TestUserConfigService_ZeroValues tests behavior when config has zero values.
func TestUserConfigService_ZeroValues(t *testing.T) {
	// Create a temporary directory for the test config file
	tempDir, err := os.MkdirTemp("", "user_config_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Write YAML with zero values
	configPath := filepath.Join(tempDir, "dashboard.yaml")
	zeroYAML := `max_faults: 0
prediction_pct: 0.0
pages_per_day: 0.0
`
	if err := os.WriteFile(configPath, []byte(zeroYAML), 0644); err != nil {
		t.Fatal(err)
	}

	// Load configuration
	svc, err := LoadDashboardConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config with zero values: %v", err)
	}

	// Verify zero values are preserved (not replaced with defaults)
	if svc.GetMaxFaults() != 0 {
		t.Errorf("Expected max_faults=0, got %d", svc.GetMaxFaults())
	}
	if svc.GetPredictionPct() != 0.0 {
		t.Errorf("Expected prediction_pct=0.0, got %f", svc.GetPredictionPct())
	}
	if svc.GetPagesPerDay() != 0.0 {
		t.Errorf("Expected pages_per_day=0.0, got %f", svc.GetPagesPerDay())
	}
}

// TestUserConfigService_FileNotFound tests behavior when config file doesn't exist.
func TestUserConfigService_FileNotFound(t *testing.T) {
	// Load from non-existent file
	svc, err := LoadDashboardConfig("/nonexistent/path/dashboard.yaml")

	if err != nil {
		t.Fatalf("Expected no error for missing file (graceful fallback), got: %v", err)
	}

	// Verify defaults are used
	if svc.GetMaxFaults() != 10 {
		t.Errorf("Expected default max_faults=10, got %d", svc.GetMaxFaults())
	}
	if svc.GetPredictionPct() != 0.15 {
		t.Errorf("Expected default prediction_pct=0.15, got %f", svc.GetPredictionPct())
	}
	if svc.GetPagesPerDay() != 25.0 {
		t.Errorf("Expected default pages_per_day=25.0, got %f", svc.GetPagesPerDay())
	}
}

// TestUserConfigService_InvalidYAML tests behavior when YAML file is malformed.
func TestUserConfigService_InvalidYAML(t *testing.T) {
	// Create a temporary directory for the test config file
	tempDir, err := os.MkdirTemp("", "user_config_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Write invalid YAML content
	configPath := filepath.Join(tempDir, "dashboard.yaml")
	invalidYAML := `max_faults: 10
invalid yaml here
  - missing value
`
	if err := os.WriteFile(configPath, []byte(invalidYAML), 0644); err != nil {
		t.Fatal(err)
	}

	// Load configuration
	svc, err := LoadDashboardConfig(configPath)

	if err != nil {
		t.Fatalf("Expected no error for invalid YAML (graceful fallback), got: %v", err)
	}

	// Verify defaults are used
	if svc.GetMaxFaults() != 10 {
		t.Errorf("Expected default max_faults=10, got %d", svc.GetMaxFaults())
	}
}

// TestUserConfigService_DefaultValues tests that all defaults are applied correctly.
func TestUserConfigService_DefaultValues(t *testing.T) {
	svc := NewUserConfigService(GetDefaultConfig())

	// Verify all default values
	if svc.GetMaxFaults() != 10 {
		t.Errorf("Expected max_faults=10, got %d", svc.GetMaxFaults())
	}
	if svc.GetPredictionPct() != 0.15 {
		t.Errorf("Expected prediction_pct=0.15, got %f", svc.GetPredictionPct())
	}
	if svc.GetPagesPerDay() != 25.0 {
		t.Errorf("Expected pages_per_day=25.0, got %f", svc.GetPagesPerDay())
	}
}

// TestUserConfigService_PartialConfig tests behavior when config has missing fields.
func TestUserConfigService_PartialConfig(t *testing.T) {
	// Create a temporary directory for the test config file
	tempDir, err := os.MkdirTemp("", "user_config_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Write partial YAML (only one field defined)
	configPath := filepath.Join(tempDir, "dashboard.yaml")
	partialYAML := `max_faults: 20
`
	if err := os.WriteFile(configPath, []byte(partialYAML), 0644); err != nil {
		t.Fatal(err)
	}

	// Load configuration
	svc, err := LoadDashboardConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load partial config: %v", err)
	}

	// Verify defined value is used, others use defaults
	if svc.GetMaxFaults() != 20 {
		t.Errorf("Expected max_faults=20, got %d", svc.GetMaxFaults())
	}
	if svc.GetPredictionPct() != 0.15 {
		t.Errorf("Expected default prediction_pct=0.15, got %f", svc.GetPredictionPct())
	}
	if svc.GetPagesPerDay() != 25.0 {
		t.Errorf("Expected default pages_per_day=25.0, got %f", svc.GetPagesPerDay())
	}
}

// TestUserConfigService_Getters tests all getter methods return correct values.
func TestUserConfigService_Getters(t *testing.T) {
	maxFaults := 10
	predictionPct := 0.15
	pagesPerDay := 25.0
	config := &DashboardConfig{
		MaxFaults:     &maxFaults,
		PredictionPct: &predictionPct,
		PagesPerDay:   &pagesPerDay,
	}
	svc := NewUserConfigService(config)

	tests := []struct {
		name     string
		getter   func() int
		expected int
	}{
		{"GetMaxFaults", svc.GetMaxFaults, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.getter()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}

	// Test float64 getters separately due to type mismatch
	if svc.GetPredictionPct() != 0.15 {
		t.Errorf("Expected prediction_pct=0.15, got %f", svc.GetPredictionPct())
	}
	if svc.GetPagesPerDay() != 25.0 {
		t.Errorf("Expected pages_per_day=25.0, got %f", svc.GetPagesPerDay())
	}
}

// TestUserConfigService_EmptyConfig tests behavior when config file is empty.
func TestUserConfigService_EmptyConfig(t *testing.T) {
	// Create a temporary directory for the test config file
	tempDir, err := os.MkdirTemp("", "user_config_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Write empty YAML content
	configPath := filepath.Join(tempDir, "dashboard.yaml")
	if err := os.WriteFile(configPath, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	// Load configuration
	svc, err := LoadDashboardConfig(configPath)

	if err != nil {
		t.Fatalf("Expected no error for empty config (graceful fallback), got: %v", err)
	}

	// Verify defaults are used
	if svc.GetMaxFaults() != 10 {
		t.Errorf("Expected default max_faults=10, got %d", svc.GetMaxFaults())
	}
	if svc.GetPredictionPct() != 0.15 {
		t.Errorf("Expected default prediction_pct=0.15, got %f", svc.GetPredictionPct())
	}
	if svc.GetPagesPerDay() != 25.0 {
		t.Errorf("Expected default pages_per_day=25.0, got %f", svc.GetPagesPerDay())
	}
}
