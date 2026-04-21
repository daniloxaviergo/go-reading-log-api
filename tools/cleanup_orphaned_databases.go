// cleanup_orphaned_databases.go
// Standalone script for cleaning up orphaned test databases
// Can be run directly or invoked via make test-clean

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go-reading-log-api-next/internal/config"
)

const (
	// Timeout for the entire cleanup operation
	cleanupTimeout = 60 * time.Second
	// Pattern for test database names
	testDBPattern = "reading_log_test_%"
)

func main() {
	// Load environment configuration
	if err := loadConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Get configuration
	cfg := config.LoadConfig()

	// Build connection string for main database
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBDatabase,
	)

	// Create connection pool
	ctx, cancel := context.WithTimeout(context.Background(), cleanupTimeout)
	defer cancel()

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating connection pool: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Get excluded database name (current test database if set)
	excludeName := os.Getenv("DB_DATABASE_TEST")
	if excludeName == "" {
		excludeName = cfg.DBDatabase + "_test"
	}

	// Run cleanup
	fmt.Printf("Starting cleanup of orphaned test databases...\n")
	fmt.Printf("Excluding current test database: %s\n", excludeName)
	fmt.Printf("Searching for databases matching pattern: %s\n\n", testDBPattern)

	droppedCount, err := cleanupOrphanedDatabases(pool, excludeName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cleanup completed with warnings: %v\n", err)
		os.Exit(0) // Graceful exit - errors are logged but don't fail the command
	}

	fmt.Printf("\n========================================\n")
	fmt.Printf("Cleanup complete!\n")
	fmt.Printf("Dropped %d orphaned test database(s)\n", droppedCount)
	fmt.Printf("========================================\n")

	os.Exit(0)
}

// loadConfig loads environment configuration from .env.test
func loadConfig() error {
	// Try to load .env.test
	if err := godotenv.Load(".env.test"); err != nil {
		// If .env.test doesn't exist, try .env
		if err := godotenv.Load(".env"); err != nil {
			// No environment file found, use defaults
			fmt.Println("No .env.test or .env file found, using environment variables")
		}
	}
	return nil
}

// cleanupOrphanedDatabases identifies and drops test databases matching the pattern
// excludeName is the current test database name to exclude from cleanup
// Returns the number of dropped databases and any error encountered
func cleanupOrphanedDatabases(pool *pgxpool.Pool, excludeName string) (int, error) {
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
			fmt.Fprintf(os.Stderr, "Warning: failed to scan database name: %v\n", err)
			continue
		}
		toDrop = append(toDrop, name)
	}

	if len(toDrop) == 0 {
		fmt.Println("No orphaned test databases found.")
		return 0, nil
	}

	fmt.Printf("Found %d orphaned test database(s) to clean up:\n", len(toDrop))

	// Drop each orphaned database
	droppedCount := 0
	for _, dbName := range toDrop {
		// Print progress
		fmt.Printf("  Dropping %s... ", dbName)

		_, dropErr := pool.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
		if dropErr != nil {
			fmt.Printf("FAILED: %v\n", dropErr)
			// Log error but continue with remaining databases
			continue
		}

		fmt.Printf("DONE\n")
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
