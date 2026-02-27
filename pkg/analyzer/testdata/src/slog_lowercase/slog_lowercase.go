package slog_lowercase

import (
	"context"
	"log/slog"
)

func examples() {
	// Rule 1: log messages must start with lowercase letter

	// Package-level functions
	slog.Info("Starting server on port 8080")  // want `log message must start with a lowercase letter`
	slog.Error("Failed to connect to database") // want `log message must start with a lowercase letter`
	slog.Warn("Warning about something")        // want `log message must start with a lowercase letter`
	slog.Debug("Debug mode enabled")            // want `log message must start with a lowercase letter`

	// Context variants
	ctx := context.Background()
	slog.InfoContext(ctx, "Starting with context")  // want `log message must start with a lowercase letter`
	slog.ErrorContext(ctx, "Failed with context")   // want `log message must start with a lowercase letter`
	slog.WarnContext(ctx, "Warning with context")   // want `log message must start with a lowercase letter`
	slog.DebugContext(ctx, "Debug with context")    // want `log message must start with a lowercase letter`

	// Logger instance methods
	logger := slog.Default()
	logger.Info("Starting via logger instance")  // want `log message must start with a lowercase letter`
	logger.Error("Failed via logger instance")   // want `log message must start with a lowercase letter`
	logger.Warn("Warning via logger instance")   // want `log message must start with a lowercase letter`
	logger.Debug("Debug via logger instance")    // want `log message must start with a lowercase letter`
	logger.InfoContext(ctx, "Starting ctx via logger instance") // want `log message must start with a lowercase letter`

	// Single uppercase letter
	slog.Info("A")     // want `log message must start with a lowercase letter`
	// All caps acronym
	slog.Info("HTTP server started") // want `log message must start with a lowercase letter`

	// Correct - starts with lowercase
	slog.Info("starting server on port 8080")
	slog.Error("failed to connect to database")
	slog.Warn("something went wrong")
	slog.Debug("debug mode enabled")
	slog.InfoContext(ctx, "starting with context")
	logger.Info("starting via logger instance")

	// Numbers and punctuation at start are ok (not uppercase letter)
	slog.Info("8080 is the port")
	slog.Info("3 retries remaining")
	slog.Info("")
}
