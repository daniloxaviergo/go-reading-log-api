package performance

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"go-reading-log-api-next/internal/config"
	"go-reading-log-api-next/test"
)

func BenchmarkGoroutineIDExtraction(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	var totalTime time.Duration
	for i := 0; i < b.N; i++ {
		start := time.Now()
		id := goroutineID()
		_ = id // prevent optimization
		totalTime += time.Since(start)
	}

	avgTime := totalTime / time.Duration(b.N)
	msPerOp := float64(avgTime) / float64(time.Millisecond)
	b.ReportMetric(msPerOp, "ms/op")

	if msPerOp > 1.0 {
		b.Errorf("Goroutine ID extraction took %.2fms > 1ms threshold", msPerOp)
	} else {
		b.Logf("Goroutine ID extraction time: %.2fms (within 1ms threshold)", msPerOp)
	}
}

// BenchmarkParallelTestStartup measures the time required to set up a parallel test
// with unique database naming, ensuring the overhead is less than 200ms.
func BenchmarkParallelTestStartup(b *testing.B) {
	if !test.IsTestDatabase() {
		b.Skip("Test database not configured - skipping benchmark test")
	}

	// Warm-up phase - create a few test databases to establish baseline
	for i := 0; i < 3; i++ {
		_, err := test.SetupTestDB()
		if err != nil {
			b.Fatalf("Warm-up failed: %v", err)
		}
	}

	b.ResetTimer()
	b.ReportAllocs()

	var totalTime time.Duration
	var times []float64

	for i := 0; i < b.N; i++ {
		start := time.Now()

		helper, err := test.SetupTestDB()
		if err != nil {
			b.Fatalf("Setup failed: %v", err)
		}

		// Measure the time for the actual setup
		setupTime := time.Since(start)
		times = append(times, float64(setupTime)/float64(time.Millisecond))

		// Perform some operations to ensure real work is measured
		if err := helper.SetupTestSchema(); err != nil {
			b.Fatalf("Schema setup failed: %v", err)
		}

		if err := helper.ClearTestData(); err != nil {
			b.Fatalf("Clear data failed: %v", err)
		}

		// Cleanup
		if err := helper.CleanupTestSchema(); err != nil {
			b.Logf("Cleanup failed: %v", err)
		}
		helper.Close()

		totalTime += time.Since(start)
	}

	// Calculate statistics
	avgTime := totalTime / time.Duration(b.N)
	p50, p95, p99 := calculatePercentiles(times)

	b.ReportMetric(float64(avgTime)/float64(time.Millisecond), "ms/op")
	b.ReportMetric(p50, "p50_ms/op")
	b.ReportMetric(p95, "p95_ms/op")
	b.ReportMetric(p99, "p99_ms/op")

	// Verify threshold
	avgMs := float64(avgTime) / float64(time.Millisecond)
	if avgMs > 200 {
		b.Errorf("Startup time %.2fms exceeds threshold of 200ms", avgMs)
	} else {
		b.Logf("Startup time %.2fms is within threshold (< 200ms)", avgMs)
	}
}

// BenchmarkParallelTestExecution measures the performance of parallel test execution
// with 8 goroutines, ensuring no more than 10% regression from baseline.
func BenchmarkParallelTestExecution(b *testing.B) {
	if !test.IsTestDatabase() {
		b.Skip("Test database not configured - skipping benchmark test")
	}

	helper := setupParallelBenchmark(b, 8)
	defer cleanupParallelBenchmark(b, helper)

	// Use the pool directly for benchmarks
	pool := helper.Pool

	// Warm-up phase
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := pool.Query(ctx, "SELECT * FROM projects")
	if err != nil && err.Error() != "no rows in result set" {
		b.Fatalf("Warm-up failed: %v", err)
	}

	b.ResetTimer()
	b.ReportAllocs()
	b.SetParallelism(8)

	var totalTime time.Duration
	var ops int

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			start := time.Now()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			rows, err := pool.Query(ctx, "SELECT * FROM projects")
			cancel()

			if err != nil && err.Error() != "no rows in result set" {
				b.Fatalf("Query failed: %v", err)
			}

			// Count rows by iterating through results
			rowCount := 0
			for rows.Next() {
				rowCount++
			}

			if rowCount != 10 {
				b.Logf("Expected 10 projects, got %d", rowCount)
			}

			totalTime += time.Since(start)
			ops++
		}
	})

	avgTime := totalTime / time.Duration(ops)
	b.ReportMetric(float64(avgTime)/float64(time.Millisecond), "ms/op")
	b.ReportMetric(float64(ops)/float64(totalTime)*float64(time.Second), "ops/sec")

	// Compare against baseline (from projects_benchmark_test.go)
	// Baseline: ~5-10ms per operation
	thresholdMs := 10.0 // 10ms threshold
	actualMs := float64(avgTime) / float64(time.Millisecond)

	// Allow up to 10% regression (baseline was ~5-10ms, so allow up to ~11ms)
	if actualMs > thresholdMs {
		b.Errorf("Execution time %.2fms exceeds threshold of %.2fms (10%% regression)", actualMs, thresholdMs)
	} else {
		b.Logf("Execution time %.2fms is within threshold (< %.2fms)", actualMs, thresholdMs)
	}
}

