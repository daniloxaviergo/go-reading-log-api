package unit

import (
	"context"
	"fmt"
	"testing"

	"go-reading-log-api-next/internal/domain/models"
	"go-reading-log-api-next/test"
)

// TestLogRepositoryIntegration tests the mock log repository
func TestLogRepositoryIntegration(t *testing.T) {
	t.Log("Log repository tests are in test package")
}

// TestMockLogRepositoryTests tests that the mock repository tests work correctly
func TestMockLogRepositoryTests(t *testing.T) {
	t.Log("Mock log repository tests are in test package")
}

// TestLogRepositoryGetByID tests the GetByID method
func TestLogRepositoryGetByID(t *testing.T) {
	mock := test.NewMockLogRepository()

	// Test adding and retrieving a log
	log := &models.Log{
		ID:        1,
		ProjectID: 1,
		StartPage: 1,
		EndPage:   10,
	}
	mock.AddLog(log)

	ctx := context.Background()
	result, err := mock.GetByID(ctx, 1)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.ID != 1 {
		t.Errorf("Expected log ID 1, got %d", result.ID)
	}
}

// TestLogRepositoryGetByIDNotFound tests GetByID for non-existent log
func TestLogRepositoryGetByIDNotFound(t *testing.T) {
	mock := test.NewMockLogRepository()

	ctx := context.Background()
	_, err := mock.GetByID(ctx, 999)

	if err == nil {
		t.Error("Expected error for non-existent log, got nil")
	}
}

// TestLogRepositoryGetByProjectID tests the GetByProjectID method
func TestLogRepositoryGetByProjectID(t *testing.T) {
	mock := test.NewMockLogRepository()

	// Add test logs for project
	log1 := &models.Log{ID: 1, ProjectID: 1}
	log2 := &models.Log{ID: 2, ProjectID: 1}
	log3 := &models.Log{ID: 3, ProjectID: 1}

	mock.AddLog(log1)
	mock.AddLog(log2)
	mock.AddLog(log3)
	mock.AddLogsForProject(1, []*models.Log{log1, log2, log3})

	ctx := context.Background()
	result, err := mock.GetByProjectID(ctx, 1)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 logs, got %d", len(result))
	}
}

// TestLogRepositoryGetByProjectIDEmpty tests GetByProjectID with no logs
func TestLogRepositoryGetByProjectIDEmpty(t *testing.T) {
	mock := test.NewMockLogRepository()

	ctx := context.Background()
	result, err := mock.GetByProjectID(ctx, 1)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected 0 logs, got %d", len(result))
	}
}

// TestLogRepositoryGetAll tests the GetAll method
func TestLogRepositoryGetAll(t *testing.T) {
	mock := test.NewMockLogRepository()

	// Add test logs
	log1 := &models.Log{ID: 1}
	log2 := &models.Log{ID: 2}
	log3 := &models.Log{ID: 3}

	mock.AddLog(log1)
	mock.AddLog(log2)
	mock.AddLog(log3)

	ctx := context.Background()
	result, err := mock.GetAll(ctx)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 logs, got %d", len(result))
	}
}

// TestLogRepositoryGetAllEmpty tests GetAll with no logs
func TestLogRepositoryGetAllEmpty(t *testing.T) {
	mock := test.NewMockLogRepository()

	ctx := context.Background()
	result, err := mock.GetAll(ctx)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected 0 logs, got %d", len(result))
	}
}

// TestLogRepositoryError handling tests
func TestLogRepositoryError(t *testing.T) {
	mock := test.NewMockLogRepository()
	mock.SetError(fmt.Errorf("database error"))

	ctx := context.Background()

	_, err := mock.GetByID(ctx, 1)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	_, err = mock.GetByProjectID(ctx, 1)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	_, err = mock.GetAll(ctx)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// TestLogRepositoryCallTracking tests call tracking
func TestLogRepositoryCallTracking(t *testing.T) {
	mock := test.NewMockLogRepository()

	ctx := context.Background()
	mock.GetByID(ctx, 1)
	mock.GetByID(ctx, 2)
	mock.GetByProjectID(ctx, 1)

	if mock.CallCount() != 2 {
		t.Errorf("Expected 2 GetByID calls, got %d", mock.CallCount())
	}

	if mock.GetByProjectIDCallCount() != 1 {
		t.Errorf("Expected 1 GetByProjectID call, got %d", mock.GetByProjectIDCallCount())
	}

	if mock.GetByIDLastCall() != 2 {
		t.Errorf("Expected last GetByID call to be ID 2, got %d", mock.GetByIDLastCall())
	}
}
