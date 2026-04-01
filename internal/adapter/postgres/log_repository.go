package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-reading-log-api-next/internal/domain/models"
)

// LogRepositoryImpl implements LogRepository interface using PostgreSQL
type LogRepositoryImpl struct {
	pool *pgxpool.Pool
}

// NewLogRepositoryImpl creates a new LogRepositoryImpl with the given connection pool
func NewLogRepositoryImpl(pool *pgxpool.Pool) *LogRepositoryImpl {
	return &LogRepositoryImpl{pool: pool}
}

// GetByID retrieves a log by its ID
func (r *LogRepositoryImpl) GetByID(ctx context.Context, id int64) (*models.Log, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultContextTimeout)
	defer cancel()

	query := `
		SELECT id, project_id, data, start_page, end_page, wday, note, text, created_at, updated_at
		FROM logs
		WHERE id = $1
	`

	var log models.Log
	var data, note, text *string
	var createdAt, updatedAt time.Time

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&log.ID,
		&log.ProjectID,
		&data,
		&log.StartPage,
		&log.EndPage,
		&log.Wday,
		&note,
		&text,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("log with ID %d not found: %w", id, err)
		}
		return nil, fmt.Errorf("failed to get log by ID %d: %w", id, err)
	}

	log.Data = data
	log.Note = note
	log.Text = text
	log.CreatedAt = &createdAt
	log.UpdatedAt = &updatedAt

	return &log, nil
}

// GetByProjectID retrieves all logs for a given project ID
func (r *LogRepositoryImpl) GetByProjectID(ctx context.Context, projectID int64) ([]*models.Log, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultContextTimeout)
	defer cancel()

	query := `
		SELECT id, project_id, data, start_page, end_page, wday, note, text, created_at, updated_at
		FROM logs
		WHERE project_id = $1
		ORDER BY id ASC
	`

	rows, err := r.pool.Query(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to query logs by project ID %d: %w", projectID, err)
	}
	defer rows.Close()

	var logs []*models.Log
	for rows.Next() {
		var log models.Log
		var data, note, text *string
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&log.ID,
			&log.ProjectID,
			&data,
			&log.StartPage,
			&log.EndPage,
			&log.Wday,
			&note,
			&text,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan log row: %w", err)
		}

		log.Data = data
		log.Note = note
		log.Text = text
		log.CreatedAt = &createdAt
		log.UpdatedAt = &updatedAt

		logs = append(logs, &log)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating logs: %w", err)
	}

	return logs, nil
}

// GetByProjectIDOrdered retrieves logs for a project ordered by data DESC
func (r *LogRepositoryImpl) GetByProjectIDOrdered(ctx context.Context, projectID int64) ([]*models.Log, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultContextTimeout)
	defer cancel()

	query := `
		SELECT id, project_id, data, start_page, end_page, wday, note, text, created_at, updated_at
		FROM logs
		WHERE project_id = $1
		ORDER BY data DESC
	`

	rows, err := r.pool.Query(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to query logs by project ID %d: %w", projectID, err)
	}
	defer rows.Close()

	var logs []*models.Log
	for rows.Next() {
		var log models.Log
		var data, note, text *string
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&log.ID,
			&log.ProjectID,
			&data,
			&log.StartPage,
			&log.EndPage,
			&log.Wday,
			&note,
			&text,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan log row: %w", err)
		}

		log.Data = data
		log.Note = note
		log.Text = text
		log.CreatedAt = &createdAt
		log.UpdatedAt = &updatedAt

		logs = append(logs, &log)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating logs: %w", err)
	}

	return logs, nil
}

// GetAll retrieves all logs
func (r *LogRepositoryImpl) GetAll(ctx context.Context) ([]*models.Log, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultContextTimeout)
	defer cancel()

	query := `
		SELECT id, project_id, data, start_page, end_page, wday, note, text, created_at, updated_at
		FROM logs
		ORDER BY id ASC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all logs: %w", err)
	}
	defer rows.Close()

	var logs []*models.Log
	for rows.Next() {
		var log models.Log
		var data, note, text *string
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&log.ID,
			&log.ProjectID,
			&data,
			&log.StartPage,
			&log.EndPage,
			&log.Wday,
			&note,
			&text,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan log row: %w", err)
		}

		log.Data = data
		log.Note = note
		log.Text = text
		log.CreatedAt = &createdAt
		log.UpdatedAt = &updatedAt

		logs = append(logs, &log)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating logs: %w", err)
	}

	return logs, nil
}
