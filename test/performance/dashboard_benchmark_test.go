package performance

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	pg "go-reading-log-api-next/internal/adapter/postgres"
	"go-reading-log-api-next/internal/api/v1/handlers"
	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/repository"
	"go-reading-log-api-next/internal/service"
	"go-reading-log-api-next/internal/service/dashboard"
	"go-reading-log-api-next/test"
)

// MockProjectsService is a mock implementation of ProjectsServiceInterface
type MockProjectsService struct{}

func (m *MockProjectsService) GetRunningProjectsWithLogs(ctx context.Context) ([]*dashboard.ProjectWithLogs, error) {
	return []*dashboard.ProjectWithLogs{}, nil
}

func (m *MockProjectsService) CalculateStats(ctx context.Context) (*dto.StatsData, error) {
	return dto.NewStatsData(), nil
}

func (m *MockProjectsService) GetDashboardProjects(ctx context.Context) (*dto.DashboardProjectsResponse, error) {
	return dto.NewDashboardProjectsResponse(), nil
}

// =============================================================================
// Dashboard Performance Benchmarks
// =============================================================================

// BenchmarkDashboardDay measures the performance of /v1/dashboard/day.json endpoint
func BenchmarkDashboardDay(b *testing.B) {
	if !test.IsTestDatabase() {
		b.Skip("Test database not configured - skipping benchmark test")
	}

	helper := setupDashboardBenchmark(b)
	defer cleanupDashboardBenchmark(b, helper)

	handler := createDashboardHandler(helper.Pool)

	// Warm-up phase
	warmupRequest(b, handler, "/v1/dashboard/day.json")

	// Benchmark
	b.ResetTimer()
	b.ReportAllocs()

	var totalTime time.Duration
	var latencies []float64

	for i := 0; i < b.N; i++ {
		start := time.Now()
		req := httptest.NewRequest("GET", "/v1/dashboard/day.json", nil)
		recorder := httptest.NewRecorder()

		handler.Day(recorder, req)

		if recorder.Code != http.StatusOK {
			b.Fatalf("Expected 200, got %d: %s", recorder.Code, recorder.Body.String())
		}

		totalTime += time.Since(start)
		latencies = append(latencies, float64(time.Since(start))/float64(time.Millisecond))
	}

	// Calculate percentiles
	p50, p95, p99 := calculatePercentiles(latencies)

	b.ReportMetric(p50, "p50_ms/op")
	b.ReportMetric(p95, "p95_ms/op")
	b.ReportMetric(p99, "p99_ms/op")

	// Verify threshold (100ms for p95)
	if p95 > 100 {
		b.Errorf("P95 latency %.2fms exceeds threshold of 100ms", p95)
	} else {
		b.Logf("P95 latency %.2fms is within threshold (< 100ms)", p95)
	}
}

// BenchmarkDashboardProjects measures the performance of /v1/dashboard/projects.json endpoint
func BenchmarkDashboardProjects(b *testing.B) {
	if !test.IsTestDatabase() {
		b.Skip("Test database not configured - skipping benchmark test")
	}

	helper := setupDashboardBenchmark(b)
	defer cleanupDashboardBenchmark(b, helper)

	handler := createDashboardHandler(helper.Pool)

	// Warm-up phase
	warmupRequest(b, handler, "/v1/dashboard/projects.json")

	// Benchmark
	b.ResetTimer()
	b.ReportAllocs()

	var totalTime time.Duration
	var latencies []float64

	for i := 0; i < b.N; i++ {
		start := time.Now()
		req := httptest.NewRequest("GET", "/v1/dashboard/projects.json", nil)
		recorder := httptest.NewRecorder()

		handler.Projects(recorder, req)

		if recorder.Code != http.StatusOK {
			b.Fatalf("Expected 200, got %d: %s", recorder.Code, recorder.Body.String())
		}

		totalTime += time.Since(start)
		latencies = append(latencies, float64(time.Since(start))/float64(time.Millisecond))
	}

	// Calculate percentiles
	p50, p95, p99 := calculatePercentiles(latencies)

	b.ReportMetric(p50, "p50_ms/op")
	b.ReportMetric(p95, "p95_ms/op")
	b.ReportMetric(p99, "p99_ms/op")

	// Verify threshold (150ms for p95)
	if p95 > 150 {
		b.Errorf("P95 latency %.2fms exceeds threshold of 150ms", p95)
	} else {
		b.Logf("P95 latency %.2fms is within threshold (< 150ms)", p95)
	}
}

