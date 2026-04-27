---
id: RDL-104
title: Fix test
status: To Do
assignee:
  - Thomas
created_date: '2026-04-27 11:50'
updated_date: '2026-04-27 13:58'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
=== RUN   TestDashboardDayEndpoint_Integration
Warning: Failed to load dashboard config from , using defaults: failed to read config file: open : no such file or directory
DEBUG: Raw JSON: {"data":{"type":"dashboard_day","attributes":{"stats":{"previous_week_pages":0,"last_week_pages":0,"per_pages":133.333,"mean_day":0,"spec_mean_day":0,"progress_geral":191.667,"total_pages":0,"pages":0,"count_pages":0,"speculate_pages":0}},"id":"1777290539"}}
    dashboard_integration_test.go:76: 
        	Error Trace:	/home/danilo/scripts/github/go-reading-log-api-next/test/dashboard_integration_test.go:76
        	Error:      	Received unexpected error:
        	            	progress_geral mismatch: got 191.667000, expected 12.500000
        	Test:       	TestDashboardDayEndpoint_Integration
--- FAIL: TestDashboardDayEndpoint_Integration (0.26s)
=== RUN   TestDashboardProjectsEndpoint_Integration
Warning: Failed to load dashboard config from , using defaults: failed to read config file: open : no such file or directory
panic: test timed out after 2s
	running tests:
		TestDashboardProjectsEndpoint_Integration (2s)

goroutine 82 [running]:
testing.(*M).startAlarm.func1()
	/usr/lib/go/src/testing/testing.go:2802 +0x34b
created by time.goFunc
	/usr/lib/go/src/time/sleep.go:215 +0x2d

goroutine 1 [chan receive]:
testing.(*T).Run(0x1880685fa6c8, {0x9dfe62?, 0x188068687b30?}, 0x9f5698)
	/usr/lib/go/src/testing/testing.go:2109 +0x4e5
testing.runTests.func1(0x1880685fa6c8)
	/usr/lib/go/src/testing/testing.go:2585 +0x37
testing.tRunner(0x1880685fa6c8, 0x188068687c58)
	/usr/lib/go/src/testing/testing.go:2036 +0xea
testing.runTests({0x9d3174, 0x17}, {0x9d66c4, 0x1c}, 0x1880684f2180, {0xe8b740, 0x22, 0x22}, {0xc273f12b44600050, 0x7741a791, ...})
	/usr/lib/go/src/testing/testing.go:2583 +0x505
testing.(*M).Run(0x1880685b0e60)
	/usr/lib/go/src/testing/testing.go:2443 +0x6ac
main.main()
	_testmain.go:112 +0x9b

goroutine 73 [select]:
github.com/jackc/puddle/v2.(*Pool[...]).initResourceValue(0xa12de0, {0xa02d98, 0x188068828bd0}, 0x18806894e580)
	/home/danilo/go/pkg/mod/github.com/jackc/puddle/v2@v2.2.2/pool.go:459 +0x157
github.com/jackc/puddle/v2.(*Pool[...]).acquire(0xa12de0, {0xa02d98, 0x188068828bd0})
	/home/danilo/go/pkg/mod/github.com/jackc/puddle/v2@v2.2.2/pool.go:396 +0x1f9
github.com/jackc/puddle/v2.(*Pool[...]).Acquire(0xa12de0, {0xa02d98, 0x188068828bd0})
	/home/danilo/go/pkg/mod/github.com/jackc/puddle/v2@v2.2.2/pool.go:347 +0x89
github.com/jackc/pgx/v5/pgxpool.(*Pool).Acquire(0x1880686761c0, {0xa02d98?, 0x188068828bd0?})
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:615 +0x154
github.com/jackc/pgx/v5/pgxpool.(*Pool).QueryRow(0x188068975ee0?, {0xa02d98, 0x188068828bd0}, {0x9e25fa, 0x2d}, {0x1880689327a0, 0x1, 0x1})
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:767 +0x4f
go-reading-log-api-next/internal/adapter/postgres.(*DashboardRepositoryImpl).GetProjectsWithLogs(0x188068602160, {0xa02c10?, 0xeb45c0?})
	/home/danilo/scripts/github/go-reading-log-api-next/internal/adapter/postgres/dashboard_repository.go:335 +0x393
go-reading-log-api-next/internal/api/v1/handlers.(*DashboardHandler).Projects(0x188068683ee0, {0xa01920, 0x188068622840}, 0x0?)
	/home/danilo/scripts/github/go-reading-log-api-next/internal/api/v1/handlers/dashboard_handler.go:133 +0x78
go-reading-log-api-next/test.TestDashboardProjectsEndpoint_Integration(0x18806890a488)
	/home/danilo/scripts/github/go-reading-log-api-next/test/dashboard_integration_test.go:112 +0x3d9
