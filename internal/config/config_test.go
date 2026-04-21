package config

import (
	"os"
	"testing"
	"time"
)

// TestLoadConfigDefaultValues tests that LoadConfig returns default values
// when no .env file or environment variables are set.
func TestLoadConfigDefaultValues(t *testing.T) {
	// Ensure no env vars are set
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASS")
	os.Unsetenv("DB_DATABASE")
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("SERVER_HOST")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("LOG_FORMAT")

	config := LoadConfig()

	// Verify database defaults
	if config.DBHost != "localhost" {
		t.Errorf("expected DB_HOST='localhost', got '%s'", config.DBHost)
	}
	if config.DBPort != 5432 {
		t.Errorf("expected DB_PORT=5432, got %d", config.DBPort)
	}
	if config.DBUser != "postgres" {
		t.Errorf("expected DB_USER='postgres', got '%s'", config.DBUser)
	}
	if config.DBPassword != "" {
		t.Errorf("expected DB_PASS='', got '%s'", config.DBPassword)
	}
	if config.DBDatabase != "reading_log" {
		t.Errorf("expected DB_DATABASE='reading_log', got '%s'", config.DBDatabase)
	}

	// Verify server defaults
	if config.ServerPort != "3000" {
		t.Errorf("expected SERVER_PORT='3000', got '%s'", config.ServerPort)
	}
	if config.ServerHost != "0.0.0.0" {
		t.Errorf("expected SERVER_HOST='0.0.0.0', got '%s'", config.ServerHost)
	}

	// Verify logging defaults
	if config.LogLevel != "info" {
		t.Errorf("expected LOG_LEVEL='info', got '%s'", config.LogLevel)
	}
	if config.LogFormat != "text" {
		t.Errorf("expected LOG_FORMAT='text', got '%s'", config.LogFormat)
	}
}

// TestLoadConfigEnvironmentVariables tests that environment variables override defaults.
func TestLoadConfigEnvironmentVariables(t *testing.T) {
	// Set environment variables
	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "custom_user")
	os.Setenv("DB_PASS", "secret123")
	os.Setenv("DB_DATABASE", "custom_db")
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("SERVER_HOST", "127.0.0.1")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("LOG_FORMAT", "json")

	config := LoadConfig()

	// Verify environment variable values are loaded
	if config.DBHost != "db.example.com" {
		t.Errorf("expected DB_HOST='db.example.com', got '%s'", config.DBHost)
	}
	if config.DBPort != 5433 {
		t.Errorf("expected DB_PORT=5433, got %d", config.DBPort)
	}
	if config.DBUser != "custom_user" {
		t.Errorf("expected DB_USER='custom_user', got '%s'", config.DBUser)
	}
	if config.DBPassword != "secret123" {
		t.Errorf("expected DB_PASS='secret123', got '%s'", config.DBPassword)
	}
	if config.DBDatabase != "custom_db" {
		t.Errorf("expected DB_DATABASE='custom_db', got '%s'", config.DBDatabase)
	}
	if config.ServerPort != "8080" {
		t.Errorf("expected SERVER_PORT='8080', got '%s'", config.ServerPort)
	}
	if config.ServerHost != "127.0.0.1" {
		t.Errorf("expected SERVER_HOST='127.0.0.1', got '%s'", config.ServerHost)
	}
	if config.LogLevel != "debug" {
		t.Errorf("expected LOG_LEVEL='debug', got '%s'", config.LogLevel)
	}
	if config.LogFormat != "json" {
		t.Errorf("expected LOG_FORMAT='json', got '%s'", config.LogFormat)
	}
}

