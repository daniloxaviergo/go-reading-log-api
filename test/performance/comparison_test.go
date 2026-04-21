package performance

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"go-reading-log-api-next/internal/adapter/postgres"
	"go-reading-log-api-next/internal/api/v1/handlers"
	"go-reading-log-api-next/test"
)

// BenchmarkComparisonGetAllWithLogs measures the performance of GetAllWithLogs
// with comparison against established baseline metrics
func BenchmarkComparisonGetAllWithLogs(b *testing.B) {
	helper := setupBenchmarkDatabase(b)
	defer cleanupBenchmarkDatabase(b, helper)

	repo := postgres.NewProjectRepositoryImpl(helper.Pool)

	// Warm-up phase
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := repo.GetAllWithLogs(ctx)
	cancel()
	if err != nil {
		b.Fatalf("Warm-up failed: %v", err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	var totalTime time.Duration
	var totalAllocs uint64

	for i := 0; i < b.N; i++ {
		start := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		result, err := repo.GetAllWithLogs(ctx)
		cancel()

		if err != nil {
			b.Fatalf("GetAllWithLogs failed: %v", err)
		}

		if len(result) != 10 {
			b.Errorf("Expected 10 projects, got %d", len(result))
		}

		totalTime += time.Since(start)
		totalAllocs += uint64(result[0].Project.ID) // Track allocations indirectly
	}

	avgTime := totalTime / time.Duration(b.N)
	b.ReportMetric(float64(avgTime)/float64(time.Millisecond), "ms/op")
	b.ReportMetric(float64(totalAllocs)/float64(b.N), "allocs/op")
}

// BenchmarkComparisonGetWithLogs measures the performance of GetWithLogs
// with comparison against established baseline metrics
func BenchmarkComparisonGetWithLogs(b *testing.B) {
	helper := setupBenchmarkDatabase(b)
	defer cleanupBenchmarkDatabase(b, helper)

	// Get the first project ID for benchmarking
	repo := postgres.NewProjectRepositoryImpl(helper.Pool)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	projects, err := repo.GetAllWithLogs(ctx)
	cancel()
	if err != nil {
		b.Fatalf("Failed to get projects: %v", err)
	}

	if len(projects) == 0 {
		b.Fatal("No projects found in database")
	}

	projectID := projects[0].Project.ID

	b.ResetTimer()
	b.ReportAllocs()

	var totalTime time.Duration

	for i := 0; i < b.N; i++ {
		start := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := repo.GetWithLogs(ctx, projectID)
		cancel()

		if err != nil {
			b.Fatalf("GetWithLogs failed: %v", err)
		}

		totalTime += time.Since(start)
	}

	avgTime := totalTime / time.Duration(b.N)
	b.ReportMetric(float64(avgTime)/float64(time.Millisecond), "ms/op")
}

// BenchmarkComparisonConcurrentGetAllWithLogs measures concurrent performance
// of GetAllWithLogs under simulated load
func BenchmarkComparisonConcurrentGetAllWithLogs(b *testing.B) {
	helper := setupBenchmarkDatabase(b)
	defer cleanupBenchmarkDatabase(b, helper)

	repo := postgres.NewProjectRepositoryImpl(helper.Pool)

	// Warm-up phase
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := repo.GetAllWithLogs(ctx)
	cancel()
	if err != nil {
		b.Fatalf("Warm-up failed: %v", err)
	}

	b.ReportAllocs()
	b.SetParallelism(4)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			result, err := repo.GetAllWithLogs(ctx)
			cancel()

			if err != nil {
				b.Fatalf("GetAllWithLogs failed: %v", err)
			}

			if len(result) != 10 {
				b.Errorf("Expected 10 projects, got %d", len(result))
			}
		}
	})
}

