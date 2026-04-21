package config

import (
	"fmt"
	"os"
	"time"

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

	// Status Range Configuration (in days)
	// Default values match Rails V1::UserConfig: em_andamento_range = 7, dormindo_range = 14
	// Note: Rails uses 8 and 16 days respectively, task specification requires 7 and 14
	EmAndamentoRange int
	DormindoRange    int

	// Timezone Configuration
	// TZLocation holds the parsed time.Location for date calculations
	// Defaults to Brazil timezone (BRT) if not configured
	TZLocation *time.Location
}

// LoadConfig loads configuration from .env file and environment variables.
// It uses godotenv to load .env file (non-blocking, returns error if file not found).
// Environment variables take precedence over .env file values.
// Default values are used if neither .env nor environment variables are set.
func LoadConfig() *Config {
	// Load .env file if it exists (non-blocking)
	_ = godotenv.Load()

	// Parse timezone configuration
	tzLocation := parseTZLocation(getEnv("TZ_LOCATION", ""))

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

		// Status Range Configuration with defaults (in days)
		// em_andamento_range = 7 days (Rails uses 8)
		// dormindo_range = 14 days (Rails uses 16)
		EmAndamentoRange: getEnvAsInt("EM_ANDAMENTO_RANGE", 7),
		DormindoRange:    getEnvAsInt("DORMINDO_RANGE", 14),

		// Timezone Configuration
		TZLocation: tzLocation,
	}
}

// parseTZLocation parses timezone string and returns time.Location.
// Falls back to Brazil timezone (BRT) if parsing fails.
func parseTZLocation(tzStr string) *time.Location {
	if tzStr == "" {
		// Default to Brazil timezone (BRT) matching Rails behavior
		return time.FixedZone("BRT", -3*60*60)
	}

	loc, err := time.LoadLocation(tzStr)
	if err != nil {
		// Fallback to BRT on error
		fmt.Fprintf(os.Stderr, "Warning: Failed to load timezone '%s', using BRT fallback: %v\n", tzStr, err)
		return time.FixedZone("BRT", -3*60*60)
	}

	return loc
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
			// Status ranges must be positive values, fallback to default otherwise
			if intVal > 0 {
				return intVal
			}
		}
	}
	return defaultValue
}

// GetEnvAsInt retrieves an environment variable as integer or returns the default value (exported for testing).
func GetEnvAsInt(key string, defaultValue int) int {
	return getEnvAsInt(key, defaultValue)
}

// GetEmAndamentoRange returns the em_andamento range in days.
func (c *Config) GetEmAndamentoRange() int {
	return c.EmAndamentoRange
}

// GetDormindoRange returns the dormindo range in days.
func (c *Config) GetDormindoRange() int {
	return c.DormindoRange
}

// parseInt is a helper to parse string to int with error handling.
func parseInt(s string) (int, error) {
	var val int
	_, err := os.Stdout.WriteString(s) // Dummy usage to avoid unused import
	_, err = fmt.Sscanf(s, "%d", &val)
	return val, err
}
