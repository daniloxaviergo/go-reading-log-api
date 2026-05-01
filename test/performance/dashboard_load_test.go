package performance

import (
	"context"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go-reading-log-api-next/internal/api/v1/handlers"
	"go-reading-log-api-next/test"
)

// =============================================================================
// Dashboard Concurrent Load Tests
// =============================================================================

// BenchmarkDashboardConcurrent10 measures performance with 10 concurrent users
func BenchmarkDashboardConcurrent10(b *testing.B) {
	if !test.IsTestDatabase() {
		b.Skip("Test database not configured - skipping benchmark test")
	}

	helper := setupDashboardBenchmark(b)
	defer cleanupDashboardBenchmark(b, helper)

	handler := createDashboardHandler(helper.Pool)

	// Warm-up phase
	warmupConcurrent(b, handler, 10)

	// Benchmark
	b.ResetTimer()
	b.ReportAllocs()

	var totalTime time.Duration
	var successCount int64
	var failCount int64

	b.SetParallelism(10)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			start := time.Now()
			req := httptest.NewRequest("GET", "/v1/dashboard/day.json", nil)
			recorder := httptest.NewRecorder()

			handler.Day(recorder, req)

			if recorder.Code == http.StatusOK {
				atomic.AddInt64(&successCount, 1)
			} else {
				atomic.AddInt64(&failCount, 1)
			}

			totalTime += time.Since(start)
		}
	})

	// Calculate metrics
	avgTime := totalTime / time.Duration(successCount+failCount)
	opsPerSec := float64(successCount+failCount) / float64(totalTime) * float64(time.Second)

	b.ReportMetric(float64(avgTime)/float64(time.Millisecond), "ms/op")
	b.ReportMetric(opsPerSec, "ops/sec")
	b.ReportMetric(float64(successCount), "successes")
	b.ReportMetric(float64(failCount), "failures")

	// Verify QPS target (>100 QPS)
	if opsPerSec < 100 {
		b.Errorf("QPS %.2f below target of 100 QPS", opsPerSec)
	} else {
		b.Logf("QPS %.2f meets target (> 100 QPS)", opsPerSec)
	}

	// Verify error rate
	total := successCount + failCount
	if total > 0 && float64(failCount)/float64(total) > 0.01 {
		b.Errorf("Error rate %.2f%% exceeds 1%% threshold", float64(failCount)/float64(total)*100)
	} else {
		b.Logf("Error rate %.2f%% is within threshold (< 1%%)", float64(failCount)/float64(total)*100)
	}

	// Verify connection pool
	verifyConnectionPool(b, helper.Pool)
}

// BenchmarkDashboardConcurrent50 measures performance with 50 concurrent users
func BenchmarkDashboardConcurrent50(b *testing.B) {
	if !test.IsTestDatabase() {
		b.Skip("Test database not configured - skipping benchmark test")
	}

	helper := setupDashboardBenchmark(b)
	defer cleanupDashboardBenchmark(b, helper)

	handler := createDashboardHandler(helper.Pool)

	// Warm-up phase
	warmupConcurrent(b, handler, 50)

	// Benchmark
	b.ResetTimer()
	b.ReportAllocs()

	var totalTime time.Duration
	var successCount int64
	var failCount int64

	b.SetParallelism(50)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			start := time.Now()
			req := httptest.NewRequest("GET", "/v1/dashboard/day.json", nil)
			recorder := httptest.NewRecorder()

			handler.Day(recorder, req)

			if recorder.Code == http.StatusOK {
				atomic.AddInt64(&successCount, 1)
			} else {
				atomic.AddInt64(&failCount, 1)
			}

			totalTime += time.Since(start)
		}
	})

	// Calculate metrics
	avgTime := totalTime / time.Duration(successCount+failCount)
	opsPerSec := float64(successCount+failCount) / float64(totalTime) * float64(time.Second)

	b.ReportMetric(float64(avgTime)/float64(time.Millisecond), "ms/op")
	b.ReportMetric(opsPerSec, "ops/sec")
	b.ReportMetric(float64(successCount), "successes")
	b.ReportMetric(float64(failCount), "failures")

	// Verify QPS target (>100 QPS)
	if opsPerSec < 100 {
		b.Errorf("QPS %.2f below target of 100 QPS", opsPerSec)
	} else {
		b.Logf("QPS %.2f meets target (> 100 QPS)", opsPerSec)
	}

	// Verify error rate
	total := successCount + failCount
	if total > 0 && float64(failCount)/float64(total) > 0.01 {
		b.Errorf("Error rate %.2f%% exceeds 1%% threshold", float64(failCount)/float64(total)*100)
	} else {
		b.Logf("Error rate %.2f%% is within threshold (< 1%%)", float64(failCount)/float64(total)*100)
	}

	// Verify connection pool
	verifyConnectionPool(b, helper.Pool)
}

