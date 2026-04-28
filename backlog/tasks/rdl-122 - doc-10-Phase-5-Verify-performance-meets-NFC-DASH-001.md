---
id: RDL-122
title: '[doc-10 Phase 5] Verify performance meets NFC-DASH-001'
status: To Do
assignee:
  - thomas
created_date: '2026-04-28 00:30'
updated_date: '2026-04-28 05:53'
labels:
  - performance
  - phase-5
  - backend
dependencies: []
documentation:
  - doc-010
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Run performance tests with production-like data volume (10,000+ logs). Verify response time is < 500ms at p95 percentile. Add database indexes if needed.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Response time < 500ms at p95 with 10,000+ logs
- [ ] #2 Database queries use appropriate indexes
- [ ] #3 Performance test results documented
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task focuses on verifying that the Go API meets the NFC-DASH-001 performance requirement of <500ms p95 response time with production-like data volume (10,000+ logs). The implementation will:

1. **Create large-scale benchmark tests** - Extend existing benchmark infrastructure to test with 10,000+ logs across multiple projects
2. **Measure p95 latency** - Track percentile latencies using the existing `calculatePercentiles` function in `test/performance/comparison_test.go`
3. **Analyze database query performance** - Use EXPLAIN ANALYZE to identify bottlenecks and verify index usage
4. **Add missing database indexes** - If performance tests reveal bottlenecks, add targeted indexes based on query analysis
5. **Document results** - Create comprehensive performance documentation with test methodology and results

**Architecture Decisions:**
- Reuse existing benchmark infrastructure (`test/performance/` directory) to maintain consistency
- Follow the pattern established in `dashboard_benchmark_test.go` for percentile calculations
- Use Go's native benchmarking framework (`testing.B`) for consistency with existing tests
- Create a dedicated large-scale dataset benchmark file to avoid cluttering existing benchmarks

**Why this approach:**
- Leverages existing, proven testing patterns in the codebase
- Provides clear separation between small-scale (current) and large-scale (10,000+ logs) benchmarks
- Enables automated regression detection through threshold verification
- Aligns with NFC-DASH-001 acceptance criteria

### 2. Files to Modify

**New Files to Create:**
- `test/performance/large_scale_benchmark_test.go` - Main benchmark file for 10,000+ logs testing with:
  - `BenchmarkLargeScaleGetAllWithLogs` - Tests GetAllWithLogs with 10,000+ logs
  - `BenchmarkLargeScaleGetWithLogs` - Tests GetWithLogs with large dataset
  - `BenchmarkLargeScaleDashboardDay` - Tests dashboard/day endpoint with large dataset
  - `BenchmarkLargeScaleConcurrent` - Concurrent load test with 10,000+ logs
  - Helper functions for large dataset setup (100+ projects, 10,000+ logs)

- `docs/performance/large-scale-benchmarks.md` - Documentation covering:
  - Test methodology and data volume
  - Performance results (p50, p95, p99 latencies)
  - Database query analysis (EXPLAIN ANALYZE output)
  - Index recommendations and rationale
  - Threshold verification results

**Existing Files to Read/Modify:**
- `test/performance/baseline.go` - May need to add large-scale benchmark entries
- `test/performance/comparison_test.go` - Reuse `calculatePercentiles` and `VerifyPerformanceThresholds` functions
- `test/test_helper.go` - May need to extend for faster large dataset creation
- `docs/benchmarks.md` - Update with large-scale benchmark documentation
- `Makefile` - Add `benchmark-large-scale` target

**Database Schema Files (if indexes needed):**
- `docs/database.sql` - Add CREATE INDEX statements for documentation
- `test/setup_comparison_data.sh` - Add index creation if needed

### 3. Dependencies

**Prerequisites:**
- PostgreSQL running and accessible (existing requirement)
- Test database configured (`.env.test` file)
- Existing benchmark infrastructure in place (already completed per RDL-090)

**Existing Tasks Completed:**
- RDL-090: Performance testing infrastructure for dashboard endpoints (provides benchmark patterns)
- RDL-028: Database indexes for logs table (provides `index_logs_on_project_id_and_data_desc`)
- RDL-029: Initial performance benchmark verification

**Blocking Issues:**
- None identified - all prerequisites are in place

**Setup Steps:**
1. Ensure PostgreSQL is running with sufficient resources for 10,000+ log dataset
2. Verify existing benchmark tests pass before adding large-scale tests
3. Confirm test database can handle large dataset creation within reasonable time

### 4. Code Patterns

**Follow Existing Patterns:**
- **Benchmark Structure**: Match `dashboard_benchmark_test.go` pattern with warm-up phase, benchmark loop, and percentile calculations
- **Percentile Calculation**: Reuse `calculatePercentiles()` from `comparison_test.go`
- **Threshold Verification**: Use `VerifyPerformanceThresholds()` for automated pass/fail
- **Connection Pool Verification**: Include `verifyConnectionPool()` calls to ensure efficient connection reuse
- **Error Handling**: Follow existing pattern of `b.Fatalf()` for critical errors, `b.Errorf()` for threshold failures

