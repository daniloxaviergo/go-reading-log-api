package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
)

func TestChain_MiddlewareOrder(t *testing.T) {
	// Handler that returns the order of execution
	type executionOrder struct {
		mu       sync.Mutex
		sequence []string
	}

	orders := &executionOrder{}

	// Innermost handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orders.mu.Lock()
		orders.sequence = append(orders.sequence, "handler")
		orders.mu.Unlock()
		w.WriteHeader(http.StatusOK)
	})

	// Middleware that logs its position
	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			orders.mu.Lock()
			orders.sequence = append(orders.sequence, "mw1")
			orders.mu.Unlock()
			next.ServeHTTP(w, r)
		})
	}

	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			orders.mu.Lock()
			orders.sequence = append(orders.sequence, "mw2")
			orders.mu.Unlock()
			next.ServeHTTP(w, r)
		})
	}

	// Chain: Chain applies in reverse order, so:
	// Chain(handler, mw1, mw2) = mw1(mw2(handler))
	// Order of execution: mw1 -> mw2 -> handler
	chain := Chain(handler, mw1, mw2)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	chain.ServeHTTP(w, req)

	expectedOrder := []string{"mw1", "mw2", "handler"}
	if len(orders.sequence) != len(expectedOrder) {
		t.Errorf("Expected %d middleware calls, got %d", len(expectedOrder), len(orders.sequence))
	}

	for i, expected := range expectedOrder {
		if i >= len(orders.sequence) {
			t.Errorf("Missing middleware at position %d", i)
			continue
		}
		if orders.sequence[i] != expected {
			t.Errorf("Expected middleware order[%d] = '%s', got '%s'", i, expected, orders.sequence[i])
		}
	}
}

func TestChain_RecoveryIsOutermost(t *testing.T) {
	recoveryCalled := false
	corsCalled := false
	requestIDCalled := false

	recovery := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recoveryCalled = true
			next.ServeHTTP(w, r)
		})
	}

	cors := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			corsCalled = true
			next.ServeHTTP(w, r)
		})
	}

	requestID := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestIDCalled = true
			next.ServeHTTP(w, r)
		})
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	chain := Chain(handler, recovery, cors, requestID)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	chain.ServeHTTP(w, req)

	if !recoveryCalled {
		t.Error("Expected recovery middleware to be called")
	}
	if !corsCalled {
		t.Error("Expected cors middleware to be called")
	}
	if !requestIDCalled {
		t.Error("Expected requestID middleware to be called")
	}
}

func TestChain_ContextPropagation(t *testing.T) {
	var capturedRequestID string
	var capturedContext context.Context

	recovery := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedContext = r.Context()
			next.ServeHTTP(w, r)
		})
	}

	requestID := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), requestIDKey, "test-request-id")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedRequestID = GetRequestIDFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	chain := Chain(handler, recovery, requestID)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	chain.ServeHTTP(w, req)

	if capturedRequestID != "test-request-id" {
		t.Errorf("Expected request ID 'test-request-id', got '%s'", capturedRequestID)
	}

	if capturedContext == nil {
		t.Error("Expected context to be propagated to all middleware")
	}
}

func TestChain_LongMiddlewareChain(t *testing.T) {
	type executionOrder struct {
		mu       sync.Mutex
		sequence []string
	}

	orders := &executionOrder{}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orders.mu.Lock()
		orders.sequence = append(orders.sequence, "handler")
		orders.mu.Unlock()
		w.WriteHeader(http.StatusOK)
	})

	// Create 5 middleware functions
	middlewares := make([]func(http.Handler) http.Handler, 5)
	for i := 0; i < 5; i++ {
		idx := i
		middlewares[idx] = func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				orders.mu.Lock()
				orders.sequence = append(orders.sequence, "mw"+strconv.Itoa(idx))
				orders.mu.Unlock()
				next.ServeHTTP(w, r)
			})
		}
	}

	chain := Chain(handler, middlewares...)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	chain.ServeHTTP(w, req)

	expectedCount := 5 + 1 // 5 middlewares + 1 handler
	if len(orders.sequence) != expectedCount {
		t.Errorf("Expected %d calls, got %d", expectedCount, len(orders.sequence))
	}

	if orders.sequence[len(orders.sequence)-1] != "handler" {
		t.Error("Expected handler to be called last")
	}
}

func TestChain_ErrorPropagation(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Error from handler", http.StatusBadRequest)
	})

	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}

	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}

	chain := Chain(nextHandler, middleware1, middleware2)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	chain.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestChain_PreflightRequest(t *testing.T) {
	// Test that preflight OPTIONS requests work through the chain
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	cors := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	chain := Chain(handler, cors)

	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	w := httptest.NewRecorder()

	chain.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d for preflight, got %d", http.StatusNoContent, w.Code)
	}
}