testing.tRunner(0x18806890a488, 0x9f5698)
	/usr/lib/go/src/testing/testing.go:2036 +0xea
created by testing.(*T).Run in goroutine 1
	/usr/lib/go/src/testing/testing.go:2101 +0x4c5

goroutine 44 [select]:
github.com/jackc/pgx/v5/pgxpool.(*Pool).backgroundHealthCheck(0x1880686761c0)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:495 +0xc8
github.com/jackc/pgx/v5/pgxpool.NewWithConfig.func5()
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:338 +0x3d
created by github.com/jackc/pgx/v5/pgxpool.NewWithConfig in goroutine 73
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:335 +0x43e

goroutine 79 [IO wait]:
internal/poll.runtime_pollWait(0x7f4a6384ca00, 0x72)
	/usr/lib/go/src/runtime/netpoll.go:351 +0x85
internal/poll.(*pollDesc).wait(0x188068946280?, 0x18806899a000?, 0x0)
	/usr/lib/go/src/internal/poll/fd_poll_runtime.go:84 +0x27
internal/poll.(*pollDesc).waitRead(...)
	/usr/lib/go/src/internal/poll/fd_poll_runtime.go:89
internal/poll.(*FD).Read(0x188068946280, {0x18806899a000, 0x2000, 0x2000})
	/usr/lib/go/src/internal/poll/fd_unix.go:165 +0x2ae
net.(*netFD).Read(0x188068946280, {0x18806899a000?, 0x42739e?, 0x50048907f?})
	/usr/lib/go/src/net/fd_posix.go:68 +0x25
net.(*conn).Read(0x18806890c108, {0x18806899a000?, 0x4bfe7c?, 0x7?})
	/usr/lib/go/src/net/net.go:196 +0x45
github.com/jackc/pgx/v5/pgconn/internal/bgreader.(*BGReader).Read(0x18806894e700, {0x18806899a000, 0x2000, 0x2000})
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgconn/internal/bgreader/bgreader.go:100 +0xcb
io.ReadAtLeast({0x9fc320, 0x18806894e700}, {0x18806899a000, 0x2000, 0x2000}, 0x5)
	/usr/lib/go/src/io/io.go:335 +0x8e
github.com/jackc/pgx/v5/pgproto3.(*chunkReader).Next(0x18806890f770, 0x5)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgproto3/chunkreader.go:80 +0x289
github.com/jackc/pgx/v5/pgproto3.(*Frontend).Receive(0x18806899c008)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgproto3/frontend.go:309 +0x3c
github.com/jackc/pgx/v5/pgconn.(*PgConn).peekMessage(0x188068952b08)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgconn/pgconn.go:595 +0x14b
github.com/jackc/pgx/v5/pgconn.(*PgConn).receiveMessage(0x188068952b08)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgconn/pgconn.go:619 +0x26
github.com/jackc/pgx/v5/pgconn.connectOne({0xa02d98, 0x1880689987e0}, 0x188068944a80, 0x18806894e680, 0x0)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgconn/pgconn.go:402 +0xe09
github.com/jackc/pgx/v5/pgconn.connectPreferred({0xa02c80, 0x188068908ae0}, 0x188068944a80, {0x1880689327d0, 0x2, 0x188068962c40?})
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgconn/pgconn.go:264 +0x2b4
github.com/jackc/pgx/v5/pgconn.ConnectConfig({0xa02c80, 0x188068908ae0}, 0x188068944a80)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgconn/pgconn.go:159 +0x165
github.com/jackc/pgx/v5.connect({0xa02c80?, 0x188068908ae0?}, 0x188068944a80)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/conn.go:278 +0x353
github.com/jackc/pgx/v5.ConnectConfig({0xa02c80, 0x188068908ae0}, 0x188068944900)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/conn.go:160 +0x15c
github.com/jackc/pgx/v5/pgxpool.NewWithConfig.func3({0xa02c80, 0x188068908ae0})
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:289 +0x1e8
github.com/jackc/puddle/v2.(*Pool[...]).initResourceValue.func1()
	/home/danilo/go/pkg/mod/github.com/jackc/puddle/v2@v2.2.2/pool.go:423 +0xdc
created by github.com/jackc/puddle/v2.(*Pool[...]).initResourceValue in goroutine 73
	/home/danilo/go/pkg/mod/github.com/jackc/puddle/v2@v2.2.2/pool.go:421 +0xf6

goroutine 81 [select]:
context.(*cancelCtx).propagateCancel.func2()
	/usr/lib/go/src/context/context.go:523 +0x9a
created by context.(*cancelCtx).propagateCancel in goroutine 79
	/usr/lib/go/src/context/context.go:522 +0x409
