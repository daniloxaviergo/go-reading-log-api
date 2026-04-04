---
id: RDL-036
title: '[doc-002 Phase 5] Perform performance comparison test'
status: To Do
assignee:
  - thomas
created_date: '2026-04-03 14:05'
updated_date: '2026-04-04 06:53'
labels:
  - phase-5
  - performance-test
  - benchmarking
dependencies: []
references:
  - 'PRD Section: Acceptance Criteria - NF1'
  - NF2 Performance
  - 'PRD Section: Technical Decisions - Decision 4'
documentation:
  - doc-002
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Run comprehensive performance comparison between Go and Rails API responses using same test data. Verify response time within 10% threshold and memory usage within 20% increase. Document any regressions and optimize identified bottlenecks.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Response time within 10% of Rails implementation
- [ ] #2 Memory usage within 20% increase threshold
- [ ] #3 Performance regression identified and resolved
- [ ] #4 Performance metrics documented
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

Perform comprehensive performance comparison between Go and Rails API endpoints using identical test data and methodology. The task has an important nuance: the Rails app exists but its directory structure suggests it's not fully functional (contains only Rails skeleton without full application code). This creates a unique challenge that requires an alternate approach.

**Performance Comparison Strategy:**

1. **Primary Method - Go-to-Go Baseline**: Since the Rails app directory appears incomplete, establish performance baselines using Go implementation in different configurations:
   - Run current Go implementation against production-like dataset
   - Compare against optimized version (if optimizations needed)
   - Use established Go benchmark results from RDL-029 as reference

2. **Secondary Method - Rails Comparison**: If Rails app can be made functional:
   - Start Rails app in docker-compose alongside Go app
   - Run identical benchmark requests against both endpoints
   - Use `wrk` or similar tool for consistent performance measurement

3. **Performance Metrics to Collect**:
   - Response time (p50, p95, p99 percentiles)
   - Memory usage (RSS in MB)
   - Requests per second (RPS)
   - Database query count
   - Error rates under load

4. **Threshold Verification**:
   - Response time within 10% of baseline
   - Memory usage within 20% of baseline
   - Identify outliers and optimize

**Approach Rationale:**
- RDL-029 demonstrated current Go implementation meets performance targets
- Rails app directory incomplete - requires verification before comparison
- Establishing Go baseline first provides useful参照 even without Rails comparison
- If Rails app can be fixed, direct comparison provides most valuable data

### 2. Files to Modify

#### New Files to Create
| File | Purpose |
|------|---------|
| `test/performance/comparison_test.go` | Comprehensive performance comparison test suite |
| `test/performance/run_comparison.sh` | Bash script to run performance comparison |
| `docs/performance-comparison.md` | Detailed performance comparison documentation |
| `test/benchmark_data/fixtures.go` | Benchmark data fixtures for consistent testing |

#### Modified Files
| File | Change Type | Reason |
|------|-------------|--------|
| `test/performance/projects_benchmark_test.go` | Add tests | Extend existing benchmarks with comparison scenarios |
| `internal/api/v1/handlers/projects_handler.go` | Read only | Verify handler implementation matches performance expectations |
| `internal/adapter/postgres/project_repository.go` | Read only | Verify query efficiency |

#### Existing Benchmarks to Use (RDL-029)
| File | Available Data |
|------|----------------|
| `test/performance/projects_benchmark_test.go` | Already created with 6 benchmark functions |

### 3. Dependencies

#### Prerequisites (From已完成 Tasks)
- ✅ **RDL-008**: Test infrastructure established
- ✅ **RDL-009**: Test coverage verified (121 tests passing)
- ✅ **RDL-029**: Query performance verified against Rails (identical performance ~0.12% difference)
- ✅ **RDL-035**: Edge cases tested and documented
- ✅ **RDL-034**: JSON response comparison test executed

#### Required Infrastructure
- PostgreSQL running with test database
- Docker available for Rails app if needed
- Benchmark tools: `go test -bench`, `wrk` or `hey` for HTTP benchmarking

