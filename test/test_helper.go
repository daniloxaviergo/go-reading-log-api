package test

import (
	"context"
	"fmt"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go-reading-log-api-next/internal/config"
	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/domain/models"
	"go-reading-log-api-next/internal/repository"

)

const testContextTimeout = 5 * time.Second

// TestHelper provides common test utilities for database setup and cleanup
type TestHelper struct {
	Config     *config.Config
	Pool       *pgxpool.Pool
	TestDBName string
}

// SetupTestDB creates a test database connection using test database configuration
// It reads DB_DATABASE_TEST env var, falling back to DB_DATABASE with '_test' suffix
func SetupTestDB() (*TestHelper, error) {
	cfg := config.LoadConfig()

	// Determine test database name
	testDBName := os.Getenv("DB_DATABASE_TEST")
	if testDBName == "" {
		testDBName = cfg.DBDatabase + "_test"
	}

	// Build connection string for test database
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		testDBName,
	)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Verify connection works
	ctx, cancel := context.WithTimeout(context.Background(), testContextTimeout)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping test database: %w", err)
	}

	return &TestHelper{
		Config:     cfg,
		Pool:       pool,
		TestDBName: testDBName,
	}, nil
}

// SetupTestDBWithConfig creates a test database connection with custom configuration
func SetupTestDBWithConfig(cfg *config.Config) (*TestHelper, error) {
	// Determine test database name
	testDBName := os.Getenv("DB_DATABASE_TEST")
	if testDBName == "" {
		testDBName = cfg.DBDatabase + "_test"
	}

	// Build connection string for test database
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		testDBName,
	)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Verify connection works
	ctx, cancel := context.WithTimeout(context.Background(), testContextTimeout)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping test database: %w", err)
	}

	return &TestHelper{
		Config:     cfg,
		Pool:       pool,
		TestDBName: testDBName,
	}, nil
}

// SetupTestSchema creates the test schema (tables) in the test database
// This should be called before running integration tests
func (h *TestHelper) SetupTestSchema() error {
	ctx, cancel := context.WithTimeout(context.Background(), testContextTimeout)
	defer cancel()

	// SQL statements to create tables if they don't exist (similar to Rails schema)
	// Execute separately because pgx.Exec doesn't support multi-statement execution
	queries := []string{
		`CREATE TABLE IF NOT EXISTS projects (
			id BIGSERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			total_page INT NOT NULL DEFAULT 0,
			started_at TIMESTAMP WITH TIME ZONE,
			page INT NOT NULL DEFAULT 0,
			reinicia BOOLEAN NOT NULL DEFAULT false,
			progress VARCHAR(255),
			status VARCHAR(255),
			logs_count INT DEFAULT 0,
			days_unread INT DEFAULT 0,
			median_day VARCHAR(255),
			finished_at TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS logs (
			id BIGSERIAL PRIMARY KEY,
			project_id BIGINT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
			data VARCHAR(255),
			start_page INT NOT NULL DEFAULT 0,
			end_page INT NOT NULL DEFAULT 0,
			wday INT NOT NULL DEFAULT 0,
			note TEXT,
			text TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_logs_project_id ON logs(project_id)`,
	}

	for _, query := range queries {
		_, err := h.Pool.Exec(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to setup test schema: %w", err)
		}
	}

	return nil
}

// CleanupTestSchema drops the test tables after tests complete
func (h *TestHelper) CleanupTestSchema() error {
	ctx, cancel := context.WithTimeout(context.Background(), testContextTimeout)
	defer cancel()

	// Drop tables - be careful with this in production!
	queries := []string{
		"DROP TABLE IF EXISTS logs CASCADE;",
		"DROP TABLE IF EXISTS projects CASCADE;",
	}

	for _, query := range queries {
		_, err := h.Pool.Exec(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to cleanup test schema: %w", err)
		}
	}

	return nil
}

// ClearTestData truncates all test data but keeps schema intact
func (h *TestHelper) ClearTestData() error {
	ctx, cancel := context.WithTimeout(context.Background(), testContextTimeout)
	defer cancel()

	queries := []string{
		"TRUNCATE TABLE logs CASCADE;",
		"TRUNCATE TABLE projects CASCADE;",
	}

	for _, query := range queries {
		_, err := h.Pool.Exec(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to clear test data: %w", err)
		}
	}

	return nil
}

// GetContext returns a context with timeout for test operations
func (h *TestHelper) GetContext() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), testContextTimeout)
	_ = cancel // Will be called when context expires
	return ctx
}

// GetContextWithTimeout returns a context with custom timeout
func (h *TestHelper) GetContextWithTimeout(timeout time.Duration) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	_ = cancel // Will be called when context expires
	return ctx
}

// Close cleans up the database connection pool
func (h *TestHelper) Close() {
	if h.Pool != nil {
		h.Pool.Close()
	}
}

// Helper function to get test context with timeout
func GetTestContext() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), testContextTimeout)
	_ = cancel
	return ctx
}