// TestLoadConfigEnvVarCaseInsensitivity tests that environment variables
// with different cases are handled correctly.
func TestLoadConfigEnvVarCaseInsensitivity(t *testing.T) {
	// Test with lowercase log level
	os.Setenv("LOG_LEVEL", "DEBUG")
	config := LoadConfig()
	if config.LogLevel != "DEBUG" {
		t.Errorf("expected LOG_LEVEL='DEBUG', got '%s'", config.LogLevel)
	}

	// Test with uppercase log level
	os.Setenv("LOG_LEVEL", "DEBUG")
	config = LoadConfig()
	if config.LogLevel != "DEBUG" {
		t.Errorf("expected LOG_LEVEL='DEBUG', got '%s'", config.LogLevel)
	}
}

// TestLoadConfigInvalidPort tests that invalid port values fall back to defaults.
func TestLoadConfigInvalidPort(t *testing.T) {
	os.Setenv("DB_PORT", "invalid")
	config := LoadConfig()

	if config.DBPort != 5432 {
		t.Errorf("expected DB_PORT=5432 (default) for invalid value, got %d", config.DBPort)
	}
}

// TestLoadConfigEmptyEnvVars tests that empty environment variables use defaults.
func TestLoadConfigEmptyEnvVars(t *testing.T) {
	os.Setenv("DB_HOST", "")
	os.Setenv("SERVER_PORT", "")

	config := LoadConfig()

	if config.DBHost != "localhost" {
		t.Errorf("expected DB_HOST='localhost' for empty env var, got '%s'", config.DBHost)
	}
	if config.ServerPort != "3000" {
		t.Errorf("expected SERVER_PORT='3000' for empty env var, got '%s'", config.ServerPort)
	}
}

// TestLoadConfigStatusRangesDefaultValues tests that status range config uses defaults.
func TestLoadConfigStatusRangesDefaultValues(t *testing.T) {
	// Ensure no env vars are set
	os.Unsetenv("EM_ANDAMENTO_RANGE")
	os.Unsetenv("DORMINDO_RANGE")

	config := LoadConfig()

	// Verify default status range values
	// em_andamento_range = 7 days
	if config.EmAndamentoRange != 7 {
		t.Errorf("expected EmAndamentoRange=7 (default), got %d", config.EmAndamentoRange)
	}
	// dormindo_range = 14 days
	if config.DormindoRange != 14 {
		t.Errorf("expected DormindoRange=14 (default), got %d", config.DormindoRange)
	}
}

// TestLoadConfigStatusRangesEnvironmentVariables tests that env vars override status range defaults.
func TestLoadConfigStatusRangesEnvironmentVariables(t *testing.T) {
	// Set environment variables for status ranges
	os.Setenv("EM_ANDAMENTO_RANGE", "10")
	os.Setenv("DORMINDO_RANGE", "21")

	config := LoadConfig()

	// Verify environment variable values are loaded
	if config.EmAndamentoRange != 10 {
		t.Errorf("expected EmAndamentoRange=10, got %d", config.EmAndamentoRange)
	}
	if config.DormindoRange != 21 {
		t.Errorf("expected DormindoRange=21, got %d", config.DormindoRange)
	}
}

// TestLoadConfigStatusRangesGetterMethods tests the getter methods for status ranges.
func TestLoadConfigStatusRangesGetterMethods(t *testing.T) {
	os.Unsetenv("EM_ANDAMENTO_RANGE")
	os.Unsetenv("DORMINDO_RANGE")

	config := LoadConfig()

	// Test getter methods return correct values
	if config.GetEmAndamentoRange() != 7 {
		t.Errorf("expected GetEmAndamentoRange()=7, got %d", config.GetEmAndamentoRange())
	}
	if config.GetDormindoRange() != 14 {
		t.Errorf("expected GetDormindoRange()=14, got %d", config.GetDormindoRange())
	}
}

