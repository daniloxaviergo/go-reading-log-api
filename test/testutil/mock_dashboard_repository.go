// Package testutil provides test utilities including mock repositories
package testutil

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"

	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/repository"
)

// MockDashboardRepository is a mock implementation of the DashboardRepository interface
type MockDashboardRepository struct {
	mock.Mock
}

// NewMockDashboardRepository creates a new MockDashboardRepository
func NewMockDashboardRepository() *MockDashboardRepository {
	return &MockDashboardRepository{}
}

// GetDailyStats mocks the GetDailyStats method
func (m *MockDashboardRepository) GetDailyStats(ctx context.Context, date time.Time) (*dto.DailyStats, error) {
	args := m.Called(ctx, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DailyStats), args.Error(1)
}

// GetProjectAggregates mocks the GetProjectAggregates method
func (m *MockDashboardRepository) GetProjectAggregates(ctx context.Context) ([]*dto.ProjectAggregate, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*dto.ProjectAggregate), args.Error(1)
}

// GetFaultsByDateRange mocks the GetFaultsByDateRange method
func (m *MockDashboardRepository) GetFaultsByDateRange(ctx context.Context, start, end time.Time) (*dto.FaultStats, error) {
	args := m.Called(ctx, start, end)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.FaultStats), args.Error(1)
}

// GetWeekdayFaults mocks the GetWeekdayFaults method
func (m *MockDashboardRepository) GetWeekdayFaults(ctx context.Context, start, end time.Time) (*dto.WeekdayFaults, error) {
	args := m.Called(ctx, start, end)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.WeekdayFaults), args.Error(1)
}

// GetLogsByDateRange mocks the GetLogsByDateRange method
func (m *MockDashboardRepository) GetLogsByDateRange(ctx context.Context, start, end time.Time) ([]*dto.LogEntry, error) {
	args := m.Called(ctx, start, end)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*dto.LogEntry), args.Error(1)
}

// GetProjectWeekdayMean mocks the GetProjectWeekdayMean method
func (m *MockDashboardRepository) GetProjectWeekdayMean(ctx context.Context, projectID int64, weekday int) (float64, error) {
	args := m.Called(ctx, projectID, weekday)
	return args.Get(0).(float64), args.Error(1)
}

// CalculatePeriodPages mocks the CalculatePeriodPages method
func (m *MockDashboardRepository) CalculatePeriodPages(ctx context.Context, start, end time.Time) (int, error) {
	args := m.Called(ctx, start, end)
	return args.Int(0), args.Error(1)
}

// GetProjectsWithLogs mocks the GetProjectsWithLogs method
func (m *MockDashboardRepository) GetProjectsWithLogs(ctx context.Context) ([]*dto.ProjectAggregateResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*dto.ProjectAggregateResponse), args.Error(1)
}

// GetProjectLogs mocks the GetProjectLogs method
func (m *MockDashboardRepository) GetProjectLogs(ctx context.Context, projectID int64, limit int) ([]*dto.LogEntry, error) {
	args := m.Called(ctx, projectID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*dto.LogEntry), args.Error(1)
}

// GetMaxByWeekday mocks the GetMaxByWeekday method
func (m *MockDashboardRepository) GetMaxByWeekday(ctx context.Context, date time.Time) (*float64, error) {
	args := m.Called(ctx, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*float64), args.Error(1)
}

// GetOverallMean mocks the GetOverallMean method
func (m *MockDashboardRepository) GetOverallMean(ctx context.Context, date time.Time) (*float64, error) {
	args := m.Called(ctx, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*float64), args.Error(1)
}

// GetPreviousPeriodMean mocks the GetPreviousPeriodMean method
func (m *MockDashboardRepository) GetPreviousPeriodMean(ctx context.Context, date time.Time) (*float64, error) {
	args := m.Called(ctx, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*float64), args.Error(1)
}

// GetPreviousPeriodSpecMean mocks the GetPreviousPeriodSpecMean method
func (m *MockDashboardRepository) GetPreviousPeriodSpecMean(ctx context.Context, date time.Time) (*float64, error) {
	args := m.Called(ctx, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*float64), args.Error(1)
}

// GetMeanByWeekday mocks the GetMeanByWeekday method
func (m *MockDashboardRepository) GetMeanByWeekday(ctx context.Context, weekday int) (*float64, error) {
	args := m.Called(ctx, weekday)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*float64), args.Error(1)
}

// GetRunningProjectsWithLogs mocks the GetRunningProjectsWithLogs method
func (m *MockDashboardRepository) GetRunningProjectsWithLogs(ctx context.Context) ([]*dto.ProjectWithLogs, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*dto.ProjectWithLogs), args.Error(1)
}

// GetPool mocks the GetPool method
func (m *MockDashboardRepository) GetPool() repository.PoolInterface {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(repository.PoolInterface)
}

// Ensure MockDashboardRepository implements DashboardRepository interface
var _ repository.DashboardRepository = (*MockDashboardRepository)(nil)
