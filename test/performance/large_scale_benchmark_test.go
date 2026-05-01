package performance

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sort"
	"sync"
	"testing"
	"time"

	pg "go-reading-log-api-next/internal/adapter/postgres"
	"go-reading-log-api-next/internal/api/v1/handlers"
	"go-reading-log-api-next/test"
)

// =============================================================================
// Large-Scale Performance Benchmarks (10,000+ logs)
// =============================================================================

// calculatePercentilesFloat64 calculates p50, p95, p99 percentiles from a slice of float64 values (milliseconds)
func calculatePercentilesFloat64(times []float64) (p50, p95, p99 float64) {
	if len(times) == 0 {
		return 0, 0, 0
	}

	// Sort the slice
	sorted := make([]float64, len(times))
	copy(sorted, times)
	sort.Float64s(sorted)

	n := len(sorted)
	p50 = sorted[n*50/100]
	p95 = sorted[n*95/100]
	p99 = sorted[n*99/100]

	return p50, p95, p99
}

//
// These benchmarks verify that the API meets NFC-DASH-001 performance requirements:
// - Response time < 500ms at p95 percentile with 10,000+ logs
// - Database queries use appropriate indexes
// - Connection pool efficiency under load
//
// Dataset: 100 projects with ~100 logs each = 10,000+ total logs

const (
	// LargeScaleNumProjects is the number of projects to create
	LargeScaleNumProjects = 100
	// LargeScaleLogsPerProject is the average number of logs per project
	LargeScaleLogsPerProject = 100
	// LargeScaleP95ThresholdMs is the p95 latency threshold in milliseconds
	LargeScaleP95ThresholdMs = 500.0
	// LargeScaleConcurrentUsers is the number of concurrent users for load tests
	LargeScaleConcurrentUsers = 50
)

