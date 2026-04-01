package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSMiddleware_PreflightRequest(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handler := CORSMiddleware(next)

	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, w.Code)
	}

	if w.Body.Len() != 0 {
		t.Errorf("Expected empty body for preflight request, got %d bytes", w.Body.Len())
	}
}

func TestCORSMiddleware_NormalRequest(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handler := CORSMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	if body != "OK" {
		t.Errorf("Expected body 'OK', got '%s'", body)
	}
}

func TestCORSMiddleware_HeadersSet(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := CORSMiddleware(next)

	tests := []struct {
		method string
		url    string
	}{
		{http.MethodGet, "/test"},
		{http.MethodPost, "/test"},
		{http.MethodPut, "/test"},
		{http.MethodDelete, "/test"},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			headers := w.Header()

			if headers.Get("Access-Control-Allow-Origin") != "*" {
				t.Errorf("Expected Access-Control-Allow-Origin header to be '*', got '%s'", headers.Get("Access-Control-Allow-Origin"))
			}

			if headers.Get("Access-Control-Allow-Methods") != "GET, POST, PUT, DELETE, OPTIONS" {
				t.Errorf("Expected Access-Control-Allow-Methods header to be 'GET, POST, PUT, DELETE, OPTIONS', got '%s'", headers.Get("Access-Control-Allow-Methods"))
			}

			if headers.Get("Access-Control-Allow-Headers") != "Content-Type, Authorization" {
				t.Errorf("Expected Access-Control-Allow-Headers header to be 'Content-Type, Authorization', got '%s'", headers.Get("Access-Control-Allow-Headers"))
			}
		})
	}
}

func TestCORSMiddleware_PropagatesRequest(t *testing.T) {
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("next handler called"))
	})

	handler := CORSMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if !nextCalled {
		t.Error("Expected next handler to be called, but it was not")
	}
}
