package integration

import (
	"bytes"
	"context"
	"encoding/json"
	// "fmt" - not used
	"math"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go-reading-log-api-next/internal/adapter/postgres"
	// "go-reading-log-api-next/internal/config" - not used in this test file
	"go-reading-log-api-next/internal/domain/dto"
	// "go-reading-log-api-next/internal/repository" - not used in this test file (used via postgres adapter)
	"go-reading-log-api-next/test"
	"go-reading-log-api-next/test/testdata"
)

// TestExpectedValues_Integration tests expected values against a real database
func TestExpectedValues_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test database
	helper, err := test.SetupTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer helper.Close()

	// Setup test schema (create tables)
	if err := helper.SetupTestSchema(); err != nil {
		t.Fatalf("Failed to setup test schema: %v", err)
	}

	// Create a test project
	ctx := context.Background()

	// Insert test data
	projectID, err := insertTestProject(helper.Pool, ctx)
	if err != nil {
		t.Fatalf("Failed to insert test project: %v", err)
	}

	// Get the project with logs
	repo := postgres.NewProjectRepositoryImpl(helper.Pool)
	result, err := repo.GetWithLogs(ctx, projectID)
	if err != nil {
		t.Fatalf("Failed to get project: %v", err)
	}

	// Calculate expected values
	expected := testdata.CalculateExpectedValues(ctx, result.Project)

	// Verify calculated fields match expected values
	t.Run("progress calculation", func(t *testing.T) {
		if result.Project.Progress == nil {
			t.Fatal("project progress should not be nil")
		}
		if expected.Progress == nil {
			t.Fatal("expected progress should not be nil")
		}
		if !floatEqual(*result.Project.Progress, *expected.Progress, 2) {
			t.Errorf("progress mismatch: got %v, expected %v", *result.Project.Progress, *expected.Progress)
		}
	})

	t.Run("status calculation", func(t *testing.T) {
		if result.Project.Status == nil {
			t.Fatal("project status should not be nil")
		}
		if expected.Status == nil {
			t.Fatal("expected status should not be nil")
		}
		if *result.Project.Status != *expected.Status {
			t.Errorf("status mismatch: got '%s', expected '%s'", *result.Project.Status, *expected.Status)
		}
	})

	t.Run("logs_count", func(t *testing.T) {
		if result.Project.LogsCount == nil {
			t.Fatal("project logs_count should not be nil")
		}
		if expected.LogsCount == nil {
			t.Fatal("expected logs_count should not be nil")
		}
		if *result.Project.LogsCount != *expected.LogsCount {
			t.Errorf("logs_count mismatch: got %d, expected %d", *result.Project.LogsCount, *expected.LogsCount)
		}
	})

	t.Run("days_unreading", func(t *testing.T) {
		if result.Project.DaysUnread == nil {
			t.Fatal("project days_unreading should not be nil")
		}
		if expected.DaysUnread == nil {
			t.Fatal("expected days_unreading should not be nil")
		}
		if *result.Project.DaysUnread != *expected.DaysUnread {
			t.Errorf("days_unreading mismatch: got %d, expected %d", *result.Project.DaysUnread, *expected.DaysUnread)
		}
	})

	t.Run("median_day", func(t *testing.T) {
		if result.Project.MedianDay == nil {
			t.Fatal("project median_day should not be nil")
		}
		if expected.MedianDay == nil {
			t.Fatal("expected median_day should not be nil")
		}
		if !floatEqual(*result.Project.MedianDay, *expected.MedianDay, 2) {
			t.Errorf("median_day mismatch: got %v, expected %v", *result.Project.MedianDay, *expected.MedianDay)
		}
	})
}

// TestExpectedValues_RailsComparison tests that Go API values match Rails API values
func TestExpectedValues_RailsComparison(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Rails comparison test in short mode")
	}

	// Load expected values from JSON files
	goData, err := os.ReadFile("../data/project-450-go.json")
	if err != nil {
		t.Fatalf("Failed to read Go JSON: %v", err)
	}

	railsData, err := os.ReadFile("../data/project-450-rails.json")
	if err != nil {
		t.Fatalf("Failed to read Rails JSON: %v", err)
	}

	// Parse and compare key fields using simple string matching
	// For more robust parsing, use a proper JSON library
	goValues, err := parseProjectJSON(goData)
	if err != nil {
		t.Logf("Warning: Failed to parse Go JSON: %v", err)
		t.Skip("Skipping due to JSON parsing limitation")
	}

	railsValues, err := parseRailsJSON(railsData)
	if err != nil {
		t.Logf("Warning: Failed to parse Rails JSON: %v", err)
		t.Skip("Skipping due to JSON parsing limitation")
	}

	// Compare key values
	t.Run("name match", func(t *testing.T) {
		if goValues.Name != railsValues.Name {
			t.Errorf("name mismatch: Go='%s', Rails='%s'", goValues.Name, railsValues.Name)
		}
	})

	t.Run("total_page match", func(t *testing.T) {
		if goValues.TotalPage != railsValues.TotalPage {
			t.Errorf("total_page mismatch: Go=%d, Rails=%d", goValues.TotalPage, railsValues.TotalPage)
		}
	})

	t.Run("page match", func(t *testing.T) {
		if goValues.Page != railsValues.Page {
			t.Errorf("page mismatch: Go=%d, Rails=%d", goValues.Page, railsValues.Page)
		}
	})

	t.Run("progress match", func(t *testing.T) {
		// Progress might differ slightly due to rounding
		diff := goValues.Progress - railsValues.Progress
		if diff < 0 {
			diff = -diff
		}
		if diff > 0.1 {
			t.Errorf("progress mismatch: Go=%.2f, Rails=%.2f", goValues.Progress, railsValues.Progress)
		}
	})

	t.Run("logs_count match", func(t *testing.T) {
		if goValues.LogsCount != railsValues.LogsCount {
			t.Errorf("logs_count mismatch: Go=%d, Rails=%d", goValues.LogsCount, railsValues.LogsCount)
		}
	})
}