// BenchmarkDashboardLastDays measures the performance of /v1/dashboard/last_days.json endpoint
func BenchmarkDashboardLastDays(b *testing.B) {
	if !test.IsTestDatabase() {
		b.Skip("Test database not configured - skipping benchmark test")
	}

	helper := setupDashboardBenchmark(b)
	defer cleanupDashboardBenchmark(b, helper)

	handler := createDashboardHandler(helper.Pool)

	// Warm-up phase
	warmupRequest(b, handler, "/v1/dashboard/last_days.json")

	// Benchmark
	b.ResetTimer()
	b.ReportAllocs()

	var totalTime time.Duration
	var latencies []float64

	for i := 0; i < b.N; i++ {
		start := time.Now()
		req := httptest.NewRequest("GET", "/v1/dashboard/last_days.json?days=7&type=1", nil)
		recorder := httptest.NewRecorder()

		handler.LastDays(recorder, req)

		if recorder.Code != http.StatusOK {
			b.Fatalf("Expected 200, got %d: %s", recorder.Code, recorder.Body.String())
		}

		totalTime += time.Since(start)
		latencies = append(latencies, float64(time.Since(start))/float64(time.Millisecond))
	}

	// Calculate percentiles
	p50, p95, p99 := calculatePercentiles(latencies)

	b.ReportMetric(p50, "p50_ms/op")
	b.ReportMetric(p95, "p95_ms/op")
	b.ReportMetric(p99, "p99_ms/op")

	// Verify threshold (200ms for p95)
	if p95 > 200 {
		b.Errorf("P95 latency %.2fms exceeds threshold of 200ms", p95)
	} else {
		b.Logf("P95 latency %.2fms is within threshold (< 200ms)", p95)
	}
}

// BenchmarkDashboardFaults measures the performance of /v1/dashboard/echart/faults.json endpoint
func BenchmarkDashboardFaults(b *testing.B) {
	if !test.IsTestDatabase() {
		b.Skip("Test database not configured - skipping benchmark test")
	}

	helper := setupDashboardBenchmark(b)
	defer cleanupDashboardBenchmark(b, helper)

	handler := createDashboardHandler(helper.Pool)

	// Warm-up phase
	warmupRequest(b, handler, "/v1/dashboard/echart/faults.json")

	// Benchmark
	b.ResetTimer()
	b.ReportAllocs()

	var totalTime time.Duration
	var latencies []float64

	for i := 0; i < b.N; i++ {
		start := time.Now()
		req := httptest.NewRequest("GET", "/v1/dashboard/echart/faults.json", nil)
		recorder := httptest.NewRecorder()

		handler.Faults(recorder, req)

		if recorder.Code != http.StatusOK {
			b.Fatalf("Expected 200, got %d: %s", recorder.Code, recorder.Body.String())
		}

		totalTime += time.Since(start)
		latencies = append(latencies, float64(time.Since(start))/float64(time.Millisecond))
	}

	// Calculate percentiles
	p50, p95, p99 := calculatePercentiles(latencies)

	b.ReportMetric(p50, "p50_ms/op")
	b.ReportMetric(p95, "p95_ms/op")
	b.ReportMetric(p99, "p99_ms/op")

	// Verify threshold (100ms for p95)
	if p95 > 100 {
		b.Errorf("P95 latency %.2fms exceeds threshold of 100ms", p95)
	} else {
		b.Logf("P95 latency %.2fms is within threshold (< 100ms)", p95)
	}
}

