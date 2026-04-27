---
id: RDL-103
title: Fix test
status: To Do
assignee:
  - thomas
created_date: '2026-04-27 10:52'
updated_date: '2026-04-27 11:29'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
fix these test

=== RUN   TestErrorScenarios
=== RUN   TestErrorScenarios/Day_Endpoint_-_Invalid_Date
    error_scenarios_test.go:86: Unknown endpoint: /v1/dashboard/day.json?date=invalid
=== RUN   TestErrorScenarios/Last_Days_-_Invalid_Type
    error_scenarios_test.go:86: Unknown endpoint: /v1/dashboard/last_days.json?type=99
panic: test timed out after 2s
	running tests:
		TestErrorScenarios (2s)
		TestErrorScenarios/Last_Days_-_Invalid_Type (2s)

goroutine 46 [running]:
testing.(*M).startAlarm.func1()
	/usr/lib/go/src/testing/testing.go:2802 +0x34b
created by time.goFunc
	/usr/lib/go/src/time/sleep.go:215 +0x2d

goroutine 1 [chan receive]:
testing.(*T).Run(0x31275e01c6c8, {0xa695d8?, 0x31275df6fb30?}, 0xa92c68)
	/usr/lib/go/src/testing/testing.go:2109 +0x4e5
testing.runTests.func1(0x31275e01c6c8)
	/usr/lib/go/src/testing/testing.go:2585 +0x37
testing.tRunner(0x31275e01c6c8, 0x31275df6fc58)
	/usr/lib/go/src/testing/testing.go:2036 +0xea
testing.runTests({0xa6c9cb, 0x17}, {0xa7a9a1, 0x28}, 0x31275de86180, {0xf8e780, 0x30, 0x30}, {0xc273edc006500fc5, 0x77484c23, ...})
	/usr/lib/go/src/testing/testing.go:2583 +0x505
testing.(*M).Run(0x31275df46e60)
	/usr/lib/go/src/testing/testing.go:2443 +0x6ac
main.main()
	_testmain.go:140 +0x9b

goroutine 8 [chan receive]:
testing.(*T).Run(0x31275e01c908, {0xa6d735?, 0x31275defe100?}, 0x31275dfc0500)
	/usr/lib/go/src/testing/testing.go:2109 +0x4e5
go-reading-log-api-next/test/integration.RunErrorScenarios(0x31275e01c908, {0xf8cee0, 0x5, 0x31275df07760?})
	/home/danilo/scripts/github/go-reading-log-api-next/test/integration/error_scenarios_test.go:36 +0xd1
go-reading-log-api-next/test/integration.TestErrorScenarios(0x31275e01c908)
	/home/danilo/scripts/github/go-reading-log-api-next/test/integration/error_scenarios_test.go:259 +0x6c
testing.tRunner(0x31275e01c908, 0xa92c68)
	/usr/lib/go/src/testing/testing.go:2036 +0xea
created by testing.(*T).Run in goroutine 1
	/usr/lib/go/src/testing/testing.go:2101 +0x4c5

goroutine 61 [sync.WaitGroup.Wait]:
sync.runtime_SemacquireWaitGroup(0x31275deff660?, 0x80?)
	/usr/lib/go/src/runtime/sema.go:114 +0x2e
sync.(*WaitGroup).Wait(0x31275de89b10)
	/usr/lib/go/src/sync/waitgroup.go:206 +0x85
go-reading-log-api-next/test.cleanupOrphanedDatabasesConcurrent(0x31275e0d8000, {0x31275dfac2d0, 0x2e})
	/home/danilo/scripts/github/go-reading-log-api-next/test/test_helper.go:551 +0x5e6
go-reading-log-api-next/test.(*TestHelper).Close.func1()
	/home/danilo/scripts/github/go-reading-log-api-next/test/test_helper.go:377 +0x1bf
go-reading-log-api-next/test.(*TestHelper).Close(0x31275df42780?)
	/home/danilo/scripts/github/go-reading-log-api-next/test/test_helper.go:403 +0x8b
runtime.Goexit()
	/usr/lib/go/src/runtime/panic.go:694 +0x5e
testing.(*common).FailNow(0x31275dfc8908)
	/usr/lib/go/src/testing/testing.go:1022 +0x4a
