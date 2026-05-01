package dashboard

import (
	"fmt"
	"time"
)

// floatPtr is a helper function to create a pointer to a float64
func floatPtr(f float64) *float64 {
	return &f
}

// Scenario represents a complete test scenario with all required data
type Scenario struct {
	Name        string
	Description string
	Projects    []*ProjectFixture
	Logs        []*LogFixture
	Expected    *ExpectedResults // Calculated values for validation
}

// ExpectedResults holds pre-calculated values for validation
type ExpectedResults struct {
	Stats      *StatsExpectations
	EchartData map[string]interface{}
	LogsCount  int
}

// StatsExpectations defines expected statistical values
type StatsExpectations struct {
	PreviousWeekPages int
	LastWeekPages     int
	PerPages          *float64
	MeanDay           float64
	SpecMeanDay       float64
	ProgressGeral     float64
}

// Pre-built Scenarios

// ScenarioZeroPages: Project with zero pages (edge case)
// Note: This scenario has no logs - use only for specific edge case testing
func ScenarioZeroPages() *Scenario {
	return &Scenario{
		Name:        "Zero Pages Edge Case",
		Description: "Project with total_page=0 for edge case testing. No logs - for edge case testing only.",
		Projects: []*ProjectFixture{
			{
				ID:        1,
				Name:      "Empty Project",
				TotalPage: 0,
				Page:      0,
				Status:    "unstarted",
			},
		},
		Logs:     []*LogFixture{},
		Expected: createZeroStats(),
	}
}

// ScenarioCompleteBook: Fully completed project
func ScenarioCompleteBook() *Scenario {
	return &Scenario{
		Name:        "Complete Book",
		Description: "Project with all pages read. 30-day data for validation.",
		Projects: []*ProjectFixture{
			{
				ID:        2,
				Name:      "Completed Book",
				TotalPage: 300,
				Page:      300,
				Status:    "finished",
			},
		},
		Logs: scenarioCompleteBookLogs(),
		Expected: &ExpectedResults{
			Stats: &StatsExpectations{
				ProgressGeral: 100.0,
			},
		},
	}
}

// scenarioCompleteBookLogs returns 30 days of logs for complete book scenario
func scenarioCompleteBookLogs() []*LogFixture {
	baseDate := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	logs := make([]*LogFixture, 30)

	for i := 0; i < 30; i++ {
		logs[i] = &LogFixture{
			ID:        int64(i + 1),
			ProjectID: 2,
			Data:      baseDate.AddDate(0, 0, -i),
			StartPage: 0,
			EndPage:   300 / 30 * (i + 1), // Evenly distribute pages
			WDay:      int(baseDate.AddDate(0, 0, -i).Weekday()),
		}
	}

	return logs
}

// ScenarioMultipleProjects: Multiple projects with varying completion
func ScenarioMultipleProjects() *Scenario {
	return &Scenario{
		Name:        "Multiple Projects",
		Description: "3 projects: unstarted, running, finished. 30-day data for validation.",
		Projects: []*ProjectFixture{
			{
				ID:        10,
				Name:      "Unstarted Project",
				TotalPage: 200,
				Page:      0,
				Status:    "unstarted",
			},
			{
				ID:        11,
				Name:      "Running Project",
				TotalPage: 200,
				Page:      50,
				Status:    "running",
			},
			{
				ID:        12,
				Name:      "Finished Project",
				TotalPage: 200,
				Page:      200,
				Status:    "finished",
			},
		},
		Logs: scenarioMultipleProjectsLogs(),
		Expected: &ExpectedResults{
			Stats: &StatsExpectations{
				PerPages:      nil,    // Null when no previous period data available (new behavior)
				ProgressGeral: 41.667, // (0 + 50 + 200) / (200 + 200 + 200) * 100 = 250/600 * 100
			},
		},
	}
}

// scenarioMultipleProjectsLogs returns 30 days of logs for multiple projects scenario
func scenarioMultipleProjectsLogs() []*LogFixture {
	baseDate := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	logs := make([]*LogFixture, 0)

	// Add logs for each project covering all 7 weekdays
	projectLogs := map[int64]int{
		10: 10, // Unstarted - 10 logs
		11: 15, // Running - 15 logs
		12: 10, // Finished - 10 logs
	}

	logID := int64(100)
	for projID, count := range projectLogs {
		for i := 0; i < count; i++ {
			wday := i % 7
			logs = append(logs, &LogFixture{
				ID:        logID,
				ProjectID: projID,
				Data:      baseDate.AddDate(0, 0, -wday),
				StartPage: i * 5,
				EndPage:   (i + 1) * 5,
				WDay:      wday,
			})
			logID++
		}
	}

	return logs
}

