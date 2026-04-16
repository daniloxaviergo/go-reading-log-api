---
id: RDL-053
title: '[doc-004 Phase 2.2] Verify parallel test performance impact'
status: To Do
assignee:
  - Thomas
created_date: '2026-04-15 12:15'
updated_date: '2026-04-16 12:02'
labels:
  - benchmark
  - performance
  - testing
dependencies: []
references:
  - 'R3.AC2: Test execution speed not significantly impacted'
documentation:
  - doc-004
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create performance benchmarks to measure the impact of the parallel test database uniqueness changes. The benchmarks should compare test execution times before and after the changes, ensuring that the overhead added by unique database naming doesn't significantly impact test speed. Specifically verify that test startup time increases by less than 200ms and overall test execution time doesn't regress by more than 10%.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Measure test execution time before/after changes
- [ ] #2 Ensure < 10% performance regression
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task requires creating performance benchmarks to measure the impact of parallel test database uniqueness changes introduced in Phase 1.4 (RDL-051). The benchmarks will compare test execution times before and after the changes, ensuring that the overhead added by unique database naming doesn't significantly impact test speed.

**Architectural Decisions:**

1. **Benchmark Structure**: Create a new benchmark test file `test/performance/parallel_test_benchmark.go` that specifically tests parallel test execution with unique database names. This follows the existing pattern in `test/performance/` directory.

2. **Performance Thresholds**: Establish baseline metrics and verify that:
   - Test startup time (database creation + connection) increases by less than 200ms
   - Overall test execution time doesn't regress by more than 10%
   - Cleanup time remains under 60 seconds even with 6,000+ orphaned databases

3. **Test Scenarios**: Implement benchmarks covering:
   - Single test execution (baseline)
   - Parallel test execution (8+ goroutines)
   - Orphaned database cleanup performance
   - Database uniqueness collision avoidance

4. **Reporting**: Generate JSON reports with percentile metrics (p50, p95, p99) to provide comprehensive performance analysis.

**Why This Approach:**
- Follows existing benchmark patterns in the codebase
- Uses Go's native benchmarking framework for reliability
- Provides detailed metrics for regression detection
- Aligns with the PRD's acceptance criteria for < 10% performance regression

---

### 2. Files to Modify

| File | Action | Description |
|------|--------|-------------|
| `test/performance/parallel_test_benchmark.go` | **Create** | New benchmark file for parallel test performance testing |
| `test/performance/comparison_test.go` | **Modify** | Add parallel test comparison benchmarks (if not exists) |
| `Makefile` | **Modify** | Add `benchmark-parallel` target |
| `docs/performance-benchmarks.md` | **Create** | Documentation for benchmark results and methodology |

**Detailed File Changes:**

**New File: `test/performance/parallel_test_benchmark.go`** (~300-400 lines)
- `BenchmarkParallelTestStartup` - Measure database creation + connection time
- `BenchmarkParallelTestExecution` - Measure test execution with 8+ parallel goroutines
- `BenchmarkParallelCleanup` - Measure cleanup performance with many orphaned databases
- `BenchmarkDatabaseUniqueness` - Verify no collisions with unique database names

**Modification: `Makefile`** (add near line 100)
```makefile
benchmark-parallel:
	@echo "$(BLUE)Running parallel performance benchmarks...$(NC)"
	@export $$(xargs < .env.test | grep -v '^#' | xargs) && \
		$(GO) test -bench=BenchmarkParallel -benchmem -count=5 $(TEST_PKG)/performance
	@echo "$(GREEN)Benchmark complete$(NC)"
```

---

### 3. Dependencies

**Prerequisites:**
- Go 1.25.7 (already in use by the project) ✓
- pgx/v5 library (already installed) ✓
- PostgreSQL running and accessible
- `.env.test` file with database credentials

**No new dependencies required.**

**Existing Dependencies Used:**
- `github.com/jackc/pgx/v5/pgxpool` - Database connection pooling
- `github.com/joho/godotenv` - Environment variable loading
- `go-reading-log-api-next/test` - Test helper functions
- `go-reading-log-api-next/internal/config` - Configuration loading

---

### 4. Code Patterns

**Follow existing patterns from the codebase:**

1. **Benchmark Structure** (from `test/performance/comparison_test.go`):
```go
func BenchmarkOperationName(b *testing.B) {
    helper := setupBenchmarkDatabase(b)
    defer cleanupBenchmarkDatabase(b, helper)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // Benchmark code here
    }
}
```

2. **Parallel Test Pattern** (from `test/performance/comparison_test.go`):
```go
b.SetParallelism(8)  // Test with 8 parallel goroutines
b.RunParallel(func(pb *testing.PB) {
    for pb.Next() {
        // Parallel test code here
    }
})
```