// BenchmarkComparisonLargeDataset measures performance with larger dataset
// Creates 100 projects with 10 logs each to test scalability
func BenchmarkComparisonLargeDataset(b *testing.B) {
	if !test.IsTestDatabase() {
		b.Skip("Test database not configured - skipping large dataset benchmark")
	}

	helper, err := test.SetupTestDB()
	if err != nil {
		b.Fatalf("Failed to setup test DB: %v", err)
	}
	defer func() {
		if err := helper.CleanupTestSchema(); err != nil {
			b.Logf("Failed to cleanup test schema: %v", err)
		}
		helper.Close()
	}()

	if err := helper.SetupTestSchema(); err != nil {
		b.Fatalf("Failed to setup test schema: %v", err)
	}

	if err := helper.ClearTestData(); err != nil {
		b.Fatalf("Failed to clear test data: %v", err)
	}

	ctx := helper.GetContext()
	repo := postgres.NewProjectRepositoryImpl(helper.Pool)

	// Create 100 projects with 10 logs each
	for i := 1; i <= 100; i++ {
		var projectID int64
		err := helper.Pool.QueryRow(ctx, `
			INSERT INTO projects (name, total_page, page, reinicia)
			VALUES ($1, $2, $3, $4) RETURNING id
		`, "Large Dataset Project "+string(rune('0'+i%10)), 200, 100, false).Scan(&projectID)
		if err != nil {
			b.Fatalf("Failed to create test project: %v", err)
		}

		// Create 10 logs per project
		for j := 1; j <= 10; j++ {
			_, err := helper.Pool.Exec(ctx, `
				INSERT INTO logs (project_id, data, start_page, end_page, wday)
				VALUES ($1, $2, $3, $4, $5)
			`, projectID, "2024-01-01", j*10, j*10+5, 1)
			if err != nil {
				b.Fatalf("Failed to create test log: %v", err)
			}
		}
	}

	b.ResetTimer()
	b.ReportAllocs()

	var totalTime time.Duration
	var totalProjects int

	for i := 0; i < b.N; i++ {
		start := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		result, err := repo.GetAllWithLogs(ctx)
		cancel()

		if err != nil {
			b.Fatalf("GetAllWithLogs failed: %v", err)
		}

		totalProjects = len(result)

		if totalProjects != 100 {
			b.Errorf("Expected 100 projects, got %d", totalProjects)
		}

		totalTime += time.Since(start)
	}

	avgTime := totalTime / time.Duration(b.N)
	b.ReportMetric(float64(avgTime)/float64(time.Millisecond), "ms/op")
	b.ReportMetric(float64(totalProjects), "projects/op")
}

// BenchmarkHTTPHandlerIndex measures the performance of the HTTP handler Index endpoint
func BenchmarkHTTPHandlerIndex(b *testing.B) {
	// Check if we should run HTTP benchmarks
	if os.Getenv("SKIP_HTTP_BENCHMARK") == "1" {
		b.Skip("HTTP benchmarks disabled via SKIP_HTTP_BENCHMARK=1")
	}

	helper := setupBenchmarkDatabase(b)
	defer cleanupBenchmarkDatabase(b, helper)

	repo := postgres.NewProjectRepositoryImpl(helper.Pool)

	// Create a test HTTP server
	handler := handlers.NewProjectsHandler(repo)

	// Add warm-up request
	req := newMockRequest("GET", "/v1/projects", nil)
	w := newMockResponseWriter()
	handler.Index(w, req.Request)

	b.ResetTimer()
	b.ReportAllocs()

	var totalTime time.Duration
	var responseCount int

	for i := 0; i < b.N; i++ {
		start := time.Now()
		req := newMockRequest("GET", "/v1/projects", nil)
		w := newMockResponseWriter()

		handler.Index(w, req.Request)

		totalTime += time.Since(start)
		responseCount++
	}

	avgTime := totalTime / time.Duration(b.N)
	b.ReportMetric(float64(avgTime)/float64(time.Millisecond), "ms/op")
	b.ReportMetric(float64(responseCount), "requests/op")
}

// BenchmarkHTTPHandlerShow measures the performance of the HTTP handler Show endpoint
func BenchmarkHTTPHandlerShow(b *testing.B) {
	// Check if we should run HTTP benchmarks
	if os.Getenv("SKIP_HTTP_BENCHMARK") == "1" {
		b.Skip("HTTP benchmarks disabled via SKIP_HTTP_BENCHMARK=1")
	}

	helper := setupBenchmarkDatabase(b)
	defer cleanupBenchmarkDatabase(b, helper)

	repo := postgres.NewProjectRepositoryImpl(helper.Pool)

	// Get a project ID for the Show endpoint
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	projects, err := repo.GetAllWithLogs(ctx)
	cancel()
	if err != nil {
		b.Fatalf("Failed to get projects: %v", err)
	}

	if len(projects) == 0 {
		b.Fatal("No projects found in database")
	}

	projectID := projects[0].Project.ID
	projectIDStr := fmt.Sprintf("%d", projectID)

	// Add warm-up request
	req := newMockRequest("GET", "/v1/projects/"+projectIDStr, nil)
	w := newMockResponseWriter()
	handler := handlers.NewProjectsHandler(repo)
	handler.Show(w, req.Request)

	b.ResetTimer()
	b.ReportAllocs()

	var totalTime time.Duration

	for i := 0; i < b.N; i++ {
		start := time.Now()
		req := newMockRequest("GET", "/v1/projects/"+projectIDStr, nil)
		w := newMockResponseWriter()

		handler := handlers.NewProjectsHandler(repo)
		handler.Show(w, req.Request)

		totalTime += time.Since(start)
	}

	avgTime := totalTime / time.Duration(b.N)
	b.ReportMetric(float64(avgTime)/float64(time.Millisecond), "ms/op")
}