// TestLoadConfigStatusRangesInvalidValues tests that invalid values fall back to defaults.
func TestLoadConfigStatusRangesInvalidValues(t *testing.T) {
	os.Setenv("EM_ANDAMENTO_RANGE", "invalid")
	os.Setenv("DORMINDO_RANGE", "-5")

	config := LoadConfig()

	// Invalid values should fall back to defaults
	if config.EmAndamentoRange != 7 {
		t.Errorf("expected EmAndamentoRange=7 (default for invalid value), got %d", config.EmAndamentoRange)
	}
	if config.DormindoRange != 14 {
		t.Errorf("expected DormindoRange=14 (default for invalid value), got %d", config.DormindoRange)
	}
}

// TestLoadConfigStatusRangesEmptyValues tests that empty env vars use defaults.
func TestLoadConfigStatusRangesEmptyValues(t *testing.T) {
	os.Setenv("EM_ANDAMENTO_RANGE", "")
	os.Setenv("DORMINDO_RANGE", "")

	config := LoadConfig()

	// Empty values should use defaults
	if config.EmAndamentoRange != 7 {
		t.Errorf("expected EmAndamentoRange=7 (default for empty value), got %d", config.EmAndamentoRange)
	}
	if config.DormindoRange != 14 {
		t.Errorf("expected DormindoRange=14 (default for empty value), got %d", config.DormindoRange)
	}
}

// TestLoadConfigTimezoneDefault tests that timezone defaults to BRT.
func TestLoadConfigTimezoneDefault(t *testing.T) {
	os.Unsetenv("TZ_LOCATION")

	config := LoadConfig()

	if config.TZLocation == nil {
		t.Fatal("expected TZLocation to be set")
	}

	// Check that it's the Brazil timezone (BRT) by using FixedZone comparison
	expectedLoc := time.FixedZone("BRT", -3*60*60)

	// Compare string representation as Zone() is not exported
	if config.TZLocation.String() != expectedLoc.String() {
		t.Errorf("expected timezone '%s', got '%s'", expectedLoc.String(), config.TZLocation.String())
	}
}

// TestLoadConfigTimezoneFromEnv tests that TZ_LOCATION env var is loaded.
func TestLoadConfigTimezoneFromEnv(t *testing.T) {
	os.Setenv("TZ_LOCATION", "America/Sao_Paulo")

	config := LoadConfig()

	if config.TZLocation == nil {
		t.Fatal("expected TZLocation to be set")
	}

	// Verify the location can be used for time calculations
	// We use a Time value's Zone() method with our Location
	now := time.Now()
	_, offset := now.In(config.TZLocation).Zone()

	// Sao Paulo is UTC-3 (like BRT)
	expectedOffset := -3 * 60 * 60

	if offset != expectedOffset {
		t.Errorf("expected timezone offset %d, got %d", expectedOffset, offset)
	}
}

// TestLoadConfigTimezoneInvalidFallback tests that invalid TZ_LOCATION falls back to BRT.
func TestLoadConfigTimezoneInvalidFallback(t *testing.T) {
	os.Setenv("TZ_LOCATION", "Invalid/Timezone")

	config := LoadConfig()

	if config.TZLocation == nil {
		t.Fatal("expected TZLocation to be set")
	}

	// Should fallback to BRT
	expectedLoc := time.FixedZone("BRT", -3*60*60)

	if config.TZLocation.String() != expectedLoc.String() {
		t.Errorf("expected timezone '%s' (fallback), got '%s'", expectedLoc.String(), config.TZLocation.String())
	}
}

// TestLoadConfigTimezoneEmptyFallback tests that empty TZ_LOCATION falls back to BRT.
func TestLoadConfigTimezoneEmptyFallback(t *testing.T) {
	os.Setenv("TZ_LOCATION", "")

	config := LoadConfig()

	if config.TZLocation == nil {
		t.Fatal("expected TZLocation to be set")
	}

	// Should fallback to BRT
	expectedLoc := time.FixedZone("BRT", -3*60*60)

	if config.TZLocation.String() != expectedLoc.String() {
		t.Errorf("expected timezone '%s' (fallback), got '%s'", expectedLoc.String(), config.TZLocation.String())
	}
}
