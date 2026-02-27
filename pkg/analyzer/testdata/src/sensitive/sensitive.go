// Package sensitive focuses on rule 4 (sensitive data) with all checks enabled.
package sensitive

import (
	"context"
	"log/slog"
)

func examples() {
	password := "secret123"
	apiKey := "abc123"
	token := "xyz"

	ctx := context.Background()
	logger := slog.Default()

	// All default sensitive patterns
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

	// Case insensitive — these also trigger lowercase rule (uppercase first letter)
	slog.Info("PASSWORD: abc")   // want `log message must start with a lowercase letter` `log message may contain sensitive data \("password"\)`
	slog.Info("Token is valid")  // want `log message must start with a lowercase letter` `log message may contain sensitive data \("token"\)`
	slog.Info("SECRET stored")   // want `log message must start with a lowercase letter` `log message may contain sensitive data \("secret"\)`

	// String concatenation
	slog.Info("user password: " + password) // want `log message may contain sensitive data \("password"\)`
	slog.Debug("api_key=" + apiKey)         // want `log message may contain sensitive data \("api_key"\)`
	slog.Info("token: " + token)            // want `log message may contain sensitive data \("token"\)`

	// Context variants
	slog.InfoContext(ctx, "auth header set") // want `log message may contain sensitive data \("auth"\)`

	// Logger instance
	logger.Info("password reset")   // want `log message may contain sensitive data \("password"\)`
	logger.Debug("token refreshed") // want `log message may contain sensitive data \("token"\)`

	// Word boundary checks — these must NOT match
	slog.Info("user authenticated successfully")
	slog.Debug("api request completed")
	slog.Info("authorization granted")
	slog.Info("validation complete")
	slog.Info("payment processed")
	slog.Info("user registered")
}