testing.(*common).Fatalf(0x31275dfc8908, {0xa6ab3d?, 0xa62282?}, {0x31275df71f40?, 0xa775dc?, 0x24?})
	/usr/lib/go/src/testing/testing.go:1228 +0x59
go-reading-log-api-next/test/integration.RunErrorScenarios.func1(0x31275dfc8908)
	/home/danilo/scripts/github/go-reading-log-api-next/test/integration/error_scenarios_test.go:86 +0x535
testing.tRunner(0x31275dfc8908, 0x31275dfc0500)
	/usr/lib/go/src/testing/testing.go:2036 +0xea
created by testing.(*T).Run in goroutine 8
	/usr/lib/go/src/testing/testing.go:2101 +0x4c5

goroutine 26 [IO wait]:
internal/poll.runtime_pollWait(0x7fdcbfa60a00, 0x72)
	/usr/lib/go/src/runtime/netpoll.go:351 +0x85
internal/poll.(*pollDesc).wait(0x31275e019080?, 0x31275e20c000?, 0x0)
	/usr/lib/go/src/internal/poll/fd_poll_runtime.go:84 +0x27
internal/poll.(*pollDesc).waitRead(...)
	/usr/lib/go/src/internal/poll/fd_poll_runtime.go:89
internal/poll.(*FD).Read(0x31275e019080, {0x31275e20c000, 0x2000, 0x2000})
	/usr/lib/go/src/internal/poll/fd_unix.go:165 +0x2ae
net.(*netFD).Read(0x31275e019080, {0x31275e20c000?, 0x31275e1259e0?, 0x4d255d?})
	/usr/lib/go/src/net/fd_posix.go:68 +0x25
net.(*conn).Read(0x31275de76588, {0x31275e20c000?, 0x31275deff360?, 0x31275dff1768?})
	/usr/lib/go/src/net/net.go:196 +0x45
github.com/jackc/pgx/v5/pgconn/internal/bgreader.(*BGReader).Read(0x31275df2b8c0, {0x31275e20c000, 0x2000, 0x2000})
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgconn/internal/bgreader/bgreader.go:100 +0xcb
io.ReadAtLeast({0xa9b7a0, 0x31275df2b8c0}, {0x31275e20c000, 0x2000, 0x2000}, 0x5)
	/usr/lib/go/src/io/io.go:335 +0x8e
github.com/jackc/pgx/v5/pgproto3.(*chunkReader).Next(0x31275e20a030, 0x5)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgproto3/chunkreader.go:80 +0x289
github.com/jackc/pgx/v5/pgproto3.(*Frontend).Receive(0x31275e163688)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgproto3/frontend.go:309 +0x3c
github.com/jackc/pgx/v5/pgconn.(*PgConn).peekMessage(0x31275e027708)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgconn/pgconn.go:595 +0x14b
github.com/jackc/pgx/v5/pgconn.(*MultiResultReader).NextResult(0x31275e027838)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgconn/pgconn.go:1552 +0x4e
github.com/jackc/pgx/v5.(*Conn).execSimpleProtocol(0x31275df84b40?, {0xaa2350?, 0x31275e148b60?}, {0x31275debc4b0?, 0x31275dff1b68?}, {0x0?, 0x424cdc?, 0x31275dff1bc0?})
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/conn.go:589 +0xb0
github.com/jackc/pgx/v5.(*Conn).exec(0x31275e208500, {0xaa2350, 0x31275e148b60}, {0x31275debc4b0, 0x46}, {0x0?, 0x87f7ed?, 0x31275df2b740?})
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/conn.go:574 +0x8ac
github.com/jackc/pgx/v5.(*Conn).Exec(0x31275e208500, {0xaa2350?, 0x31275e148b60?}, {0x31275debc4b0, 0x46}, {0x0, 0x0, 0x0})
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/conn.go:481 +0x114
github.com/jackc/pgx/v5/pgxpool.(*Conn).Exec(0x31275e0402a0?, {0xaa2350?, 0x31275e148b60?}, {0x31275debc4b0?, 0x31275df4d520?}, {0x0?, 0x4f4ada?, 0x31275df4d520?})
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/conn.go:87 +0x3c
github.com/jackc/pgx/v5/pgxpool.(*Pool).Exec(0xa6f227?, {0xaa2350, 0x31275e148b60}, {0x31275debc4b0, 0x46}, {0x0, 0x0, 0x0})
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:726 +0xf7
go-reading-log-api-next/test.cleanupOrphanedDatabasesConcurrent.func1({0x31275deae810, 0x2e})
	/home/danilo/scripts/github/go-reading-log-api-next/test/test_helper.go:541 +0x566
