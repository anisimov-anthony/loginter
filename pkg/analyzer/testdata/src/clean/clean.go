// Package clean contains log messages that comply with all rules.
// No diagnostics should be reported here.
package clean

import (
	"context"
	"log/slog"
)

func examples() {
	ctx := context.Background()
	logger := slog.Default()

	// All messages: lowercase, English, no special chars, no sensitive data

	slog.Info("server started on port 8080")
	slog.Error("failed to connect to database")
	slog.Warn("something went wrong")
	slog.Debug("debug mode enabled")

	slog.InfoContext(ctx, "starting with context")
	slog.ErrorContext(ctx, "failed with context")
	slog.WarnContext(ctx, "warning with context")
	slog.DebugContext(ctx, "debug with context")

	logger.Info("logger started")
	logger.Error("logger error occurred")
	logger.Warn("logger warning")
	logger.Debug("logger debug info")

	slog.Info("processing 42 items")
	slog.Info("server at 127.0.0.1:8080")
	slog.Info("retry 3 of 5")
	slog.Info("key=value pair")
	slog.Info("well-known host")
	slog.Info("path/to/file loaded")
	slog.Info("value (default 42)")
	slog.Info("host: localhost")
	slog.Info("user authenticated successfully")
	slog.Info("authorization granted")
	slog.Info("api request completed")
	slog.Info("payment processed")
	slog.Info("user registered successfully")
	slog.Info("")
}