// ScenarioFaultsByWeekday: Faults distributed across all weekdays
func ScenarioFaultsByWeekday() *Scenario {
	distribution := map[int]int{
		0: 3, // Sunday
		1: 2, // Monday
		2: 1, // Tuesday
		3: 2, // Wednesday
		4: 1, // Thursday
		5: 3, // Friday
		6: 4, // Saturday
	}

	totalFaults := 0
	for _, count := range distribution {
		totalFaults += count
	}

	var faults []*LogFixture
	logID := int64(200)
	// Use a base date within the last 6 months from today
	// Today is 2026-04-27, so use 2026-02-01 as base (within 6-month range)
	// 2026-02-01 is a Sunday (weekday 0)
	baseDate := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)

	// Ensure we have at least 30 days of data by extending the date range
	for weekday, count := range distribution {
		for i := 0; i < count; i++ {
			// Find a date that actually falls on the target weekday
			// Start from base date and add days to reach the target weekday
			// baseDate is Sunday (0), so to get weekday N, add N days
			offsetToWeekday := weekday // Days from Sunday to reach target weekday
			weekOffset := i * 7        // Spread across different weeks
			targetDate := baseDate.AddDate(0, 0, offsetToWeekday+weekOffset)

			faults = append(faults, &LogFixture{
				ID:        int64(logID),
				ProjectID: 1,
				Data:      targetDate,
				StartPage: 0,
				EndPage:   10,
				WDay:      weekday,
				Note:      stringPtr(fmt.Sprintf("Fault on weekday %d", weekday)),
			})
			logID++
		}
	}

	return &Scenario{
		Name:        "Faults by Weekday",
		Description: "Faults distributed across all 7 weekdays with 30-day data range",
		Projects: []*ProjectFixture{
			{ID: 1, Name: "Faulty Project", TotalPage: 200, Page: 50},
		},
		Logs:     faults,
		Expected: createFaultStats(totalFaults),
	}
}

// ScenarioLastDays: Logs for last_days endpoint testing
func ScenarioLastDays() *Scenario {
	baseDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	var logs []*LogFixture
	for i := 0; i < 30; i++ {
		logs = append(logs, &LogFixture{
			ID:        int64(i + 1),
			ProjectID: 1,
			Data:      baseDate.AddDate(0, 0, -i),
			StartPage: i * 5,
			EndPage:   (i + 1) * 5,
			WDay:      int(baseDate.AddDate(0, 0, -i).Weekday()),
		})
	}

	return &Scenario{
		Name:        "Last Days Trend",
		Description: "Logs for last_days endpoint testing. 30 days of data.",
		Projects: []*ProjectFixture{
			{ID: 1, Name: "Trend Project", TotalPage: 200, Page: 50},
		},
		Logs: logs,
		Expected: &ExpectedResults{
			LogsCount: 30,
		},
	}
}

// ScenarioSpeculateActual: Logs for speculate vs actual chart
func ScenarioSpeculateActual() *Scenario {
	baseDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	var logs []*LogFixture
	// Extended to 30 days for validation
	for i := 0; i < 30; i++ {
		logs = append(logs, &LogFixture{
			ID:        int64(i + 1),
			ProjectID: 1,
			Data:      baseDate.AddDate(0, 0, -i),
			StartPage: i * 10,
			EndPage:   (i + 1) * 10,
			WDay:      int(baseDate.AddDate(0, 0, -i).Weekday()),
		})
	}

	return &Scenario{
		Name:        "Speculate Actual Test",
		Description: "Logs for speculate vs actual chart. 30 days of data.",
		Projects: []*ProjectFixture{
			{ID: 1, Name: "Speculate Project", TotalPage: 200, Page: 50},
		},
		Logs: logs,
		Expected: &ExpectedResults{
			LogsCount: 30,
		},
	}
}

// ScenarioMeanProgress: Logs for mean progress chart
func ScenarioMeanProgress() *Scenario {
	baseDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	var logs []*LogFixture

	for i := 0; i < 30; i++ {
		logs = append(logs, &LogFixture{
			ID:        int64(i + 1),
			ProjectID: 1,
			Data:      baseDate.AddDate(0, 0, -i),
			StartPage: (i * 5) % 100,
			EndPage:   ((i + 1) * 5) % 100,
			WDay:      int(baseDate.AddDate(0, 0, -i).Weekday()),
		})
	}

	return &Scenario{
		Name:        "Mean Progress Test",
		Description: "Varying daily progress for visual map testing. 30 days of data.",
		Projects: []*ProjectFixture{
			{ID: 1, Name: "Progress Project", TotalPage: 200, Page: 50},
		},
		Logs: logs,
		Expected: &ExpectedResults{
			LogsCount: 30,
		},
	}
}

