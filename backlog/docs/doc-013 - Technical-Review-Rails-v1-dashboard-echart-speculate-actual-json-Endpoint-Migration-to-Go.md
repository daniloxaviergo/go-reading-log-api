# Technical Review: Rails `/v1/dashboard/echart/speculate_actual.json` Endpoint Migration to Go

**Date:** 2026-05-10  
**Status:** ⚠️ **Critical Implementation Gaps Identified**  
**Priority:** High - Blocks dashboard visualization functionality

---

## Executive Summary

The current Go implementation of the `SpeculateActual` endpoint **does not match** the Rails behavior. There are fundamental differences in:

1. **Response format** - Go uses JSON:API envelope, Rails uses flat `{ echart: {...} }`
2. **Data calculation** - Go calculates faults, Rails calculates reading pages by weekday
3. **Date range** - Go uses 30 days, Rails uses 15 days (14 days ago to today)
4. **Series data** - Go uses fault counts, Rails uses pages read with speculative mean
5. **X-axis format** - Go uses dates, Rails uses `'10-05 (Sat)'` format
6. **Mark elements** - Go missing markPoint/markLine configurations

---

## 1. Ambiguous Technical Requirements

### 1.1 Response Format Clarification

**Current Go Implementation:**
```json
{
  "data": {
    "type": "dashboard_echart_speculate_actual",
    "id": "123456",
    "attributes": {
      "echart": { ... }
    }
  }
}
```

**Rails Expected Format:**
```json
{
  "echart": {
    "tooltip": { "trigger": "axis" },
    "grid": { ... },
    "xAxis": { ... },
    "yAxis": { ... },
    "series": [ ... ]
  }
}
```

**Recommendation:** Change to flat `{ echart: {...} }` format to match Rails.

### 1.2 "Speculate" vs "Faults" Terminology

**Ambiguity:** The Rails implementation calculates **pages read** (reading progress), not faults. The current Go implementation calculates **fault counts**.

**Rails Calculation Flow:**
1. Get logs from 14 days ago to today (15 days total)
2. Group by weekday (0-6)
3. Calculate mean pages per 7-day interval for each weekday
4. Apply 15% speculative multiplier (`per_extimativa` config)
5. Return actual pages vs speculative mean

**Current Go Implementation:**
- Uses `GetFaultsByDateRange` (fault counting logic)
- Calculates predicted faults: `faults * (1 + 0.15)`
- **Incorrect domain logic for reading progress tracking**

**Clarification Needed:** Should this endpoint track:
- [ ] Reading pages (Rails behavior)
- [ ] Fault counts (current Go implementation)
- [ ] Both (dual-series chart)

**Recommendation:** Align with Rails behavior - track **pages read**, not faults.

### 1.3 Date Range Specification

| Aspect | Rails | Go (Current) |
|--------|-------|--------------|
| Start Date | 14 days ago | 30 days ago |
| End Date | Today | Today |
| Total Days | 15 | 30 |

**Recommendation:** Change to 15-day range to match Rails:
```go
endDate := time.Now()
startDate := endDate.AddDate(0, 0, -14) // 15 days including today
```

---

## 2. Missing Repository Methods

### 2.1 Required Methods Not Present

The Rails implementation requires these repository methods that are **missing** in Go:

| Method | Purpose | Status |
|--------|---------|--------|
| `GetLogsGroupedByWeekday(ctx, startDate, endDate)` | Get logs grouped by weekday (0-6) with page sums | ❌ Missing |
| `GetMeanPagesByWeekday(ctx, weekday, startDate, endDate)` | Calculate mean pages for specific weekday | ❌ Missing |
| `GetFirstLogDate(ctx) (time.Time, error)` | Get earliest log date for interval calculation | ❌ Missing |

### 2.2 Current Repository Methods (Available)

```go
// Available but not suitable for this endpoint
GetLogsByDateRange(ctx, start, end) ([]*dto.LogEntry, error)
GetFaultsByDateRange(ctx, start, end) (*dto.FaultStats, error)
GetMeanByWeekday(ctx, weekday) (*float64, error) // Different algorithm
```

