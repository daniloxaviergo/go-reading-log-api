package dashboard

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DashboardFixtures manages all dashboard-related test data
type DashboardFixtures struct {
	pool *pgxpool.Pool

	// Pre-created data for reuse across tests
	projects map[int64]*ProjectFixture
	logs     map[int64][]*LogFixture
}

// ProjectFixture represents a single project with its associated data
type ProjectFixture struct {
	ID        int64
	Name      string
	TotalPage int
	Page      int
	StartedAt *time.Time
	Status    string

	// Associated logs
	Logs []*LogFixture
}

// LogFixture represents a log entry with full control over fields
type LogFixture struct {
	ID        int64
	ProjectID int64
	Data      time.Time // Explicit date for reproducibility
	StartPage int
	EndPage   int
	Note      *string
	WDay      int // Weekday (0-6)
}

// NewDashboardFixtures creates a new fixture manager
func NewDashboardFixtures(pool *pgxpool.Pool) *DashboardFixtures {
	return &DashboardFixtures{
		pool:     pool,
		projects: make(map[int64]*ProjectFixture),
		logs:     make(map[int64][]*LogFixture),
	}
}

// LoadScenario loads a complete scenario into the test database
func (fm *DashboardFixtures) LoadScenario(scenario *Scenario) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Begin transaction
	tx, err := fm.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	// Insert projects
	for _, proj := range scenario.Projects {
		if err := fm.insertProject(ctx, tx, proj); err != nil {
			return fmt.Errorf("failed to insert project %d: %w", proj.ID, err)
		}
	}

	// Insert logs
	for _, log := range scenario.Logs {
		if err := fm.insertLog(ctx, tx, log); err != nil {
			return fmt.Errorf("failed to insert log %d: %w", log.ID, err)
		}
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Cache the loaded fixtures
	for _, proj := range scenario.Projects {
		fm.projects[proj.ID] = proj
	}
	for _, log := range scenario.Logs {
		fm.logs[log.ProjectID] = append(fm.logs[log.ProjectID], log)
	}

	return nil
}

func (fm *DashboardFixtures) insertProject(ctx context.Context, tx pgx.Tx, proj *ProjectFixture) error {
	query := `
		INSERT INTO projects (id, name, total_page, page, started_at, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			total_page = EXCLUDED.total_page,
			page = EXCLUDED.page,
			started_at = EXCLUDED.started_at,
			status = EXCLUDED.status,
			updated_at = NOW()
	`

	_, err := tx.Exec(ctx, query,
		proj.ID,
		proj.Name,
		proj.TotalPage,
		proj.Page,
		proj.StartedAt,
		proj.Status,
	)
	if err != nil {
		return fmt.Errorf("failed to insert project: %w", err)
	}

	// Also update derived fields if needed
	// Note: data column is VARCHAR, so we need to cast it for date operations
	updateQuery := `
		UPDATE projects SET
			logs_count = COALESCE((SELECT COUNT(*) FROM logs WHERE project_id = $1), 0),
			days_unread = CASE WHEN started_at IS NOT NULL THEN 
				COALESCE((SELECT EXTRACT(DAY FROM NOW() - MAX(data::timestamp))::INT FROM logs WHERE project_id = $1), 0)
			ELSE 0 END,
			median_day = CASE WHEN logs_count > 0 AND days_unread > 0 THEN
				(page::float / days_unread)::varchar(255)
			ELSE '0' END
		WHERE id = $1
	`

	_, err = tx.Exec(ctx, updateQuery, proj.ID)
	if err != nil {
		return fmt.Errorf("failed to update project derived fields: %w", err)
	}

	return nil
}

func (fm *DashboardFixtures) insertLog(ctx context.Context, tx pgx.Tx, log *LogFixture) error {
	query := `
		INSERT INTO logs (id, project_id, data, start_page, end_page, wday, note, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		ON CONFLICT (id) DO UPDATE SET
			project_id = EXCLUDED.project_id,
			data = EXCLUDED.data,
			start_page = EXCLUDED.start_page,
			end_page = EXCLUDED.end_page,
			wday = EXCLUDED.wday,
			note = EXCLUDED.note,
			updated_at = NOW()
	`

	_, err := tx.Exec(ctx, query,
		log.ID,
		log.ProjectID,
		log.Data.Format(time.RFC3339),
		log.StartPage,
		log.EndPage,
		log.WDay,
		log.Note,
	)
	if err != nil {
		return fmt.Errorf("failed to insert log: %w", err)
	}

	return nil
}

// GetProject retrieves a loaded project fixture
func (fm *DashboardFixtures) GetProject(id int64) (*ProjectFixture, bool) {
	proj, ok := fm.projects[id]
	return proj, ok
}

// GetLogsForProject retrieves logs for a specific project
func (fm *DashboardFixtures) GetLogsForProject(projectID int64) ([]*LogFixture, bool) {
	logs, ok := fm.logs[projectID]
	return logs, ok
}

// ClearAll removes all test data from the database
func (fm *DashboardFixtures) ClearAll() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	queries := []string{
		"TRUNCATE TABLE logs CASCADE;",
		"TRUNCATE TABLE projects CASCADE;",
	}

	for _, query := range queries {
		_, err := fm.pool.Exec(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to clear data: %w", err)
		}
	}

	fm.projects = make(map[int64]*ProjectFixture)
	fm.logs = make(map[int64][]*LogFixture)

	return nil
}

// CreateTestProject creates a test project and returns it
func (fm *DashboardFixtures) CreateTestProject(id int64, name string, totalPage, page int, status string) (*ProjectFixture, error) {
	proj := &ProjectFixture{
		ID:        id,
		Name:      name,
		TotalPage: totalPage,
		Page:      page,
		Status:    status,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := fm.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	if err := fm.insertProject(ctx, tx, proj); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}

	fm.projects[id] = proj
	return proj, nil
}

// CreateTestLog creates a test log and returns it
func (fm *DashboardFixtures) CreateTestLog(id, projectID int64, data time.Time, startPage, endPage, wday int, note string) (*LogFixture, error) {
	log := &LogFixture{
		ID:        id,
		ProjectID: projectID,
		Data:      data,
		StartPage: startPage,
		EndPage:   endPage,
		WDay:      wday,
	}

	if note != "" {
		log.Note = &note
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := fm.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	if err := fm.insertLog(ctx, tx, log); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}

	fm.logs[projectID] = append(fm.logs[projectID], log)
	return log, nil
}

// GenerateRandomProject creates a project with random data for stress testing
func GenerateRandomProject(id int64, seed int) *ProjectFixture {
	idInt := int(id)
	return &ProjectFixture{
		ID:        id,
		Name:      fmt.Sprintf("Random Project %d", id),
		TotalPage: ((idInt*17 + seed) % 500) + 100,
		Page:      ((idInt*23 + seed) % 400),
		Status:    "running",
	}
}

// GenerateRandomLog creates a log with random data for stress testing
func GenerateRandomLog(id, projectID int64, baseDate time.Time, seed int) *LogFixture {
	idInt := int(id)
	dayOffset := (idInt*7 + seed) % 30
	hourOffset := (idInt*13 + seed) % 24

	return &LogFixture{
		ID:        id,
		ProjectID: projectID,
		Data:      baseDate.AddDate(0, 0, -dayOffset).Add(time.Duration(hourOffset) * time.Hour),
		StartPage: ((idInt*11 + seed) % 200),
		EndPage:   (((idInt*11 + seed) % 200) + ((idInt*17+seed)%50 + 1)),
		WDay:      int(baseDate.AddDate(0, 0, -dayOffset).Weekday()),
	}
}