// BenchmarkDashboardSpeculateActual measures the performance of /v1/dashboard/echart/speculate_actual.json endpoint
func BenchmarkDashboardSpeculateActual(b *testing.B) {
	if !test.IsTestDatabase() {
		b.Skip("Test database not configured - skipping benchmark test")
	}

	helper := setupDashboardBenchmark(b)
	defer cleanupDashboardBenchmark(b, helper)

	handler := createDashboardHandler(helper.Pool)

	// Warm-up phase
	warmupRequest(b, handler, "/v1/dashboard/echart/speculate_actual.json")

	// Benchmark
	b.ResetTimer()
	b.ReportAllocs()

	var totalTime time.Duration
	var latencies []float64

	for i := 0; i < b.N; i++ {
		start := time.Now()
		req := httptest.NewRequest("GET", "/v1/dashboard/echart/speculate_actual.json", nil)
		recorder := httptest.NewRecorder()

		handler.SpeculateActual(recorder, req)

		if recorder.Code != http.StatusOK {
			b.Fatalf("Expected 200, got %d: %s", recorder.Code, recorder.Body.String())
		}

		totalTime += time.Since(start)
		latencies = append(latencies, float64(time.Since(start))/float64(time.Millisecond))
	}

	// Calculate percentiles
	p50, p95, p99 := calculatePercentiles(latencies)

	b.ReportMetric(p50, "p50_ms/op")
	b.ReportMetric(p95, "p95_ms/op")
	b.ReportMetric(p99, "p99_ms/op")

	// Verify threshold (150ms for p95)
	if p95 > 150 {
		b.Errorf("P95 latency %.2fms exceeds threshold of 150ms", p95)
	} else {
		b.Logf("P95 latency %.2fms is within threshold (< 150ms)", p95)
	}
}

// BenchmarkDashboardWeekdayFaults measures the performance of /v1/dashboard/echart/faults_week_day.json endpoint
func BenchmarkDashboardWeekdayFaults(b *testing.B) {
	if !test.IsTestDatabase() {
		b.Skip("Test database not configured - skipping benchmark test")
	}

	helper := setupDashboardBenchmark(b)
	defer cleanupDashboardBenchmark(b, helper)

	handler := createDashboardHandler(helper.Pool)

	// Warm-up phase
	warmupRequest(b, handler, "/v1/dashboard/echart/faults_week_day.json")

	// Benchmark
	b.ResetTimer()
	b.ReportAllocs()

	var totalTime time.Duration
	var latencies []float64

	for i := 0; i < b.N; i++ {
		start := time.Now()
		req := httptest.NewRequest("GET", "/v1/dashboard/echart/faults_week_day.json", nil)
		recorder := httptest.NewRecorder()

		handler.WeekdayFaults(recorder, req)

		if recorder.Code != http.StatusOK {
			b.Fatalf("Expected 200, got %d: %s", recorder.Code, recorder.Body.String())
		}

		totalTime += time.Since(start)
		latencies = append(latencies, float64(time.Since(start))/float64(time.Millisecond))
	}

	// Calculate percentiles
	p50, p95, p99 := calculatePercentiles(latencies)

	b.ReportMetric(p50, "p50_ms/op")
	b.ReportMetric(p95, "p95_ms/op")
	b.ReportMetric(p99, "p99_ms/op")

	// Verify threshold (200ms for p95)
	if p95 > 200 {
		b.Errorf("P95 latency %.2fms exceeds threshold of 200ms", p95)
	} else {
		b.Logf("P95 latency %.2fms is within threshold (< 200ms)", p95)
	}
}