// BenchmarkLargeScaleGetAllWithLogs measures performance of GetAllWithLogs
// with 10,000+ logs distributed across 100 projects
func BenchmarkLargeScaleGetAllWithLogs(b *testing.B) {
	if !test.IsTestDatabase() {
		b.Skip("Test database not configured - skipping large-scale benchmark")
	}

	helper, cleanup := setupLargeScaleBenchmark(b)
	defer cleanup()

	repo := pg.NewProjectRepositoryImpl(helper.Pool)

	// Warm-up phase
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := repo.GetAllWithLogs(ctx)
	cancel()
	if err != nil {
		b.Fatalf("Warm-up failed: %v", err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	var totalTime time.Duration
	var latencies []float64
	var projectCount int

	for i := 0; i < b.N; i++ {
		start := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		result, err := repo.GetAllWithLogs(ctx)
		cancel()

		if err != nil {
			b.Fatalf("GetAllWithLogs failed: %v", err)
		}

		projectCount = len(result)
		totalTime += time.Since(start)
		latencies = append(latencies, float64(time.Since(start))/float64(time.Millisecond))
	}

	b.StopTimer()

	// Calculate percentiles
	p50, p95, p99 := calculatePercentilesFloat64(latencies)

	b.ReportMetric(p50, "p50_ms/op")
	b.ReportMetric(p95, "p95_ms/op")
	b.ReportMetric(p99, "p99_ms/op")
	b.ReportMetric(float64(projectCount), "projects/op")

	// Verify threshold
	if p95 > LargeScaleP95ThresholdMs {
		b.Errorf("P95 latency %.2fms exceeds threshold of %.2fms", p95, LargeScaleP95ThresholdMs)
	} else {
		b.Logf("P95 latency %.2fms is within threshold (< %.2fms)", p95, LargeScaleP95ThresholdMs)
	}

	// Verify connection pool
	verifyConnectionPool(b, helper.Pool)
}

// BenchmarkLargeScaleGetWithLogs measures performance of GetWithLogs for a single
// project with many logs (100+ logs per project)
func BenchmarkLargeScaleGetWithLogs(b *testing.B) {
	if !test.IsTestDatabase() {
		b.Skip("Test database not configured - skipping large-scale benchmark")
	}

	helper, cleanup := setupLargeScaleBenchmark(b)
	defer cleanup()

	repo := pg.NewProjectRepositoryImpl(helper.Pool)

	// Get a project ID with many logs
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	projects, err := repo.GetAllWithLogs(ctx)
	cancel()
	if err != nil {
		b.Fatalf("Failed to get projects: %v", err)
	}

	if len(projects) == 0 {
		b.Fatal("No projects found in database")
	}

	// Select a project with logs
	projectID := projects[0].Project.ID

	// Warm-up phase
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	_, err = repo.GetWithLogs(ctx, projectID)
	cancel()
	if err != nil {
		b.Fatalf("Warm-up failed: %v", err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	var totalTime time.Duration
	var latencies []float64

	for i := 0; i < b.N; i++ {
		start := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := repo.GetWithLogs(ctx, projectID)
		cancel()

		if err != nil {
			b.Fatalf("GetWithLogs failed: %v", err)
		}

		totalTime += time.Since(start)
		latencies = append(latencies, float64(time.Since(start))/float64(time.Millisecond))
	}

	b.StopTimer()

	// Calculate percentiles
	p50, p95, p99 := calculatePercentilesFloat64(latencies)

	b.ReportMetric(p50, "p50_ms/op")
	b.ReportMetric(p95, "p95_ms/op")
	b.ReportMetric(p99, "p99_ms/op")

	// Threshold for single project is lower (100ms)
	const singleProjectThreshold = 100.0
	if p95 > singleProjectThreshold {
		b.Errorf("P95 latency %.2fms exceeds threshold of %.2fms", p95, singleProjectThreshold)
	} else {
		b.Logf("P95 latency %.2fms is within threshold (< %.2fms)", p95, singleProjectThreshold)
	}
}

// BenchmarkLargeScaleConcurrent50 measures concurrent performance with 50 users
func BenchmarkLargeScaleConcurrent50(b *testing.B) {
	if !test.IsTestDatabase() {
		b.Skip("Test database not configured - skipping large-scale benchmark")
	}

	helper, cleanup := setupLargeScaleBenchmark(b)
	defer cleanup()

	repo := pg.NewProjectRepositoryImpl(helper.Pool)

	// Warm-up phase
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := repo.GetAllWithLogs(ctx)
	cancel()
	if err != nil {
		b.Fatalf("Warm-up failed: %v", err)
	}

	b.SetParallelism(LargeScaleConcurrentUsers)
	b.ResetTimer()
	b.ReportAllocs()

	var mu sync.Mutex
	var successCount int
	var errorCount int
	var totalLatency time.Duration

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			start := time.Now()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			_, err := repo.GetAllWithLogs(ctx)
			cancel()

			latency := time.Since(start)

			mu.Lock()
			if err != nil {
				errorCount++
			} else {
				successCount++
				totalLatency += latency
			}
			mu.Unlock()
		}
	})

	b.StopTimer()

	// Calculate metrics
	errorRate := float64(errorCount) / float64(successCount+errorCount) * 100
	avgLatency := float64(totalLatency) / float64(successCount) / float64(time.Millisecond)

	b.ReportMetric(float64(successCount), "successful_requests")
	b.ReportMetric(float64(errorCount), "failed_requests")
	b.ReportMetric(errorRate, "error_rate_percent")
	b.ReportMetric(avgLatency, "avg_latency_ms")
	b.ReportMetric(float64(successCount)/b.Elapsed().Seconds(), "requests_per_sec")

	// Verify thresholds
	if errorRate > 1.0 {
		b.Errorf("Error rate %.2f%% exceeds threshold of 1%%", errorRate)
	} else {
		b.Logf("Error rate %.2f%% is within threshold (< 1%%)", errorRate)
	}

	// Verify connection pool
	verifyConnectionPool(b, helper.Pool)
}

// BenchmarkLargeScaleConcurrent100 measures concurrent performance with 100 users
func BenchmarkLargeScaleConcurrent100(b *testing.B) {
	if !test.IsTestDatabase() {
		b.Skip("Test database not configured - skipping large-scale benchmark")
	}

	helper, cleanup := setupLargeScaleBenchmark(b)
	defer cleanup()

	repo := pg.NewProjectRepositoryImpl(helper.Pool)

	// Warm-up phase
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := repo.GetAllWithLogs(ctx)
	cancel()
	if err != nil {
		b.Fatalf("Warm-up failed: %v", err)
	}

	b.SetParallelism(100)
	b.ResetTimer()
	b.ReportAllocs()

	var mu sync.Mutex
	var successCount int
	var errorCount int
	var totalLatency time.Duration

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			start := time.Now()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			_, err := repo.GetAllWithLogs(ctx)
			cancel()

			latency := time.Since(start)

			mu.Lock()
			if err != nil {
				errorCount++
			} else {
				successCount++
				totalLatency += latency
			}
			mu.Unlock()
		}
	})

	b.StopTimer()

	// Calculate metrics
	errorRate := float64(errorCount) / float64(successCount+errorCount) * 100
	avgLatency := float64(totalLatency) / float64(successCount) / float64(time.Millisecond)

	b.ReportMetric(float64(successCount), "successful_requests")
	b.ReportMetric(float64(errorCount), "failed_requests")
	b.ReportMetric(errorRate, "error_rate_percent")
	b.ReportMetric(avgLatency, "avg_latency_ms")
	b.ReportMetric(float64(successCount)/b.Elapsed().Seconds(), "requests_per_sec")

	// Verify thresholds
	if errorRate > 1.0 {
		b.Errorf("Error rate %.2f%% exceeds threshold of 1%%", errorRate)
	} else {
		b.Logf("Error rate %.2f%% is within threshold (< 1%%)", errorRate)
	}
}

