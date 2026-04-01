package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestContextWithTimeout tests the ContextWithTimeout helper function
func TestContextWithTimeout(t *testing.T) {
	parentCtx := context.Background()
	timeout := 5 * time.Second

	ctx, cancel := ContextWithTimeout(parentCtx, timeout)
	defer cancel()

	// Verify context has deadline
	_, ok := ctx.Deadline()
	if !ok {
		t.Error("Expected context to have deadline")
	}

	// Verify context can be used for operations
	select {
	case <-ctx.Done():
		t.Error("Context should not be done initially")
	default:
		// Expected - context should be active
	}
}

// TestContextWithTimeout_Cancel tests that cancel works correctly
func TestContextWithTimeout_Cancel(t *testing.T) {
	ctx, cancel := ContextWithTimeout(context.Background(), 10*time.Second)

	// Cancel the context
	cancel()

	// Verify context is done after cancel
	select {
	case <-ctx.Done():
		// Expected
	default:
		t.Error("Context should be done after cancel")
	}
}

// TestContextWithTimeout_Timeout tests that context times out correctly
func TestContextWithTimeout_Timeout(t *testing.T) {
	ctx, cancel := ContextWithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// Wait for timeout
	time.Sleep(10 * time.Millisecond)

	select {
	case <-ctx.Done():
		// Expected - context should be timed out
	default:
		t.Error("Context should be done after timeout")
	}
}

// TestChain tests the Chain middleware function
func TestChain(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Chain middleware
	middleware := []func(http.Handler) http.Handler{
		RecoveryMiddleware,
		CORSMiddleware,
		RequestIDMiddleware,
	}

	chain := Chain(handler, middleware...)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	chain.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	if w.Body.String() != "OK" {
		t.Errorf("Expected body 'OK', got '%s'", w.Body.String())
	}
}

// TestChain_PanicRecovery tests that the middleware chain properly handles panics
func TestChain_PanicRecovery(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	middleware := []func(http.Handler) http.Handler{
		RecoveryMiddleware,
	}

	chain := Chain(handler, middleware...)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	chain.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d after panic, got %d", http.StatusInternalServerError, w.Code)
	}
}

// TestChain_CORS tests CORS middleware
func TestChain_CORS(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := []func(http.Handler) http.Handler{
		CORSMiddleware,
	}

	chain := Chain(handler, middleware...)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	chain.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Verify CORS headers are set
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("Expected CORS header 'Access-Control-Allow-Origin: *'")
	}
}

// TestChain_RequestID tests RequestID middleware
func TestChain_RequestID(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := GetRequestIDFromContext(r.Context())
		if reqID == "" {
			t.Error("Expected request ID in context")
		}
		w.WriteHeader(http.StatusOK)
	})

	middleware := []func(http.Handler) http.Handler{
		RequestIDMiddleware,
	}

	chain := Chain(handler, middleware...)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	chain.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

// TestChain_Logging tests LoggingMiddleware
func TestChain_Logging(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := []func(http.Handler) http.Handler{
		LoggingMiddleware,
	}

	chain := Chain(handler, middleware...)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	chain.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

// TestChain_Order tests that middleware order is correct
func TestChain_Order(t *testing.T) {
	order := []string{}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
	})

	middleware := []func(http.Handler) http.Handler{
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, "mw1")
				next.ServeHTTP(w, r)
			})
		},
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, "mw2")
				next.ServeHTTP(w, r)
			})
		},
	}

	chain := Chain(handler, middleware...)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	chain.ServeHTTP(w, req)

	// Verify order: middleware apply in reverse order (innermost first)
	// Chain applies: mw2(mw1(handler)) (iterates backwards over [mw1, mw2])
	// So order should be: mw1, mw2, handler
	if len(order) != 3 {
		t.Errorf("Expected 3 middleware calls, got %d", len(order))
	}

	if order[0] != "mw1" {
		t.Errorf("Expected 'mw1' first, got '%s'", order[0])
	}
	if order[1] != "mw2" {
		t.Errorf("Expected 'mw2' second, got '%s'", order[1])
	}
	if order[2] != "handler" {
		t.Errorf("Expected 'handler' last, got '%s'", order[2])
	}
}
