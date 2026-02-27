package slog_sensitive

import (
	"context"
	"log/slog"
)

func examples() {
	// Rule 4: log messages must not contain sensitive data

	password := "secret123"
	apiKey := "abc123"
	token := "xyz"
	secret := "mysecret"

	// String literals with sensitive keywords
	slog.Info("user password: secret123")       // want `log message may contain sensitive data \("password"\)`
	slog.Debug("api_key=abc123")                // want `log message may contain sensitive data \("api_key"\)`
	slog.Info("token: xyz")                     // want `log message may contain sensitive data \("token"\)`
	slog.Info("secret value here")              // want `log message may contain sensitive data \("secret"\)`
	slog.Info("private_key loaded")             // want `log message may contain sensitive data \("private_key"\)`
	slog.Info("access_key set")                 // want `log message may contain sensitive data \("access_key"\)`
	slog.Info("api_secret configured")          // want `log message may contain sensitive data \("api_secret"\)`
	slog.Info("user credential found")          // want `log message may contain sensitive data \("credential"\)`
	slog.Info("auth header set")                // want `log message may contain sensitive data \("auth"\)`
	slog.Info("passwd changed")                 // want `log message may contain sensitive data \("passwd"\)`
	slog.Info("apikey sent to server")          // want `log message may contain sensitive data \("apikey"\)`

	// Case insensitive
	slog.Info("PASSWORD: abc")            // want `log message may contain sensitive data \("password"\)`
	slog.Info("Token is valid")           // want `log message may contain sensitive data \("token"\)`
	slog.Info("SECRET stored")            // want `log message may contain sensitive data \("secret"\)`

	// String concatenation
	slog.Info("user password: " + password) // want `log message may contain sensitive data \("password"\)`
	slog.Debug("api_key=" + apiKey)         // want `log message may contain sensitive data \("api_key"\)`
	slog.Info("token: " + token)            // want `log message may contain sensitive data \("token"\)`
	slog.Info("secret: " + secret)          // want `log message may contain sensitive data \("secret"\)`

	// Context variants
	ctx := context.Background()
	slog.InfoContext(ctx, "auth header set") // want `log message may contain sensitive data \("auth"\)`

	// Logger instance
	logger := slog.Default()
	logger.Info("password reset")   // want `log message may contain sensitive data \("password"\)`
	logger.Debug("token refreshed") // want `log message may contain sensitive data \("token"\)`

	// Correct - word boundary checks: "auth" in "authentication" should NOT match
	slog.Info("user authenticated successfully")
	slog.Debug("api request completed")
	slog.Info("authorization granted")
	slog.Info("validation complete")
	slog.Info("payment processed")
	slog.Info("user registered")
}