FAIL	go-reading-log-api-next/test	2.011s
?   	go-reading-log-api-next/test/fixtures	[no test files]
=== RUN   TestFixtureValidator_WeekdayCoverage
--- PASS: TestFixtureValidator_WeekdayCoverage (0.00s)
=== RUN   TestFixtureValidator_WeekdayCoverage_Missing
--- PASS: TestFixtureValidator_WeekdayCoverage_Missing (0.00s)
=== RUN   TestFixtureValidator_WeekdayCoverage_NoLogs
--- PASS: TestFixtureValidator_WeekdayCoverage_NoLogs (0.00s)
=== RUN   TestFixtureValidator_DataRange
--- PASS: TestFixtureValidator_DataRange (0.00s)
=== RUN   TestFixtureValidator_DataRange_Insufficient
--- PASS: TestFixtureValidator_DataRange_Insufficient (0.00s)
=== RUN   TestFixtureValidator_DataRange_DuplicateDates
--- PASS: TestFixtureValidator_DataRange_DuplicateDates (0.00s)
=== RUN   TestFixtureValidator_Combined
--- PASS: TestFixtureValidator_Combined (0.00s)
=== RUN   TestFixtureValidator_ProjectConsistency
--- PASS: TestFixtureValidator_ProjectConsistency (0.00s)
=== RUN   TestFixtureValidator_DateRange
--- PASS: TestFixtureValidator_DateRange (0.00s)
=== RUN   TestFixtureValidator_DateRange_Narrow
--- PASS: TestFixtureValidator_DateRange_Narrow (0.00s)
=== RUN   TestValidateScenario
--- PASS: TestValidateScenario (0.00s)
=== RUN   TestMustValidateScenario
--- PASS: TestMustValidateScenario (0.00s)
=== RUN   TestMustValidateScenario_Panic
--- PASS: TestMustValidateScenario_Panic (0.00s)
=== RUN   TestFixtureValidator_Warnings
--- PASS: TestFixtureValidator_Warnings (0.00s)
PASS
ok  	go-reading-log-api-next/test/fixtures/dashboard	(cached)
=== RUN   TestErrorScenarios
=== RUN   TestErrorScenarios/Day_Endpoint_-_Invalid_Date
=== RUN   TestErrorScenarios/Last_Days_-_Invalid_Type
panic: test timed out after 2s
	running tests:
		TestErrorScenarios (2s)
		TestErrorScenarios/Last_Days_-_Invalid_Type (2s)

goroutine 48 [running]:
testing.(*M).startAlarm.func1()
	/usr/lib/go/src/testing/testing.go:2802 +0x34b
created by time.goFunc
	/usr/lib/go/src/time/sleep.go:215 +0x2d

goroutine 1 [chan receive]:
testing.(*T).Run(0x3c0f9aff66c8, {0xa696f8?, 0x3c0f9afd3b30?}, 0xa92d88)
	/usr/lib/go/src/testing/testing.go:2109 +0x4e5
testing.runTests.func1(0x3c0f9aff66c8)
	/usr/lib/go/src/testing/testing.go:2585 +0x37
testing.tRunner(0x3c0f9aff66c8, 0x3c0f9afd3c58)
	/usr/lib/go/src/testing/testing.go:2036 +0xea
testing.runTests({0xa6caeb, 0x17}, {0xa7aac1, 0x28}, 0x3c0f9af9c168, {0xf8e780, 0x30, 0x30}, {0xc273f12b474fbdad, 0x7741cc2e, ...})
	/usr/lib/go/src/testing/testing.go:2583 +0x505
testing.(*M).Run(0x3c0f9af9af00)
	/usr/lib/go/src/testing/testing.go:2443 +0x6ac
main.main()
	_testmain.go:140 +0x9b

goroutine 21 [chan receive]:
testing.(*T).Run(0x3c0f9aff6908, {0xa6d855?, 0x3c0f9afac0e0?}, 0x3c0f9af12410)
	/usr/lib/go/src/testing/testing.go:2109 +0x4e5
go-reading-log-api-next/test/integration.RunErrorScenarios(0x3c0f9aff6908, {0xf8cee0, 0x5, 0x3c0f9ae9f760?})
	/home/danilo/scripts/github/go-reading-log-api-next/test/integration/error_scenarios_test.go:46 +0xd1
go-reading-log-api-next/test/integration.TestErrorScenarios(0x3c0f9aff6908)
	/home/danilo/scripts/github/go-reading-log-api-next/test/integration/error_scenarios_test.go:269 +0x6c
testing.tRunner(0x3c0f9aff6908, 0xa92d88)
	/usr/lib/go/src/testing/testing.go:2036 +0xea
created by testing.(*T).Run in goroutine 1
	/usr/lib/go/src/testing/testing.go:2101 +0x4c5

goroutine 73 [select]:
go-reading-log-api-next/test.cleanupOrphanedDatabasesConcurrent(0x3c0f9af6e1c0, {0x3c0f9af60120, 0x2e})
	/home/danilo/scripts/github/go-reading-log-api-next/test/test_helper.go:565 +0x6d6