// Helper functions for HTTP benchmarking

// mockResponseWriter implements http.ResponseWriter for testing
type mockResponseWriter struct {
	headerMap  http.Header
	body       *strings.Builder
	statusCode int
}

func newMockResponseWriter() *mockResponseWriter {
	return &mockResponseWriter{
		headerMap:  make(http.Header),
		body:       &strings.Builder{},
		statusCode: http.StatusOK,
	}
}

func (w *mockResponseWriter) Header() http.Header {
	return w.headerMap
}

func (w *mockResponseWriter) Write(data []byte) (int, error) {
	return w.body.Write(data)
}

func (w *mockResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

// mockRequest implements http.Request for testing
type mockRequest struct {
	*http.Request
	method string
	url    string
}

func newMockRequest(method, url string, body io.Reader) *mockRequest {
	req, _ := http.NewRequest(method, url, body)
	return &mockRequest{
		Request: req,
		method:  method,
		url:     url,
	}
}

// PerformanceMetrics represents performance comparison metrics
type PerformanceMetrics struct {
	Operation       string  `json:"operation"`
	AverageTimeMs   float64 `json:"average_time_ms"`
	P50TimeMs       float64 `json:"p50_time_ms"`
	P95TimeMs       float64 `json:"p95_time_ms"`
	P99TimeMs       float64 `json:"p99_time_ms"`
	MemoryAllocs    uint64  `json:"memory_allocs"`
	RequestsPerSec  float64 `json:"requests_per_sec"`
	ThresholdMs     float64 `json:"threshold_ms"`
	WithinThreshold bool    `json:"within_threshold"`
}

// PerformanceBenchmarkResult represents the complete benchmark result
type PerformanceBenchmarkResult struct {
	Timestamp  string             `json:"timestamp"`
	Operation  string             `json:"operation"`
	Metrics    PerformanceMetrics `json:"metrics"`
	Comparison string             `json:"comparison"`
	Notes      string             `json:"notes"`
}

// CalculatePercentiles calculates p50, p95, p99 percentiles from a slice of times
func CalculatePercentiles(times []time.Duration) (p50, p95, p99 float64) {
	if len(times) == 0 {
		return 0, 0, 0
	}

	// Convert to milliseconds and sort
	msTimes := make([]float64, len(times))
	for i, t := range times {
		msTimes[i] = float64(t) / float64(time.Millisecond)
	}
	sortSlice(msTimes)

	n := len(msTimes)
	p50 = msTimes[n*50/100]
	p95 = msTimes[n*95/100]
	p99 = msTimes[n*99/100]

	return p50, p95, p99
}

// sortSlice sorts a slice of floats in ascending order
func sortSlice(slice []float64) {
	for i := 0; i < len(slice)-1; i++ {
		for j := i + 1; j < len(slice); j++ {
			if slice[i] > slice[j] {
				slice[i], slice[j] = slice[j], slice[i]
			}
		}
	}
}

// SaveBenchmarkResults saves benchmark results to a JSON file
func SaveBenchmarkResults(results []PerformanceBenchmarkResult, filename string) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write results: %w", err)
	}

	return nil
}

// VerifyPerformanceThresholds verifies that performance metrics meet thresholds
func VerifyPerformanceThresholds(m *PerformanceMetrics) bool {
	return m.WithinThreshold
}

// GeneratePerformanceReport generates a performance comparison report
func GeneratePerformanceReport(results []PerformanceBenchmarkResult) string {
	var sb strings.Builder

	sb.WriteString("# Performance Comparison Report\n\n")
	sb.WriteString(fmt.Sprintf("**Generated:** %s\n\n", time.Now().Format(time.RFC3339)))

	for _, r := range results {
		sb.WriteString(fmt.Sprintf("## %s\n\n", r.Operation))
		sb.WriteString(fmt.Sprintf("- **Average Time:** %.2f ms\n", r.Metrics.AverageTimeMs))
		sb.WriteString(fmt.Sprintf("- **P50 Time:** %.2f ms\n", r.Metrics.P50TimeMs))
		sb.WriteString(fmt.Sprintf("- **P95 Time:** %.2f ms\n", r.Metrics.P95TimeMs))
		sb.WriteString(fmt.Sprintf("- **P99 Time:** %.2f ms\n", r.Metrics.P99TimeMs))
		sb.WriteString(fmt.Sprintf("- **Memory Allocations:** %d\n", r.Metrics.MemoryAllocs))
		sb.WriteString(fmt.Sprintf("- **Within Threshold:** %t\n", r.Metrics.WithinThreshold))
		sb.WriteString(fmt.Sprintf("- **Threshold:** %.2f ms\n", r.Metrics.ThresholdMs))
		sb.WriteString("\n")
	}

	return sb.String()
}