created by go-reading-log-api-next/test.cleanupOrphanedDatabasesConcurrent in goroutine 61
	/home/danilo/scripts/github/go-reading-log-api-next/test/test_helper.go:511 +0x47b

goroutine 37 [select]:
github.com/jackc/pgx/v5/pgxpool.(*Pool).backgroundHealthCheck(0x31275e0d8000)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:495 +0xc8
github.com/jackc/pgx/v5/pgxpool.NewWithConfig.func5()
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:338 +0x3d
created by github.com/jackc/pgx/v5/pgxpool.NewWithConfig in goroutine 61
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:335 +0x43e

goroutine 27 [select]:
github.com/jackc/pgx/v5/pgxpool.(*Pool).backgroundHealthCheck(0x31275e0402a0)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:495 +0xc8
github.com/jackc/pgx/v5/pgxpool.NewWithConfig.func5()
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:338 +0x3d
created by github.com/jackc/pgx/v5/pgxpool.NewWithConfig in goroutine 26
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:335 +0x43e
FAIL	go-reading-log-api-next/test/integration	2.014s
testing: warning: no tests to run
PASS
ok  	go-reading-log-api-next/test/performance	(cached) [no tests to run]
?   	go-reading-log-api-next/test/testutil	[no test files]
=== RUN   TestDashboardRepository_GetDailyStats
--- PASS: TestDashboardRepository_GetDailyStats (0.31s)
=== RUN   TestDashboardRepository_GetDailyStats_EmptyDate
panic: test timed out after 2s
	running tests:
		TestDashboardRepository_GetDailyStats_EmptyDate (2s)

goroutine 83 [running]:
testing.(*M).startAlarm.func1()
	/usr/lib/go/src/testing/testing.go:2802 +0x34b
created by time.goFunc
	/usr/lib/go/src/time/sleep.go:215 +0x2d

goroutine 1 [chan receive]:
testing.(*T).Run(0x27d9f824c6c8, {0x9e641f?, 0x27d9f8259b30?}, 0x9f80c0)
	/usr/lib/go/src/testing/testing.go:2109 +0x4e5
testing.runTests.func1(0x27d9f824c6c8)
	/usr/lib/go/src/testing/testing.go:2585 +0x37
testing.tRunner(0x27d9f824c6c8, 0x27d9f8259c58)
	/usr/lib/go/src/testing/testing.go:2036 +0xea
testing.runTests({0x9d5d2e, 0x17}, {0x9dd368, 0x21}, 0x27d9f80c4180, {0xe95620, 0x7e, 0x7e}, {0xc273edc006b41208, 0x773f76c8, ...})
	/usr/lib/go/src/testing/testing.go:2583 +0x505
testing.(*M).Run(0x27d9f8204e60)
	/usr/lib/go/src/testing/testing.go:2443 +0x6ac
main.main()
	_testmain.go:296 +0x9b

goroutine 24 [sync.WaitGroup.Wait]:
sync.runtime_SemacquireWaitGroup(0x27d9f84104a0?, 0xa0?)
	/usr/lib/go/src/runtime/sema.go:114 +0x2e
sync.(*WaitGroup).Wait(0x27d9f840a580)
	/usr/lib/go/src/sync/waitgroup.go:206 +0x85
go-reading-log-api-next/test.cleanupOrphanedDatabasesConcurrent(0x27d9f81c21c0, {0x27d9f81ae1b0, 0x2e})
	/home/danilo/scripts/github/go-reading-log-api-next/test/test_helper.go:551 +0x5e6
go-reading-log-api-next/test.(*TestHelper).Close.func1()
	/home/danilo/scripts/github/go-reading-log-api-next/test/test_helper.go:377 +0x1bf
go-reading-log-api-next/test.(*TestHelper).Close(0x27d9f81944d0?)
	/home/danilo/scripts/github/go-reading-log-api-next/test/test_helper.go:403 +0x8b