#### External Dependencies
| Tool | Purpose | Current Status |
|------|---------|----------------|
| `go test -bench` | Go benchmarking | Built-in, available |
| `wrk` | HTTP benchmarking | Install via Makefile or docker-compose |
| `pprof` | Memory profiling | Built-in with Go |

### 4. Code Patterns

#### Go Benchmarking Pattern (Existing from RDL-029)
```go
func BenchmarkEndpoint(b *testing.B) {
    helper := setupBenchmarkDatabase(b)
    defer cleanupBenchmarkDatabase(b, helper)
    
    repo := postgres.NewProjectRepositoryImpl(helper.Pool)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        _, err := repo.GetAllWithLogs(ctx)
        cancel()
        if err != nil {
            b.Fatalf("Benchmark failed: %v", err)
        }
    }
}
```

#### HTTP Benchmarking Pattern (for Rails comparison)
```bash
# Run wrk benchmark
wrk -t4 -c100 -d30s http://localhost:3000/api/v1/projects

# Output format: Requests/sec, Latency, Transfer/sec
```

#### Performance Metrics to Log
```go
// Log benchmark results in consistent format
b.ReportAllocs()                    // Memory allocations
b.ReportMetric(float64(b.N), "op/benchmark")  // Operations count
```

### 5. Testing Strategy

#### Phase 1: baseline measurements (Go implementation only)

1. **Repository-level benchmarks** (already exist from RDL-029):
   - `BenchmarkGetAllWithLogs` - Baseline performance
   - `BenchmarkGetWithLogs` - Single record performance
   - `BenchmarkGetAllWithLogsConcurrent` - Concurrent performance
   - `BenchmarkGetAllWithLogsLargeDataset` - Scalability

2. **HTTP handler benchmarks** (new):
   - `BenchmarkProjectsHandler_Index` - List all projects endpoint
   - `BenchmarkProjectsHandler_Show` - Show single project endpoint
   - `BenchmarkLogsHandler_Index` - List logs endpoint

3. **Load tests** (new):
   - 10 concurrent users, 30 seconds
   - 50 concurrent users, 30 seconds
   - 100 concurrent users, 30 seconds

#### Phase 2: Rails comparison (if rails app functional)

1. **Setup comparison environment**:
   ```bash
   # Start both apps
   docker-compose up -d go-api rails-api
   
   # Wait for both to be ready
   sleep 10
   ```

2. **Run identical benchmarks**:
   ```bash
   # Run wrk against both endpoints
   wrk -t4 -c100 -d30s http://localhost:3000/api/v1/projects > go_results.txt
   wrk -t4 -c100 -d30s http://localhost:3001/api/v1/projects > rails_results.txt
   ```

3. **Compare results**:
   - Calculate percentage difference: `(Go - Rails) / Rails * 100`
   - Verify within 10% threshold
   - Document any deviations

#### Phase 3: Memory profiling (new)

1. **Profile Go memory usage**:
   ```go
   // Add to benchmark
   b.SetProfileRate(1)  // Profile all allocations
   runtime.GC()          // Force garbage collection
   ```

2. **Memory threshold verification**:
   - Current memory usage as baseline
   - Verify <20% increase under load
   - Profile with `go tool pprof`

### 6. Risks and Considerations

#### Risk 1: Rails App Incomplete
**Description**: The `rails-app/` directory exists but may lack complete application code or dependencies.

**Mitigation**:
- First verify Rails app can start: `docker-compose up rails`
- If Rails app incomplete, proceed with Go-to-Go baseline comparison
- Document limitation in results report

**Contingency**: If Rails app cannot be fixed, document Go performance baseline and recommend performance testing as part of future Rails migration stages.

#### Risk 2: Benchmark Variance
**Description**: Performance benchmarks can vary due to system load, caching, or other factors.

**Mitigation**:
- Run each benchmark multiple times (min 5 iterations)
- Report median/percentile values (not just mean)
- Use statistical significance testing if differences are marginal
- Clear system cache between runs: `echo 3 | sudo tee /proc/sys/vm/drop_caches`

