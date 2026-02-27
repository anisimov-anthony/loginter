// Package violations contains log calls that intentionally violate all four rules.
// The e2e test verifies that the linter reports each of them.
package violations

import (
	"context"
	"log/slog"
)

func allRules() {
	ctx := context.Background()

	// Rule 1: message must start with lowercase letter
	slog.Info("Starting server on port 8080")
	slog.Error("Failed to connect to database")
	slog.Warn("Warning: disk space low")
	slog.Debug("Debug mode enabled")
	slog.InfoContext(ctx, "Starting with context")

	// Rule 2: message must be in English only
	slog.Info("запуск сервера")
	slog.Error("ошибка подключения")
	slog.Warn("启动服务器")

	// Rule 3: message must not contain special characters or emoji
	slog.Info("server started!")
	slog.Error("connection failed!!!")
	slog.Warn("notify @admin now")
	slog.Info("issue #123 found")
	slog.Debug("step 1; step 2")
	slog.Info("run `command`")
	slog.Info("client & server")
	slog.Info("server started 🚀")
	slog.Warn("alert ⚠")
	slog.Info("loading\u2026")

	// Rule 4: message must not contain sensitive data keywords
	slog.Info("user password: secret123")
	slog.Debug("api_key=abc123")
	slog.Info("token: xyz")
	slog.Info("secret value here")
	slog.Info("private_key loaded")
	slog.Info("access_key set")
	slog.Info("api_secret configured")
	slog.Info("user credential found")
	slog.Info("auth header set")
	slog.Info("passwd changed")
	slog.Info("apikey sent")
}