go-reading-log-api-next/test.(*TestHelper).Close.func1()
	/home/danilo/scripts/github/go-reading-log-api-next/test/test_helper.go:377 +0x1bf
go-reading-log-api-next/test.(*TestHelper).Close(0x0?)
	/home/danilo/scripts/github/go-reading-log-api-next/test/test_helper.go:403 +0x8b
go-reading-log-api-next/test/integration.RunErrorScenarios.func1(0x3c0f9af42488)
	/home/danilo/scripts/github/go-reading-log-api-next/test/integration/error_scenarios_test.go:101 +0x8de
testing.tRunner(0x3c0f9af42488, 0x3c0f9af12410)
	/usr/lib/go/src/testing/testing.go:2036 +0xea
created by testing.(*T).Run in goroutine 21
	/usr/lib/go/src/testing/testing.go:2101 +0x4c5

goroutine 52 [IO wait]:
internal/poll.runtime_pollWait(0x7f4f0a26f800, 0x72)
	/usr/lib/go/src/runtime/netpoll.go:351 +0x85
internal/poll.(*pollDesc).wait(0x3c0f9b0a6900?, 0x3c0f9b0fc000?, 0x0)
	/usr/lib/go/src/internal/poll/fd_poll_runtime.go:84 +0x27
internal/poll.(*pollDesc).waitRead(...)
	/usr/lib/go/src/internal/poll/fd_poll_runtime.go:89
internal/poll.(*FD).Read(0x3c0f9b0a6900, {0x3c0f9b0fc000, 0x2000, 0x2000})
	/usr/lib/go/src/internal/poll/fd_unix.go:165 +0x2ae
net.(*netFD).Read(0x3c0f9b0a6900, {0x3c0f9b0fc000?, 0x3c0f9b0ea2d0?, 0xfb7c60?})
	/usr/lib/go/src/net/fd_posix.go:68 +0x25
net.(*conn).Read(0x3c0f9ae1a1f8, {0x3c0f9b0fc000?, 0x9?, 0x3c0f9b0bf730?})
	/usr/lib/go/src/net/net.go:196 +0x45
github.com/jackc/pgx/v5/pgconn/internal/bgreader.(*BGReader).Read(0x3c0f9aecea40, {0x3c0f9b0fc000, 0x2000, 0x2000})
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgconn/internal/bgreader/bgreader.go:100 +0xcb
io.ReadAtLeast({0xa9b8e0, 0x3c0f9aecea40}, {0x3c0f9b0fc000, 0x2000, 0x2000}, 0x5)
	/usr/lib/go/src/io/io.go:335 +0x8e
github.com/jackc/pgx/v5/pgproto3.(*chunkReader).Next(0x3c0f9b0ea840, 0x5)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgproto3/chunkreader.go:80 +0x289
github.com/jackc/pgx/v5/pgproto3.(*Frontend).Receive(0x3c0f9b0d3688)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgproto3/frontend.go:309 +0x3c
github.com/jackc/pgx/v5/pgconn.(*PgConn).peekMessage(0x3c0f9b0af608)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgconn/pgconn.go:595 +0x14b
github.com/jackc/pgx/v5/pgconn.(*MultiResultReader).NextResult(0x3c0f9b0af738)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgconn/pgconn.go:1552 +0x4e
github.com/jackc/pgx/v5.(*Conn).execSimpleProtocol(0x3c0f9b0ac5a0?, {0xaa2490?, 0x3c0f9aed1ce0?}, {0x3c0f9ae600f0?, 0x3c0f9b0bfb30?}, {0x0?, 0x424cdc?, 0x3c0f9b0bfb88?})
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/conn.go:589 +0xb0
github.com/jackc/pgx/v5.(*Conn).exec(0x3c0f9b0ee640, {0xaa2490, 0x3c0f9aed1ce0}, {0x3c0f9ae600f0, 0x46}, {0x0?, 0x87f7ed?, 0x3c0f9aece8c0?})
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/conn.go:574 +0x8ac
github.com/jackc/pgx/v5.(*Conn).Exec(0x3c0f9b0ee640, {0xaa2490?, 0x3c0f9aed1ce0?}, {0x3c0f9ae600f0, 0x46}, {0x0, 0x0, 0x0})
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/conn.go:481 +0x114
github.com/jackc/pgx/v5/pgxpool.(*Conn).Exec(0x3c0f9b0a8380?, {0xaa2490?, 0x3c0f9aed1ce0?}, {0x3c0f9ae600f0?, 0x3c0f9aefa0d0?}, {0x0?, 0x4f4ada?, 0x3c0f9aefa0d0?})
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/conn.go:87 +0x3c
github.com/jackc/pgx/v5/pgxpool.(*Pool).Exec(0xa6f347?, {0xaa2490, 0x3c0f9aed1ce0}, {0x3c0f9ae600f0, 0x46}, {0x0, 0x0, 0x0})
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:726 +0xf7
go-reading-log-api-next/test.cleanupOrphanedDatabasesConcurrent.func1({0x3c0f9b256090, 0x2e})
	/home/danilo/scripts/github/go-reading-log-api-next/test/test_helper.go:548 +0x5e6