go-reading-log-api-next/test/unit.TestDashboardRepository_GetDailyStats_EmptyDate(0x27d9f81b86c8)
	/home/danilo/scripts/github/go-reading-log-api-next/test/unit/dashboard_repository_test.go:86 +0x1d4
testing.tRunner(0x27d9f81b86c8, 0x9f80c0)
	/usr/lib/go/src/testing/testing.go:2036 +0xea
created by testing.(*T).Run in goroutine 1
	/usr/lib/go/src/testing/testing.go:2101 +0x4c5

goroutine 30 [select]:
github.com/jackc/pgx/v5/pgxpool.(*Pool).backgroundHealthCheck(0x27d9f81c21c0)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:495 +0xc8
github.com/jackc/pgx/v5/pgxpool.NewWithConfig.func5()
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:338 +0x3d
created by github.com/jackc/pgx/v5/pgxpool.NewWithConfig in goroutine 24
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:335 +0x43e

goroutine 73 [IO wait]:
internal/poll.runtime_pollWait(0x7fb1aac99c00, 0x72)
	/usr/lib/go/src/runtime/netpoll.go:351 +0x85
internal/poll.(*pollDesc).wait(0x27d9f8452380?, 0x27d9f8502000?, 0x0)
	/usr/lib/go/src/internal/poll/fd_poll_runtime.go:84 +0x27
internal/poll.(*pollDesc).waitRead(...)
	/usr/lib/go/src/internal/poll/fd_poll_runtime.go:89
internal/poll.(*FD).Read(0x27d9f8452380, {0x27d9f8502000, 0x2000, 0x2000})
	/usr/lib/go/src/internal/poll/fd_unix.go:165 +0x2ae
net.(*netFD).Read(0x27d9f8452380, {0x27d9f8502000?, 0x27d9f8416d20?, 0x4cff7d?})
	/usr/lib/go/src/net/fd_posix.go:68 +0x25
net.(*conn).Read(0x27d9f84140d8, {0x27d9f8502000?, 0x27d9f84101c0?, 0x27d9f8255768?})
	/usr/lib/go/src/net/net.go:196 +0x45
github.com/jackc/pgx/v5/pgconn/internal/bgreader.(*BGReader).Read(0x27d9f84583c0, {0x27d9f8502000, 0x2000, 0x2000})
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgconn/internal/bgreader/bgreader.go:100 +0xcb
io.ReadAtLeast({0x9ff2e0, 0x27d9f84583c0}, {0x27d9f8502000, 0x2000, 0x2000}, 0x5)
	/usr/lib/go/src/io/io.go:335 +0x8e
github.com/jackc/pgx/v5/pgproto3.(*chunkReader).Next(0x27d9f8417260, 0x5)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgproto3/chunkreader.go:80 +0x289
github.com/jackc/pgx/v5/pgproto3.(*Frontend).Receive(0x27d9f8435688)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgproto3/frontend.go:309 +0x3c
github.com/jackc/pgx/v5/pgconn.(*PgConn).peekMessage(0x27d9f845c588)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgconn/pgconn.go:595 +0x14b
github.com/jackc/pgx/v5/pgconn.(*MultiResultReader).NextResult(0x27d9f845c6b8)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgconn/pgconn.go:1552 +0x4e
github.com/jackc/pgx/v5.(*Conn).execSimpleProtocol(0x27d9f82012c0?, {0xa05c98?, 0x27d9f84390a0?}, {0x27d9f8464050?, 0x27d9f8255b68?}, {0x0?, 0x424c3c?, 0x27d9f8255bc0?})
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/conn.go:589 +0xb0
github.com/jackc/pgx/v5.(*Conn).exec(0x27d9f844adc0, {0xa05c98, 0x27d9f84390a0}, {0x27d9f8464050, 0x46}, {0x0?, 0x7dc3ed?, 0x27d9f8458240?})
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/conn.go:574 +0x8ac
github.com/jackc/pgx/v5.(*Conn).Exec(0x27d9f844adc0, {0xa05c98?, 0x27d9f84390a0?}, {0x27d9f8464050, 0x46}, {0x0, 0x0, 0x0})
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/conn.go:481 +0x114
github.com/jackc/pgx/v5/pgxpool.(*Conn).Exec(0x27d9f84540e0?, {0xa05c98?, 0x27d9f84390a0?}, {0x27d9f8464050?, 0x27d9f84401a0?}, {0x0?, 0x4f1c7a?, 0x27d9f84401a0?})
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/conn.go:87 +0x3c
github.com/jackc/pgx/v5/pgxpool.(*Pool).Exec(0x9d806f?, {0xa05c98, 0x27d9f84390a0}, {0x27d9f8464050, 0x46}, {0x0, 0x0, 0x0})
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:726 +0xf7
go-reading-log-api-next/test.cleanupOrphanedDatabasesConcurrent.func1({0x27d9f8444120, 0x2e})
	/home/danilo/scripts/github/go-reading-log-api-next/test/test_helper.go:541 +0x566