**Gap Analysis:**
- `GetLogsByDateRange` returns individual logs, not grouped/ aggregated data
- `GetMeanByWeekday` uses different algorithm (project-level, not global)
- No method for weekday grouping with cumulative mean calculation

### 2.3 Recommended Repository Methods

```go
// Add to DashboardRepository interface

// GetWeekdayPagesGrouped returns pages read grouped by weekday for date range
// Returns map[weekday]float64 where key is 0-6 (Sunday-Saturday)
GetWeekdayPagesGrouped(ctx context.Context, start, end time.Time) (map[int]float64, error)

// GetWeekdayMeanWithIntervals calculates mean pages per 7-day interval for each weekday
// Algorithm: total_pages / count_reads where count_reads = floor((day - begin_data) / 7)
// Returns map[weekday]float64
GetWeekdayMeanWithIntervals(ctx context.Context, start, end, beginData time.Time) (map[int]float64, error)

// GetFirstLogDate returns the date of the first log entry in the system
GetFirstLogDate(ctx context.Context) (time.Time, error)

// GetLogsByDateRangeGrouped returns logs with weekday grouping and page sums
// More efficient than GetLogsByDateRange + client-side grouping
GetLogsByDateRangeGrouped(ctx context.Context, start, end time.Time) ([]*dto.WeekdayLogGroup, error)
```

### 2.4 New DTOs Required

```go
// WeekdayLogGroup represents aggregated log data by weekday
type WeekdayLogGroup struct {
    Weekday    int                  // 0-6 (Sunday-Saturday)
    TotalPages int                  // Sum of pages for this weekday
    LogCount   int                  // Number of log entries
    Dates      []time.Time          // Individual dates with data
}

// SpeculateActualData represents the comparative data structure
type SpeculateActualData struct {
    Date       string  `json:"date"`        // YYYY-MM-DD
    Pages      int     `json:"pages"`       // Actual pages read
    Mean       float64 `json:"mean"`        // Mean pages for weekday
    SpecMean   float64 `json:"spec_mean"`   // Mean * (1 + per_extimativa)
}
```

---

## 3. Current Go Implementation Issues

### 3.1 Handler Implementation Problems

**File:** `internal/api/v1/handlers/dashboard_handler.go:559-622`

**Issues:**

1. **Wrong data source:** Uses `GetFaultsByDateRange` instead of log page data
2. **Wrong date range:** 30 days instead of 15 days
3. **Wrong series names:** 'Actual'/'Speculated' instead of 'Pages'/'Mean'
4. **Missing mark elements:** No markPoint (max/min) or markLine (pages_per_day)
5. **Wrong X-axis format:** Uses raw dates instead of `'10-05 (Sat)'`
6. **JSON:API envelope:** Wraps response in envelope instead of flat `{ echart: ... }`
7. **Missing grid config:** Rails includes `grid: { left: '3%', right: '4%', ... }`
8. **Missing smooth flag:** Rails uses `smooth: true` for line curves

### 3.2 ECharts DTO Limitations

**File:** `internal/domain/dto/dashboard_response.go`

**Missing ECharts Configurations:**

```go
// Series struct missing:
MarkPoint map[string]interface{} `json:"markPoint,omitempty"`  // ❌ Missing
MarkLine  map[string]interface{} `json:"markLine,omitempty"`   // ❌ Missing

// Axis struct missing:
BoundaryGap interface{} `json:"boundaryGap,omitempty"` // ❌ Missing (should be false)

// Grid struct should be populated:
Grid *Grid `json:"grid,omitempty"` // ✅ Present but not used
```

### 3.3 Service Layer Missing

**No service layer exists** for SpeculateActual calculations. All logic is in handler, violating Clean Architecture.

**Recommended Structure:**
```
internal/service/dashboard/
├── speculate_actual_service.go    # Business logic
├── speculate_actual_service_test.go
```

---

## 4. Test Coverage Gaps

### 4.1 Missing Test Files

| File | Purpose | Status |
|------|---------|--------|
| `internal/api/v1/handlers/dashboard_handler_speculate_actual_test.go` | Handler integration tests | ❌ Missing |
| `internal/service/dashboard/speculate_actual_service_test.go` | Service layer unit tests | ❌ Missing |
| `internal/adapter/postgres/speculate_actual_repository_test.go` | Repository integration tests | ❌ Missing |

