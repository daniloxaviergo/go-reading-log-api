---
id: RDL-093
title: Fix the test broken
status: To Do
assignee:
  - workflow
created_date: '2026-04-22 17:44'
updated_date: '2026-04-22 17:44'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
```go
=== RUN   TestSetupRoutes_MiddlewareChain
2026/04/22 14:42:25 INFO Request completed method=GET path=/healthz status=200 duration=16.428µs request_id=27ba93e3-75d6-4829-86fa-741667936659
--- PASS: TestSetupRoutes_MiddlewareChain (0.00s)
PASS
ok  	go-reading-log-api-next/internal/api/v1	(cached)
=== RUN   TestDashboardHandler_Day
DEBUG: Raw JSON: {"data":{"type":"dashboard_day","attributes":{"total_pages":100,"log_count":5},"id":"1705314600"}}
--- PASS: TestDashboardHandler_Day (0.00s)
=== RUN   TestDashboardHandler_Day_EmptyData
DEBUG: Raw JSON: {"data":{"type":"dashboard_day","attributes":{"total_pages":0,"log_count":0},"id":"1705746600"}}
--- PASS: TestDashboardHandler_Day_EmptyData (0.00s)
=== RUN   TestDashboardHandler_Day_InvalidDate
--- PASS: TestDashboardHandler_Day_InvalidDate (0.00s)
=== RUN   TestDashboardHandler_Projects
DEBUG: Raw JSON: {"data":[{"type":"dashboard_projects","attributes":[{"project_id":1,"project_name":"Project 1","total_pages":100,"log_count":5,"progress":100},{"project_id":2,"project_name":"Project 2","total_pages":200,"log_count":3,"progress":100}],"id":"1776879793"}]}
--- PASS: TestDashboardHandler_Projects (0.00s)
=== RUN   TestDashboardHandler_Projects_Empty
DEBUG: Raw JSON: {"data":[{"type":"dashboard_projects","attributes":[],"id":"1776879793"}]}
--- PASS: TestDashboardHandler_Projects_Empty (0.00s)
=== RUN   TestDashboardHandler_Faults
DEBUG: Raw JSON: {"data":{"type":"dashboard_echart_faults","attributes":{"title":"Fault Percentage","tooltip":{"formatter":"{a} \u003cbr/\u003e{b} : {c}%"},"series":[{"name":"Faults","type":"gauge","data":[80],"itemStyle":{"color":"#f44336"}}]},"id":"1776879793"}}
    dashboard_handler_test.go:275: 
        	Error Trace:	/home/danilo/scripts/github/go-reading-log-api-next/internal/api/v1/handlers/dashboard_handler_test.go:275
        	Error:      	Not equal: 
        	            	expected: "Faults Gauge"
        	            	actual  : "Fault Percentage"
        	            	
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -1 +1 @@
        	            	-Faults Gauge
        	            	+Fault Percentage
        	Test:       	TestDashboardHandler_Faults
--- FAIL: TestDashboardHandler_Faults (0.00s)
=== RUN   TestDashboardHandler_SpeculateActual
DEBUG: Raw JSON: {"data":{"type":"dashboard_echart_speculate_actual","attributes":{"title":"Speculated vs Actual Faults","tooltip":{"trigger":"axis"},"legend":{"show":true,"data":["Actual","Speculated"]},"series":[{"name":"Actual","type":"line","data":[50],"lineStyle":{"type":"solid","width":2}},{"name":"Speculated","type":"line","data":[57],"lineStyle":{"type":"dashed","width":2}}],"xAxis":{"type":"category","name":"Date"},"yAxis":{"type":"value","name":"Fault Count"}},"id":"1776879793"}}
--- PASS: TestDashboardHandler_SpeculateActual (0.00s)
=== RUN   TestDashboardHandler_WeekdayFaults
Validation error: weekday 2 is missing from output
    dashboard_handler_test.go:372: 
        	Error Trace:	/home/danilo/scripts/github/go-reading-log-api-next/internal/api/v1/handlers/dashboard_handler_test.go:372
        	Error:      	Not equal: 
        	            	expected: 200
        	            	actual  : 400
        	Test:       	TestDashboardHandler_WeekdayFaults
    dashboard_handler_test.go:379: 
        	Error Trace:	/home/danilo/scripts/github/go-reading-log-api-next/internal/api/v1/handlers/dashboard_handler_test.go:379
        	Error:      	Should be true
        	Test:       	TestDashboardHandler_WeekdayFaults
--- FAIL: TestDashboardHandler_WeekdayFaults (0.00s)
=== RUN   TestDashboardHandler_MeanProgress
--- FAIL: TestDashboardHandler_MeanProgress (0.00s)
panic: 
	assert: mock: I don't know what to return because the method call was unexpected.
		Either do Mock.On("GetLogsByDateRange").Return(...) first, or remove the GetLogsByDateRange() call.
		This method was unexpected:
			GetLogsByDateRange(context.backgroundCtx,time.Time,time.Time)
			0: context.backgroundCtx{emptyCtx:context.emptyCtx{}}
			1: time.Date(2026, time.March, 24, 0, 0, 0, 0, time.Local)
			2: time.Date(2026, time.April, 22, 0, 0, 0, 0, time.Local)
		at: [/home/danilo/scripts/github/go-reading-log-api-next/internal/api/v1/handlers/dashboard_handler_test.go:46 /home/danilo/scripts/github/go-reading-log-api-next/internal/service/dashboard/mean_progress_service.go:81 /home/danilo/scripts/github/go-reading-log-api-next/internal/service/dashboard/mean_progress_service.go:139 /home/danilo/scripts/github/go-reading-log-api-next/internal/api/v1/handlers/dashboard_handler.go:487 /home/danilo/scripts/github/go-reading-log-api-next/internal/api/v1/handlers/dashboard_handler_test.go:423] [recovered, repanicked]

goroutine 44 [running]:
testing.tRunner.func1.2({0x6ae920, 0xc0003040b0})
	/usr/lib/go/src/testing/testing.go:1872 +0x237
testing.tRunner.func1()
	/usr/lib/go/src/testing/testing.go:1875 +0x35b
panic({0x6ae920?, 0xc0003040b0?})
	/usr/lib/go/src/runtime/panic.go:783 +0x132
github.com/stretchr/testify/mock.(*Mock).fail(0xc00022ad70, {0x73982b?, 0x8?}, {0xc000229a40?, 0x3?, 0x3?})
	/home/danilo/go/pkg/mod/github.com/stretchr/testify@v1.11.1/mock/mock.go:359 +0x125
github.com/stretchr/testify/mock.(*Mock).MethodCalled(0xc00022ad70, {0x7ec031, 0x12}, {0xc0002dd320, 0x3, 0x3})
	/home/danilo/go/pkg/mod/github.com/stretchr/testify@v1.11.1/mock/mock.go:527 +0x77b
github.com/stretchr/testify/mock.(*Mock).Called(0xc00022ad70, {0xc0002dd320, 0x3, 0x3})
	/home/danilo/go/pkg/mod/github.com/stretchr/testify@v1.11.1/mock/mock.go:491 +0x125
go-reading-log-api-next/internal/api/v1/handlers.(*MockDashboardRepository).GetLogsByDateRange(0xc00022ad70, {0x7997b8, 0x9e11e0}, {0x0, 0xee153f530, 0x9c0580}, {0x0, 0xee17a30b0, 0x9c0580})
	/home/danilo/scripts/github/go-reading-log-api-next/internal/api/v1/handlers/dashboard_handler_test.go:46 +0x146
go-reading-log-api-next/internal/service/dashboard.(*MeanProgressService).GetMeanProgressData(0xc000070db0, {0x7997b8, 0x9e11e0})
	/home/danilo/scripts/github/go-reading-log-api-next/internal/service/dashboard/mean_progress_service.go:81 +0x7b
go-reading-log-api-next/internal/service/dashboard.(*MeanProgressService).GenerateChartConfig(0x14000000000000?, {0x7997b8?, 0x9e11e0?})
	/home/danilo/scripts/github/go-reading-log-api-next/internal/service/dashboard/mean_progress_service.go:139 +0x32
go-reading-log-api-next/internal/api/v1/handlers.(*DashboardHandler).MeanProgress(0x7997b8?, {0x798a08, 0xc000229940}, 0xc00021ded0?)
	/home/danilo/scripts/github/go-reading-log-api-next/internal/api/v1/handlers/dashboard_handler.go:487 +0x95
go-reading-log-api-next/internal/api/v1/handlers.TestDashboardHandler_MeanProgress(0xc000203dc0)
	/home/danilo/scripts/github/go-reading-log-api-next/internal/api/v1/handlers/dashboard_handler_test.go:423 +0x3d4
testing.tRunner(0xc000203dc0, 0x741ec0)
	/usr/lib/go/src/testing/testing.go:1934 +0xea
created by testing.(*T).Run in goroutine 1
	/usr/lib/go/src/testing/testing.go:1997 +0x465
FAIL	go-reading-log-api-next/internal/api/v1/handlers	0.008s
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
