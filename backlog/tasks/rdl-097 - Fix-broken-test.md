---
id: RDL-097
title: Fix broken test
status: To Do
assignee:
  - catarina
created_date: '2026-04-23 18:15'
updated_date: '2026-04-23 18:15'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Fix the test go test -v -timeout=5s ./...

```go
ok  	go-reading-log-api-next/internal/validation	(cached)
=== RUN   TestDashboardDayEndpoint_Integration
Warning: Failed to load dashboard config from , using defaults: failed to read config file: open : no such file or directory
DEBUG: Raw JSON: {"data":{"type":"dashboard_day","attributes":{"total_pages":0,"log_count":0},"id":"1776967972"}}
    dashboard_integration_test.go:72: 
        	Error Trace:	/home/danilo/scripts/github/go-reading-log-api-next/test/dashboard_integration_test.go:72
        	Error:      	Not equal: 
        	            	expected: 133.333
        	            	actual  : 0
        	Test:       	TestDashboardDayEndpoint_Integration
    dashboard_integration_test.go:76: 
        	Error Trace:	/home/danilo/scripts/github/go-reading-log-api-next/test/dashboard_integration_test.go:76
        	Error:      	Received unexpected error:
        	            	progress_geral mismatch: got 0.000000, expected 12.500000
        	Test:       	TestDashboardDayEndpoint_Integration
--- FAIL: TestDashboardDayEndpoint_Integration (0.10s)
=== RUN   TestDashboardProjectsEndpoint_Integration
Warning: Failed to load dashboard config from , using defaults: failed to read config file: open : no such file or directory
DEBUG: Raw JSON: {"data":[{"type":"dashboard_projects","attributes":[{"project_id":10,"project_name":"Unstarted Project","total_pages":0,"log_count":0,"progress":0},{"project_id":11,"project_name":"Running Project","total_pages":50,"log_count":2,"progress":100},{"project_id":12,"project_name":"Finished Project","total_pages":0,"log_count":0,"progress":0}],"id":"1776967972"}]}
    dashboard_integration_test.go:121: 
        	Error Trace:	/home/danilo/scripts/github/go-reading-log-api-next/test/dashboard_integration_test.go:121
        	Error:      	Expected value not to be nil.
        	Test:       	TestDashboardProjectsEndpoint_Integration
--- FAIL: TestDashboardProjectsEndpoint_Integration (0.11s)
=== RUN   TestDashboardLastDaysEndpoint_Integration
Warning: Failed to load dashboard config from , using defaults: failed to read config file: open : no such file or directory
=== RUN   TestDashboardLastDaysEndpoint_Integration/type_1
DEBUG: Raw JSON: {"data":{"type":"dashboard_last_days","attributes":{"avg_per_day":0,"days":7,"end_date":"2026-04-23T15:12:52-03:00","start_date":"2026-04-17T15:12:52-03:00","total_faults":0,"type":"1"},"id":"1776967972"}}
    dashboard_integration_test.go:202: 
        	Error Trace:	/home/danilo/scripts/github/go-reading-log-api-next/test/dashboard_integration_test.go:202
        	Error:      	Should NOT be empty, but was []
        	Test:       	TestDashboardLastDaysEndpoint_Integration/type_1
=== RUN   TestDashboardLastDaysEndpoint_Integration/type_2
DEBUG: Raw JSON: {"data":{"type":"dashboard_last_days","attributes":{"avg_per_day":0,"days":7,"end_date":"2026-04-23T15:12:52-03:00","start_date":"2026-04-17T15:12:52-03:00","total_faults":0,"type":"2"},"id":"1776967972"}}
    dashboard_integration_test.go:202: 
        	Error Trace:	/home/danilo/scripts/github/go-reading-log-api-next/test/dashboard_integration_test.go:202
        	Error:      	Should NOT be empty, but was []
        	Test:       	TestDashboardLastDaysEndpoint_Integration/type_2
=== RUN   TestDashboardLastDaysEndpoint_Integration/type_3
DEBUG: Raw JSON: {"data":{"type":"dashboard_last_days","attributes":{"avg_per_day":0,"days":7,"end_date":"2026-04-23T15:12:52-03:00","start_date":"2026-04-17T15:12:52-03:00","total_faults":0,"type":"3"},"id":"1776967972"}}
    dashboard_integration_test.go:202: 
        	Error Trace:	/home/danilo/scripts/github/go-reading-log-api-next/test/dashboard_integration_test.go:202
        	Error:      	Should NOT be empty, but was []
        	Test:       	TestDashboardLastDaysEndpoint_Integration/type_3
=== RUN   TestDashboardLastDaysEndpoint_Integration/type_4
DEBUG: Raw JSON: {"data":{"type":"dashboard_last_days","attributes":{"avg_per_day":0,"days":7,"end_date":"2026-04-23T15:12:52-03:00","start_date":"2026-04-17T15:12:52-03:00","total_faults":0,"type":"4"},"id":"1776967972"}}
    dashboard_integration_test.go:202: 
        	Error Trace:	/home/danilo/scripts/github/go-reading-log-api-next/test/dashboard_integration_test.go:202
        	Error:      	Should NOT be empty, but was []
        	Test:       	TestDashboardLastDaysEndpoint_Integration/type_4
=== RUN   TestDashboardLastDaysEndpoint_Integration/type_5
DEBUG: Raw JSON: {"data":{"type":"dashboard_last_days","attributes":{"avg_per_day":0,"days":7,"end_date":"2026-04-23T15:12:52-03:00","start_date":"2026-04-17T15:12:52-03:00","total_faults":0,"type":"5"},"id":"1776967972"}}
    dashboard_integration_test.go:202: 
        	Error Trace:	/home/danilo/scripts/github/go-reading-log-api-next/test/dashboard_integration_test.go:202
        	Error:      	Should NOT be empty, but was []
        	Test:       	TestDashboardLastDaysEndpoint_Integration/type_5
=== RUN   TestDashboardLastDaysEndpoint_Integration/invalid_type
DEBUG: Raw JSON: {"data":{"type":"dashboard_last_days","attributes":{"avg_per_day":0,"days":7,"end_date":"2026-04-23T15:12:52-03:00","start_date":"2026-04-17T15:12:52-03:00","total_faults":0,"type":"99"},"id":"1776967972"}}
    dashboard_integration_test.go:214: 
        	Error Trace:	/home/danilo/scripts/github/go-reading-log-api-next/test/dashboard_integration_test.go:214
        	Error:      	Not equal: 
        	            	expected: 422
        	            	actual  : 200
        	Test:       	TestDashboardLastDaysEndpoint_Integration/invalid_type
--- FAIL: TestDashboardLastDaysEndpoint_Integration (0.21s)
    --- FAIL: TestDashboardLastDaysEndpoint_Integration/type_1 (0.00s)
    --- FAIL: TestDashboardLastDaysEndpoint_Integration/type_2 (0.00s)
    --- FAIL: TestDashboardLastDaysEndpoint_Integration/type_3 (0.00s)
    --- FAIL: TestDashboardLastDaysEndpoint_Integration/type_4 (0.00s)
    --- FAIL: TestDashboardLastDaysEndpoint_Integration/type_5 (0.00s)
    --- FAIL: TestDashboardLastDaysEndpoint_Integration/invalid_type (0.00s)
=== RUN   TestDashboardFaultsChart_Integration
Warning: Failed to load dashboard config from , using defaults: failed to read config file: open : no such file or directory
DEBUG: Raw JSON: {"data":{"type":"dashboard_echart_faults","attributes":{"title":"Faults Gauge","tooltip":{"formatter":"{a} \u003cbr/\u003e{b} : {c}%"},"series":[{"name":"Faults","type":"gauge","data":[0],"itemStyle":{"color":"#4caf50"}}]},"id":"1776967972"}}
    dashboard_integration_test.go:260: 
        	Error Trace:	/home/danilo/scripts/github/go-reading-log-api-next/test/dashboard_integration_test.go:260
        	Error:      	Expected value not to be nil.
        	Test:       	TestDashboardFaultsChart_Integration
--- FAIL: TestDashboardFaultsChart_Integration (0.17s)
panic: runtime error: invalid memory address or nil pointer dereference [recovered, repanicked]
[signal SIGSEGV: segmentation violation code=0x1 addr=0x30 pc=0x84aa60]

goroutine 77 [running]:
testing.tRunner.func1.2({0x918c40, 0xe43e10})
	/usr/lib/go/src/testing/testing.go:1872 +0x237
testing.tRunner.func1()
	/usr/lib/go/src/testing/testing.go:1875 +0x35b
panic({0x918c40?, 0xe43e10?})
	/usr/lib/go/src/runtime/panic.go:783 +0x132
go-reading-log-api-next/test.TestDashboardFaultsChart_Integration(0xc000103180)
	/home/danilo/scripts/github/go-reading-log-api-next/test/dashboard_integration_test.go:261 +0x4c0
testing.tRunner(0xc000103180, 0x9f0508)
	/usr/lib/go/src/testing/testing.go:1934 +0xea
created by testing.(*T).Run in goroutine 1
	/usr/lib/go/src/testing/testing.go:1997 +0x465
FAIL	go-reading-log-api-next/test	0.591s
?   	go-reading-log-api-next/test/fixtures	[no test files]
?   	go-reading-log-api-next/test/fixtures/dashboard	[no test files]
=== RUN   TestErrorScenarios
=== RUN   TestErrorScenarios/Day_Endpoint_-_Invalid_Date
    error_scenarios_test.go:86: Unknown endpoint: /v1/dashboard/day.json?date=invalid
=== RUN   TestErrorScenarios/Last_Days_-_Invalid_Type
    error_scenarios_test.go:86: Unknown endpoint: /v1/dashboard/last_days.json?type=99
=== RUN   TestErrorScenarios/Projects_Endpoint_-_Empty_Database
DEBUG: Raw JSON: {"data":[{"type":"dashboard_projects","attributes":[],"id":"1776967972"}]}
=== RUN   TestErrorScenarios/Day_Endpoint_-_Empty_Database
DEBUG: Raw JSON: {"data":{"type":"dashboard_day","attributes":{"total_pages":0,"log_count":0},"id":"1776967972"}}
=== RUN   TestErrorScenarios/Mean_Progress_-_Empty_Database
DEBUG: Raw JSON: {"data":{"type":"dashboard_echart_mean_progress","attributes":{"title":"Mean Progress","tooltip":{"trigger":"axis"},"legend":{"show":true,"data":["Progress"]},"series":[{"name":"Progress","type":"line","data":[],"itemStyle":{"color":[]}}],"xAxis":{"type":"category","name":"Date"},"yAxis":{"type":"value","name":"Progress (%)"}},"id":"1776967973"}}
    error_scenarios_test.go:166: 
        	Error Trace:	/home/danilo/scripts/github/go-reading-log-api-next/test/integration/error_scenarios_test.go:166
        	            				/home/danilo/scripts/github/go-reading-log-api-next/test/integration/error_scenarios_test.go:90
        	Error:      	Expected value not to be nil.
        	Test:       	TestErrorScenarios/Mean_Progress_-_Empty_Database
--- FAIL: TestErrorScenarios (0.82s)
```
<!-- SECTION:DESCRIPTION:END -->

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