// Helper function to get test context with custom timeout
func GetTestContextWithTimeout(timeout time.Duration) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	_ = cancel
	return ctx
}

// Helper to check if we're running against a test database (safety check)
func IsTestDatabase() bool {
	testDB := os.Getenv("DB_DATABASE_TEST")
	dbName := os.Getenv("DB_DATABASE")
	return testDB != "" || dbName != "reading_log"
}

// Helper to generate a unique test name for test data
func TestName(t *testing.T) string {
	return fmt.Sprintf("%s_%d", t.Name(), time.Now().UnixMilli())
}

// MockProjectRepository is a mock implementation of repository.ProjectRepository
// Used for testing without database dependency
type MockProjectRepository struct {
	Projects         map[int64]*models.Project
	WithLogsMap      map[int64]*dto.ProjectResponse
	GetByIDCalls     []int64
	GetAllCalls      int
	GetWithLogsCalls []int64
	Err              error
}

// NewMockProjectRepository creates a new MockProjectRepository
func NewMockProjectRepository() *MockProjectRepository {
	return &MockProjectRepository{
		Projects:    make(map[int64]*models.Project),
		WithLogsMap: make(map[int64]*dto.ProjectResponse),
	}
}

// AddProject adds a project to the mock repository
func (m *MockProjectRepository) AddProject(project *models.Project) {
	m.Projects[project.ID] = project
}

// AddProjectWithLogs adds a project with logs response to the mock repository
func (m *MockProjectRepository) AddProjectWithLogs(response *dto.ProjectResponse) {
	m.WithLogsMap[response.ID] = response
}

// SetError sets a generic error to return on all operations
func (m *MockProjectRepository) SetError(err error) {
	m.Err = err
}

// GetByID retrieves a project by ID
func (m *MockProjectRepository) GetByID(ctx context.Context, id int64) (*models.Project, error) {
	m.GetByIDCalls = append(m.GetByIDCalls, id)

	if m.Err != nil {
		return nil, m.Err
	}

	if project, ok := m.Projects[id]; ok {
		return project, nil
	}

	return nil, fmt.Errorf("project with ID %d not found", id)
}

// GetAll retrieves all projects
func (m *MockProjectRepository) GetAll(ctx context.Context) ([]*models.Project, error) {
	m.GetAllCalls++

	if m.Err != nil {
		return nil, m.Err
	}

	projects := make([]*models.Project, 0, len(m.Projects))
	for _, project := range m.Projects {
		projects = append(projects, project)
	}

	// Sort by ID for deterministic ordering
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].ID < projects[j].ID
	})

	return projects, nil
}

// GetAllWithLogs retrieves all projects with their associated logs
func (m *MockProjectRepository) GetAllWithLogs(ctx context.Context) ([]*repository.ProjectWithLogs, error) {
	m.GetAllCalls++

	if m.Err != nil {
		return nil, m.Err
	}

	result := make([]*repository.ProjectWithLogs, 0, len(m.Projects))
	for id, project := range m.Projects {
		var response *dto.ProjectResponse

		// Check if we have a pre-configured response with logs
		if resp, ok := m.WithLogsMap[id]; ok {
			response = resp
		} else {
			// Build response from project
			response = &dto.ProjectResponse{
				ID:         project.ID,
				Name:       project.Name,
				StartedAt:  formatTimePtr(project.StartedAt),
				Progress:   project.Progress,
				TotalPage:  project.TotalPage,
				Page:       project.Page,
				Status:     project.Status,
				LogsCount:  project.LogsCount,
				DaysUnread: project.DaysUnread,
				MedianDay:  formatTimePtr(project.MedianDay),
				FinishedAt: formatTimePtr(project.FinishedAt),
			}
		}

		result = append(result, &repository.ProjectWithLogs{
			Project: response,
			Logs:    []*dto.LogResponse{}, // Empty logs for now
		})
	}

	// Sort by project ID for deterministic ordering
	sort.Slice(result, func(i, j int) bool {
		return result[i].Project.ID < result[j].Project.ID
	})

	return result, nil
}

// GetWithLogs retrieves a project with its associated logs
func (m *MockProjectRepository) GetWithLogs(ctx context.Context, id int64) (*repository.ProjectWithLogs, error) {
	m.GetWithLogsCalls = append(m.GetWithLogsCalls, id)

	if m.Err != nil {
		return nil, m.Err
	}

	if response, ok := m.WithLogsMap[id]; ok {
		return &repository.ProjectWithLogs{
			Project: response,
			Logs:    []*dto.LogResponse{},
		}, nil
	}

	if project, ok := m.Projects[id]; ok {
		response := &dto.ProjectResponse{
			ID:         project.ID,
			Name:       project.Name,
			StartedAt:  formatTimePtr(project.StartedAt),
			Progress:   project.Progress,
			TotalPage:  project.TotalPage,
			Page:       project.Page,
			Status:     project.Status,
			LogsCount:  project.LogsCount,
			DaysUnread: project.DaysUnread,
			MedianDay:  formatTimePtr(project.MedianDay),
			FinishedAt: formatTimePtr(project.FinishedAt),
		}
		return &repository.ProjectWithLogs{
			Project: response,
			Logs:    []*dto.LogResponse{},
		}, nil
	}

	return nil, fmt.Errorf("project with ID %d not found", id)
}