### 4.2 Test Scenarios Not Covered

**Based on Rails behavior, tests should cover:**

1. **Date Range Calculations**
   - 15-day range (14 days ago to today)
   - Edge case: Less than 15 days of logs exist
   - Edge case: No logs in range

2. **Weekday Grouping**
   - Correct mapping of dates to weekdays (0-6)
   - Sunday = 0, Saturday = 6
   - DST boundary handling

3. **Mean Calculations**
   - Mean = total_pages / count_reads (7-day intervals)
   - count_reads = floor((day - begin_data) / 7)
   - Division by zero handling

4. **Speculative Mean**
   - spec_mean = mean * (1 + per_extimativa)
   - Config value retrieval (default 15%)
   - Rounding to 3 decimals

5. **X-axis Format**
   - Format: `'10-05 (Sat)'` (DD-MM (DDD))
   - Timezone handling
   - Locale-specific day names

6. **Mark Elements**
   - markPoint: max/min values present
   - markLine: pages_per_day from config present

7. **Response Format**
   - Flat `{ echart: {...} }` structure
   - No JSON:API envelope
   - Content-Type: `application/json`

### 4.3 Existing Test Coverage

**Current:** `dashboard_handler_test.go` has tests for:
- ✅ `Day` endpoint
- ✅ `Projects` endpoint
- ✅ `ProjectsWithLogs` endpoint
- ✅ `LastDays` endpoint
- ✅ `Faults` endpoint
- ✅ `WeekdayFaults` endpoint
- ✅ `MeanProgress` endpoint
- ✅ `YearlyTotal` endpoint
- ❌ **`SpeculateActual` endpoint** - No tests exist

---

## 5. Edge Cases Requiring Special Handling

### 5.1 Empty Data Scenarios

| Scenario | Expected Behavior |
|----------|-------------------|
| No logs in system | Return empty series arrays, valid chart config |
| Logs exist but none in 15-day range | Return empty series arrays |
| Only 1 day of logs in range | Single data point in series |
| First log < 7 days ago | count_reads = 1 for all weekdays |
| per_extimativa config = 0 | spec_mean = mean (no multiplier) |

### 5.2 Division by Zero

**Scenarios:**
- `count_reads = 0` when `(day - begin_data) < 7 days`
- `prev_period_mean = 0` in ratio calculations

**Rails Behavior:** Returns `nil`/`nil` for mean when denominator is zero

**Go Implementation:** Should return `*float64` (nullable pointer) and handle `nil` in calculations

### 5.3 Timezone Handling

**Rails:** Uses `Time.zone.now` (application timezone)

**Go:** Uses `time.Now()` (system timezone)

**Risk:** X-axis dates may differ if server timezone ≠ Rails app timezone

**Recommendation:** Use UTC consistently or configure timezone in app settings

### 5.4 Weekday Boundary

**PostgreSQL:** `EXTRACT(DOW FROM timestamp)` returns 0=Sunday, 6=Saturday

**Go:** `time.Weekday()` returns 0=Sunday, 6=Saturday

**Status:** ✅ Compatible - No conversion needed

### 5.5 Date Formatting

**Rails:** `day.to_date.strftime('%d-%m (%a)')`

**Go Equivalent:**
```go
day.Format("02-01 (Mon)") // Produces "10-05 (Sat)"
```

**Note:** Month format is `01` not `05` - verify Rails uses `%m` (month) not `%d` (day) in second position

---

## 6. Technical Feasibility Assessment

### 6.1 Implementation Complexity

| Component | Complexity | Effort Estimate |
|-----------|------------|-----------------|
| Repository methods | Medium | 4-6 hours |
| Service layer | Medium | 3-4 hours |
| Handler rewrite | Low | 2-3 hours |
| DTO extensions | Low | 1-2 hours |
| Test coverage | High | 6-8 hours |
| **Total** | | **16-23 hours** |

### 6.2 Risk Assessment

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Data mismatch with Rails | High | Critical | Side-by-side comparison tests |
| Timezone inconsistencies | Medium | High | Explicit UTC usage |
| Mean calculation errors | Medium | High | Unit tests with known inputs |
| Config retrieval failures | Low | Medium | Default fallback values |
| Performance on large datasets | Low | Medium | Indexed queries on `data` column |