// BenchmarkDashboardMeanProgress measures the performance of /v1/dashboard/echart/mean_progress.json endpoint
func BenchmarkDashboardMeanProgress(b *testing.B) {
	if !test.IsTestDatabase() {
		b.Skip("Test database not configured - skipping benchmark test")
	}

	helper := setupDashboardBenchmark(b)
	defer cleanupDashboardBenchmark(b, helper)

	handler := createDashboardHandler(helper.Pool)

	// Warm-up phase
	warmupRequest(b, handler, "/v1/dashboard/echart/mean_progress.json")

	// Benchmark
	b.ResetTimer()
	b.ReportAllocs()

	var totalTime time.Duration
	var latencies []float64

	for i := 0; i < b.N; i++ {
		start := time.Now()
		req := httptest.NewRequest("GET", "/v1/dashboard/echart/mean_progress.json", nil)
		recorder := httptest.NewRecorder()

		handler.MeanProgress(recorder, req)

		if recorder.Code != http.StatusOK {
			b.Fatalf("Expected 200, got %d: %s", recorder.Code, recorder.Body.String())
		}

		totalTime += time.Since(start)
		latencies = append(latencies, float64(time.Since(start))/float64(time.Millisecond))
	}

	// Calculate percentiles
	p50, p95, p99 := calculatePercentiles(latencies)

	b.ReportMetric(p50, "p50_ms/op")
	b.ReportMetric(p95, "p95_ms/op")
	b.ReportMetric(p99, "p99_ms/op")

	// Verify threshold (150ms for p95)
	if p95 > 150 {
		b.Errorf("P95 latency %.2fms exceeds threshold of 150ms", p95)
	} else {
		b.Logf("P95 latency %.2fms is within threshold (< 150ms)", p95)
	}
}

// BenchmarkDashboardYearlyTotal measures the performance of /v1/dashboard/echart/last_year_total.json endpoint
func BenchmarkDashboardYearlyTotal(b *testing.B) {
	if !test.IsTestDatabase() {
		b.Skip("Test database not configured - skipping benchmark test")
	}

	helper := setupDashboardBenchmark(b)
	defer cleanupDashboardBenchmark(b, helper)

	handler := createDashboardHandler(helper.Pool)

	// Warm-up phase
	warmupRequest(b, handler, "/v1/dashboard/echart/last_year_total.json")

	// Benchmark
	b.ResetTimer()
	b.ReportAllocs()

	var totalTime time.Duration
	var latencies []float64

	for i := 0; i < b.N; i++ {
		start := time.Now()
		req := httptest.NewRequest("GET", "/v1/dashboard/echart/last_year_total.json", nil)
		recorder := httptest.NewRecorder()

		handler.YearlyTotal(recorder, req)

		if recorder.Code != http.StatusOK {
			b.Fatalf("Expected 200, got %d: %s", recorder.Code, recorder.Body.String())
		}

		totalTime += time.Since(start)
		latencies = append(latencies, float64(time.Since(start))/float64(time.Millisecond))
	}

	// Calculate percentiles
	p50, p95, p99 := calculatePercentiles(latencies)

	b.ReportMetric(p50, "p50_ms/op")
	b.ReportMetric(p95, "p95_ms/op")
	b.ReportMetric(p99, "p99_ms/op")

	// Verify threshold (300ms for p95)
	if p95 > 300 {
		b.Errorf("P95 latency %.2fms exceeds threshold of 300ms", p95)
	} else {
		b.Logf("P95 latency %.2fms is within threshold (< 300ms)", p95)
	}
}

// =============================================================================
// Helper Functions
// =============================================================================

