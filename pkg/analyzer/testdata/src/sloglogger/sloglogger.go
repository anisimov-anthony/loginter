// Package sloglogger tests rules applied to *slog.Logger instance method calls.
package sloglogger

import (
	"context"
	"log/slog"
)

func examples() {
	logger := slog.Default()
	ctx := context.Background()

	// Rule 1: lowercase
	logger.Info("Starting server")  // want `log message must start with a lowercase letter`
	logger.Error("Failed to connect") // want `log message must start with a lowercase letter`
	logger.Warn("Warning here")     // want `log message must start with a lowercase letter`
	logger.Debug("Debug info")      // want `log message must start with a lowercase letter`

	// Rule 1 with context variants
	logger.InfoContext(ctx, "Starting with context")  // want `log message must start with a lowercase letter`
	logger.ErrorContext(ctx, "Failed with context")   // want `log message must start with a lowercase letter`
	logger.WarnContext(ctx, "Warning with context")   // want `log message must start with a lowercase letter`
	logger.DebugContext(ctx, "Debug with context")    // want `log message must start with a lowercase letter`

	// Rule 2: English only
	logger.Info("запуск сервера")  // want `log message must be in English only`
	logger.Error("ошибка системы") // want `log message must be in English only`

	// Rule 3: no special chars
	logger.Info("server started!") // want `log message must not contain special characters or emoji`
	logger.Warn("alert ⚠")        // want `log message must not contain special characters or emoji`

	// Rule 4: no sensitive data
	logger.Info("auth header set")   // want `log message may contain sensitive data \("auth"\)`
	logger.Debug("token refreshed")  // want `log message may contain sensitive data \("token"\)`

	// Correct
	logger.Info("starting server")
	logger.Error("failed to connect")
	logger.Warn("something went wrong")
	logger.Debug("debug mode enabled")
	logger.InfoContext(ctx, "starting with context")
}