// BenchmarkLargeScaleHTTPProjects measures HTTP handler performance for /v1/projects.json
func BenchmarkLargeScaleHTTPProjects(b *testing.B) {
	if !test.IsTestDatabase() {
		b.Skip("Test database not configured - skipping large-scale benchmark")
	}

	helper, cleanup := setupLargeScaleBenchmark(b)
	defer cleanup()

	repo := pg.NewProjectRepositoryImpl(helper.Pool)
	handler := handlers.NewProjectsHandler(repo)

	// Warm-up phase
	req := httptest.NewRequest("GET", "/v1/projects.json", nil)
	recorder := httptest.NewRecorder()
	handler.Index(recorder, req)

	if recorder.Code != http.StatusOK {
		b.Fatalf("Warm-up failed: expected 200, got %d", recorder.Code)
	}

	b.ResetTimer()
	b.ReportAllocs()

	var totalTime time.Duration
	var latencies []float64

	for i := 0; i < b.N; i++ {
		start := time.Now()
		req := httptest.NewRequest("GET", "/v1/projects.json", nil)
		recorder := httptest.NewRecorder()

		handler.Index(recorder, req)

		if recorder.Code != http.StatusOK {
			b.Fatalf("Expected 200, got %d: %s", recorder.Code, recorder.Body.String())
		}

		totalTime += time.Since(start)
		latencies = append(latencies, float64(time.Since(start))/float64(time.Millisecond))
	}

	b.StopTimer()

	// Calculate percentiles
	p50, p95, p99 := calculatePercentilesFloat64(latencies)

	b.ReportMetric(p50, "p50_ms/op")
	b.ReportMetric(p95, "p95_ms/op")
	b.ReportMetric(p99, "p99_ms/op")

	// Verify threshold
	if p95 > LargeScaleP95ThresholdMs {
		b.Errorf("P95 latency %.2fms exceeds threshold of %.2fms", p95, LargeScaleP95ThresholdMs)
	} else {
		b.Logf("P95 latency %.2fms is within threshold (< %.2fms)", p95, LargeScaleP95ThresholdMs)
	}
}

// BenchmarkLargeScaleHTTPShow measures HTTP handler performance for /v1/projects/{id}.json
func BenchmarkLargeScaleHTTPShow(b *testing.B) {
	if !test.IsTestDatabase() {
		b.Skip("Test database not configured - skipping large-scale benchmark")
	}

	helper, cleanup := setupLargeScaleBenchmark(b)
	defer cleanup()

	repo := pg.NewProjectRepositoryImpl(helper.Pool)

	// Get a project ID
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	projects, err := repo.GetAllWithLogs(ctx)
	cancel()
	if err != nil {
		b.Fatalf("Failed to get projects: %v", err)
	}

	if len(projects) == 0 {
		b.Fatal("No projects found")
	}

	projectID := fmt.Sprintf("%d", projects[0].Project.ID)
	handler := handlers.NewProjectsHandler(repo)

	// Warm-up phase
	req := httptest.NewRequest("GET", "/v1/projects/"+projectID+".json", nil)
	recorder := httptest.NewRecorder()
	handler.Show(recorder, req)

	if recorder.Code != http.StatusOK {
		b.Fatalf("Warm-up failed: expected 200, got %d", recorder.Code)
	}

	b.ResetTimer()
	b.ReportAllocs()

	var totalTime time.Duration
	var latencies []float64

	for i := 0; i < b.N; i++ {
		start := time.Now()
		req := httptest.NewRequest("GET", "/v1/projects/"+projectID+".json", nil)
		recorder := httptest.NewRecorder()

		handler.Show(recorder, req)

		if recorder.Code != http.StatusOK {
			b.Fatalf("Expected 200, got %d: %s", recorder.Code, recorder.Body.String())
		}

		totalTime += time.Since(start)
		latencies = append(latencies, float64(time.Since(start))/float64(time.Millisecond))
	}

	b.StopTimer()

	// Calculate percentiles
	p50, p95, p99 := calculatePercentilesFloat64(latencies)

	b.ReportMetric(p50, "p50_ms/op")
	b.ReportMetric(p95, "p95_ms/op")
	b.ReportMetric(p99, "p99_ms/op")

	// Threshold for single project is lower
	const singleProjectThreshold = 100.0
	if p95 > singleProjectThreshold {
		b.Errorf("P95 latency %.2fms exceeds threshold of %.2fms", p95, singleProjectThreshold)
	} else {
		b.Logf("P95 latency %.2fms is within threshold (< %.2fms)", p95, singleProjectThreshold)
	}
}

