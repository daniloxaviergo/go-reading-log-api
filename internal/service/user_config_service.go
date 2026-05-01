package service

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// DashboardConfig holds the dashboard configuration values loaded from YAML.
type DashboardConfig struct {
	MaxFaults     *int     `yaml:"max_faults"`
	PredictionPct *float64 `yaml:"prediction_pct"`
	PagesPerDay   *float64 `yaml:"pages_per_day"`
}

// UserConfigService provides access to dashboard configuration with file-based
// loading and hardcoded defaults as fallback.
type UserConfigService struct {
	config *DashboardConfig
}

// LoadDashboardConfig loads the dashboard configuration from the specified YAML file path.
// If the file doesn't exist or contains invalid YAML, it returns a default configuration
// with warning logs. The service will always return a valid configuration.
func LoadDashboardConfig(path string) (*UserConfigService, error) {
	// Try to load from file
	config, err := loadFromFile(path)
	if err != nil {
		// Log warning but continue with defaults (graceful fallback)
		fmt.Printf("Warning: Failed to load dashboard config from %s, using defaults: %v\n", path, err)
		return NewUserConfigService(GetDefaultConfig()), nil
	}

	// Merge with defaults for any nil fields (missing values)
	defaultConfig := GetDefaultConfig()
	if config.MaxFaults == nil {
		config.MaxFaults = defaultConfig.MaxFaults
	}
	if config.PredictionPct == nil {
		config.PredictionPct = defaultConfig.PredictionPct
	}
	if config.PagesPerDay == nil {
		config.PagesPerDay = defaultConfig.PagesPerDay
	}

	return NewUserConfigService(config), nil
}

// loadFromFile reads and parses the YAML configuration file.
func loadFromFile(path string) (*DashboardConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config DashboardConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &config, nil
}

// GetDefaultConfig returns the default configuration values.
func GetDefaultConfig() *DashboardConfig {
	maxFaults := 10
	predictionPct := 0.15
	pagesPerDay := 25.0
	return &DashboardConfig{
		MaxFaults:     &maxFaults,
		PredictionPct: &predictionPct,
		PagesPerDay:   &pagesPerDay,
	}
}

// NewUserConfigService creates a new UserConfigService with the given configuration.
func NewUserConfigService(config *DashboardConfig) *UserConfigService {
	return &UserConfigService{
		config: config,
	}
}

// GetMaxFaults returns the maximum number of faults allowed before an alert is triggered.
func (s *UserConfigService) GetMaxFaults() int {
	if s.config.MaxFaults != nil {
		return *s.config.MaxFaults
	}
	return 0
}

// GetPredictionPct returns the prediction percentage for speculative calculations.
func (s *UserConfigService) GetPredictionPct() float64 {
	if s.config.PredictionPct != nil {
		return *s.config.PredictionPct
	}
	return 0.0
}

// GetPagesPerDay returns the default pages per day target.
func (s *UserConfigService) GetPagesPerDay() float64 {
	if s.config.PagesPerDay != nil {
		return *s.config.PagesPerDay
	}
	return 0.0
}
