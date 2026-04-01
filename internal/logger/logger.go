package logger

import (
	"log/slog"
	"os"
)

// Initialize creates and configures a structured logger using slog.
// level specifies the minimum log level (debug, info, warn, error).
// format specifies the output format (text or json).
// Returns a configured *slog.Logger that can be used throughout the application.
func Initialize(level, format string) *slog.Logger {
	// Parse log level
	lvl := parseLogLevel(level)

	// Create appropriate handler based on format
	var handler slog.Handler
	switch format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: lvl,
		})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: lvl,
		})
	}

	// Create and return logger with custom options if needed
	return slog.New(handler)
}

// parseLogLevel converts a string level to slog.Level.
// Supports case-insensitive: debug, info, warn, error.
// Falls back to info level if invalid level is provided.
func parseLogLevel(level string) slog.Level {
	switch level {
	case "debug", "DEBUG", "Debug":
		return slog.LevelDebug
	case "info", "INFO", "Info":
		return slog.LevelInfo
	case "warn", "WARN", "Warn":
		return slog.LevelWarn
	case "error", "ERROR", "Error":
		return slog.LevelError
	default:
		// Log warning about invalid level (using default logger before config is ready)
		// For now, just return info as default
		return slog.LevelInfo
	}
}
