package fixtures

import (
	"fmt"
	"time"
)

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
	PerPages          float64
	MeanDay           float64
	SpecMeanDay       float64
	ProgressGeral     float64
}

// Pre-built Scenarios

// ScenarioZeroPages: Project with zero pages (edge case)
func ScenarioZeroPages() *Scenario {
	return &Scenario{
		Name:        "Zero Pages Edge Case",
		Description: "Project with total_page=0 for edge case testing",
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
		Description: "Project with all pages read",
		Projects: []*ProjectFixture{
			{
				ID:        2,
				Name:      "Completed Book",
				TotalPage: 300,
				Page:      300,
				Status:    "finished",
			},
		},
		Logs: []*LogFixture{
			{
				ID:        1,
				ProjectID: 2,
				Data:      time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
				StartPage: 0,
				EndPage:   300,
				WDay:      1, // Monday
			},
		},
		Expected: &ExpectedResults{
			Stats: &StatsExpectations{
				ProgressGeral: 100.0,
			},
		},
	}
}

// ScenarioMultipleProjects: Multiple projects with varying completion
func ScenarioMultipleProjects() *Scenario {
	return &Scenario{
		Name:        "Multiple Projects",
		Description: "3 projects: unstarted, running, finished",
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
		Logs: []*LogFixture{
			// Logs for running project
			{ID: 100, ProjectID: 11, Data: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC), StartPage: 0, EndPage: 25, WDay: 1},
			{ID: 101, ProjectID: 11, Data: time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC), StartPage: 25, EndPage: 50, WDay: 2},
		},
		Expected: &ExpectedResults{
			Stats: &StatsExpectations{
				ProgressGeral: 12.5, // (0 + 25 + 200) / (200 + 200 + 200) * 100
			},
		},
	}
}

// ScenarioFaultsByWeekday: Faults distributed across all weekdays
func ScenarioFaultsByWeekday() *Scenario {
	// Define fault distribution (example: more faults on weekends)
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

	// Generate fault fixtures
	var faults []*LogFixture
	logID := int64(200)
	baseDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	for weekday, count := range distribution {
		// Create multiple faults for this weekday
		for i := 0; i < count; i++ {
			targetDate := baseDate.AddDate(0, 0, -weekday)

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
		Description: "Faults distributed across all 7 weekdays",
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
		Description: "Logs for last_days endpoint testing",
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
	for i := 0; i < 15; i++ {
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
		Description: "Logs for speculate vs actual chart",
		Projects: []*ProjectFixture{
			{ID: 1, Name: "Speculate Project", TotalPage: 200, Page: 50},
		},
		Logs: logs,
		Expected: &ExpectedResults{
			LogsCount: 15,
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
		Description: "Varying daily progress for visual map testing",
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
	// Create logs for each week over 52 weeks
	for w := 0; w < 52; w++ {
		weekStart := baseDate.AddDate(0, 0, w*7)

		// Create multiple logs per week
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
		Description: "Logs spanning 52 weeks for yearly trend",
		Projects: []*ProjectFixture{
			{ID: 1, Name: "Yearly Project", TotalPage: 200, Page: 50},
		},
		Logs: logs,
		Expected: &ExpectedResults{
			LogsCount: 364, // 52 weeks * 7 days
		},
	}
}

// ScenarioEmptyData: Empty database state for error handling tests
func ScenarioEmptyData() *Scenario {
	return &Scenario{
		Name:        "Empty Data State",
		Description: "Empty database for testing zero-value responses",
		Projects:    []*ProjectFixture{},
		Logs:        []*LogFixture{},
		Expected: &ExpectedResults{
			Stats: &StatsExpectations{
				PreviousWeekPages: 0,
				LastWeekPages:     0,
				PerPages:          0.0,
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
		Description: "5 projects with different states and completion levels",
		Projects: []*ProjectFixture{
			{ID: 100, Name: "Unstarted", TotalPage: 200, Page: 0, Status: "unstarted"},
			{ID: 101, Name: "Running 25%", TotalPage: 200, Page: 50, Status: "running"},
			{ID: 102, Name: "Running 50%", TotalPage: 200, Page: 100, Status: "running"},
			{ID: 103, Name: "Sleeping", TotalPage: 200, Page: 75, Status: "sleeping"},
			{ID: 104, Name: "Finished", TotalPage: 200, Page: 200, Status: "finished"},
		},
		Logs: []*LogFixture{
			// Logs for running projects
			{ID: 200, ProjectID: 101, Data: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC), StartPage: 0, EndPage: 25, WDay: 1},
			{ID: 201, ProjectID: 101, Data: time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC), StartPage: 25, EndPage: 50, WDay: 2},
			{ID: 202, ProjectID: 102, Data: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC), StartPage: 0, EndPage: 50, WDay: 1},
			{ID: 203, ProjectID: 102, Data: time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC), StartPage: 50, EndPage: 75, WDay: 2},
			{ID: 204, ProjectID: 102, Data: time.Date(2024, 1, 17, 10, 0, 0, 0, time.UTC), StartPage: 75, EndPage: 100, WDay: 3},
		},
		Expected: &ExpectedResults{
			Stats: &StatsExpectations{
				ProgressGeral: 28.75, // (0 + 50 + 100 + 75 + 200) / (200 * 5) * 100
			},
		},
	}
}

// Helper functions

func createZeroStats() *ExpectedResults {
	return &ExpectedResults{
		Stats: &StatsExpectations{
			PreviousWeekPages: 0,
			LastWeekPages:     0,
			PerPages:          0.0,
			MeanDay:           0.0,
			SpecMeanDay:       0.0,
			ProgressGeral:     0.0,
		},
	}
}

func createFaultStats(totalFaults int) *ExpectedResults {
	return &ExpectedResults{
		Stats: &StatsExpectations{
			PreviousWeekPages: totalFaults * 10, // 10 pages per fault
			LastWeekPages:     totalFaults * 10,
			PerPages:          133.333,
			MeanDay:           9.67,
			SpecMeanDay:       11.12,
			ProgressGeral:     25.0, // 50/200
		},
	}
}

func stringPtr(s string) *string {
	return &s
}
