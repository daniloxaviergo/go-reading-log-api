// cleanup.go
// Utility functions for test database cleanup
// Exports functions that can be imported by test files

package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go-reading-log-api-next/internal/config"
)

const (
	// Timeout for the entire cleanup operation
	cleanupTimeout = 60 * time.Second
	// Pattern for test database names
	testDBPattern = "reading_log_test_%"
)

// CleanupOrphanedDatabases identifies and drops test databases matching the pattern
// excludeName is the current test database name to exclude from cleanup
// Returns the number of dropped databases and any error encountered
func CleanupOrphanedDatabases(pool *pgxpool.Pool, excludeName string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cleanupTimeout)
	defer cancel()

	// Query for orphaned databases matching the pattern
	query := `
		SELECT datname 
		FROM pg_database 
		WHERE datname LIKE $1
		AND datname != $2
		AND pg_catalog.pg_get_userbyid(datdba) = current_user
	`

	rows, err := pool.Query(ctx, query, testDBPattern, excludeName)
	if err != nil {
		return 0, fmt.Errorf("failed to query test databases: %w", err)
	}
	defer rows.Close()

	var toDrop []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			// Log error but continue with remaining databases
			continue
		}
		toDrop = append(toDrop, name)
	}

	if len(toDrop) == 0 {
		return 0, nil
	}

	// Drop each orphaned database
	droppedCount := 0
	for _, dbName := range toDrop {
		_, dropErr := pool.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
		if dropErr != nil {
			// Log error but continue with remaining databases
			continue
		}
		droppedCount++
	}

	return droppedCount, nil
}

// ValidateDatabases runs a quick validation to ensure database connectivity
func ValidateDatabases() error {
	cfg := config.LoadConfig()

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBDatabase,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}
	defer pool.Close()

	// Ping to verify connection
	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}