created by go-reading-log-api-next/test.cleanupOrphanedDatabasesConcurrent in goroutine 24
	/home/danilo/scripts/github/go-reading-log-api-next/test/test_helper.go:511 +0x47b

goroutine 74 [select]:
github.com/jackc/pgx/v5/pgxpool.(*Pool).backgroundHealthCheck(0x27d9f84540e0)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:495 +0xc8
github.com/jackc/pgx/v5/pgxpool.NewWithConfig.func5()
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:338 +0x3d
created by github.com/jackc/pgx/v5/pgxpool.NewWithConfig in goroutine 73
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:335 +0x43e
FAIL	go-reading-log-api-next/test/unit	2.012s
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task involves fixing two critical test failures that cause test timeouts and panics:

**Issue A: Error Scenarios Test - Unknown Endpoint Panic**
The `RunErrorScenarios` function in `error_scenarios_test.go` uses exact string matching for endpoints in its switch statement. When tests use endpoints with query parameters (e.g., `/v1/dashboard/day.json?date=invalid`), they don't match the case statements (e.g., `/v1/dashboard/day.json`), causing the test to panic with "Unknown endpoint".

**Solution**: Modify the endpoint matching logic to extract the path portion (before query parameters) using `strings.Split(endpoint, "?")[0]` before comparison.

**Issue B: Test Timeout in cleanupOrphanedDatabasesConcurrent**
The `cleanupOrphanedDatabasesConcurrent` function in `test_helper.go` causes tests to hang indefinitely. The root cause is that goroutines can block indefinitely on the semaphore channel (`sem <- struct{}{}`) without respecting context cancellation. When the 10-second context times out, goroutines waiting on the semaphore never exit, causing `wg.Wait()` to block forever.

**Solution**: Refactor `cleanupOrphanedDatabasesConcurrent` to use select statements with context channels, ensuring goroutines can exit cleanly when the context times out, even while waiting on the semaphore.

### 2. Files to Modify

**File 1: `/home/danilo/scripts/github/go-reading-log-api-next/test/integration/error_scenarios_test.go`**
- **Location**: Lines 60-90 (switch statement in `RunErrorScenarios`)
- **Changes**:
  - Add helper function `extractPath(endpoint string) string` to extract path without query parameters
  - Modify all case statements to use `extractPath(scenario.Endpoint)` instead of `scenario.Endpoint`
  - This ensures endpoints like `/v1/dashboard/day.json?date=invalid` match the case `/v1/dashboard/day.json`

**File 2: `/home/danilo/scripts/github/go-reading-log-api-next/test/test_helper.go`**
- **Location**: Lines 505-560 (function `cleanupOrphanedDatabasesConcurrent`)
- **Changes**:
  - Refactor the goroutine launch loop to use select statements for context-aware semaphore acquisition
  - Add a `done` channel to signal goroutines to exit when context cancels
  - Wrap the semaphore acquisition in a select statement that listens for both semaphore slot and context cancellation
  - Ensure proper cleanup even when goroutines exit early due to timeout
  - Reduce timeout from 10 seconds to 5 seconds to prevent test blocking

### 3. Dependencies

**Prerequisites**:
- No external dependencies required
- Existing test infrastructure must be functional
- PostgreSQL must be running for integration tests

**Related Tasks**:
- None - this is a standalone bug fix
- This task unblocks all dashboard-related tests

### 4. Code Patterns

**Following Existing Patterns**:

1. **Path Extraction Pattern** (new helper):
```go
func extractPath(endpoint string) string {
    if idx := strings.Index(endpoint, "?"); idx != -1 {
        return endpoint[:idx]
    }
    return endpoint
}
```

2. **Context-Aware Semaphore Pattern** (for cleanup):
```go
select {
case sem <- struct{}{}:
    // Acquired semaphore slot
    defer func() { <-sem }()
    // Perform operation
case <-ctx.Done():
    // Context cancelled, exit gracefully
    return
}
```

3. **Error Handling Pattern**:
- Continue processing other databases even if one fails
- Aggregate errors using `errors.Join()` for visibility
- Log errors but don't block test completion

4. **Naming Conventions**:
- Helper functions: `extractPath` (camelCase)
- Test functions: `TestErrorScenarios` (PascalCase)
- Variables: `toDrop`, `sem`, `wg` (camelCase)

### 5. Testing Strategy

**Unit Tests**:
- No new unit tests required - this is a fix to existing test infrastructure
- Verify `extractPath` helper works correctly with edge cases:
  - `/v1/dashboard/day.json` → `/v1/dashboard/day.json`
  - `/v1/dashboard/day.json?date=invalid` → `/v1/dashboard/day.json`
  - `/v1/dashboard/day.json?date=invalid&type=test` → `/v1/dashboard/day.json`

**Integration Tests**:
- Run `go test -v ./test/integration -run TestErrorScenarios` to verify:
  - `Day Endpoint - Invalid Date` test passes (no more "Unknown endpoint" panic)
  - `Last Days - Invalid Type` test passes (no more timeout)
  - All other error scenarios complete successfully

**Regression Tests**:
- Run full test suite: `go test ./...`
- Verify no new timeouts or panics
- Verify test execution time is reasonable (< 30 seconds for full suite)

**Edge Cases to Cover**:
1. Empty query string: `/v1/dashboard/day.json?`
2. Multiple query parameters: `/v1/dashboard/last_days.json?type=99&days=5`
3. No query parameters: `/v1/dashboard/projects.json`
4. Cleanup with no orphaned databases
5. Cleanup with many orphaned databases (> 100)

### 6. Risks and Considerations

**Known Issues**:
- **Risk**: The current implementation may still have edge cases where goroutines don't exit cleanly
- **Mitigation**: Add explicit timeout handling with select statements and ensure all goroutines have exit paths

**Trade-offs**:
- **Trade-off**: Reducing timeout from 10s to 5s may cause cleanup to be incomplete on slow systems
- **Mitigation**: Accept that some orphaned databases may remain; they will be cleaned up on next test run

**Deployment Considerations**:
- No deployment changes required
- This is a test infrastructure fix only
- No impact on production code

**Blocking Issues**:
- None identified - this is a straightforward fix
- If issues persist after implementation, consider simplifying the cleanup logic to sequential drops (slower but more reliable)

**Additional Notes**:
- After fixing, verify that `go fmt` and `go vet` pass with no errors
- Ensure test output is clean with no panics or timeouts
- Document the fix in the task description for future reference
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Progress

### Issue A: Error Scenarios Test - Unknown Endpoint Panic
**Status:** ✅ FIXED

Added `extractPath` helper function and updated all case statements in the switch to use it. This ensures endpoints with query parameters (e.g., `/v1/dashboard/day.json?date=invalid`) match the case statements correctly.

**Changes:**
- Added `extractPath(endpoint string) string` helper function
- Modified all 8 case statements to use `extractPath(scenario.Endpoint)`
- Added `strings` import

### Issue B: Test Timeout in cleanupOrphanedDatabasesConcurrent
**Status:** ✅ FIXED

Refactored `cleanupOrphanedDatabasesConcurrent` to use context-aware semaphore acquisition. Goroutines can now exit cleanly when the context times out.

**Changes:**
- Replaced blocking `sem <- struct{}{}` with select statement that listens for both semaphore slot and context cancellation
- Changed `dropCtx` to use parent context `ctx` instead of `context.Background()` for proper cancellation propagation
- Added timeout-protected wait using a `done` channel to prevent `wg.Wait()` from blocking forever
- Goroutines now exit gracefully when context times out, even while waiting on semaphore

---

**Next Steps:**
1. Run tests to verify both fixes
2. Run go fmt and go vet
3. Check Definition of Done items
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
