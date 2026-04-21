---
id: RDL-079
title: >-
  [doc-008 Phase 1] Create DashboardRepository interface and PostgreSQL
  implementation
status: To Do
assignee:
  - thomas
created_date: '2026-04-21 15:49'
updated_date: '2026-04-21 16:33'
labels:
  - phase-1
  - repository
  - database
dependencies: []
references:
  - REQ-DASH-002
  - AC-DASH-002
  - 'Decision 6: Repository Pattern Extension'
documentation:
  - doc-008
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Define internal/repository/dashboard_repository.go with all dashboard query methods and implement in internal/adapter/postgres/dashboard_repository.go. Include GetDailyStats, GetProjectAggregates, GetFaultsByDateRange, GetWeekdayFaults methods using pgx for efficient database access.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Interface defines all required dashboard query methods
- [ ] #2 PostgreSQL implementation uses pgx for efficient queries
- [ ] #3 Connection pooling configured correctly
- [ ] #4 Unit tests verify each repository method independently
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task involves creating a **Dashboard Repository** layer following Clean Architecture principles. The approach will mirror the existing Project and Log repository patterns while extending them with dashboard-specific aggregation methods.

**Architecture Decision**: Extend the existing repository interface pattern rather than creating a separate service layer for Phase 2, as suggested by Decision 6 in doc-008. This maintains architectural simplicity while providing the necessary data access methods.

**Key Design Considerations**:
- Use `pgx` directly for efficient queries (already in use throughout the codebase)
- Follow existing naming conventions (`DashboardRepository` interface + `DashboardRepositoryImpl` implementation)
- Maintain consistent error handling patterns
- Use context timeouts matching existing code (15 seconds default)

**Methods to Implement**:
| Method | Purpose | Query Type |
|--------|---------|------------|
| `GetDailyStats` | Daily page statistics with weekday breakdown | Aggregation with GROUP BY |
| `GetProjectAggregates` | Project-level sums and counts | JOIN aggregation |
| `GetFaultsByDateRange` | Fault counting within date range | COUNT with WHERE |
| `GetWeekdayFaults` | Fault distribution by weekday | Aggregation with EXTRACT |

### 2. Files to Modify

#### New Files to Create:

```
internal/
├── repository/
│   └── dashboard_repository.go          # Interface definition
├── adapter/
│   └── postgres/
│       └── dashboard_repository.go      # Implementation
└── domain/
    └── dto/
        └── dashboard_response.go         # Response DTOs (may already exist)
```

#### Modified Files:

```
internal/api/v1/routes.go                # Add dashboard route registrations
cmd/server.go                            # Wire up dashboard repository (if needed)
```

### 3. Dependencies

**Prerequisites for Implementation**:
- [x] Go 1.25.7 environment ready
- [x] PostgreSQL connection pool available (`pgxpool.Pool`)
- [x] Existing `repository` and `postgres` packages accessible
- [x] Domain models (`models.Project`, `models.Log`) available

**External Dependencies** (already in go.mod):
- `github.com/jackc/pgx/v5` - Database driver
- `github.com/jackc/pgx/v5/pgxpool` - Connection pooling

**No New Dependencies Required**

### 4. Code Patterns

**Pattern 1: Repository Interface Definition**
```go
// internal/repository/dashboard_repository.go
package repository

import (
    "context"
    "time"
)

type DashboardRepository interface {
    GetDailyStats(ctx context.Context, date time.Time) (*DailyStats, error)
    GetProjectAggregates(ctx context.Context) ([]*ProjectAggregate, error)
    GetFaultsByDateRange(ctx context.Context, start, end time.Time) (int, error)
    GetWeekdayFaults(ctx context.Context, start, end time.Time) (map[int]int, error)
}
```

**Pattern 2: Implementation Structure**
```go
// internal/adapter/postgres/dashboard_repository.go
package postgres

import (
    "context"
    "fmt"
    "time"
    
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

const dashboardContextTimeout = 15 * time.Second

type DashboardRepositoryImpl struct {
    pool *pgxpool.Pool
}

func NewDashboardRepositoryImpl(pool *pgxpool.Pool) *DashboardRepositoryImpl {
    return &DashboardRepositoryImpl{pool: pool}
}
```