// formatTimePtr converts a time.Time pointer to a string pointer for JSON serialization
func formatTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(time.RFC3339)
	return &s
}

// CallCount returns the number of times GetByID was called
func (m *MockProjectRepository) CallCount() int {
	return len(m.GetByIDCalls)
}

// GetAllCallCount returns the number of times GetAll was called
func (m *MockProjectRepository) GetAllCallCount() int {
	return m.GetAllCalls
}

// GetWithLogsCallCount returns the number of times GetWithLogs was called
func (m *MockProjectRepository) GetWithLogsCallCount() int {
	return len(m.GetWithLogsCalls)
}

// GetByIDLastCall returns the last project ID requested
func (m *MockProjectRepository) GetByIDLastCall() int64 {
	if len(m.GetByIDCalls) == 0 {
		return 0
	}
	return m.GetByIDCalls[len(m.GetByIDCalls)-1]
}

// MockLogRepository is a mock implementation of repository.LogRepository
// Used for testing without database dependency
type MockLogRepository struct {
	Logs                map[int64]*models.Log
	ByProjectID         map[int64][]*models.Log
	GetByIDCalls        []int64
	GetByProjectIDCalls []int64
	GetAllCalls         int
	Err                 error
}

// NewMockLogRepository creates a new MockLogRepository
func NewMockLogRepository() *MockLogRepository {
	return &MockLogRepository{
		Logs:        make(map[int64]*models.Log),
		ByProjectID: make(map[int64][]*models.Log),
	}
}

// AddLog adds a log to the mock repository
func (m *MockLogRepository) AddLog(log *models.Log) {
	m.Logs[log.ID] = log
}

// AddLogsForProject adds logs for a specific project
func (m *MockLogRepository) AddLogsForProject(projectID int64, logs []*models.Log) {
	m.ByProjectID[projectID] = logs
}

// SetError sets a generic error to return on all operations
func (m *MockLogRepository) SetError(err error) {
	m.Err = err
}

// GetByID retrieves a log by ID
func (m *MockLogRepository) GetByID(ctx context.Context, id int64) (*models.Log, error) {
	m.GetByIDCalls = append(m.GetByIDCalls, id)

	if m.Err != nil {
		return nil, m.Err
	}

	if log, ok := m.Logs[id]; ok {
		return log, nil
	}

	return nil, fmt.Errorf("log with ID %d not found", id)
}

// GetByProjectID retrieves all logs for a project
func (m *MockLogRepository) GetByProjectID(ctx context.Context, projectID int64) ([]*models.Log, error) {
	m.GetByProjectIDCalls = append(m.GetByProjectIDCalls, projectID)

	if m.Err != nil {
		return nil, m.Err
	}

	if logs, ok := m.ByProjectID[projectID]; ok {
		return logs, nil
	}

	return []*models.Log{}, nil
}

// GetByProjectIDOrdered retrieves all logs for a project ordered
func (m *MockLogRepository) GetByProjectIDOrdered(ctx context.Context, projectID int64) ([]*models.Log, error) {
	m.GetByProjectIDCalls = append(m.GetByProjectIDCalls, projectID)

	if m.Err != nil {
		return nil, m.Err
	}

	if logs, ok := m.ByProjectID[projectID]; ok {
		return logs, nil
	}

	return []*models.Log{}, nil
}

// GetAll retrieves all logs
func (m *MockLogRepository) GetAll(ctx context.Context) ([]*models.Log, error) {
	m.GetAllCalls++

	if m.Err != nil {
		return nil, m.Err
	}

	logs := make([]*models.Log, 0, len(m.Logs))
	for _, log := range m.Logs {
		logs = append(logs, log)
	}

	// Sort by ID for deterministic ordering
	sort.Slice(logs, func(i, j int) bool {
		return logs[i].ID < logs[j].ID
	})

	return logs, nil
}

// CallCount returns the number of times GetByID was called
func (m *MockLogRepository) CallCount() int {
	return len(m.GetByIDCalls)
}

// GetByProjectIDCallCount returns the number of times GetByProjectID was called
func (m *MockLogRepository) GetByProjectIDCallCount() int {
	return len(m.GetByProjectIDCalls)
}

// GetAllCallCount returns the number of times GetAll was called
func (m *MockLogRepository) GetAllCallCount() int {
	return m.GetAllCalls
}

// GetByIDLastCall returns the last log ID requested
func (m *MockLogRepository) GetByIDLastCall() int64 {
	if len(m.GetByIDCalls) == 0 {
		return 0
	}
	return m.GetByIDCalls[len(m.GetByIDCalls)-1]
}