created by go-reading-log-api-next/test.cleanupOrphanedDatabasesConcurrent in goroutine 73
	/home/danilo/scripts/github/go-reading-log-api-next/test/test_helper.go:511 +0x485

goroutine 53 [sync.WaitGroup.Wait]:
sync.runtime_SemacquireWaitGroup(0x3c0f9aed1730?, 0xe0?)
	/usr/lib/go/src/runtime/sema.go:114 +0x2e
sync.(*WaitGroup).Wait(0x3c0f9b2105d0)
	/usr/lib/go/src/sync/waitgroup.go:206 +0x85
go-reading-log-api-next/test.cleanupOrphanedDatabasesConcurrent.func2()
	/home/danilo/scripts/github/go-reading-log-api-next/test/test_helper.go:561 +0x25
created by go-reading-log-api-next/test.cleanupOrphanedDatabasesConcurrent in goroutine 73
	/home/danilo/scripts/github/go-reading-log-api-next/test/test_helper.go:560 +0x676

goroutine 44 [select]:
github.com/jackc/pgx/v5/pgxpool.(*Pool).backgroundHealthCheck(0x3c0f9b0a8380)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:495 +0xc8
github.com/jackc/pgx/v5/pgxpool.NewWithConfig.func5()
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:338 +0x3d
created by github.com/jackc/pgx/v5/pgxpool.NewWithConfig in goroutine 52
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:335 +0x43e

goroutine 79 [select]:
github.com/jackc/pgx/v5/pgxpool.(*Pool).backgroundHealthCheck(0x3c0f9af6e1c0)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:495 +0xc8
github.com/jackc/pgx/v5/pgxpool.NewWithConfig.func5()
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:338 +0x3d
created by github.com/jackc/pgx/v5/pgxpool.NewWithConfig in goroutine 73
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:335 +0x43e
FAIL	go-reading-log-api-next/test/integration	2.008s
testing: warning: no tests to run
PASS
ok  	go-reading-log-api-next/test/performance	0.003s [no tests to run]
?   	go-reading-log-api-next/test/testutil	[no test files]
=== RUN   TestDashboardRepository_GetDailyStats
--- PASS: TestDashboardRepository_GetDailyStats (0.21s)
=== RUN   TestDashboardRepository_GetDailyStats_EmptyDate
panic: test timed out after 2s
	running tests:
		TestDashboardRepository_GetDailyStats_EmptyDate (2s)

goroutine 68 [running]:
testing.(*M).startAlarm.func1()
	/usr/lib/go/src/testing/testing.go:2802 +0x34b
created by time.goFunc
	/usr/lib/go/src/time/sleep.go:215 +0x2d

goroutine 1 [chan receive]:
testing.(*T).Run(0x1c61dcf566c8, {0x9e653f?, 0x1c61dcf63b30?}, 0x9f81e0)
	/usr/lib/go/src/testing/testing.go:2109 +0x4e5
testing.runTests.func1(0x1c61dcf566c8)
	/usr/lib/go/src/testing/testing.go:2585 +0x37
testing.tRunner(0x1c61dcf566c8, 0x1c61dcf63c58)
	/usr/lib/go/src/testing/testing.go:2036 +0xea
testing.runTests({0x9d5e4e, 0x17}, {0x9dd488, 0x21}, 0x1c61dcdce180, {0xe95620, 0x7e, 0x7e}, {0xc273f12b474ab5b4, 0x773d7afe, ...})
	/usr/lib/go/src/testing/testing.go:2583 +0x505
testing.(*M).Run(0x1c61dcf0ee60)
	/usr/lib/go/src/testing/testing.go:2443 +0x6ac
main.main()
	_testmain.go:296 +0x9b

goroutine 43 [select]:
go-reading-log-api-next/test.cleanupOrphanedDatabasesConcurrent(0x1c61dcf84380, {0x1c61dcdf68a0, 0x2e})
	/home/danilo/scripts/github/go-reading-log-api-next/test/test_helper.go:565 +0x6d6
go-reading-log-api-next/test.(*TestHelper).Close.func1()
	/home/danilo/scripts/github/go-reading-log-api-next/test/test_helper.go:377 +0x1bf
go-reading-log-api-next/test.(*TestHelper).Close(0x1c61dce07020?)
	/home/danilo/scripts/github/go-reading-log-api-next/test/test_helper.go:403 +0x8b
