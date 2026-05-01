package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-reading-log-api-next/internal/domain/dto"
)

// PoolInterface defines the interface for database pool operations
// This is defined here to avoid circular dependencies
type PoolInterface interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	AcquireFunc(ctx context.Context, fn func(*pgxpool.Conn) error) error
	Close()
	Config() *pgxpool.Config
	Reset()
}

// DashboardRepository defines the interface for dashboard query operations
// These methods provide aggregated data for dashboard views
type DashboardRepository interface {
	// GetDailyStats returns daily page statistics with weekday breakdown
	GetDailyStats(ctx context.Context, date time.Time) (*dto.DailyStats, error)

	// GetProjectAggregates returns project-level sums and counts
	GetProjectAggregates(ctx context.Context) ([]*dto.ProjectAggregate, error)

	// GetFaultsByDateRange returns the count of faults within a date range
	GetFaultsByDateRange(ctx context.Context, start, end time.Time) (*dto.FaultStats, error)

	// GetWeekdayFaults returns fault distribution by weekday (0-6 = Sunday-Saturday)
	GetWeekdayFaults(ctx context.Context, start, end time.Time) (*dto.WeekdayFaults, error)

	// GetLogsByDateRange returns log entries within a date range with page calculations
	GetLogsByDateRange(ctx context.Context, start, end time.Time) ([]*dto.LogEntry, error)

	// GetProjectWeekdayMean calculates the mean pages for a project on a specific weekday
	GetProjectWeekdayMean(ctx context.Context, projectID int64, weekday int) (float64, error)

	// CalculatePeriodPages calculates total pages within a date range
	CalculatePeriodPages(ctx context.Context, start, end time.Time) (int, error)

	// GetProjectsWithLogs returns all projects with eager-loaded logs (first 4 per project)
	GetProjectsWithLogs(ctx context.Context) ([]*dto.ProjectAggregateResponse, error)

	// GetProjectLogs returns logs for a specific project ordered by date DESC
	GetProjectLogs(ctx context.Context, projectID int64, limit int) ([]*dto.LogEntry, error)

	// GetMaxByWeekday returns the maximum pages read in a single day for the target weekday
	GetMaxByWeekday(ctx context.Context, date time.Time) (*float64, error)

	// GetOverallMean calculates the overall mean across all weekdays for the target date
	GetOverallMean(ctx context.Context, date time.Time) (*float64, error)

	// GetPreviousPeriodMean returns the mean for the same weekday 7 days prior
	GetPreviousPeriodMean(ctx context.Context, date time.Time) (*float64, error)

	// GetPreviousPeriodSpecMean returns the speculative mean for the same weekday 7 days prior
	GetPreviousPeriodSpecMean(ctx context.Context, date time.Time) (*float64, error)

	// GetMeanByWeekday calculates the mean pages per 7-day interval for a specific weekday
	// Algorithm: total_pages / count_reads where count_reads = floor((log_data - begin_data) / 7 days)
	// Returns nil for no data or zero intervals (consistent with GetMaxByWeekday)
	GetMeanByWeekday(ctx context.Context, weekday int) (*float64, error)

	// GetRunningProjectsWithLogs returns all projects with eager-loaded logs
	// The service layer filters by running status using the 7-day threshold
	// Returns projects with first 4 logs per project ordered by data DESC
	GetRunningProjectsWithLogs(ctx context.Context) ([]*dto.ProjectWithLogs, error)

	// GetPool returns the underlying database connection pool
	GetPool() PoolInterface
}