**Pattern 3: Query Execution with Error Wrapping**
```go
func (r *DashboardRepositoryImpl) GetDailyStats(ctx context.Context, date time.Time) (*DailyStats, error) {
    ctx, cancel := context.WithTimeout(ctx, dashboardContextTimeout)
    defer cancel()
    
    query := `
        SELECT 
            COALESCE(SUM(start_page), 0) as total_pages,
            COUNT(*) as log_count
        FROM logs
        WHERE data::date = $1
    `
    
    var stats DailyStats
    err := r.pool.QueryRow(ctx, query, date).Scan(
        &stats.TotalPages,
        &stats.LogCount,
    )
    if err != nil {
        if err == pgx.ErrNoRows {
            return &DailyStats{}, nil  // Return zero values instead of error
        }
        return nil, fmt.Errorf("failed to get daily stats: %w", err)
    }
    
    return &stats, nil
}
```

**Pattern 4: Map-Based Return for Weekday Data**
```go
func (r *DashboardRepositoryImpl) GetWeekdayFaults(ctx context.Context, start, end time.Time) (map[int]int, error) {
    ctx, cancel := context.WithTimeout(ctx, dashboardContextTimeout)
    defer cancel()
    
    query := `
        SELECT 
            EXTRACT(DOW FROM data)::int as weekday,
            COUNT(*) as fault_count
        FROM logs
        WHERE data BETWEEN $1 AND $2
        GROUP BY EXTRACT(DOW FROM data)
        ORDER BY weekday
    `
    
    rows, err := r.pool.Query(ctx, query, start, end)
    if err != nil {
        return nil, fmt.Errorf("failed to query weekday faults: %w", err)
    }
    defer rows.Close()
    
    result := make(map[int]int)
    for rows.Next() {
        var weekday int
        var count int
        if err := rows.Scan(&weekday, &count); err != nil {
            return nil, fmt.Errorf("failed to scan weekday fault: %w", err)
        }
        result[weekday] = count
    }
    
    // Ensure all 7 days are present (0-6)
    for i := 0; i < 7; i++ {
        if _, exists := result[i]; !exists {
            result[i] = 0
        }
    }
    
    return result, nil
}
```

### 5. Testing Strategy

**Unit Tests Structure**:
```go
// test/dashboard_repository_test.go
package test

import (
    "context"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    
    "go-reading-log-api-next/internal/adapter/postgres"
    "go-reading-log-api-next/internal/repository"
)

func TestDashboardRepository_GetDailyStats(t *testing.T) {
    // Setup test database
    helper, err := SetupTestDB()
    require.NoError(t, err)
    defer helper.Close()
    
    repo := postgres.NewDashboardRepositoryImpl(helper.Pool)
    
    // Create test data
    testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
    err = createTestLogs(helper.Pool, []testLog{
        {Data: testDate, StartPage: 10, EndPage: 20},
        {Data: testDate, StartPage: 20, EndPage: 30},
    })
    require.NoError(t, err)
    
    // Execute
    stats, err := repo.GetDailyStats(context.Background(), testDate)
    
    // Verify
    assert.NoError(t, err)
    assert.Equal(t, 20, stats.TotalPages)  // (20-10) + (30-20) = 20
    assert.Equal(t, 2, stats.LogCount)
}

func TestDashboardRepository_GetWeekdayFaults(t *testing.T) {
    helper, err := SetupTestDB()
    require.NoError(t, err)
    defer helper.Close()
    
    repo := postgres.NewDashboardRepositoryImpl(helper.Pool)
    
    // Create test data spanning multiple weekdays
    startOfWeek := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC) // Monday
    endOfWeek := startOfWeek.AddDate(0, 0, 7)
    
    err = createTestLogs(helper.Pool, []testLog{
        {Data: startOfWeek, StartPage: 0, EndPage: 10}, // Monday
        {Data: startOfWeek.AddDate(0, 0, 1), StartPage: 0, EndPage: 10}, // Tuesday
        {Data: startOfWeek.AddDate(0, 0, 1), StartPage: 0, EndPage: 10}, // Tuesday again
    })
    require.NoError(t, err)
    
    stats, err := repo.GetWeekdayFaults(context.Background(), startOfWeek, endOfWeek)
    
    assert.NoError(t, err)
    assert.Equal(t, 1, stats[1])  // Monday
    assert.Equal(t, 2, stats[2])  // Tuesday
    assert.Equal(t, 0, stats[0])  // Sunday (no data)
}
```