go-reading-log-api-next/test/unit.TestDashboardRepository_GetDailyStats_EmptyDate(0x1c61dcf57d48)
	/home/danilo/scripts/github/go-reading-log-api-next/test/unit/dashboard_repository_test.go:86 +0x1d4
testing.tRunner(0x1c61dcf57d48, 0x9f81e0)
	/usr/lib/go/src/testing/testing.go:2036 +0xea
created by testing.(*T).Run in goroutine 1
	/usr/lib/go/src/testing/testing.go:2101 +0x4c5

goroutine 49 [select]:
github.com/jackc/pgx/v5/pgxpool.(*Pool).backgroundHealthCheck(0x1c61dcf84380)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:495 +0xc8
github.com/jackc/pgx/v5/pgxpool.NewWithConfig.func5()
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:338 +0x3d
created by github.com/jackc/pgx/v5/pgxpool.NewWithConfig in goroutine 43
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:335 +0x43e

goroutine 62 [IO wait]:
internal/poll.runtime_pollWait(0x7fb08b8d3a00, 0x72)
	/usr/lib/go/src/runtime/netpoll.go:351 +0x85
internal/poll.(*pollDesc).wait(0x1c61dceb2680?, 0x1c61dcef8000?, 0x0)
	/usr/lib/go/src/internal/poll/fd_poll_runtime.go:84 +0x27
internal/poll.(*pollDesc).waitRead(...)
	/usr/lib/go/src/internal/poll/fd_poll_runtime.go:89
internal/poll.(*FD).Read(0x1c61dceb2680, {0x1c61dcef8000, 0x2000, 0x2000})
	/usr/lib/go/src/internal/poll/fd_unix.go:165 +0x2ae
net.(*netFD).Read(0x1c61dceb2680, {0x1c61dcef8000?, 0x1c61dce9d380?, 0xebb8a0?})
	/usr/lib/go/src/net/fd_posix.go:68 +0x25
net.(*conn).Read(0x1c61dce82108, {0x1c61dcef8000?, 0x9?, 0x1c61dcedb730?})
	/usr/lib/go/src/net/net.go:196 +0x45
github.com/jackc/pgx/v5/pgconn/internal/bgreader.(*BGReader).Read(0x1c61dceae6c0, {0x1c61dcef8000, 0x2000, 0x2000})
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgconn/internal/bgreader/bgreader.go:100 +0xcb
io.ReadAtLeast({0x9ff400, 0x1c61dceae6c0}, {0x1c61dcef8000, 0x2000, 0x2000}, 0x5)
	/usr/lib/go/src/io/io.go:335 +0x8e
github.com/jackc/pgx/v5/pgproto3.(*chunkReader).Next(0x1c61dce9d8c0, 0x5)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgproto3/chunkreader.go:80 +0x289
github.com/jackc/pgx/v5/pgproto3.(*Frontend).Receive(0x1c61dcec1b08)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgproto3/frontend.go:309 +0x3c
github.com/jackc/pgx/v5/pgconn.(*PgConn).peekMessage(0x1c61dceceb08)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgconn/pgconn.go:595 +0x14b
github.com/jackc/pgx/v5/pgconn.(*MultiResultReader).NextResult(0x1c61dcecec38)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgconn/pgconn.go:1552 +0x4e
github.com/jackc/pgx/v5.(*Conn).execSimpleProtocol(0x1c61dcf0ad20?, {0xa05db8?, 0x1c61dce99030?}, {0x1c61dceca0a0?, 0x1c61dcedbb30?}, {0x0?, 0x424c3c?, 0x1c61dcedbb88?})
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/conn.go:589 +0xb0
github.com/jackc/pgx/v5.(*Conn).exec(0x1c61dceb77c0, {0xa05db8, 0x1c61dce99030}, {0x1c61dceca0a0, 0x46}, {0x0?, 0x7dc3ed?, 0x1c61dceae540?})
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/conn.go:574 +0x8ac
github.com/jackc/pgx/v5.(*Conn).Exec(0x1c61dceb77c0, {0xa05db8?, 0x1c61dce99030?}, {0x1c61dceca0a0, 0x46}, {0x0, 0x0, 0x0})
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/conn.go:481 +0x114
github.com/jackc/pgx/v5/pgxpool.(*Conn).Exec(0x1c61dcec61c0?, {0xa05db8?, 0x1c61dce99030?}, {0x1c61dceca0a0?, 0x1c61dcea00d0?}, {0x0?, 0x4f1c7a?, 0x1c61dcea00d0?})
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/conn.go:87 +0x3c
github.com/jackc/pgx/v5/pgxpool.(*Pool).Exec(0x9d818f?, {0xa05db8, 0x1c61dce99030}, {0x1c61dceca0a0, 0x46}, {0x0, 0x0, 0x0})
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:726 +0xf7
go-reading-log-api-next/test.cleanupOrphanedDatabasesConcurrent.func1({0x1c61dceac270, 0x2e})
	/home/danilo/scripts/github/go-reading-log-api-next/test/test_helper.go:548 +0x5e6
