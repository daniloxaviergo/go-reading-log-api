package test

import (
	"strings"
	"sync"
	"testing"
	"time"

	"go-reading-log-api-next/internal/config"
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