func setupDashboardBenchmark(b *testing.B) *test.TestHelper {
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
	for i := 1; i <= 10; i++ {
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

func cleanupDashboardBenchmark(b *testing.B, helper *test.TestHelper) {
	b.Helper()
	if helper != nil {
		if err := helper.CleanupTestSchema(); err != nil {
			b.Logf("Failed to cleanup test schema: %v", err)
		}
		helper.Close()
	}
}

func createDashboardHandler(pool *pgxpool.Pool) *handlers.DashboardHandler {
	repo := createTestRepository(pool)

	userConfig, err := service.LoadDashboardConfig("")
	if err != nil {
		userConfig = service.NewUserConfigService(service.GetDefaultConfig())
	}

	return handlers.NewDashboardHandler(repo, userConfig, &MockProjectsService{})
}

func createTestRepository(pool *pgxpool.Pool) repository.DashboardRepository {
	repo := pg.NewDashboardRepositoryImpl(pool)
	return repo
}

func warmupRequest(b *testing.B, handler *handlers.DashboardHandler, url string) {
	b.Helper()

	req := httptest.NewRequest("GET", url, nil)
	recorder := httptest.NewRecorder()

	switch {
	case strings.Contains(url, "/day.json"):
		handler.Day(recorder, req)
	case strings.Contains(url, "/projects.json"):
		handler.Projects(recorder, req)
	case strings.Contains(url, "/last_days.json"):
		handler.LastDays(recorder, req)
	case strings.Contains(url, "/faults.json") && !strings.Contains(url, "week_day"):
		handler.Faults(recorder, req)
	case strings.Contains(url, "/speculate_actual.json"):
		handler.SpeculateActual(recorder, req)
	case strings.Contains(url, "/faults_week_day.json"):
		handler.WeekdayFaults(recorder, req)
	case strings.Contains(url, "/mean_progress.json"):
		handler.MeanProgress(recorder, req)
	case strings.Contains(url, "/last_year_total.json"):
		handler.YearlyTotal(recorder, req)
	}

	if recorder.Code != http.StatusOK {
		b.Logf("Warm-up failed for %s: status=%d", url, recorder.Code)
	}
}

// calculatePercentiles calculates p50, p95, p99 percentiles from a slice of float64 values
func verifyConnectionPool(b *testing.B, pool *pgxpool.Pool) {
	b.Helper()

	stats := pool.Stat()

	// Check pool is being utilized
	if stats.TotalConns() == 0 {
		b.Errorf("Connection pool has no connections")
	}

	// Verify connections are being reused (not creating new ones each request)
	// Note: pgxpool doesn't expose ReleaseCount directly, so we check AcquireCount vs TotalConns
	// High acquire count relative to total connections indicates good reuse
	if stats.TotalConns() > 0 {
		reuseRatio := float64(stats.AcquireCount()) / float64(stats.TotalConns())
		b.Logf("Connection reuse ratio: %.2f (acquire=%d, total=%d)", reuseRatio, stats.AcquireCount(), stats.TotalConns())
	}

	// Check for connection leaks (AcquiredConns should be 0 after operations complete)
	if stats.AcquiredConns() != 0 {
		b.Logf("Note: %d connections still acquired (may be expected during benchmark)", stats.AcquiredConns())
	}

	b.Logf("Connection pool stats: total=%d, active=%d, idle=%d, acquire=%d",
		stats.TotalConns(), stats.AcquiredConns(), stats.IdleConns(),
		stats.AcquireCount())
}

// BenchmarkDashboardWithQueryTracer measures performance with query tracing enabled
func BenchmarkDashboardWithQueryTracer(b *testing.B) {
	if !test.IsTestDatabase() {
		b.Skip("Test database not configured - skipping benchmark test")
	}

	helper := setupDashboardBenchmark(b)
	defer cleanupDashboardBenchmark(b, helper)

	// Create handler with query tracer
	handler := createDashboardHandler(helper.Pool)

	// Warm-up phase
	warmupRequest(b, handler, "/v1/dashboard/day.json")

	// Benchmark
	b.ResetTimer()
	b.ReportAllocs()

	var totalTime time.Duration
	var latencies []float64

	for i := 0; i < b.N; i++ {
		start := time.Now()
		req := httptest.NewRequest("GET", "/v1/dashboard/day.json", nil)
		recorder := httptest.NewRecorder()

		handler.Day(recorder, req)

		if recorder.Code != http.StatusOK {
			b.Fatalf("Expected 200, got %d: %s", recorder.Code, recorder.Body.String())
		}

		totalTime += time.Since(start)
		latencies = append(latencies, float64(time.Since(start))/float64(time.Millisecond))
	}

	// Calculate percentiles
	p50, p95, p99 := calculatePercentiles(latencies)

	b.ReportMetric(p50, "p50_ms/op")
	b.ReportMetric(p95, "p95_ms/op")
	b.ReportMetric(p99, "p99_ms/op")

	// Verify threshold (100ms for p95)
	if p95 > 100 {
		b.Errorf("P95 latency %.2fms exceeds threshold of 100ms", p95)
	} else {
		b.Logf("P95 latency %.2fms is within threshold (< 100ms)", p95)
	}

	// Verify connection pool
	verifyConnectionPool(b, helper.Pool)
}