created by go-reading-log-api-next/test.cleanupOrphanedDatabasesConcurrent in goroutine 43
	/home/danilo/scripts/github/go-reading-log-api-next/test/test_helper.go:511 +0x485

goroutine 63 [sync.WaitGroup.Wait]:
sync.runtime_SemacquireWaitGroup(0x1c61dce98850?, 0x60?)
	/usr/lib/go/src/runtime/sema.go:114 +0x2e
sync.(*WaitGroup).Wait(0x1c61dce9e3d0)
	/usr/lib/go/src/sync/waitgroup.go:206 +0x85
go-reading-log-api-next/test.cleanupOrphanedDatabasesConcurrent.func2()
	/home/danilo/scripts/github/go-reading-log-api-next/test/test_helper.go:561 +0x25
created by go-reading-log-api-next/test.cleanupOrphanedDatabasesConcurrent in goroutine 43
	/home/danilo/scripts/github/go-reading-log-api-next/test/test_helper.go:560 +0x676

goroutine 64 [select]:
github.com/jackc/pgx/v5/pgxpool.(*Pool).backgroundHealthCheck(0x1c61dcec61c0)
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:495 +0xc8
github.com/jackc/pgx/v5/pgxpool.NewWithConfig.func5()
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:338 +0x3d
created by github.com/jackc/pgx/v5/pgxpool.NewWithConfig in goroutine 62
	/home/danilo/go/pkg/mod/github.com/jackc/pgx/v5@v5.9.1/pgxpool/pool.go:335 +0x43e
FAIL	go-reading-log-api-next/test/unit	2.008s
?   	go-reading-log-api-next/tools	[no test files]
FAIL
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task involves fixing multiple failing tests caused by two distinct root causes:

**Issue A: Dashboard Day Endpoint Calculation Mismatch**
- The `progress_geral` calculation in `Day` handler uses incorrect formula
- Currently: `(sum of all log end_page values) / (sum of project total_page) * 100`
- Should be: `(sum of project Page field) / (sum of project total_page) * 100`
- The handler is aggregating log data instead of using project state

**Issue B: Test Timeout in Cleanup Operations**
- `cleanupOrphanedDatabasesConcurrent` function causes test timeouts
- Root cause: Creating new connection pools inside each goroutine exhausts connections
- The function spawns goroutines that each call `pgxpool.New()`, causing connection pool exhaustion when many orphaned databases exist
- Semaphore limits concurrent drops to 5, but each goroutine still creates its own pool, leading to resource exhaustion

**Solution Strategy:**
1. Fix `Day` handler to calculate `progress_geral` using project `Page` field from database
2. Refactor `cleanupOrphanedDatabasesConcurrent` to reuse a single connection pool for all DROP operations
3. Update test scenario expectations to match correct calculation

### 2. Files to Modify

**Primary Files:**
- `internal/api/v1/handlers/dashboard_handler.go`
  - Fix `Day` handler's `progress_geral` calculation to query project `page` field
  - Add SQL query to fetch project page values for progress calculation

- `test/test_helper.go`
  - Refactor `cleanupOrphanedDatabasesConcurrent` to use single connection pool
  - Move pool creation outside goroutine loop
  - Ensure proper connection reuse and cleanup

**Secondary Files:**
- `test/fixtures/dashboard/scenarios.go`
  - Update `ScenarioMultipleProjects` expected `progress_geral` value
  - Current expected: 12.5 (incorrect)
  - New expected: Calculate based on actual fixture data (Page values: 0 + 50 + 200 = 250, total capacity: 600, progress: 41.67)

**Files to Read (for context):**
- `test/dashboard_integration_test.go` - Understand test validation logic
- `internal/service/dashboard/day_service.go` - Understand alternative calculation approach
- `internal/adapter/postgres/dashboard_repository.go` - Understand repository methods

### 3. Dependencies

**Prerequisites:**
- Database must be running and accessible for integration tests
- Test database cleanup from previous runs should be performed before testing
- No external dependencies - all changes are internal to the codebase

**Blocking Issues:**
- None identified - all required code exists and can be modified in place

**Setup Steps:**
1. Run `make test-clean` to clear any orphaned test databases before testing
2. Ensure PostgreSQL is running: `pg_isready -h localhost -p 5432`
3. Verify `.env.test` file exists with correct database credentials

### 4. Code Patterns

**Following Existing Patterns:**

**Handler Pattern (dashboard_handler.go):**
```go
// Current pattern for database queries in handlers
query := `SELECT ... FROM projects WHERE id = $1`
var value int
err := h.repo.GetPool().QueryRow(ctx, query, projectID).Scan(&value)
```

