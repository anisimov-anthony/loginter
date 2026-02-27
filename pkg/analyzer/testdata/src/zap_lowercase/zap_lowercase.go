package zap_lowercase

import "go.uber.org/zap"

func examples() {
	// Rule 1: log messages must start with lowercase letter

	logger := zap.NewNop()
	sugar := logger.Sugar()

	// Logger methods
	logger.Info("Starting server on port 8080")  // want `log message must start with a lowercase letter`
	logger.Error("Failed to connect to database") // want `log message must start with a lowercase letter`
	logger.Warn("Warning about something")        // want `log message must start with a lowercase letter`
	logger.Debug("Debug mode enabled")            // want `log message must start with a lowercase letter`
	logger.Fatal("Fatal system error")            // want `log message must start with a lowercase letter`
	logger.Panic("Panic occurred")                // want `log message must start with a lowercase letter`
	logger.DPanic("DPanic in development")        // want `log message must start with a lowercase letter`

	// SugaredLogger w-methods (structured)
	sugar.Infow("Starting server", "port", 8080)  // want `log message must start with a lowercase letter`
	sugar.Errorw("Failed to connect", "err", "x") // want `log message must start with a lowercase letter`
	sugar.Warnw("Warning here", "key", "val")     // want `log message must start with a lowercase letter`
	sugar.Debugw("Debug info", "key", "val")      // want `log message must start with a lowercase letter`
	sugar.Fatalw("Fatal error", "key", "val")     // want `log message must start with a lowercase letter`

	// SugaredLogger f-methods (printf-style)
	sugar.Infof("Starting on port %d", 8080)  // want `log message must start with a lowercase letter`
	sugar.Errorf("Failed: %s", "err")         // want `log message must start with a lowercase letter`
	sugar.Warnf("Warning: %s", "msg")         // want `log message must start with a lowercase letter`
	sugar.Debugf("Debug: %v", "val")          // want `log message must start with a lowercase letter`

	// All caps acronym
	logger.Info("HTTP server started") // want `log message must start with a lowercase letter`

	// Correct - starts with lowercase
	logger.Info("starting server on port 8080")
	logger.Error("failed to connect to database")
	logger.Warn("something went wrong")
	logger.Debug("debug mode enabled")
	sugar.Infow("starting server", "port", 8080)
	sugar.Infof("starting on port %d", 8080)

	// Numbers and punctuation at start are ok
	logger.Info("8080 is the port")
	logger.Info("3 retries remaining")
	logger.Info("")
}