#### Risk 3: Over-Optimization
**Description**: Optimizing for benchmark may hurt production performance (e.g., over-caching).

**Mitigation**:
- Keep benchmarks close to production patterns
- Document optimization trade-offs
- Test with realistic dataset sizes
- Monitor memory usage alongside speed

#### Risk 4: Threshold Misinterpretation
**Description**: 10% performance difference may be statistically insignificant.

**Mitigation**:
- Report confidence intervals for all measurements
- Use appropriate statistical tests (t-test, etc.)
- Consider business impact: 10% on 1ms = 0.1ms vs 10% on 1000ms = 100ms
- Focus on absolute improvements where meaningful

#### Considerations

1. **Test Dataset Size**:
   - Small (10 projects, 50 logs): Baseline measurements
   - Medium (100 projects, 500 logs): Performance regression detection
   - Large (1000 projects, 5000 logs): Scalability verification

2. **Performance Goal Realism**:
   - Rails uses Ruby (slower language) but has mature ORM
   - Go uses compiled code but may have less optimized queries
   - Target: 10% difference is achievable with proper optimization
   - Expected: Go should be faster or equal with best practices

3. **Definition of "Performance"**:
   - Response time (wall clock) - Primary metric
   - Memory usage - Secondary metric
   - Error rate - Secondary metric (should be < 1%)
   - Throughput (RPS) - Tertiary metric

4. **Reporting Requirements**:
   - Raw benchmark data for all test runs
   - Statistical summary (mean, median, p95, p99)
   - Comparison with threshold
   - Root cause analysis for any regressions
   - Optimization recommendations

### 7. Acceptance Criteria Mapping

| AC | Method | Verification |
|----|--------|--------------|
| AC1: Response time within 10% of Rails | Compare Go vs Rails benchmarks | Calculate percentage difference |
| AC2: Memory usage within 20% increase | Profile memory with pprof | Compare RSS before/after |
| AC3: Performance regression identified and resolved | Benchmark before/after changes | Document fixes applied |
| AC4: Performance metrics documented | Generate `docs/performance-comparison.md` | Report includes all metrics |

### 8. Implementation Steps

1. **Verify Rails App Functionality** (1-2 hours)
   ```bash
   cd /home/danilo/scripts/github/go-reading-log-api-next
   docker-compose up rails -d
   docker-compose logs rails
   # Verify Rails app starts and responds to /api/v1/projects
   ```

2. **Generate Performance Baseline** (2-3 hours)
   ```bash
   # Run Go benchmarks
   go test -bench=Benchmark -benchmem ./test/performance/
   # Save results to baseline_results.txt
   ```

3. **Run Rails Comparison** (2-3 hours, if rails functional)
   ```bash
   # Start both apps
   docker-compose up -d go-api rails-api
   
   # Run wrk benchmarks
   wrk -t4 -c100 -d30s http://localhost:3000/api/v1/projects > go_results.txt
   wrk -t4 -c100 -d30s http://localhost:3001/api/v1/projects > rails_results.txt
   
   # Analyze results
   # Calculate percentage difference
   ```

4. **Analyze Results and Identify Issues** (1-2 hours)
   - Compare measured values against thresholds
   - Identify any violations
   - Profile code for hotspots if needed
   - Fix identified issues

5. **Document Findings** (1-2 hours)
   - Generate `docs/performance-comparison.md`
   - Include raw benchmark data
   - Include statistical analysis
   - Include optimization recommendations

6. **Final Verification** (1 hour)
   - Run all tests pass: `make test`
   - Format and vet: `make fmt && make vet`
   - Verify DOD items checked off

**Total Estimated Time**: 8-14 hours

### 9. Success Criteria

The task is successful when:
1. Performance comparison executed (Go vs Rails or Go baseline)
2. All metrics within thresholds OR regressions identified and fixed
3. Documentation complete in `docs/performance-comparison.md`
4. All acceptance criteria checked off in task
5. All Definition of Done items verified

