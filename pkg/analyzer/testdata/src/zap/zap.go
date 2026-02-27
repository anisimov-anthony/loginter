// Package zap tests all four rules applied to go.uber.org/zap Logger and SugaredLogger.
package zap

import "go.uber.org/zap"

func examples() {
	logger := zap.NewNop()
	sugar := logger.Sugar()

	// Rule 1: lowercase — Logger methods
	logger.Info("Starting server on port 8080")   // want `log message must start with a lowercase letter`
	logger.Error("Failed to connect to database")  // want `log message must start with a lowercase letter`
	logger.Warn("Warning about something")         // want `log message must start with a lowercase letter`
	logger.Debug("Debug mode enabled")             // want `log message must start with a lowercase letter`
	logger.Fatal("Fatal error occurred")           // want `log message must start with a lowercase letter`

	// Rule 1: lowercase — SugaredLogger methods
	sugar.Infow("Starting server", "port", 8080)  // want `log message must start with a lowercase letter`
	sugar.Infof("Starting on port %d", 8080)      // want `log message must start with a lowercase letter`
	sugar.Errorw("Failed to connect", "err", "x") // want `log message must start with a lowercase letter`

	// Rule 2: English only
	logger.Info("запуск сервера")         // want `log message must be in English only`
	logger.Error("ошибка подключения")    // want `log message must be in English only`
	sugar.Infow("запуск", "port", 8080)   // want `log message must be in English only`
	sugar.Infof("запуск %d", 8080)        // want `log message must be in English only`

	// Rule 3: no special chars
	logger.Info("server started!")        // want `log message must not contain special characters or emoji`
	logger.Error("failed 😱")            // want `log message must not contain special characters or emoji`
	sugar.Infow("started!", "port", 8080) // want `log message must not contain special characters or emoji`
	sugar.Infof("started! port %d", 8080) // want `log message must not contain special characters or emoji`

	// Rule 4: sensitive data
	logger.Info("auth header set")                // want `log message may contain sensitive data \("auth"\)`
	logger.Info("token: xyz")                     // want `log message may contain sensitive data \("token"\)`
	logger.Info("user password: abc")             // want `log message may contain sensitive data \("password"\)`
	sugar.Infow("password changed", "user", "bob") // want `log message may contain sensitive data \("password"\)`
	sugar.Debugf("token: %s", "val")               // want `log message may contain sensitive data \("token"\)`

	// Correct — no diagnostics
	logger.Info("starting server on port 8080")
	logger.Error("failed to connect to database")
	logger.Warn("something went wrong")
	logger.Debug("debug mode enabled")
	sugar.Infow("starting server", "port", 8080)
	sugar.Infof("starting on port %d", 8080)
	logger.Info("user authenticated successfully")
	logger.Info("processing 42 items")
}
