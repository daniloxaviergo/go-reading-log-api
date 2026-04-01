package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecoveryMiddleware_PanicCaught(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})
	recoveryHandler := RecoveryMiddleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Recovery middleware failed to catch panic: %v", r)
			}
		}()

		recoveryHandler.ServeHTTP(w, req)
	}()

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d after panic, got %d", http.StatusInternalServerError, w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Internal Server Error") {
		t.Errorf("Expected 'Internal Server Error' in body, got '%s'", body)
	}
}

func TestRecoveryMiddleware_NoPanic(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	recoveryHandler := RecoveryMiddleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	recoveryHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d for normal request, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	if body != "OK" {
		t.Errorf("Expected body 'OK', got '%s'", body)
	}
}

func TestRecoveryMiddleware_MultiplePanics(t *testing.T) {
	// First request - should catch panic
	handler1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})
	recoveryHandler := RecoveryMiddleware(handler1)

	req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	w1 := httptest.NewRecorder()

	recoveryHandler.ServeHTTP(w1, req1)

	if w1.Code != http.StatusInternalServerError {
		t.Errorf("First request: expected status code %d, got %d", http.StatusInternalServerError, w1.Code)
	}

	// Second request - should work again (middleware should reset)
	handler2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	recoveryHandler = RecoveryMiddleware(handler2)
	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	w2 := httptest.NewRecorder()

	recoveryHandler.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("Second request: expected status code %d, got %d", http.StatusOK, w2.Code)
	}
}