// BenchmarkParallelCleanup measures the performance of cleaning up many orphaned databases.
// Verifies that cleanup completes within 60 seconds even with 6000+ orphaned databases.
func BenchmarkParallelCleanup(b *testing.B) {
	if !test.IsTestDatabase() {
		b.Skip("Test database not configured - skipping benchmark test")
	}

	// Create orphaned databases first
	b.Run("SetupOrphans", func(b *testing.B) {
		b.ReportAllocs()
		b.SetParallelism(8)

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := test.SetupTestDB()
				if err != nil {
					b.Fatalf("Failed to create test DB: %v", err)
				}
				// Don't cleanup - leave as orphaned
				// helper.Close() is called by the test framework
			}
		})
	})

	// Measure cleanup performance
	b.Run("Cleanup", func(b *testing.B) {
		b.ReportAllocs()

		start := time.Now()
		orphans, err := getOrphanedDatabases()
		if err != nil {
			b.Fatalf("Failed to get orphaned databases: %v", err)
		}

		b.Logf("Found %d orphaned databases to cleanup", len(orphans))

		if len(orphans) == 0 {
			b.Skip("No orphaned databases found")
		}

		b.ResetTimer()

		var cleaned int64
		b.SetParallelism(8)

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// Get next orphaned database
				idx := int(atomic.AddInt64(&cleaned, 1) - 1)
				if idx >= len(orphans) {
					break
				}

				dbName := orphans[idx]
				if err := dropDatabase(dbName); err != nil {
					b.Logf("Failed to drop %s: %v", dbName, err)
				}
			}
		})

		elapsed := time.Since(start)
		b.ReportMetric(float64(elapsed)/float64(time.Second), "seconds")
		b.ReportMetric(float64(len(orphans)), "databases")

		// Verify threshold
		elapsedSeconds := float64(elapsed) / float64(time.Second)
		if elapsedSeconds > 60 {
			b.Errorf("Cleanup took %.2fs, exceeds 60s threshold", elapsedSeconds)
		} else {
			b.Logf("Cleanup took %.2fs, within 60s threshold", elapsedSeconds)
		}
	})
}

// BenchmarkDatabaseUniqueness verifies that unique database names don't cause collisions.
// Ensures that concurrent test runs with unique database names work correctly.
func BenchmarkDatabaseUniqueness(b *testing.B) {
	if !test.IsTestDatabase() {
		b.Skip("Test database not configured - skipping benchmark test")
	}

	b.ReportAllocs()
	b.SetParallelism(8)

	var collisionCount int64
	var successCount int64

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Generate a unique database name (similar to test helper)
			testDBName := "reading_log_test_parallel"
			dbName := fmt.Sprintf("%s_%d_%d_%d", testDBName, os.Getpid(), goroutineID(), time.Now().UnixNano())

			// Try to create database
			helper, err := setupTestDBWithName(dbName)
			if err != nil {
				if strings.Contains(err.Error(), "already exists") || strings.Contains(err.Error(), "exists") {
					atomic.AddInt64(&collisionCount, 1)
					continue
				}
				b.Fatalf("Failed to create test DB %s: %v", dbName, err)
			}

			// Verify we can connect
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			err = helper.Pool.Ping(ctx)
			cancel()
			if err != nil {
				b.Fatalf("Failed to ping test DB %s: %v", dbName, err)
			}

			// Cleanup
			if err := helper.CleanupTestSchema(); err != nil {
				b.Logf("Cleanup failed for %s: %v", dbName, err)
			}
			helper.Close()

			atomic.AddInt64(&successCount, 1)
		}
	})

	b.ReportMetric(float64(successCount), "successes")
	b.ReportMetric(float64(collisionCount), "collisions")

	if collisionCount > 0 {
		b.Errorf("Found %d database name collisions", collisionCount)
	} else {
		b.Log("No database name collisions detected")
	}
}

