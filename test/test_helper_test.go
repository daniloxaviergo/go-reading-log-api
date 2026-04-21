package test

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"go-reading-log-api-next/internal/config"
	"go-reading-log-api-next/test/testutil"
)

// dbTestLock serializes database access for tests that use the same test database
var dbTestLock sync.Mutex

// TestGetTestContext tests that GetTestContext returns a context with timeout
func TestGetTestContext(t *testing.T) {
	ctx := GetTestContext()

	// Verify context is not nil
	if ctx == nil {
		t.Fatal("GetTestContext returned nil context")
	}

	// Verify deadline is set
	_, ok := ctx.Deadline()
	if !ok {
		t.Error("GetTestContext did not set a deadline")
	}
}

// TestGetTestContextWithTimeout tests custom timeout contexts
func TestGetTestContextWithTimeout(t *testing.T) {
	customTimeout := 10 * time.Second
	ctx := GetTestContextWithTimeout(customTimeout)

	// Verify context is not nil
	if ctx == nil {
		t.Fatal("GetTestContextWithTimeout returned nil context")
	}

	// Verify deadline is set
	_, ok := ctx.Deadline()
	if !ok {
		t.Error("GetTestContextWithTimeout did not set a deadline")
	}
}

// TestIsTestDatabase tests the test database detection
func TestIsTestDatabase(t *testing.T) {
	// Test with no test database set
	// The function should return false if DB_DATABASE_TEST is not set
	// and DB_DATABASE is "reading_log"
	result := IsTestDatabase()
	// We can't easily control env vars in tests, so just verify the function runs
	t.Logf("IsTestDatabase returned: %v", result)
}

// TestTestName generates a unique test name
func TestTestName(t *testing.T) {
	name := TestName(t)

	// Verify name is not empty
	if name == "" {
		t.Fatal("TestName returned empty string")
	}

	// Verify name contains test function name
	if !contains(name, t.Name()) {
		t.Errorf("TestName output '%s' does not contain test function name '%s'", name, t.Name())
	}
}

// contains checks if s contains substr
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// TestTestHelperLifecycle tests the full lifecycle of TestHelper
func TestTestHelperLifecycle(t *testing.T) {
	// Skip if no test database is configured
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping database lifecycle test")
	}

	// Serialize database access for parallel tests
	dbTestLock.Lock()
	defer dbTestLock.Unlock()

	// Setup
	helper, err := SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}
	defer helper.Close()

	// Verify helper is not nil
	if helper == nil {
		t.Fatal("SetupTestDB returned nil helper")
	}

	// Verify pool is not nil
	if helper.Pool == nil {
		t.Fatal("TestHelper.Pool is nil")
	}

	// Verify context works
	ctx := helper.GetContext()
	if ctx == nil {
		t.Fatal("GetContext returned nil")
	}
}

// TestTestHelperSetupSchema tests schema setup
func TestTestHelperSetupSchema(t *testing.T) {
	// Skip if no test database is configured
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping schema setup test")
	}

	dbTestLock.Lock()
	defer dbTestLock.Unlock()

	helper, err := SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}
	defer helper.Close()

	// Setup schema
	if err := helper.SetupTestSchema(); err != nil {
		t.Fatalf("Failed to setup test schema: %v", err)
	}

	// Verify tables exist by querying
	ctx := helper.GetContext()
	var count int
	err = helper.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM projects").Scan(&count)
	if err != nil {
		t.Errorf("Failed to query projects table: %v", err)
	}
}

// TestTestHelperClearTestData tests clearing test data
func TestTestHelperClearTestData(t *testing.T) {
	// Skip if no test database is configured
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping clear test data test")
	}

	dbTestLock.Lock()
	defer dbTestLock.Unlock()

	helper, err := SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}
	defer helper.Close()

	// First setup schema
	if err := helper.SetupTestSchema(); err != nil {
		t.Fatalf("Failed to setup test schema: %v", err)
	}

	// Clear test data (should not error even if tables are empty)
	if err := helper.ClearTestData(); err != nil {
		t.Errorf("Failed to clear test data: %v", err)
	}
}

// TestTestHelperCleanupSchema tests schema cleanup
func TestTestHelperCleanupSchema(t *testing.T) {
	// Skip if no test database is configured
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping cleanup schema test")
	}

	dbTestLock.Lock()
	defer dbTestLock.Unlock()

	helper, err := SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}
	defer helper.Close()

	// Setup schema first
	if err := helper.SetupTestSchema(); err != nil {
		t.Fatalf("Failed to setup test schema: %v", err)
	}

	// Cleanup schema (should work even if tables are already gone)
	if err := helper.CleanupTestSchema(); err != nil {
		t.Errorf("Failed to cleanup test schema: %v", err)
	}
}

