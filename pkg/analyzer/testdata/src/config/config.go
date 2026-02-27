// Package config tests that only the enabled rules produce diagnostics.
// In this test only the lowercase check is enabled, so English/special/sensitive
// violations must NOT produce diagnostics.
package config

import "log/slog"

func examples() {
	// Lowercase violations — SHOULD be reported (check enabled)
	slog.Info("Starting server")  // want `log message must start with a lowercase letter`
	slog.Error("Failed to connect") // want `log message must start with a lowercase letter`
	slog.Warn("Warning here")     // want `log message must start with a lowercase letter`

	// English violation — should NOT be reported (check disabled)
	slog.Info("запуск сервера")

	// Special chars violation — should NOT be reported (check disabled)
	slog.Info("server started!")
	slog.Info("client & server")

	// Sensitive data violation — should NOT be reported (check disabled)
	slog.Info("auth header set")
	slog.Info("token: xyz")

	// Correct
	slog.Info("starting server")
	slog.Error("failed to connect")
}