3. **Metrics Reporting**:
```go
avgTime := totalTime / time.Duration(b.N)
b.ReportMetric(float64(avgTime)/float64(time.Millisecond), "ms/op")
b.ReportAllocs()
```

4. **Context Usage**: All database operations use context with timeout:
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

**Naming Conventions:**
- Benchmark functions: `BenchmarkParallelTestStartup`, `BenchmarkParallelTestExecution`
- Helper functions: `setupParallelBenchmark`, `cleanupParallelBenchmark`
- Variable names: `numParallel`, `testDBName`, `benchmarkResults`

---

### 5. Testing Strategy

**Benchmark Tests to Implement:**

| Benchmark | Purpose | Threshold | Metrics |
|-----------|---------|-----------|---------|
| `BenchmarkParallelTestStartup` | Measure DB creation + connection time | < 200ms | avg, p95 |
| `BenchmarkParallelTestExecution` | Measure test execution with 8 goroutines | < 10% regression | ops/sec |
| `BenchmarkParallelCleanup` | Measure cleanup with 6000+ orphans | < 60s | total time |
| `BenchmarkDatabaseUniqueness` | Verify no collisions with unique names | 0 collisions | count |

**Test Execution:**
```bash
# Run parallel benchmarks
make benchmark-parallel

# Or run directly
go test -bench=BenchmarkParallel -benchmem -count=5 ./test/performance

# Generate report
go tool pprof -http=:8080 profile.out
```

**Acceptance Verification:**
1. Run benchmarks 3 times each for statistical reliability
2. Calculate average and p95 metrics
3. Compare against thresholds:
   - Startup time: < 200ms ✓
   - Execution regression: < 10% ✓
   - Cleanup time: < 60s ✓
4. Generate JSON report with all metrics

---

### 6. Risks and Considerations

**Blocking Issues:**
- None identified. Implementation uses existing patterns from `test/performance/` directory.

**Potential Pitfalls:**
1. **Test Database Cleanup Interference**: Parallel cleanup might interfere with ongoing tests
   - *Mitigation*: Use unique database names per test, cleanup only orphaned databases

2. **PostgreSQL Connection Pool Limits**: 8+ parallel tests might exhaust connection pool
   - *Mitigation*: Monitor connection count, use connection pooling with appropriate max connections

3. **Disk I/O Contention**: Multiple database creates/drops might cause I/O contention
   - *Mitigation*: Run benchmarks on isolated test environment, not shared development database

4. **Statistical Variance**: Benchmark results might vary between runs
   - *Mitigation*: Run each benchmark multiple times, use p95 for threshold comparison

**Trade-offs:**
- Using 8 parallel goroutines balances realistic parallelism with resource constraints
- 60-second timeout for cleanup prevents hanging but might fail on very large orphan sets
- p95 metrics chosen over p99 for more stable threshold comparisons

**Deployment Considerations:**
- No migration required (no schema changes)
- No downtime required
- Benchmarks can be run in CI/CD pipeline
- Results should be stored for regression tracking

**Verification Checklist:**
- [ ] `make benchmark-parallel` runs without error
- [ ] All benchmarks complete within expected time
- [ ] JSON report is generated with all metrics
- [ ] Metrics meet acceptance thresholds
- [ ] No test database collisions occur
- [ ] Cleanup completes within 60 seconds

---

### 7. Implementation Steps

**Step 1: Create `test/performance/parallel_test_benchmark.go`**
- Import required packages (context, testing, time, etc.)
- Implement `setupParallelBenchmark` helper
- Implement `BenchmarkParallelTestStartup`
- Implement `BenchmarkParallelTestExecution` with 8 goroutines
- Implement `BenchmarkParallelCleanup`
- Implement `BenchmarkDatabaseUniqueness`

**Step 2: Update `Makefile`**
- Add `benchmark-parallel` target
- Add colorized output for consistency

**Step 3: Generate and Verify Results**
- Run benchmarks 3 times each
- Calculate average metrics
- Verify against thresholds
- Generate JSON report

**Step 4: Documentation**
- Update `docs/performance-benchmarks.md` with methodology
- Document threshold definitions
- Include example benchmark results
<!-- SECTION:PLAN:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 All unit tests pass use testing-expert subagent for test execution and verification
- [ ] #2 All integration tests pass use testing-expert subagent for test execution and verification
- [ ] #3 go fmt and go vet pass with no errors
- [ ] #4 Clean Architecture layers properly followed
- [ ] #5 Error responses consistent with existing patterns
- [ ] #6 HTTP status codes correct for response type
- [x] #7 Database queries optimized with proper indexes
- [x] #8 Documentation updated in QWEN.md
- [x] #9 New code paths include error path tests
- [x] #10 HTTP handlers test both success and error responses
- [x] #11 Integration tests verify actual database interactions
- [x] #12 Tests use testing-expert subagent for test execution and verification
<!-- DOD:END -->