**Naming Conventions:**
- Benchmark functions: `BenchmarkLargeScale<Operation>` (e.g., `BenchmarkLargeScaleGetAllWithLogs`)
- Helper functions: `setupLargeScaleBenchmark()`, `cleanupLargeScaleBenchmark()`
- Documentation files: `docs/performance/<topic>.md`

**Integration Patterns:**
- Import existing performance utilities from `test/performance/` package
- Use `test.TestHelper` for database setup/cleanup
- Follow Clean Architecture - benchmarks test repository and handler layers

**Code Quality:**
- Run `go fmt ./test/performance/...` after implementation
- Run `go vet ./test/performance/...` for static analysis
- Ensure all benchmarks include `b.ReportAllocs()` for memory profiling

### 5. Testing Strategy

**Benchmark Tests (test/performance/large_scale_benchmark_test.go):**

1. **Dataset Setup** - Create realistic test data:
   - 100+ projects with varying total_page values
   - 10,000+ logs distributed across projects
   - Realistic date ranges (6+ months of data)
   - Varied wday values to test weekday calculations

2. **Benchmark Functions:**
   - `BenchmarkLargeScaleGetAllWithLogs` - Full dataset query (target: <500ms p95)
   - `BenchmarkLargeScaleGetWithLogs` - Single project with many logs (target: <100ms p95)
   - `BenchmarkLargeScaleDashboardDay` - Dashboard endpoint (target: <500ms p95)
   - `BenchmarkLargeScaleConcurrent50` - 50 concurrent users (target: >100 QPS, <1% errors)
   - `BenchmarkLargeScaleConcurrent100` - 100 concurrent users (target: >100 QPS, <1% errors)

3. **Edge Cases to Cover:**
   - Empty result sets (no logs for date range)
   - Single project with all 10,000 logs
   - Even distribution across 100 projects
   - Extreme date ranges (all data vs. single day)

4. **Verification Steps:**
   - Run benchmarks 3 times (`-count=3`) for statistical stability
   - Verify p95 latency <500ms for all endpoints
   - Verify connection pool efficiency (acquire count vs. total connections)
   - Verify error rate <1% under load

**Database Query Analysis:**
- Run `EXPLAIN ANALYZE` on critical queries with large dataset
- Document query plans and index usage
- Identify missing indexes if query time >100ms

**Acceptance Criteria Verification:**
- [ ] Response time <500ms at p95 with 10,000+ logs (all endpoints)
- [ ] Database queries use appropriate indexes (verified via EXPLAIN ANALYZE)
- [ ] Performance test results documented in `docs/performance/large-scale-benchmarks.md`

### 6. Risks and Considerations

**Potential Risks:**

1. **Database Performance** - Large dataset (10,000+ logs) may cause slow test setup:
   - **Mitigation**: Use batch inserts, optimize setup queries, consider parallel data creation
   - **Fallback**: If setup takes >5 minutes, reduce dataset to 5,000 logs with documentation

2. **Memory Usage** - Loading 10,000+ logs into memory during benchmarks:
   - **Mitigation**: Monitor memory with `b.ReportAllocs()`, consider streaming results if needed
   - **Fallback**: Implement cursor-based pagination for repository methods

3. **Missing Indexes** - Current indexes may not be sufficient for large datasets:
   - **Mitigation**: Run EXPLAIN ANALYZE before implementation to identify gaps
   - **Plan**: Prepare index creation scripts for common query patterns:
     - `CREATE INDEX CONCURRENTLY ON logs(data::date)` for date range queries
     - `CREATE INDEX CONCURRENTLY ON logs(project_id, data DESC)` (already exists, verify usage)

4. **Test Environment Variability** - Benchmark results vary by hardware:
   - **Mitigation**: Document hardware specs in results, use relative thresholds
   - **Consideration**: Establish baseline measurements for CI/CD integration

5. **Connection Pool Exhaustion** - High concurrency may exhaust pool:
   - **Mitigation**: Verify pool configuration, increase max connections if needed
   - **Monitor**: Track `AcquireCount` vs `TotalConns` in benchmark results

**Trade-offs:**
- **Setup Time vs. Dataset Size**: Larger datasets provide more realistic testing but increase setup time
- **Precision vs. Speed**: Running benchmarks 3x for statistical stability increases total test time
- **Index Overhead**: Additional indexes improve read performance but add write overhead

**Deployment Considerations:**
- Index creation with `CREATE INDEX CONCURRENTLY` for production (PostgreSQL 12+)
- Schedule index creation during maintenance windows if using standard `CREATE INDEX`
- Monitor query performance after index deployment

**Success Metrics:**
- All benchmarks pass with p95 <500ms
- No N+1 query patterns detected
- Connection pool utilization <80% under load
- Error rate <1% during concurrent tests
<!-- SECTION:PLAN:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 All unit tests pass
- [ ] #2 All integration tests pass execution and verification
- [ ] #3 go fmt and go vet pass with no errors
- [ ] #4 Clean Architecture layers properly followed
- [ ] #5 Error responses consistent with existing patterns
- [ ] #6 HTTP status codes correct for response type
- [ ] #7 Documentation updated in QWEN.md
- [ ] #8 New code paths include error path tests
- [ ] #9 HTTP handlers test both success and error responses
- [ ] #10 Integration tests verify actual database interactions
<!-- DOD:END -->