**Integration Test Approach**:
- Use `test.SetupTestDB()` for database setup/teardown
- Create test fixtures with known data
- Verify calculations match expected values
- Test edge cases (empty results, null values)

### 6. Risks and Considerations

**Risk 1: Timestamp Column Type Mismatch**
- **Issue**: The `logs.data` column is VARCHAR in the schema but queried as date
- **Mitigation**: Use `data::date` cast in SQL queries (consistent with existing code)
- **Status**: LOW - Already handled by existing `parseLogDate` pattern

**Risk 2: NULL Handling in Aggregations**
- **Issue**: SUM returns NULL for empty sets, not 0
- **Mitigation**: Use `COALESCE(SUM(...), 0)` in all aggregation queries
- **Status**: MEDIUM - Requires careful review of all COUNT/SUM operations

**Risk 3: Weekday Index Offset**
- **Issue**: PostgreSQL EXTRACT(DOW) returns 0-6 (Sunday-Saturday), which may differ from Go's time.Weekday
- **Mitigation**: Document the mapping clearly; use explicit conversion
- **Status**: LOW - Well-documented behavior, easy to verify

**Risk 4: Performance with Large Datasets**
- **Issue**: Aggregation queries on large tables could be slow
- **Mitigation**: Ensure indexes exist on `logs.data` column (already created per RDL-028)
- **Status**: LOW - Index already exists for date-based queries

**Consideration: Consistency with Rails API**
- Must match Rails calculation logic exactly (per PRD Decision 4)
- Verify fault counting includes all faults (not just closed)
- Verify date range calculations match Rails `Time.zone` behavior
- **Status**: HIGH - Critical for feature parity

**Consideration: Error Response Format**
- Must return consistent error format matching existing handlers
- Use standard `fmt.Errorf("failed to ...: %w", err)` wrapping
- Return HTTP 500 for unexpected errors, 404 for not found where appropriate
- **Status**: MEDIUM - Needs alignment with handler layer

**Implementation Checklist**:
- [ ] Create `internal/repository/dashboard_repository.go` interface
- [ ] Create `internal/adapter/postgres/dashboard_repository.go` implementation
- [ ] Implement `GetDailyStats` method
- [ ] Implement `GetProjectAggregates` method  
- [ ] Implement `GetFaultsByDateRange` method
- [ ] Implement `GetWeekdayFaults` method
- [ ] Add unit tests for each method
- [ ] Add integration tests with test database
- [ ] Update `internal/api/v1/routes.go` to register dashboard routes
- [ ] Run `go fmt` and `go vet` to verify code quality
- [ ] Verify all tests pass with `go test ./...`
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Progress - Task RDL-079

### Status: In Progress

I've analyzed the codebase structure and understand the patterns to follow:

**Existing Patterns Identified:**
1. `ProjectRepository` interface in `internal/repository/project_repository.go`
2. `ProjectRepositoryImpl` implementation in `internal/adapter/postgres/project_repository.go`
3. `LogRepository` interface in `internal/repository/log_repository.go`
4. `LogRepositoryImpl` implementation in `internal/adapter/postgres/log_repository.go`

**Key Design Decisions:**
- Use 15-second context timeout (matching existing code)
- Follow JSON:API response patterns via DTOs
- Use pgxpool for connection pooling (already configured)
- Implement COALESCE for NULL handling in aggregations

### Files to Create:

```
internal/
├── repository/
│   └── dashboard_repository.go          # Interface definition
├── adapter/
│   └── postgres/
│       └── dashboard_repository.go      # Implementation
└── domain/
    └── dto/
        └── dashboard_response.go         # Response DTOs
```

### Implementation Plan:
1. Create `dashboard_repository.go` interface with all required methods
2. Create `dashboard_response.go` DTOs for response structures
3. Implement `dashboard_repository.go` using pgx
4. Add unit tests for each method
5. Add integration tests with test database
6. Verify code quality with go fmt and go vet

**Next Step:** Creating the repository interface definition.
<!-- SECTION:NOTES:END -->

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
