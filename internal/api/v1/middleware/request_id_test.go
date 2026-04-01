package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestIDMiddleware_GeneratesUniqueIDs(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := RequestIDMiddleware(next)

	var ids []string
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		requestID := w.Header().Get("X-Request-ID")
		ids = append(ids, requestID)

		if requestID == "" {
			t.Errorf("Expected non-empty request ID, got empty string")
		}
	}

	// Verify all IDs are unique
	seen := make(map[string]bool)
	for _, id := range ids {
		if seen[id] {
			t.Errorf("Duplicate request ID found: %s", id)
		}
		seen[id] = true
	}
}

func TestRequestIDMiddleware_ContextPropagation(t *testing.T) {
	var capturedID string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedID = GetRequestIDFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	handler := RequestIDMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if capturedID == "" {
		t.Error("Expected non-empty request ID in context")
	}

	// Verify the ID in context matches the header
	headerID := w.Header().Get("X-Request-ID")
	if capturedID != headerID {
		t.Errorf("Request ID in context (%s) does not match header (%s)", capturedID, headerID)
	}
}

func TestRequestIDMiddleware_ResponseHeader(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := RequestIDMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	requestID := w.Header().Get("X-Request-ID")

	if requestID == "" {
		t.Error("Expected X-Request-ID header to be set")
	}

	// Verify UUID format (8-4-4-4-12 hex chars)
	parts := strings.Split(requestID, "-")
	if len(parts) != 5 {
		t.Errorf("Expected UUID format with 5 parts, got %d parts", len(parts))
	}

	expectedLengths := []int{8, 4, 4, 4, 12}
	for i, length := range expectedLengths {
		if len(parts[i]) != length {
			t.Errorf("Expected UUID part %d to have length %d, got %d", i, length, len(parts[i]))
		}
	}
}

func TestRequestIDMiddleware_NoExistingContext(t *testing.T) {
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	handler := RequestIDMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if !nextCalled {
		t.Error("Expected next handler to be called")
	}
}

func TestGetRequestIDFromContext_EmptyContext(t *testing.T) {
	ctx := context.Background()
	id := GetRequestIDFromContext(ctx)
	if id != "" {
		t.Errorf("Expected empty string for empty context, got '%s'", id)
	}
}

func TestGetRequestIDFromContext_ValidID(t *testing.T) {
	id := "12345678-1234-1234-1234-123456789abc"
	ctx := context.WithValue(context.Background(), requestIDKey, id)
	retrieved := GetRequestIDFromContext(ctx)
	if retrieved != id {
		t.Errorf("Expected '%s', got '%s'", id, retrieved)
	}
}

func TestGetRequestIDFromContext_WrongType(t *testing.T) {
	ctx := context.WithValue(context.Background(), requestIDKey, 12345)
	id := GetRequestIDFromContext(ctx)
	if id != "" {
		t.Errorf("Expected empty string for non-string value, got '%s'", id)
	}
}