**Context Timeout Pattern:**
```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()
```

**Error Wrapping Pattern:**
```go
return fmt.Errorf("failed to calculate progress: %w", err)
```

**Math Rounding Pattern (3 decimal places):**
```go
math.Round(value*1000) / 1000
```

**Cleanup Pattern (test_helper.go):**
```go
// Single pool reuse pattern (to be implemented)
pool := pgxpool.New(ctx, connStr)
defer pool.Close()
for _, dbName := range toDrop {
    pool.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
}
```

**Naming Conventions:**
- Use snake_case for database columns and JSON fields
- Use camelCase for Go struct fields
- Function names: PascalCase for exported, camelCase for unexported
- Test functions: `Test<FunctionName>_<Scenario>`

### 5. Testing Strategy

**Unit Tests:**
- No new unit tests required - existing unit tests should pass after fixes
- Verify `TestDashboardRepository_GetDailyStats` passes (currently timing out)

**Integration Tests:**
- `TestDashboardDayEndpoint_Integration` - Verify `progress_geral` calculation matches expected value
- `TestDashboardProjectsEndpoint_Integration` - Verify no timeout occurs
- `TestErrorScenarios` - Verify no timeout occurs in cleanup

**Edge Cases to Cover:**
1. Empty database state (zero projects)
2. Single project with zero pages
3. Multiple projects with varying completion states
4. Cleanup with many orphaned databases (stress test)

**Testing Approach:**
```bash
# Run specific failing tests
go test -v -run TestDashboardDayEndpoint_Integration ./test
go test -v -run TestDashboardProjectsEndpoint_Integration ./test
go test -v -run TestErrorScenarios ./test/integration

# Run full test suite
go test -v ./...
```

**Validation Criteria:**
- All tests complete within 30 seconds (no timeouts)
- `progress_geral` matches expected value (41.67 for ScenarioMultipleProjects)
- No connection pool exhaustion errors
- Cleanup completes without blocking

### 6. Risks and Considerations

**Known Issues:**
1. **Calculation Discrepancy Risk:** The expected value in `ScenarioMultipleProjects` (12.5) appears to be incorrect based on the fixture data. The comment says `(0 + 25 + 200) / (200 + 200 + 200) * 100` which equals 41.67, not 12.5. This needs clarification:
   - Option A: Update expected value to 41.67 (matches actual fixture data)
   - Option B: Change fixture data to match expected value 12.5
   - **Recommendation:** Option A - fixture data is the source of truth

2. **Connection Pool Exhaustion:** The current cleanup implementation creates a new pool for each database drop. With many orphaned databases, this can exhaust PostgreSQL connections.
   - **Mitigation:** Reuse single pool for all DROP operations
   - **Fallback:** Add connection timeout and retry logic if pool acquisition fails

3. **Race Conditions:** Parallel tests create unique databases with PID+goroutine+timestamp naming
   - **Verification:** Ensure cleanup doesn't drop active test databases
   - **Safety:** The `excludeName` parameter prevents dropping current test database

**Deployment Considerations:**
- No deployment changes required - changes are test-only
- Production database cleanup logic remains unchanged
- Test database cleanup is isolated to test environment

**Rollback Plan:**
- If issues arise, revert to previous version using git
- Test database cleanup can be manually run: `make test-clean`

**Performance Impact:**
- Cleanup operation should complete in < 10 seconds (reduced from potential timeout)
- Single pool reuse reduces connection overhead
- No impact on production performance

**Acceptance Criteria Mapping:**
- [x] All unit tests pass
- [x] All integration tests pass execution and verification
- [x] go fmt and go vet pass with no errors
- [x] Error responses consistent with existing patterns
- [x] HTTP status codes correct for response type
- [ ] Documentation updated in QWEN.md (not applicable for test fixes)
- [x] Integration tests verify actual database interactions
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
### Implementation Progress

**Status**: In Progress - Starting implementation

**Root Causes Identified:**
1. **Dashboard Day Calculation Issue**: The `progress_geral` calculation in `dashboard_handler.go` uses sum of log `end_page` values instead of project `Page` field from database
2. **Test Timeout Issue**: `cleanupOrphanedDatabasesConcurrent` in `test_helper.go` creates a new connection pool for each goroutine, causing connection pool exhaustion

**Plan:**
1. Fix `Day` handler to calculate `progress_geral` using project `Page` field from database
2. Refactor `cleanupOrphanedDatabasesConcurrent` to reuse a single connection pool for all DROP operations
3. Update test scenario expected value for `progress_geral` to match correct calculation

**Next Steps:**
- Modify `internal/api/v1/handlers/dashboard_handler.go` - Fix progress_geral calculation
- Modify `test/test_helper.go` - Refactor cleanup to use single pool
- Modify `test/fixtures/dashboard/scenarios.go` - Update expected value
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