// ScenarioYearlyTotal: Logs for yearly total chart
func ScenarioYearlyTotal() *Scenario {
	baseDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	var logs []*LogFixture
	for w := 0; w < 52; w++ {
		weekStart := baseDate.AddDate(0, 0, w*7)

		for d := 0; d < 7; d++ {
			logs = append(logs, &LogFixture{
				ID:        int64(w*7 + d + 1),
				ProjectID: 1,
				Data:      weekStart.AddDate(0, 0, d),
				StartPage: (w * 10) % 100,
				EndPage:   ((w + 1) * 10) % 100,
				WDay:      int(weekStart.AddDate(0, 0, d).Weekday()),
			})
		}
	}

	return &Scenario{
		Name:        "Yearly Total Test",
		Description: "Logs spanning 52 weeks for yearly trend. Extensive data for validation.",
		Projects: []*ProjectFixture{
			{ID: 1, Name: "Yearly Project", TotalPage: 200, Page: 50},
		},
		Logs: logs,
		Expected: &ExpectedResults{
			LogsCount: 364,
		},
	}
}

// ScenarioEmptyData: Empty database state for error handling tests
func ScenarioEmptyData() *Scenario {
	return &Scenario{
		Name:        "Empty Data State",
		Description: "Empty database for testing zero-value responses. No logs - for error handling tests only.",
		Projects:    []*ProjectFixture{},
		Logs:        []*LogFixture{},
		Expected: &ExpectedResults{
			Stats: &StatsExpectations{
				PreviousWeekPages: 0,
				LastWeekPages:     0,
				PerPages:          floatPtr(0.0),
				MeanDay:           0.0,
				SpecMeanDay:       0.0,
				ProgressGeral:     0.0,
			},
		},
	}
}

// ScenarioMultipleProjectsExtended: Extended multiple projects scenario
func ScenarioMultipleProjectsExtended() *Scenario {
	return &Scenario{
		Name:        "Multiple Projects Extended",
		Description: "5 projects with different states and completion levels. 30-day data for validation.",
		Projects: []*ProjectFixture{
			{ID: 100, Name: "Unstarted", TotalPage: 200, Page: 0, Status: "unstarted"},
			{ID: 101, Name: "Running 25%", TotalPage: 200, Page: 50, Status: "running"},
			{ID: 102, Name: "Running 50%", TotalPage: 200, Page: 100, Status: "running"},
			{ID: 103, Name: "Sleeping", TotalPage: 200, Page: 75, Status: "sleeping"},
			{ID: 104, Name: "Finished", TotalPage: 200, Page: 200, Status: "finished"},
		},
		Logs: scenarioMultipleProjectsExtendedLogs(),
		Expected: &ExpectedResults{
			Stats: &StatsExpectations{
				ProgressGeral: 28.75,
			},
		},
	}
}

// scenarioMultipleProjectsExtendedLogs returns 30 days of logs for extended scenario
func scenarioMultipleProjectsExtendedLogs() []*LogFixture {
	baseDate := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	logs := make([]*LogFixture, 0)

	// Add logs for each project covering all 7 weekdays
	projectLogs := map[int64]int{
		100: 5, // Unstarted - 5 logs
		101: 8, // Running 25% - 8 logs
		102: 8, // Running 50% - 8 logs
		103: 7, // Sleeping - 7 logs
		104: 5, // Finished - 5 logs
	}

	logID := int64(200)
	for projID, count := range projectLogs {
		for i := 0; i < count; i++ {
			wday := i % 7
			logs = append(logs, &LogFixture{
				ID:        logID,
				ProjectID: projID,
				Data:      baseDate.AddDate(0, 0, -wday),
				StartPage: i * 5,
				EndPage:   (i + 1) * 5,
				WDay:      wday,
			})
			logID++
		}
	}

	return logs
}

func createZeroStats() *ExpectedResults {
	return &ExpectedResults{
		Stats: &StatsExpectations{
			PreviousWeekPages: 0,
			LastWeekPages:     0,
			PerPages:          floatPtr(0.0),
			MeanDay:           0.0,
			SpecMeanDay:       0.0,
			ProgressGeral:     0.0,
		},
	}
}

func createFaultStats(totalFaults int) *ExpectedResults {
	return &ExpectedResults{
		Stats: &StatsExpectations{
			PreviousWeekPages: totalFaults * 10,
			LastWeekPages:     totalFaults * 10,
			PerPages:          floatPtr(280.0),
			MeanDay:           9.67,
			SpecMeanDay:       11.12,
			ProgressGeral:     25.0,
		},
	}
}

func stringPtr(s string) *string {
	return &s
}
