package test

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go-reading-log-api-next/internal/config"
	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/domain/models"
	"go-reading-log-api-next/internal/repository"
)

const testContextTimeout = 30 * time.Second

// TestHelper provides common test utilities for database setup and cleanup
type TestHelper struct {
	Config     *config.Config
	Pool       *pgxpool.Pool
	TestDBName string
}

// getGoroutineID extracts the goroutine ID from the runtime stack trace
// This is used to ensure unique database names for parallel tests running in the same process
func getGoroutineID() uint64 {
	// Get the goroutine stack trace
	// We use a buffer size of 32 to capture the goroutine ID
	buf := make([]byte, 32)
	runtime.Stack(buf, false)

	// Parse the goroutine ID from the stack trace
	// Format: goroutine 123 [running]:
	// We look for "goroutine " and extract the following number
	str := string(buf)
	start := strings.Index(str, "goroutine ")
	if start == -1 {
		return 0
	}
	start += len("goroutine ")

	// Find the end of the number
	end := start
	for end < len(str) && str[end] >= '0' && str[end] <= '9' {
		end++
	}

	if end <= start {
		return 0
	}

	// Convert to uint64
	var id uint64
	for i := start; i < end; i++ {
		id = id*10 + uint64(str[i]-'0')
	}
	return id
}

// SetupTestDB creates a test database connection using test database configuration
// It reads DB_DATABASE_TEST env var, falling back to DB_DATABASE with '_test' suffix
// For parallel tests, a unique database name is created based on the test name
func SetupTestDB() (*TestHelper, error) {
	// Load .env.test for test-specific configuration
	_ = godotenv.Load(".env.test")

	// Override DB_HOST to localhost for local testing (overrides .env file)
	// Environment variable DB_HOST takes precedence, so we must explicitly set it
	if os.Getenv("DB_HOST") == "postgres" {
		os.Setenv("DB_HOST", "localhost")
	}

	cfg := config.LoadConfig()

	// Determine test database name
	testDBName := os.Getenv("DB_DATABASE_TEST")
	if testDBName == "" {
		testDBName = cfg.DBDatabase + "_test"
	}

	// Check if we're running in parallel mode and need unique database
	// Use a unique database name based on PID, goroutine ID, and timestamp for parallel test isolation
	goroutineID := getGoroutineID()
	testDBName = fmt.Sprintf("%s_%d_%d_%d", testDBName, os.Getpid(), goroutineID, time.Now().UnixNano())

	// Connect to the main database to create the test database if needed
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBDatabase)
	mainPool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create main connection pool: %w", err)
	}
	defer mainPool.Close()

	// Create the test database if it doesn't exist
	ctx, cancel := context.WithTimeout(context.Background(), testContextTimeout)
	defer cancel()
	_, err = mainPool.Exec(ctx, fmt.Sprintf(`CREATE DATABASE "%s"`, testDBName))
	if err != nil && !strings.Contains(err.Error(), "already exists") && !strings.Contains(err.Error(), "exists") {
		return nil, fmt.Errorf("failed to create test database: %w", err)
	}

	// Build connection string for test database
	connStr = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
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
	ctx, cancel = context.WithTimeout(context.Background(), testContextTimeout)
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
	// Load .env.test for test-specific configuration
	_ = godotenv.Load(".env.test")

	// Override DB_HOST to localhost for local testing (overrides .env file)
	// Environment variable DB_HOST takes precedence, so we must explicitly set it
	if os.Getenv("DB_HOST") == "postgres" {
		os.Setenv("DB_HOST", "localhost")
	}

	// Determine test database name
	testDBName := os.Getenv("DB_DATABASE_TEST")
	if testDBName == "" {
		testDBName = cfg.DBDatabase + "_test"
	}

	// Use unique database name for parallel tests
	// Include goroutine ID to ensure uniqueness when multiple goroutines run in same process
	goroutineID := getGoroutineID()
	testDBName = fmt.Sprintf("%s_%d_%d_%d", testDBName, os.Getpid(), goroutineID, time.Now().UnixNano())

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

