package unit

import (
	"testing"
	"time"

	"go-reading-log-api-next/internal/domain/models"
)

// TestProject_ParseLogDate tests the ParseLogDate function with multiple formats
func TestProject_ParseLogDate(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		wantOk bool
	}{
		{
			name:   "YYYY-MM-DD format",
			input:  "2024-01-15",
			wantOk: true,
		},
		{
			name:   "RFC3339 format",
			input:  "2024-01-15T10:30:00Z",
			wantOk: true,
		},
		{
			name:   "Standard datetime format",
			input:  "2024-01-15 10:30:00",
			wantOk: true,
		},
		{
			name:   "Invalid format",
			input:  "not-a-date",
			wantOk: false,
		},
		{
			name:   "Empty string",
			input:  "",
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTime, gotOk := models.ParseLogDate(tt.input)

			if gotOk != tt.wantOk {
				t.Errorf("ParseLogDate(%q) = _, %v; want ok=%v", tt.input, gotOk, tt.wantOk)
			}

			if gotOk && tt.wantOk {
				// Verify the parsed time is reasonable
				if gotTime.IsZero() {
					t.Errorf("ParseLogDate(%q) returned zero time", tt.input)
				}
			}
		})
	}
}

// TestProject_ParseLogDate_Timezone tests timezone-aware date parsing
func TestProject_ParseLogDate_Timezone(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		timezone *time.Location
		wantDate string
		wantOk   bool
	}{
		{
			name:     "UTC timezone",
			input:    "2024-01-15T10:30:00Z",
			timezone: time.UTC,
			wantDate: "2024-01-15",
			wantOk:   true,
		},
		{
			name:     "Brazil timezone (BRT)",
			input:    "2024-01-15T10:30:00Z",
			timezone: time.FixedZone("BRT", -3*60*60),
			wantDate: "2024-01-15",
			wantOk:   true,
		},
		{
			name:     "Japan timezone (JST)",
			input:    "2024-01-15T10:30:00Z",
			timezone: time.FixedZone("JST", 9*60*60),
			wantDate: "2024-01-15",
			wantOk:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTime, gotOk := models.ParseLogDateWithTimezone(tt.input, tt.timezone)

			if gotOk != tt.wantOk {
				t.Errorf("ParseLogDateWithTimezone(%q) = _, %v; want ok=%v", tt.input, gotOk, tt.wantOk)
			}

			if gotOk && tt.wantOk {
				// Verify the date part matches expected
				gotDate := gotTime.Format("2006-01-02")
				if gotDate != tt.wantDate {
					t.Errorf("Expected date %s, got %s", tt.wantDate, gotDate)
				}
			}
		})
	}
}
