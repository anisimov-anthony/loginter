// Package basic tests all four rules applied to slog package-level functions
// and *slog.Logger methods. Used for both Run and RunWithSuggestedFixes.
package basic

import (
	"context"
	"log/slog"
)

func examples() {
	ctx := context.Background()

	// Rule 1: lowercase
	slog.Info("Starting server on port 8080")    // want `log message must start with a lowercase letter`
	slog.Error("Failed to connect to database")  // want `log message must start with a lowercase letter`
	slog.InfoContext(ctx, "Starting with context") // want `log message must start with a lowercase letter`

	// Rule 2: English only
	slog.Info("запуск сервера")      // want `log message must be in English only`
	slog.Warn("ошибка подключения")  // want `log message must be in English only`

	// Rule 3: no special chars
	slog.Info("server started!")     // want `log message must not contain special characters or emoji`
	slog.Error("failed 😱")         // want `log message must not contain special characters or emoji`
	slog.Info("client & server")     // want `log message must not contain special characters or emoji`

	// Rule 4: sensitive data
	slog.Info("auth header set")     // want `log message may contain sensitive data \("auth"\)`
	slog.Info("token: xyz")          // want `log message may contain sensitive data \("token"\)`
	slog.Info("user password: abc")  // want `log message may contain sensitive data \("password"\)`

	// Correct — no diagnostics expected
	slog.Info("starting server on port 8080")
	slog.Error("failed to connect to database")
	slog.Warn("something went wrong")
	slog.Debug("debug info")
	slog.InfoContext(ctx, "starting with context")
	slog.Info("user authenticated successfully")
	slog.Info("processing 42 items")
	slog.Info("key=value pair")
	slog.Info("well-known host")
}