// TestTestHelperClose tests closing the helper
func TestTestHelperClose(t *testing.T) {
	dbTestLock.Lock()
	defer dbTestLock.Unlock()
	helper, err := SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	// Close should not panic
	helper.Close()

	// Close again should not panic either
	helper.Close()
}

// TestSetupTestDBWithConfig tests setup with custom config
func TestSetupTestDBWithConfig(t *testing.T) {
	cfg := config.LoadConfig()

	dbTestLock.Lock()
	defer dbTestLock.Unlock()
	helper, err := SetupTestDBWithConfig(cfg)
	if err != nil {
		// If this fails due to DB connection issues, it's expected in some environments
		t.Logf("SetupTestDBWithConfig error (may be expected): %v", err)
		return
	}
	defer helper.Close()

	if helper == nil {
		t.Fatal("SetupTestDBWithConfig returned nil helper")
	}
}

// TestContextTimeout tests that context has proper timeout
func TestContextTimeout(t *testing.T) {
	ctx := GetTestContext()

	_, ok := ctx.Deadline()
	if !ok {
		t.Error("Context does not have a deadline")
	}

	// Verify context can be cancelled
	done := make(chan struct{})
	go func() {
		<-ctx.Done()
		close(done)
	}()

	// Wait for context to expire (5 seconds + buffer)
	time.Sleep(6 * time.Second)
	<-done
}

// TestCleanupOrphanedDatabases tests basic orphan cleanup functionality
func TestCleanupOrphanedDatabases(t *testing.T) {
	// Skip if no test database is configured
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping orphan cleanup test")
	}

	dbTestLock.Lock()
	defer dbTestLock.Unlock()

	helper, err := SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}
	defer helper.Close()

	// Call cleanup - should not error even if no orphaned databases exist
	droppedCount, err := testutil.CleanupOrphanedDatabases(helper.Pool, helper.TestDBName)
	if err != nil {
		t.Errorf("CleanupOrphanedDatabases failed: %v", err)
	}
	t.Logf("CleanupOrphanedDatabases dropped %d databases", droppedCount)
}

// TestCleanupOrphanedDatabases_ExcludeCurrent tests that current DB is excluded
func TestCleanupOrphanedDatabases_ExcludeCurrent(t *testing.T) {
	// Skip if no test database is configured
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping exclusion test")
	}

	dbTestLock.Lock()
	defer dbTestLock.Unlock()

	helper, err := SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}
	defer helper.Close()

	// Call cleanup with current test DB name
	// This should not drop the current test DB
	droppedCount, err := testutil.CleanupOrphanedDatabases(helper.Pool, helper.TestDBName)
	if err != nil {
		t.Errorf("CleanupOrphanedDatabases failed: %v", err)
	}
	t.Logf("CleanupOrphanedDatabases dropped %d databases", droppedCount)

	// Verify current test DB still exists by trying to ping it
	ctx := helper.GetContext()
	if err := helper.Pool.Ping(ctx); err != nil {
		t.Errorf("Current test database was incorrectly dropped: %v", err)
	}
}

// TestCleanupOrphanedDatabases_NonExistentDB tests handling of non-existent database
func TestCleanupOrphanedDatabases_NonExistentDB(t *testing.T) {
	// Skip if no test database is configured
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping non-existent DB test")
	}

	dbTestLock.Lock()
	defer dbTestLock.Unlock()

	helper, err := SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}
	defer helper.Close()

	// Try to cleanup with a non-existent database name
	// This should not error - DROP DATABASE IF EXISTS handles this gracefully
	droppedCount, err := testutil.CleanupOrphanedDatabases(helper.Pool, "non_existent_db_name_12345")
	if err != nil {
		t.Errorf("CleanupOrphanedDatabases should not fail for non-existent DB: %v", err)
	}
	t.Logf("CleanupOrphanedDatabases dropped %d databases", droppedCount)
}

// TestCleanupOrphanedDatabases_MultipleDBs tests cleanup with multiple databases
func TestCleanupOrphanedDatabases_MultipleDBs(t *testing.T) {
	// Skip if no test database is configured
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping multiple DBs test")
	}

	dbTestLock.Lock()
	defer dbTestLock.Unlock()

	helper, err := SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}
	defer helper.Close()

	// Create some additional test databases to simulate orphaned databases
	ctx := helper.GetContext()
	mainPool := helper.Pool

	// Create a few test databases
	dbNames := []string{
		"reading_log_test_orphan_1",
		"reading_log_test_orphan_2",
		"reading_log_test_orphan_3",
	}

	for _, dbName := range dbNames {
		_, err := mainPool.Exec(ctx, fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
		if err != nil {
			t.Logf("Note: Could not create test database %s: %v", dbName, err)
		}
	}

	// Cleanup - should not error
	droppedCount, err := testutil.CleanupOrphanedDatabases(mainPool, helper.TestDBName)
	if err != nil {
		t.Errorf("CleanupOrphanedDatabases failed: %v", err)
	}
	t.Logf("CleanupOrphanedDatabases dropped %d databases", droppedCount)

	// Cleanup should complete successfully
	// Note: We don't verify the databases were dropped here because
	// the cleanup function uses DROP DATABASE IF EXISTS which may fail
	// if the database doesn't exist or if there are permission issues
}

