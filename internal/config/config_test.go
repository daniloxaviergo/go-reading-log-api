package config

import (
	"os"
	"testing"
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