// BenchmarkDashboardConcurrent100 measures performance with 100 concurrent users
func BenchmarkDashboardConcurrent100(b *testing.B) {
	if !test.IsTestDatabase() {
		b.Skip("Test database not configured - skipping benchmark test")
	}

	helper := setupDashboardBenchmark(b)
	defer cleanupDashboardBenchmark(b, helper)

	handler := createDashboardHandler(helper.Pool)

	// Warm-up phase
	warmupConcurrent(b, handler, 100)

	// Benchmark
	b.ResetTimer()
	b.ReportAllocs()

	var totalTime time.Duration
	var successCount int64
	var failCount int64

	b.SetParallelism(100)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			start := time.Now()
			req := httptest.NewRequest("GET", "/v1/dashboard/day.json", nil)
			recorder := httptest.NewRecorder()

			handler.Day(recorder, req)

			if recorder.Code == http.StatusOK {
				atomic.AddInt64(&successCount, 1)
			} else {
				atomic.AddInt64(&failCount, 1)
			}

			totalTime += time.Since(start)
		}
	})

	// Calculate metrics
	avgTime := totalTime / time.Duration(successCount+failCount)
	opsPerSec := float64(successCount+failCount) / float64(totalTime) * float64(time.Second)

	b.ReportMetric(float64(avgTime)/float64(time.Millisecond), "ms/op")
	b.ReportMetric(opsPerSec, "ops/sec")
	b.ReportMetric(float64(successCount), "successes")
	b.ReportMetric(float64(failCount), "failures")

	// Verify QPS target (>100 QPS)
	if opsPerSec < 100 {
		b.Errorf("QPS %.2f below target of 100 QPS", opsPerSec)
	} else {
		b.Logf("QPS %.2f meets target (> 100 QPS)", opsPerSec)
	}

	// Verify error rate
	total := successCount + failCount
	if total > 0 && float64(failCount)/float64(total) > 0.01 {
		b.Errorf("Error rate %.2f%% exceeds 1%% threshold", float64(failCount)/float64(total)*100)
	} else {
		b.Logf("Error rate %.2f%% is within threshold (< 1%%)", float64(failCount)/float64(total)*100)
	}

	// Verify connection pool
	verifyConnectionPool(b, helper.Pool)
}

