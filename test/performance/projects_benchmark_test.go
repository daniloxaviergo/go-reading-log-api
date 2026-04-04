package performance

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"go-reading-log-api-next/internal/adapter/postgres"
	"go-reading-log-api-next/test"
)

// setupBenchmarkDatabase creates a clean test database for benchmarking
// It sets up 10 projects with 5 logs each for consistent benchmarking
func setupBenchmarkDatabase(b *testing.B) *test.TestHelper {
	b.Helper()

	if !test.IsTestDatabase() {
		b.Skip("Test database not configured - skipping benchmark test")
	}

	helper, err := test.SetupTestDB()
	if err != nil {
		b.Fatalf("Failed to setup test DB: %v", err)
	}

	// Setup schema
	if err := helper.SetupTestSchema(); err != nil {
		helper.Close()
		b.Fatalf("Failed to setup test schema: %v", err)
	}

	// Clear any existing data
	if err := helper.ClearTestData(); err != nil {
		helper.Close()
		b.Fatalf("Failed to clear test data: %v", err)
	}

	ctx := helper.GetContext()

	// Create 10 projects with 5 logs each for consistent benchmarking
	for i := 1; i <= 10; i++ {
		var projectID int64
		err := helper.Pool.QueryRow(ctx, `
			INSERT INTO projects (name, total_page, page, reinicia)
			VALUES ($1, $2, $3, $4) RETURNING id
		`, "Benchmark Project "+string(rune('0'+i)), 200, 100, false).Scan(&projectID)
		if err != nil {
			helper.Close()
			b.Fatalf("Failed to create test project: %v", err)
		}

		// Create 5 logs per project
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

// cleanupBenchmarkDatabase cleans up the benchmark database
func cleanupBenchmarkDatabase(b *testing.B, helper *test.TestHelper) {
	b.Helper()
	if helper != nil {
		if err := helper.CleanupTestSchema(); err != nil {
			b.Logf("Failed to cleanup test schema: %v", err)
		}
		helper.Close()
	}
}

// BenchmarkGetAllWithLogs measures the performance of GetAllWithLogs repository method
// This benchmark runs 100 iterations to get stable timing measurements
func BenchmarkGetAllWithLogs(b *testing.B) {
	helper := setupBenchmarkDatabase(b)
	defer cleanupBenchmarkDatabase(b, helper)

	repo := postgres.NewProjectRepositoryImpl(helper.Pool)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := repo.GetAllWithLogs(ctx)
		cancel()

		if err != nil {
			b.Fatalf("GetAllWithLogs failed: %v", err)
		}
	}
}

// BenchmarkGetWithLogs measures the performance of GetWithLogs repository method
func BenchmarkGetWithLogs(b *testing.B) {
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
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := repo.GetWithLogs(ctx, projectID)
		cancel()

		if err != nil {
			b.Fatalf("GetWithLogs failed: %v", err)
		}
	}
}

// BenchmarkGetAllWithLogsConcurrent measures concurrent performance of GetAllWithLogs
func BenchmarkGetAllWithLogsConcurrent(b *testing.B) {
	helper := setupBenchmarkDatabase(b)
	defer cleanupBenchmarkDatabase(b, helper)

	repo := postgres.NewProjectRepositoryImpl(helper.Pool)

	b.ReportAllocs()
	b.SetParallelism(4)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			_, err := repo.GetAllWithLogs(ctx)
			cancel()

			if err != nil {
				b.Fatalf("GetAllWithLogs failed: %v", err)
			}
		}
	})
}

// BenchmarkGetWithLogsConcurrent measures concurrent performance of GetWithLogs
func BenchmarkGetWithLogsConcurrent(b *testing.B) {
	helper := setupBenchmarkDatabase(b)
	defer cleanupBenchmarkDatabase(b, helper)

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

	b.ReportAllocs()
	b.SetParallelism(4)

	projectIDs := make([]int64, len(projects))
	for i, p := range projects {
		projectIDs[i] = p.Project.ID
	}

	// Use atomic counter for round-robin access to project IDs
	var counter uint64

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			// Get current index atomically and increment
			idx := int(atomic.AddUint64(&counter, 1) % uint64(len(projectIDs)))
			_, err := repo.GetWithLogs(ctx, projectIDs[idx])
			cancel()

			if err != nil {
				b.Fatalf("GetWithLogs failed: %v", err)
			}
		}
	})
}

// BenchmarkGetAllWithLogsLargeDataset measures performance with larger dataset
// Creates 100 projects with 10 logs each to test scalability
func BenchmarkGetAllWithLogsLargeDataset(b *testing.B) {
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
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		result, err := repo.GetAllWithLogs(ctx)
		cancel()

		if err != nil {
			b.Fatalf("GetAllWithLogs failed: %v", err)
		}

		if len(result) != 100 {
			b.Errorf("Expected 100 projects, got %d", len(result))
		}
	}
}

// PerformanceNoteGetAllWithLogs documents the performance characteristics of GetAllWithLogs
/*
GetAllWithLogs Performance Analysis:
- Uses a single LEFT OUTER JOIN query to fetch all projects with logs
- Groups results in Go memory to avoid N+1 queries
- Performance scales linearly with number of projects and logs

Query: SELECT ... FROM projects p LEFT OUTER JOIN logs l ON p.id = l.project_id ORDER BY p.id ASC, l.data DESC

Optimizations applied:
- Single query with JOIN instead of separate queries per project
- Result set grouping in memory after query execution
- Proper index on logs(project_id) for efficient JOIN
*/

// PerformanceNoteGetWithLogs documents the performance characteristics of GetWithLogs
/*
GetWithLogs Performance Analysis:
- Uses two queries: one for project, one for logs
- Logs are fetched with a simple WHERE clause on project_id
- Query for logs: SELECT ... FROM logs WHERE project_id = $1 ORDER BY data DESC

Optimizations applied:
- Separate queries allow caching of project data if needed
- Efficient index usage on logs(project_id) for WHERE clause
- Proper ORDER BY for sorted log retrieval
*/