### 6.3 Dependencies

**Required:**
- ✅ PostgreSQL connection pool (existing)
- ✅ UserConfig service (existing)
- ✅ DTO infrastructure (existing)
- ⚠️ New repository methods (to be implemented)
- ⚠️ Service layer (to be implemented)

**No external dependencies needed.**

---

## 7. Recommendations

### 7.1 Immediate Actions (Priority 1)

1. **Clarify Domain Logic**
   - Confirm with stakeholders: Should endpoint track pages or faults?
   - Document decision in PRD

2. **Implement Missing Repository Methods**
   - `GetWeekdayPagesGrouped`
   - `GetFirstLogDate`
   - `GetWeekdayMeanWithIntervals`

3. **Extend ECharts DTOs**
   - Add `MarkPoint` and `MarkLine` to `Series` struct
   - Add `Grid` configuration support

4. **Create Service Layer**
   - `SpeculateActualService` with business logic
   - Dependency injection pattern (consistent with existing)

### 7.2 Implementation Steps (Priority 2)

1. **Rewrite Handler**
   - Change response format to `{ echart: {...} }`
   - Update date range to 15 days
   - Add markPoint/markLine configurations
   - Fix X-axis formatting

2. **Add Comprehensive Tests**
   - Unit tests for service layer
   - Integration tests for repository
   - Handler tests with mock data
   - Side-by-side Rails comparison tests

3. **Documentation**
   - Update API documentation
   - Add example curl requests
   - Document calculation algorithms

### 7.3 Long-term Improvements (Priority 3)

1. **Performance Optimization**
   - Add database indexes on `logs.data` and `logs.wday`
   - Consider materialized views for complex aggregations

2. **Monitoring**
   - Add metrics for calculation performance
   - Log cache hits/misses for config values

3. **Refactoring**
   - Consider query builder pattern for complex SQL
   - Extract weekday calculation logic to utility package

---

## 8. Acceptance Criteria

### 8.1 Functional Requirements

- [ ] Response format matches Rails: `{ echart: {...} }`
- [ ] Date range is 15 days (14 days ago to today)
- [ ] Series names are 'Pages' and 'Mean'
- [ ] X-axis format is '10-05 (Sat)'
- [ ] markPoint includes max and min
- [ ] markLine includes pages_per_day from config
- [ ] Config value per_extimativa = 10% (or configurable)
- [ ] Weekday grouping uses 0-6 (Sunday-Saturday)
- [ ] Mean calculation matches Rails algorithm
- [ ] Speculative mean = mean * (1 + per_extimativa)

### 8.2 Non-Functional Requirements

- [ ] Response time < 500ms for 15-day range
- [ ] Handles empty data gracefully
- [ ] Handles division by zero
- [ ] Timezone-safe calculations
- [ ] Test coverage > 80%
- [ ] No breaking changes to existing endpoints

---

## 9. Related Tasks

- **RDL-111**: Updated StatsData DTO with nullable fields ✅
- **RDL-118**: Implemented null handling in service layer ✅
- **RDL-119**: Added comprehensive test coverage for null validation ✅
- **RDL-XXX**: [New] Implement SpeculateActual repository methods
- **RDL-XXX**: [New] Implement SpeculateActual service layer
- **RDL-XXX**: [New] Rewrite SpeculateActual handler
- **RDL-XXX**: [New] Add ECharts markPoint/markLine support

---

## 10. Appendix

### 10.1 Rails Source Reference

**Controller:** `rails-app/app/controllers/v1/dashboard/echart/speculate_actual_controller.rb`

```ruby
class V1::Dashboard::Echart::SpeculateActualController < ApplicationController
  def index
    previous_week = 15.days.ago.to_date
    logs  = Log.range_data(previous_week)
    wdays = V1::GroupLog.new(logs).by_wday

    hoje, previous_week = (Time.zone.now.end_of_day), 14.days.ago.to_date
    spec_efec = V1::Dashboard::SpeculateActual.new(previous_week, hoje, wdays)
    echart = ::V1::Dashboard::Echart::SpeculateActual.new(spec_efec.comparative)

    render json: { echart: echart.graph }
  end
end
```