// BenchmarkDashboardConcurrentAllEndpoints measures performance across all endpoints with concurrent users
func BenchmarkDashboardConcurrentAllEndpoints(b *testing.B) {
	if !test.IsTestDatabase() {
		b.Skip("Test database not configured - skipping benchmark test")
	}

	helper := setupDashboardBenchmark(b)
	defer cleanupDashboardBenchmark(b, helper)

	handler := createDashboardHandler(helper.Pool)

	// Endpoints to test
	endpoints := []string{
		"/v1/dashboard/day.json",
		"/v1/dashboard/projects.json",
		"/v1/dashboard/last_days.json?days=7&type=1",
		"/v1/dashboard/echart/faults.json",
		"/v1/dashboard/echart/speculate_actual.json",
		"/v1/dashboard/echart/faults_week_day.json",
		"/v1/dashboard/echart/mean_progress.json",
		"/v1/dashboard/echart/last_year_total.json",
	}

	// Warm-up phase
	warmupConcurrentAllEndpoints(b, handler, endpoints, 50)

	// Benchmark
	b.ResetTimer()
	b.ReportAllocs()

	var totalTime time.Duration
	var successCount int64
	var failCount int64
	var endpointCounts [8]int64

	b.SetParallelism(50)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Select a random endpoint
			endpoint := endpoints[int(time.Now().UnixNano())%len(endpoints)]

			start := time.Now()
			req := httptest.NewRequest("GET", endpoint, nil)
			recorder := httptest.NewRecorder()

			// Route to appropriate handler
			switch {
			case strings.Contains(endpoint, "/day.json"):
				handler.Day(recorder, req)
			case strings.Contains(endpoint, "/projects.json"):
				handler.Projects(recorder, req)
			case strings.Contains(endpoint, "/last_days.json"):
				handler.LastDays(recorder, req)
			case strings.Contains(endpoint, "/faults.json") && !strings.Contains(endpoint, "week_day"):
				handler.Faults(recorder, req)
			case strings.Contains(endpoint, "/speculate_actual.json"):
				handler.SpeculateActual(recorder, req)
			case strings.Contains(endpoint, "/faults_week_day.json"):
				handler.WeekdayFaults(recorder, req)
			case strings.Contains(endpoint, "/mean_progress.json"):
				handler.MeanProgress(recorder, req)
			case strings.Contains(endpoint, "/last_year_total.json"):
				handler.YearlyTotal(recorder, req)
			}

			if recorder.Code == http.StatusOK {
				atomic.AddInt64(&successCount, 1)
				// Track per-endpoint counts
				for i, e := range endpoints {
					if e == endpoint {
						atomic.AddInt64(&endpointCounts[i], 1)
						break
					}
				}
			} else {
				atomic.AddInt64(&failCount, 1)
			}

			totalTime += time.Since(start)
		}
	})

	// Calculate metrics
	avgTime := totalTime / time.Duration(successCount+failCount)
	opsPerSec := float64(successCount+failCount) / float64(totalTime) * float64(time.Second)

	b.ReportMetric(float64(avgTime)/float64(time.Millisecond), "ms/op")
	b.ReportMetric(opsPerSec, "ops/sec")
	b.ReportMetric(float64(successCount), "successes")
	b.ReportMetric(float64(failCount), "failures")

	// Print per-endpoint statistics
	b.Log("Per-endpoint request counts:")
	for i, count := range endpointCounts {
		b.Logf("  %s: %d requests", endpoints[i], count)
	}

	// Verify QPS target (>100 QPS)
	if opsPerSec < 100 {
		b.Errorf("QPS %.2f below target of 100 QPS", opsPerSec)
	} else {
		b.Logf("QPS %.2f meets target (> 100 QPS)", opsPerSec)
	}

	// Verify error rate
	total := successCount + failCount
	if total > 0 && float64(failCount)/float64(total) > 0.01 {
		b.Errorf("Error rate %.2f%% exceeds 1%% threshold", float64(failCount)/float64(total)*100)
	} else {
		b.Logf("Error rate %.2f%% is within threshold (< 1%%)", float64(failCount)/float64(total)*100)
	}

	// Verify connection pool
	verifyConnectionPool(b, helper.Pool)
}

// =============================================================================
// Helper Functions
// =============================================================================

func warmupConcurrent(b *testing.B, handler *handlers.DashboardHandler, concurrentUsers int) {
	b.Helper()

	var wg sync.WaitGroup
	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for i := 0; i < concurrentUsers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := httptest.NewRequest("GET", "/v1/dashboard/day.json", nil)
			recorder := httptest.NewRecorder()
			handler.Day(recorder, req)
		}()
	}

	wg.Wait()
}

