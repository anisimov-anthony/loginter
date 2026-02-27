// Package clean contains log calls that comply with all four rules.
// The e2e test verifies that the linter reports zero diagnostics here.
package clean

import (
	"context"
	"log/slog"
)

func allClean() {
	ctx := context.Background()

	// All lowercase, English, no special chars, no sensitive keywords
	slog.Info("server started on port 8080")
	slog.Error("failed to connect to database")
	slog.Warn("disk space is low")
	slog.Debug("debug mode enabled")
	slog.InfoContext(ctx, "starting with context")

	// Numbers and punctuation at start are fine
	slog.Info("8080 is the port")
	slog.Info("3 retries remaining")

	// Allowed punctuation
	slog.Info("host: localhost")
	slog.Info("path/to/file")
	slog.Info("key=value pair")
	slog.Info("well-known endpoint")
	slog.Info("value (default 42)")
	slog.Info("processing items in batch", "count", 42)

	// Word boundaries — these must NOT trigger sensitive rule
	slog.Info("user authenticated successfully")
	slog.Info("authorization granted")
	slog.Info("api request completed")
	slog.Info("payment processed")
}