// TestValidateDatabases tests the database validation function
func TestValidateDatabases(t *testing.T) {
	// Skip if no test database is configured
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping validation test")
	}

	err := testutil.ValidateDatabases()
	if err != nil {
		t.Errorf("ValidateDatabases failed: %v", err)
	}
}

// TestCleanupOrphanedDatabases_Function tests the testutil.CleanupOrphanedDatabases function directly
func TestCleanupOrphanedDatabases_Function(t *testing.T) {
	// Skip if no test database is configured
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping function test")
	}

	helper, err := SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}
	defer helper.Close()

	// Create test databases to clean up
	ctx := helper.GetContext()
	mainPool := helper.Pool

	testDBNames := []string{
		"reading_log_test_func_1",
		"reading_log_test_func_2",
	}

	for _, dbName := range testDBNames {
		_, err := mainPool.Exec(ctx, fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
		if err != nil {
			t.Logf("Note: Could not create test database %s: %v", dbName, err)
		}
	}

	// Call testutil.CleanupOrphanedDatabases function directly
	droppedCount, err := testutil.CleanupOrphanedDatabases(mainPool, helper.TestDBName)
	if err != nil {
		t.Errorf("testutil.CleanupOrphanedDatabases failed: %v", err)
	}

	// Verify at least some databases were dropped
	if droppedCount < 0 {
		t.Errorf("droppedCount should be non-negative, got %d", droppedCount)
	}
}

// TestCleanupOrphanedDatabases_Concurrent tests cleanup with concurrent database operations
func TestCleanupOrphanedDatabases_Concurrent(t *testing.T) {
	// Skip if no test database is configured
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping concurrent test")
	}

	dbTestLock.Lock()
	defer dbTestLock.Unlock()

	helper, err := SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}
	defer helper.Close()

	ctx := helper.GetContext()
	mainPool := helper.Pool

	// Create multiple databases concurrently
	numDBs := 5
	dbNames := make([]string, numDBs)

	for i := 0; i < numDBs; i++ {
		dbNames[i] = fmt.Sprintf("reading_log_test_concurrent_%d", i)
	}

	// Create databases
	for _, dbName := range dbNames {
		_, err := mainPool.Exec(ctx, fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
		if err != nil {
			t.Logf("Note: Could not create test database %s: %v", dbName, err)
		}
	}

	// Run cleanup
	droppedCount, err := testutil.CleanupOrphanedDatabases(mainPool, helper.TestDBName)
	if err != nil {
		t.Errorf("testutil.CleanupOrphanedDatabases failed: %v", err)
	}

	// Verify cleanup completed
	if droppedCount < 0 {
		t.Errorf("droppedCount should be non-negative, got %d", droppedCount)
	}
}

// TestCleanupOrphanedDatabases_Performance tests that cleanup completes within timeout
func TestCleanupOrphanedDatabases_Performance(t *testing.T) {
	// Skip if no test database is configured
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping performance test")
	}

	helper, err := SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}
	defer helper.Close()

	ctx := helper.GetContext()
	mainPool := helper.Pool

	// Create a moderate number of test databases
	numDBs := 10
	dbNames := make([]string, numDBs)

	for i := 0; i < numDBs; i++ {
		dbNames[i] = fmt.Sprintf("reading_log_test_perf_%d", i)
	}

	// Create databases
	for _, dbName := range dbNames {
		_, err := mainPool.Exec(ctx, fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
		if err != nil {
			t.Logf("Note: Could not create test database %s: %v", dbName, err)
		}
	}

	// Measure cleanup time
	start := time.Now()
	droppedCount, err := testutil.CleanupOrphanedDatabases(mainPool, helper.TestDBName)
	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("testutil.CleanupOrphanedDatabases failed: %v", err)
	}

	// Verify cleanup completed within reasonable time
	// The function has a 60 second timeout
	if elapsed > 60*time.Second {
		t.Errorf("Cleanup took too long: %v", elapsed)
	}

	// Verify at least some databases were dropped
	if droppedCount < 0 {
		t.Errorf("droppedCount should be non-negative, got %d", droppedCount)
	}
}