func warmupConcurrentAllEndpoints(b *testing.B, handler *handlers.DashboardHandler, endpoints []string, concurrentUsers int) {
	b.Helper()

	var wg sync.WaitGroup
	_, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for i := 0; i < concurrentUsers; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			endpoint := endpoints[idx%len(endpoints)]
			req := httptest.NewRequest("GET", endpoint, nil)
			recorder := httptest.NewRecorder()

			switch {
			case strings.Contains(endpoint, "/day.json"):
				handler.Day(recorder, req)
			case strings.Contains(endpoint, "/projects.json"):
				handler.Projects(recorder, req)
			case strings.Contains(endpoint, "/last_days.json"):
				handler.LastDays(recorder, req)
			case strings.Contains(endpoint, "/faults.json") && !strings.Contains(endpoint, "week_day"):
				handler.Faults(recorder, req)
			case strings.Contains(endpoint, "/speculate_actual.json"):
				handler.SpeculateActual(recorder, req)
			case strings.Contains(endpoint, "/faults_week_day.json"):
				handler.WeekdayFaults(recorder, req)
			case strings.Contains(endpoint, "/mean_progress.json"):
				handler.MeanProgress(recorder, req)
			case strings.Contains(endpoint, "/last_year_total.json"):
				handler.YearlyTotal(recorder, req)
			}
		}(i)
	}

	wg.Wait()
}

// BenchmarkDashboardSustainedLoad measures sustained load performance
func BenchmarkDashboardSustainedLoad(b *testing.B) {
	if !test.IsTestDatabase() {
		b.Skip("Test database not configured - skipping benchmark test")
	}

	helper := setupDashboardBenchmark(b)
	defer cleanupDashboardBenchmark(b, helper)

	handler := createDashboardHandler(helper.Pool)

	// Sustained load test parameters
	const (
		duration     = 10 * time.Second
		concurrent   = 20
		requestsPerS = 50
	)

	b.Logf("Starting sustained load test: %v duration, %d concurrent, %d requests/s",
		duration, concurrent, requestsPerS)

	// Warm-up phase
	warmupConcurrent(b, handler, concurrent)

	// Sustained load benchmark
	b.ResetTimer()
	b.ReportAllocs()

	var totalTime time.Duration
	var successCount int64
	var failCount int64

	// Create worker pool
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), duration+10*time.Second)
	defer cancel()

	// Calculate total requests based on duration and rate
	totalRequests := int32(requestsPerS * duration.Seconds())
	var requestsMade int32 = 0

	for i := 0; i < concurrent; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				default:
					// Rate limiting - sleep to achieve target requests per second
					if requestsMade >= totalRequests {
						return
					}

					start := time.Now()
					req := httptest.NewRequest("GET", "/v1/dashboard/day.json", nil)
					recorder := httptest.NewRecorder()

					handler.Day(recorder, req)

					if recorder.Code == http.StatusOK {
						atomic.AddInt64(&successCount, 1)
					} else {
						atomic.AddInt64(&failCount, 1)
					}

					totalTime += time.Since(start)
					atomic.AddInt32(&requestsMade, 1)

					// Sleep to rate limit
					elapsed := time.Since(start)
					targetDuration := time.Duration(float64(time.Second) / float64(requestsPerS/concurrent))
					if elapsed < targetDuration {
						time.Sleep(targetDuration - elapsed)
					}
				}
			}
		}(i)
	}

	wg.Wait()
	b.StopTimer()

	// Calculate metrics
	actualDuration := duration.Seconds()
	opsPerSec := float64(successCount+failCount) / actualDuration

	b.ReportMetric(float64(totalTime)/float64(time.Millisecond), "total_ms")
	b.ReportMetric(opsPerSec, "ops/sec")
	b.ReportMetric(float64(successCount), "successes")
	b.ReportMetric(float64(failCount), "failures")

	// Verify QPS target (>100 QPS)
	if opsPerSec < 100 {
		b.Errorf("QPS %.2f below target of 100 QPS", opsPerSec)
	} else {
		b.Logf("QPS %.2f meets target (> 100 QPS)", opsPerSec)
	}

	// Verify error rate
	total := successCount + failCount
	if total > 0 && float64(failCount)/float64(total) > 0.01 {
		b.Errorf("Error rate %.2f%% exceeds 1%% threshold", float64(failCount)/float64(total)*100)
	} else {
		b.Logf("Error rate %.2f%% is within threshold (< 1%%)", float64(failCount)/float64(total)*100)
	}

	// Verify connection pool
	verifyConnectionPool(b, helper.Pool)
}

// newSeededRNG creates a new random number generator with the given seed
func newSeededRNG(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed))
}