// Close cleans up the database connection pool using defer to ensure cleanup runs even on panic
// The cleanup must complete within 1 second of test completion to not block test results
func (h *TestHelper) Close() {
	if h.Pool != nil {
		// Capture the pool and config before deferring
		pool := h.Pool
		testDBName := h.TestDBName
		cfg := h.Config

		// Defer cleanup to run when function returns (even on panic)
		// This ensures the database is dropped even if the test panics
		defer func() {
			// Use a separate connection pool for DROP DATABASE to avoid issues
			// with closing the pool being dropped
			connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
				cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBDatabase)

			// First, cleanup orphaned databases from previous sessions
			// Use 60 second timeout for potentially large cleanup (6000+ databases)
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			mainPool, err := pgxpool.New(ctx, connStr)
			if err == nil {
				// Call cleanupOrphanedDatabases to drop old test databases
				// Pass current testDBName to exclude it from cleanup
				_ = cleanupOrphanedDatabases(mainPool, testDBName)
				mainPool.Close()
			}

			// Then drop the current test database
			// Use 1 second timeout for immediate drop
			ctx2, cancel2 := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel2()

			mainPool2, err := pgxpool.New(ctx2, connStr)
			if err == nil {
				// Use DROP DATABASE IF EXISTS to handle missing databases gracefully
				// Log errors but don't propagate them to avoid blocking test results
				_, dropErr := mainPool2.Exec(ctx2, fmt.Sprintf("DROP DATABASE IF EXISTS %s", testDBName))
				if dropErr != nil {
					// Log the error but don't fail the test
					// This ensures cleanup doesn't block test results
					_ = dropErr
				}
				mainPool2.Close()
			}
		}()

		// Close the connection pool (runs before defer block completes)
		pool.Close()
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

// cleanupOrphanedDatabases identifies and drops test databases older than 24 hours
// excludeName is the current test database name to exclude from cleanup
// This function queries pg_database for databases matching the pattern reading_log_test_%
// and drops each identified orphan. The cleanup must complete in under 1 minute for 6,000+ databases.
func cleanupOrphanedDatabases(pool *pgxpool.Pool, excludeName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Query for orphaned databases matching the pattern
	// Use pg_get_userbyid(datdba) to filter by current user for safety
	query := `
		SELECT datname 
		FROM pg_database 
		WHERE datname LIKE $1
		AND datname != $2
		AND pg_catalog.pg_get_userbyid(datdba) = current_user
	`

	rows, err := pool.Query(ctx, query, "reading_log_test_%", excludeName)
	if err != nil {
		return fmt.Errorf("failed to query test databases: %w", err)
	}
	defer rows.Close()

	var toDrop []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			continue
		}
		toDrop = append(toDrop, name)
	}

	// Drop each orphaned database
	for _, dbName := range toDrop {
		_, dropErr := pool.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
		if dropErr != nil {
			// Log errors but don't fail the cleanup
			// This ensures cleanup doesn't block test results
			_ = dropErr
		}
	}

	return nil
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
				MedianDay:  project.MedianDay,
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
			MedianDay:  project.MedianDay,
			FinishedAt: formatTimePtr(project.FinishedAt),
		}
		return &repository.ProjectWithLogs{
			Project: response,
			Logs:    []*dto.LogResponse{},
		}, nil
	}

	return nil, fmt.Errorf("project with ID %d not found", id)
}

// Create inserts a new project into the mock repository
func (m *MockProjectRepository) Create(ctx context.Context, project *models.Project) (*models.Project, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	// Generate a new ID (use the next available ID + 1)
	var maxID int64
	for id := range m.Projects {
		if id > maxID {
			maxID = id
		}
	}
	project.ID = maxID + 1
	m.Projects[project.ID] = project

	return project, nil
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

// Create inserts a new log into the mock repository
func (m *MockLogRepository) Create(ctx context.Context, log *models.Log) (*models.Log, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	// Generate a new ID (use the next available ID + 1)
	var maxID int64
	for id := range m.Logs {
		if id > maxID {
			maxID = id
		}
	}
	log.ID = maxID + 1
	m.Logs[log.ID] = log

	// Also add to ByProjectID map
	if _, ok := m.ByProjectID[log.ProjectID]; !ok {
		m.ByProjectID[log.ProjectID] = make([]*models.Log, 0)
	}
	m.ByProjectID[log.ProjectID] = append(m.ByProjectID[log.ProjectID], log)

	return log, nil
}