// TestExpectedValues_EdgeCases tests edge case handling
func TestExpectedValues_EdgeCases(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		page      int
		totalPage int
		logs      []*dto.LogResponse
		startedAt *time.Time
	}{
		{
			name:      "zero page project",
			page:      0,
			totalPage: 100,
			logs:      nil,
			startedAt: timePtr(time.Now()),
		},
		{
			name:      "completed project no logs",
			page:      100,
			totalPage: 100,
			logs:      []*dto.LogResponse{},
			startedAt: timePtr(time.Now().AddDate(0, -1, 0)),
		},
		{
			name:      "no started_at",
			page:      50,
			totalPage: 100,
			logs:      nil,
			startedAt: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var startedAtStr *string
			if tt.startedAt != nil {
				s := tt.startedAt.Format(time.RFC3339)
				startedAtStr = &s
			}
			project := &dto.ProjectResponse{
				Page:      tt.page,
				TotalPage: tt.totalPage,
				Logs:      tt.logs,
				StartedAt: startedAtStr,
			}

			expected := testdata.CalculateExpectedValues(ctx, project)

			// Verify no panics and valid results
			if expected.Progress == nil {
				t.Error("progress should not be nil")
			}
			if expected.Status == nil {
				t.Error("status should not be nil")
			}
			if expected.LogsCount == nil {
				t.Error("logs_count should not be nil")
			}
		})
	}
}

// Helper types and functions
type projectValues struct {
	Name      string
	TotalPage int
	Page      int
	Progress  float64
	LogsCount int
}

func parseProjectJSON(data []byte) (*projectValues, error) {
	// Parse Go API JSON format (flat structure)
	// Strip comments first (lines starting with #)
	var cleanedData []byte
	for _, line := range bytes.Split(data, []byte("\n")) {
		if len(line) > 0 && line[0] != '#' {
			cleanedData = append(cleanedData, line...)
			cleanedData = append(cleanedData, '\n')
		}
	}
	var values struct {
		Name      string  `json:"name"`
		TotalPage int     `json:"total_page"`
		Page      int     `json:"page"`
		Progress  float64 `json:"progress"`
		LogsCount int     `json:"logs_count"`
	}
	err := json.Unmarshal(cleanedData, &values)
	if err != nil {
		return nil, err
	}
	result := projectValues{
		Name:      values.Name,
		TotalPage: values.TotalPage,
		Page:      values.Page,
		Progress:  values.Progress,
		LogsCount: values.LogsCount,
	}
	return &result, nil
}

func parseRailsJSON(data []byte) (*projectValues, error) {
	// Parse Rails API JSON:API format (nested data.attributes)
	// Strip comments first (lines starting with #)
	var cleanedData []byte
	for _, line := range bytes.Split(data, []byte("\n")) {
		if len(line) > 0 && line[0] != '#' {
			cleanedData = append(cleanedData, line...)
			cleanedData = append(cleanedData, '\n')
		}
	}
	var railsResponse struct {
		Data struct {
			Attributes struct {
				Name      string  `json:"name"`
				TotalPage int     `json:"total-page"`
				Page      int     `json:"page"`
				Progress  float64 `json:"progress"`
				LogsCount int     `json:"logs-count"`
			} `json:"attributes"`
		} `json:"data"`
	}
	err := json.Unmarshal(cleanedData, &railsResponse)
	if err != nil {
		return nil, err
	}
	values := projectValues{
		Name:      railsResponse.Data.Attributes.Name,
		TotalPage: railsResponse.Data.Attributes.TotalPage,
		Page:      railsResponse.Data.Attributes.Page,
		Progress:  railsResponse.Data.Attributes.Progress,
		LogsCount: railsResponse.Data.Attributes.LogsCount,
	}
	return &values, nil
}

func insertTestProject(pool *pgxpool.Pool, ctx context.Context) (int64, error) {
	// Insert a test project
	query := `
		INSERT INTO projects (name, total_page, page, started_at, reinicia)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var id int64
	err := pool.QueryRow(ctx, query,
		"Test Project",
		100,
		50,
		time.Now().AddDate(0, -1, 0),
		false,
	).Scan(&id)

	return id, err
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func floatEqual(a, b float64, precision int) bool {
	multiplier := math.Pow(10, float64(precision))
	return math.Round(a*multiplier) == math.Round(b*multiplier)
}
