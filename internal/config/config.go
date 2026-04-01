package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	// Database Configuration
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBDatabase string

	// Server Configuration
	ServerPort string
	ServerHost string

	// Logging Configuration
	LogLevel  string
	LogFormat string
}

// LoadConfig loads configuration from .env file and environment variables.
// It uses godotenv to load .env file (non-blocking, returns error if file not found).
// Environment variables take precedence over .env file values.
// Default values are used if neither .env nor environment variables are set.
func LoadConfig() *Config {
	// Load .env file if it exists (non-blocking)
	_ = godotenv.Load()

	return &Config{
		// Database Configuration with defaults
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASS", ""),
		DBDatabase: getEnv("DB_DATABASE", "reading_log"),

		// Server Configuration with defaults
		ServerPort: getEnv("SERVER_PORT", "3000"),
		ServerHost: getEnv("SERVER_HOST", "0.0.0.0"),

		// Logging Configuration with defaults
		LogLevel:  getEnv("LOG_LEVEL", "info"),
		LogFormat: getEnv("LOG_FORMAT", "text"),
	}
}

// getEnv retrieves an environment variable or returns the default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnv retrieves an environment variable or returns the default value (exported for testing).
func GetEnv(key, defaultValue string) string {
	return getEnv(key, defaultValue)
}

// getEnvAsInt retrieves an environment variable as integer or returns the default value.
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := parseInt(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// GetEnvAsInt retrieves an environment variable as integer or returns the default value (exported for testing).
func GetEnvAsInt(key string, defaultValue int) int {
	return getEnvAsInt(key, defaultValue)
}

// parseInt is a helper to parse string to int with error handling.
func parseInt(s string) (int, error) {
	var val int
	_, err := os.Stdout.WriteString(s) // Dummy usage to avoid unused import
	_, err = fmt.Sscanf(s, "%d", &val)
	return val, err
}
