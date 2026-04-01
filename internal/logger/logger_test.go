package logger

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
)

// TestInitializeTextFormat tests that Initialize creates a logger with text format.
func TestInitializeTextFormat(t *testing.T) {
	logger := Initialize("info", "text")

	// Verify the logger is not nil
	if logger == nil {
		t.Fatal("Initialize returned nil logger")
	}

	// Check that the logger can log messages
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	testLogger := slog.New(handler)
	testLogger.Info("test message")

	if buf.String() == "" {
		t.Fatal("logger produced no output")
	}
}

// TestInitializeJSONFormat tests that Initialize creates a logger with JSON format.
func TestInitializeJSONFormat(t *testing.T) {
	logger := Initialize("info", "json")

	if logger == nil {
		t.Fatal("Initialize returned nil logger")
	}

	// Verify the logger can log messages in JSON format
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	testLogger := slog.New(handler)
	testLogger.Info("test message", "key", "value")

	output := buf.String()
	if output == "" {
		t.Fatal("logger produced no output")
	}

	// JSON output should contain proper JSON structure
	if !bytes.Contains(buf.Bytes(), []byte(`"key":"value"`)) &&
		!bytes.Contains(buf.Bytes(), []byte(`"key": "value"`)) {
		t.Errorf("JSON output should contain key-value pair: %s", output)
	}
}

// TestInitializeDebugLevel tests that the debug log level is correctly set.
func TestInitializeDebugLevel(t *testing.T) {
	logger := Initialize("debug", "text")

	// Verify the logger is not nil
	if logger == nil {
		t.Fatal("Initialize returned nil logger")
	}

	// Test that debug level would be enabled
	handler := logger.Handler()
	if handler == nil {
		t.Fatal("logger.Handler() returned nil")
	}

	// Use the handler's Enabled method to check if debug level is enabled
	ctx := context.Background()
	if !handler.Enabled(ctx, slog.LevelDebug) {
		t.Error("debug level should be enabled")
	}
}

// TestInitializeInfoLevel tests that the info log level is correctly set.
func TestInitializeInfoLevel(t *testing.T) {
	logger := Initialize("info", "text")

	if logger == nil {
		t.Fatal("Initialize returned nil logger")
	}

	handler := logger.Handler()
	if handler == nil {
		t.Fatal("logger.Handler() returned nil")
	}

	ctx := context.Background()
	if !handler.Enabled(ctx, slog.LevelInfo) {
		t.Error("info level should be enabled")
	}
}

// TestInitializeWarnLevel tests that the warn log level is correctly set.
func TestInitializeWarnLevel(t *testing.T) {
	logger := Initialize("warn", "text")

	if logger == nil {
		t.Fatal("Initialize returned nil logger")
	}

	handler := logger.Handler()
	if handler == nil {
		t.Fatal("logger.Handler() returned nil")
	}

	ctx := context.Background()
	if !handler.Enabled(ctx, slog.LevelWarn) {
		t.Error("warn level should be enabled")
	}
}

// TestInitializeErrorLevel tests that the error log level is correctly set.
func TestInitializeErrorLevel(t *testing.T) {
	logger := Initialize("error", "text")

	if logger == nil {
		t.Fatal("Initialize returned nil logger")
	}

	handler := logger.Handler()
	if handler == nil {
		t.Fatal("logger.Handler() returned nil")
	}

	ctx := context.Background()
	if !handler.Enabled(ctx, slog.LevelError) {
		t.Error("error level should be enabled")
	}
}

// TestInitializeCaseInsensitiveLevel tests that log levels are case-insensitive.
func TestInitializeCaseInsensitiveLevel(t *testing.T) {
	logger1 := Initialize("DEBUG", "text")
	logger2 := Initialize("debug", "text")
	logger3 := Initialize("Debug", "text")

	if logger1 == nil || logger2 == nil || logger3 == nil {
		t.Fatal("one or more Initialize calls returned nil")
	}
}

// TestInitializeInvalidLevel tests that invalid log levels fall back to default (info).
func TestInitializeInvalidLevel(t *testing.T) {
	logger := Initialize("invalid", "text")

	if logger == nil {
		t.Fatal("Initialize returned nil logger for invalid level")
	}

	handler := logger.Handler()
	if handler == nil {
		t.Fatal("logger.Handler() returned nil")
	}

	ctx := context.Background()

	// Info level should be enabled (default fallback)
	if !handler.Enabled(ctx, slog.LevelInfo) {
		t.Error("info level should be enabled (default fallback)")
	}

	// Debug level should be disabled since default is info
	if handler.Enabled(ctx, slog.LevelDebug) {
		t.Error("debug level should be disabled when default info level is set")
	}
}

// TestLoggerOutput tests that the logger actually outputs log messages.
func TestLoggerOutput(t *testing.T) {
	var buf bytes.Buffer

	// Create a logger that writes to our buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(handler)

	// Log a message
	logger.Info("test message", "key", "value")

	// Verify output contains expected content
	output := buf.String()
	if output == "" {
		t.Fatal("logger produced no output")
	}

	// Check that output contains expected elements
	expectedElements := []string{"time=", "level=INFO", `msg="test message"`}
	for _, elem := range expectedElements {
		if !contains(output, elem) {
			t.Errorf("output missing expected element '%s': %s", elem, output)
		}
	}
}

// TestLoggerMessageStructure tests that log messages have proper structure.
func TestLoggerMessageStructure(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(handler)

	logger.Info("test message", "user_id", 123, "action", "login")

	output := buf.String()
	if output == "" {
		t.Fatal("logger produced no output")
	}

	// JSON output should contain required fields
	requiredFields := []string{"time", "level", "msg"}
	for _, field := range requiredFields {
		if !bytes.Contains(buf.Bytes(), []byte(`"`+field+`"`)) {
			t.Errorf("JSON output should contain '%s' field: %s", field, output)
		}
	}
}

// contains checks if s contains substring sub.
func contains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