**Business Logic:** `rails-app/app/classes/v1/dashboard/speculate_actual.rb`

```ruby
class V1::Dashboard::SpeculateActual
  attr_reader :comparative
  
  def initialize(date_start, date_end, wdays)
    config = ::V1::UserConfig.new.get

    begin_data = ::Log.order(data: :asc).first.data.to_date
    logs       = ::Log.range_data(date_start, date_end)
    date_end   = date_end.to_date

    @comparative = (date_start..date_end).each_with_object({}) do |day, means|
      equal_day  = lambda { |log| log[:data].to_date == day }
      read_pages = logs.select(&equal_day).map { |l| l.read_pages }

      wdays[day.wday] << read_pages.sum
      reads_wday  = wdays[day.wday].flatten
      total_pages = reads_wday.sum.to_f
      count_reads = (begin_data..day).step(7).map { |d| d }.size.to_f

      mean = (total_pages / count_reads).round(3)
      spec_mean = ((mean * config.per_extimativa) + mean).round(3)

      means[day.to_s] = {
        pages:     read_pages.sum,
        mean:      mean,
        spec_mean: spec_mean
      }
    end
  end
end
```

**ECharts Config:** `rails-app/app/classes/v1/dashboard/echart/speculate_actual.rb`

```ruby
class V1::Dashboard::Echart::SpeculateActual
  
  def initialize(comparative)
    @config = ::V1::UserConfig.new.get
    @comparative = comparative
  end

  def graph
    {
      tooltip: { trigger: 'axis' },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      yAxis: { type: 'value' },
      xAxis: { type: 'category', boundaryGap: false, data: xAxis.compact },
      series: [
        {
          name: 'Pages',
          type: 'line',
          smooth: true,
          itemStyle: { normal: { areaStyle: { type: 'default' } } },
          data: series_data_page,
          markPoint: { data: [{ type: 'max' }, { type: 'min' }] },
          markLine: { data: [{ name: 'pages_per_day', yAxis: @config.pages_per_day }] }
        },
        {
          name: 'Mean',
          type: 'line',
          smooth: true,
          itemStyle: { normal: { areaStyle: { type: 'default' } } },
          data: series_data_mean
        }
      ]
    }
  end

  private 

  def series_data_page
    @comparative.values.map { |data| data[:pages] }
  end

  def series_data_mean
    @comparative.values.map { |data| data[:spec_mean] }
  end

  def xAxis
    @comparative.keys.map { |day| day.to_date.strftime('%d-%m (%a)') }
  end
end
```

### 10.2 Expected Response Example

```json
{
  "echart": {
    "tooltip": {
      "trigger": "axis"
    },
    "grid": {
      "left": "3%",
      "right": "4%",
      "bottom": "3%",
      "containLabel": true
    },
    "yAxis": {
      "type": "value"
    },
    "xAxis": {
      "type": "category",
      "boundaryGap": false,
      "data": ["10-05 (Sat)", "11-05 (Sun)", "12-05 (Mon)", "..."]
    },
    "series": [
      {
        "name": "Pages",
        "type": "line",
        "smooth": true,
        "itemStyle": {
          "normal": {
            "areaStyle": {
              "type": "default"
            }
          }
        },
        "data": [0, 45, 30, 52, 0, 38, 25],
        "markPoint": {
          "data": [
            { "type": "max", "name": "" },
            { "type": "min", "name": "" }
          ]
        },
        "markLine": {
          "lineStyle": {
            "normal": {
              "color": "#333"
            }
          },
          "data": [
            { "name": "pages_per_day", "yAxis": 35.5 }
          ]
        }
      },
      {
        "name": "Mean",
        "type": "line",
        "smooth": true,
        "itemStyle": {
          "normal": {
            "areaStyle": {
              "type": "default"
            }
          }
        },
        "data": [0, 51.75, 34.5, 59.8, 0, 43.7, 28.75]
      }
    ]
  }
}
```

---

*Last updated: 2026-05-10*  
*Author: Technical Lead*  
*Review status: Pending stakeholder confirmation on domain logic*