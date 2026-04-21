package testdata

import (
	"time"

	"go-reading-log-api-next/internal/domain/dto"
)

// Project450ExpectedLog contains expected log data for Project 450
// Based on: test/data/project-450-go-logs.json and test/data/project-450-rails-logs.json
var Project450ExpectedLogs = []ExpectedLog{
	{
		ID:        1,
		ProjectID: 450,
		Data:      "2026-02-19T00:00:00Z",
		StartPage: 0,
		EndPage:   25,
		Note:      stringPtr("First reading session"),
	},
	{
		ID:        2,
		ProjectID: 450,
		Data:      "2026-02-20T00:00:00Z",
		StartPage: 25,
		EndPage:   50,
		Note:      stringPtr("Continued reading"),
	},
	{
		ID:        3,
		ProjectID: 450,
		Data:      "2026-02-21T00:00:00Z",
		StartPage: 50,
		EndPage:   75,
		Note:      stringPtr("Progressing well"),
	},
}

// Project450LogCount is the expected number of logs for Project 450
const Project450LogCount = 38

// Project450ExpectedValues contains all expected values for Project 450
var Project450ExpectedValues = ExpectedProject{
	ID:        450,
	Name:      "História da Igreja VIII.1",
	TotalPage: 691,
	Page:      691,
	StartedAt: timeToString(time.Date(2026, 2, 19, 0, 0, 0, 0, time.UTC)),
	Progress:  floatPtr(100.0),
	Status:    stringPtr("finished"),
	LogsCount: intPtr(Project450LogCount),
	DaysUnread: func() *int {
		// Calculated from last log date to current date
		// Using BRT timezone to match Rails behavior
		lastLogDate := time.Date(2026, 4, 18, 0, 0, 0, 0, time.UTC)
		now := time.Now()
		tzLocation := time.FixedZone("BRT", -3*60*60)
		nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, tzLocation)
		lastReadDate := time.Date(lastLogDate.Year(), lastLogDate.Month(), lastLogDate.Day(), 0, 0, 0, 0, tzLocation)
		days := int(nowDate.Sub(lastReadDate).Hours() / 24)
		if days < 0 {
			days = 0
		}
		return &days
	}(),
	MedianDay: func() *float64 {
		// median_day = page / days_reading
		// Using BRT timezone
		startedAt := time.Date(2026, 2, 19, 0, 0, 0, 0, time.UTC)
		now := time.Now()
		tzLocation := time.FixedZone("BRT", -3*60*60)
		nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, tzLocation)
		startedAtDate := time.Date(startedAt.Year(), startedAt.Month(), startedAt.Day(), 0, 0, 0, 0, tzLocation)
		daysReading := int(nowDate.Sub(startedAtDate).Hours() / 24)
		if daysReading <= 0 {
			result := 0.0
			return &result
		}
		medianDay := float64(691) / float64(daysReading)
		rounded := medianDay // No rounding needed for this specific calculation
		return &rounded
	}(),
	FinishedAt: nil, // null for completed projects
	Logs:       Project450ExpectedLogs,
}

// GetProject450Logs converts ExpectedLog to dto.LogResponse for testing
func GetProject450Logs() []*dto.LogResponse {
	logs := make([]*dto.LogResponse, len(Project450ExpectedLogs))
	for i, el := range Project450ExpectedLogs {
		dataTime, _ := time.Parse(time.RFC3339, el.Data)
		logs[i] = &dto.LogResponse{
			ID:        el.ID,
			Data:      &dataTime,
			StartPage: el.StartPage,
			EndPage:   el.EndPage,
			Note:      el.Note,
		}
	}
	return logs
}

// GetProject450ExpectedValues returns a copy of Project450ExpectedValues
func GetProject450ExpectedValues() ExpectedProject {
	return Project450ExpectedValues
}