### 10. Existing Assets to Leverage

#### From RDL-029 (Query Performance)
- Benchmark infrastructure already built
- Database setup functions already written
- Performance analysis patterns documented

#### Test Infrastructure
- `test/test_helper.go` - Database setup/teardown
- `test/integration/` - Integration test patterns
- `docker-compose.yml` - Service orchestration

#### Make Commands
- `make test` - Run all tests
- `make test-coverage` - Coverage report
- `make docker-up` - Start services
- `make docker-down` - Stop services

### 11. Sample Implementation (Performance Test)

```go
// test/performance/comparison_test.go
func BenchmarkComparisonGoVsRails(b *testing.B) {
    // Setup test data
    helper := setupBenchmarkDatabase(b)
    defer cleanupBenchmarkDatabase(b, helper)
    
    // Run Go benchmarks
    goTimes := make([]float64, b.N)
    for i := 0; i < b.N; i++ {
        start := time.Now()
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        _, err := repo.GetAllWithLogs(ctx)
        cancel()
        if err != nil {
            b.Fatalf("Go benchmark failed: %v", err)
        }
        goTimes[i] = time.Since(start).Seconds()
    }
    
    // Compare with Rails (if available)
    // railsTimes := fetchRailsTimes()
    
    // Calculate statistics
    goMean := mean(goTimes)
    // railsMean := mean(railsTimes)
    // difference := ((goMean - railsMean) / railsMean) * 100
    
    b.ReportMetric(goMean*1000, "ms/op")  // Convert to milliseconds
}
```

### 12. Documentation Template

```markdown
# Performance Comparison Report

## Summary
- Go API Response Time: X ms (p95: Y ms)
- Rails API Response Time: X ms (p95: Y ms)
- Performance Difference: +Z% (within threshold: YES/NO)
- Memory Usage: X MB (within threshold: YES/NO)

## Details

### methodology
- Tool: wrk/Go benchmark
- Concurrent users: 100
- Duration: 30 seconds
- Dataset: 100 projects, 500 logs

### Results

| Endpoint | Go (ms) | Rails (ms) | Difference |within Threshold| |
|----------|---------|------------|------------|----------------|
| GET /api/v1/projects | X | Y | Z% | YES/NO |
| GET /api/v1/projects/:id | X | Y | Z% | YES/NO |

### Optimizations Applied
- [ ] Query optimization
- [ ] Index usage
- [ ] Connection pooling
- [ ] Caching

### Recommendations
1. Add database index on project_id if missing
2. Consider connection pool size adjustment
3. Profile memory usage under high load
```

### 13. Post-Implementation Tasks

1. Update `docs/performance-comparison.md` with latest results
2. Update `QWEN.md` with performance characteristics
3. Consider adding performance benchmark as CI/CD step
4. Document performance regression detection process
5. Schedule periodic performance re-evaluation

## Final Notes

The primary challenge is the Rails app directory completeness. Begin by verifying Rails app functionality before investing significant time in comparison testing. If Rails app is not functional, focus on comprehensive Go performance testing and establish baseline metrics for future comparison.

The task builds on existing benchmark infrastructure from RDL-029 and should integrate seamlessly with current test patterns.
<!-- SECTION:PLAN:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 All unit tests pass use testing-expert subagent for test execution and verification
- [ ] #2 All integration tests pass use testing-expert subagent for test execution and verification
- [ ] #3 go fmt and go vet pass with no errors
- [ ] #4 Clean Architecture layers properly followed
- [ ] #5 Error responses consistent with existing patterns
- [ ] #6 HTTP status codes correct for response type
- [ ] #7 Database queries optimized with proper indexes
- [ ] #8 Documentation updated in QWEN.md
- [ ] #9 New code paths include error path tests
- [ ] #10 HTTP handlers test both success and error responses
- [ ] #11 Integration tests verify actual database interactions
- [ ] #12 Tests use testing-expert subagent for test execution and verification
<!-- DOD:END -->