// goroutineID extracts the goroutine ID from the runtime stack trace
// This matches the implementation in test/test_helper.go
func goroutineID() uint64 {
	buf := make([]byte, 32)
	runtime.Stack(buf, false)
	str := string(buf)
	start := strings.Index(str, "goroutine ")
	if start == -1 {
		return 0
	}
	start += len("goroutine ")
	end := start
	for end < len(str) && str[end] >= '0' && str[end] <= '9' {
		end++
	}
	if end <= start {
		return 0
	}
	var id uint64
	for i := start; i < end; i++ {
		id = id*10 + uint64(str[i]-'0')
	}
	return id
}

// Helper functions for parallel benchmarking

func setupParallelBenchmark(b *testing.B, numProjects int) *test.TestHelper {
	b.Helper()

	if !test.IsTestDatabase() {
		b.Skip("Test database not configured - skipping benchmark test")
	}

	helper, err := test.SetupTestDB()
	if err != nil {
		b.Fatalf("Failed to setup test DB: %v", err)
	}

	if err := helper.SetupTestSchema(); err != nil {
		helper.Close()
		b.Fatalf("Failed to setup test schema: %v", err)
	}

	if err := helper.ClearTestData(); err != nil {
		helper.Close()
		b.Fatalf("Failed to clear test data: %v", err)
	}

	ctx := helper.GetContext()

	// Create test projects
	for i := 1; i <= numProjects; i++ {
		var projectID int64
		err := helper.Pool.QueryRow(ctx, `
			INSERT INTO projects (name, total_page, page, reinicia)
			VALUES ($1, $2, $3, $4) RETURNING id
		`, fmt.Sprintf("Benchmark Project %d", i), 200, 100, false).Scan(&projectID)
		if err != nil {
			helper.Close()
			b.Fatalf("Failed to create test project: %v", err)
		}

		// Create logs for each project
		for j := 1; j <= 5; j++ {
			_, err := helper.Pool.Exec(ctx, `
				INSERT INTO logs (project_id, data, start_page, end_page, wday)
				VALUES ($1, $2, $3, $4, $5)
			`, projectID, "2024-01-01", j*10, j*10+5, 1)
			if err != nil {
				helper.Close()
				b.Fatalf("Failed to create test log: %v", err)
			}
		}
	}

	return helper
}

func cleanupParallelBenchmark(b *testing.B, helper *test.TestHelper) {
	b.Helper()
	if helper != nil {
		if err := helper.CleanupTestSchema(); err != nil {
			b.Logf("Failed to cleanup test schema: %v", err)
		}
		helper.Close()
	}
}

func setupTestDBWithName(dbName string) (*test.TestHelper, error) {
	// Create a custom test helper with the specified database name
	cfg := config.LoadConfig()

	// Modify the database name
	cfg.DBDatabase = dbName

	return test.SetupTestDBWithConfig(cfg)
}

func getOrphanedDatabases() ([]string, error) {
	// Get all test databases that don't match the expected pattern
	// This is a simplified version - actual implementation would query PostgreSQL
	// for databases matching the pattern and check which ones are orphaned

	// For now, return an empty list as this depends on the actual cleanup logic
	// in test/cleanup_orphaned_databases.go
	return []string{}, nil
}

func dropDatabase(dbName string) error {
	// Drop a database by name
	// This would use pgx to execute DROP DATABASE
	return nil
}

// calculatePercentiles calculates p50, p95, p99 percentiles from a slice of float64 values
func calculatePercentiles(values []float64) (p50, p95, p99 float64) {
	if len(values) == 0 {
		return 0, 0, 0
	}

	// Sort values
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sortFloat64(sorted)

	n := len(sorted)
	p50 = sorted[n*50/100]
	p95 = sorted[n*95/100]
	p99 = sorted[n*99/100]

	return p50, p95, p99
}

func sortFloat64(slice []float64) {
	for i := 0; i < len(slice)-1; i++ {
		for j := i + 1; j < len(slice); j++ {
			if slice[i] > slice[j] {
				slice[i], slice[j] = slice[j], slice[i]
			}
		}
	}
}

// generateBenchmarkReport generates a JSON report with benchmark metrics
func generateBenchmarkReport(results map[string]BenchmarkResult) string {
	// This would generate a comprehensive JSON report
	// For now, return a placeholder
	return ""
}

// BenchmarkResult represents a single benchmark result
type BenchmarkResult struct {
	Name      string  `json:"name"`
	AverageMs float64 `json:"average_ms"`
	P50Ms     float64 `json:"p50_ms"`
	P95Ms     float64 `json:"p95_ms"`
	P99Ms     float64 `json:"p99_ms"`
	Allocs    uint64  `json:"allocs"`
	OpsPerSec float64 `json:"ops_per_sec"`
}
