package performance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// BaselineStore manages baseline measurements for performance comparisons
type BaselineStore struct {
	filePath string
	data     *BaselineData
}

// BaselineData represents the stored baseline metrics
type BaselineData struct {
	Version     string                    `json:"version"`
	CreatedAt   time.Time                 `json:"created_at"`
	CommitSHA   string                    `json:"commit_sha"`
	Benchmarks  map[string]BenchmarkStats `json:"benchmarks"`
	Environment EnvironmentInfo           `json:"environment"`
}

// BenchmarkStats stores statistics for a single benchmark
type BenchmarkStats struct {
	AverageMs       float64 `json:"average_ms"`
	P50Ms           float64 `json:"p50_ms"`
	P95Ms           float64 `json:"p95_ms"`
	P99Ms           float64 `json:"p99_ms"`
	AllocsPerOp     uint64  `json:"allocs_per_op"`
	OpsPerSec       float64 `json:"ops_per_sec"`
	ThresholdMs     float64 `json:"threshold_ms"`
	WithinThreshold bool    `json:"within_threshold"`
}

// EnvironmentInfo captures the environment where benchmarks were run
type EnvironmentInfo struct {
	GoVersion  string `json:"go_version"`
	OS         string `json:"os"`
	Arch       string `json:"arch"`
	CPU        string `json:"cpu"`
	MemoryMB   int64  `json:"memory_mb"`
	PostgreSQL string `json:"postgresql_version"`
}

// NewBaselineStore creates a new baseline store
func NewBaselineStore(filePath string) (*BaselineStore, error) {
	store := &BaselineStore{
		filePath: filePath,
		data: &BaselineData{
			Version:    "1.0.0",
			CreatedAt:  time.Now(),
			Benchmarks: make(map[string]BenchmarkStats),
		},
	}

	// Load existing baseline if available
	if err := store.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load baseline: %w", err)
	}

	return store, nil
}

// load reads the baseline from disk
func (bs *BaselineStore) load() error {
	data, err := os.ReadFile(bs.filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &bs.data)
}

// Save writes the baseline to disk
func (bs *BaselineStore) Save() error {
	// Ensure directory exists
	dir := filepath.Dir(bs.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := json.MarshalIndent(bs.data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return os.WriteFile(bs.filePath, data, 0644)
}

// UpdateBenchmark updates a benchmark's stats in the baseline
func (bs *BaselineStore) UpdateBenchmark(name string, stats BenchmarkStats) {
	bs.data.Benchmarks[name] = stats
}

// GetBenchmark returns stats for a specific benchmark, or false if not found
func (bs *BaselineStore) GetBenchmark(name string) (BenchmarkStats, bool) {
	stats, ok := bs.data.Benchmarks[name]
	return stats, ok
}

// GetComparison calculates the percentage difference between current and baseline
func (bs *BaselineStore) GetComparison(currentMs, baselineMs float64) float64 {
	if baselineMs == 0 {
		return 0
	}
	return ((currentMs - baselineMs) / baselineMs) * 100
}

// VerifyThreshold checks if regression is within acceptable threshold
func (bs *BaselineStore) VerifyThreshold(currentMs, baselineMs, thresholdPercent float64) bool {
	regression := bs.GetComparison(currentMs, baselineMs)
	return regression <= thresholdPercent
}

// GetBaselinePath returns the default baseline file path
func GetBaselinePath() string {
	// Use a relative path that works from project root
	return filepath.Join("test", "performance", "baseline.json")
}

// LoadOrCreateBaseline loads existing baseline or creates new one
func LoadOrCreateBaseline(filePath string) (*BaselineStore, error) {
	store, err := NewBaselineStore(filePath)
	if err != nil {
		return nil, err
	}

	// If no baseline exists, create initial one
	if len(store.data.Benchmarks) == 0 {
		if err := store.Save(); err != nil {
			return nil, fmt.Errorf("failed to create initial baseline: %w", err)
		}
	}

	return store, nil
}
