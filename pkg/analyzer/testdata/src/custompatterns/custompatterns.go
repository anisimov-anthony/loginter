// Package custompatterns tests that custom sensitive patterns work correctly.
// Only sensitive check is enabled with custom patterns: "ssn", "credit_card".
// Default patterns (password, token, etc.) also apply via AllSensitivePatterns.
package custompatterns

import "log/slog"

func examples() {
	// Custom pattern: ssn
	slog.Info("user ssn: 123-45-6789") // want `log message may contain sensitive data \("ssn"\)`

	// Custom pattern: credit_card
	slog.Info("credit_card: 4111111111111111") // want `log message may contain sensitive data \("credit_card"\)`

	// Default pattern still works
	slog.Info("user password: secret123") // want `log message may contain sensitive data \("password"\)`
	slog.Info("auth header set")          // want `log message may contain sensitive data \("auth"\)`

	// Non-sensitive — no diagnostics
	slog.Info("payment processed")
	slog.Info("user registered")
	slog.Info("starting server")

	// Uppercase message — but only sensitive check enabled, not lowercase
	slog.Info("Starting server")

	// Special chars — but only sensitive check enabled
	slog.Info("server started!")
}