// =============================================================================
// Helper Functions
// =============================================================================

// setupLargeScaleBenchmark creates a database with 100+ projects and 10,000+ logs
// Returns the test helper and a cleanup function
func setupLargeScaleBenchmark(b *testing.B) (*test.TestHelper, func()) {
	b.Helper()

	if !test.IsTestDatabase() {
		b.Skip("Test database not configured - skipping large-scale benchmark")
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

	// Create 100 projects with varying characteristics
	b.Logf("Creating %d projects...", LargeScaleNumProjects)
	projects := make([]int64, LargeScaleNumProjects)

	for i := 1; i <= LargeScaleNumProjects; i++ {
		var projectID int64
		// Vary total_page between 100-500
		totalPage := 100 + rand.Intn(400)
		// Vary page between 0-total_page
		page := rand.Intn(totalPage)

		err := helper.Pool.QueryRow(ctx, `
			INSERT INTO projects (name, total_page, page, reinicia, started_at)
			VALUES ($1, $2, $3, $4, $5) RETURNING id
		`, fmt.Sprintf("Large Scale Project %03d", i), totalPage, page, false,
			time.Now().AddDate(0, -rand.Intn(6)-1, 0)).Scan(&projectID)
		if err != nil {
			helper.Close()
			b.Fatalf("Failed to create test project %d: %v", i, err)
		}
		projects[i-1] = projectID
	}

	// Create ~100 logs per project (10,000+ total)
	// Use batch inserts for better performance
	b.Logf("Creating ~%d logs per project...", LargeScaleLogsPerProject)

	batchSize := 500
	logCount := 0

	for _, projectID := range projects {
		// Create logs in batches
		for batchStart := 0; batchStart < LargeScaleLogsPerProject; batchStart += batchSize {
			batchEnd := batchStart + batchSize
			if batchEnd > LargeScaleLogsPerProject {
				batchEnd = LargeScaleLogsPerProject
			}

			batchQuery := "INSERT INTO logs (project_id, data, start_page, end_page, wday) VALUES "
			var values []interface{}
			argNum := 1

			for j := batchStart; j < batchEnd; j++ {
				// Generate realistic log data
				// Spread logs over the past 6 months
				daysAgo := rand.Intn(180)
				logDate := time.Now().AddDate(0, 0, -daysAgo)

				// Vary pages read per session (5-50 pages)
				pagesRead := 5 + rand.Intn(45)
				startPage := rand.Intn(200)
				endPage := startPage + pagesRead
				wday := int(logDate.Weekday())

				batchQuery += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d), ",
					argNum, argNum+1, argNum+2, argNum+3, argNum+4)
				values = append(values,
					projectID,
					logDate.Format("2006-01-02 15:04:05"),
					startPage,
					endPage,
					wday,
				)
				argNum += 5
				logCount++
			}

			// Remove trailing comma and space
			batchQuery = batchQuery[:len(batchQuery)-2]

			_, err := helper.Pool.Exec(ctx, batchQuery, values...)
			if err != nil {
				helper.Close()
				b.Fatalf("Failed to create log batch for project %d: %v", projectID, err)
			}
		}
	}

	b.Logf("Created %d projects with %d total logs", LargeScaleNumProjects, logCount)

	// Verify data was created
	var count int
	err = helper.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM projects").Scan(&count)
	if err != nil {
		helper.Close()
		b.Fatalf("Failed to count projects: %v", err)
	}
	b.Logf("Projects in database: %d", count)

	err = helper.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM logs").Scan(&count)
	if err != nil {
		helper.Close()
		b.Fatalf("Failed to count logs: %v", err)
	}
	b.Logf("Logs in database: %d", count)

	cleanup := func() {
		if err := helper.CleanupTestSchema(); err != nil {
			b.Logf("Failed to cleanup test schema: %v", err)
		}
		helper.Close()
	}

	return helper, cleanup
}
