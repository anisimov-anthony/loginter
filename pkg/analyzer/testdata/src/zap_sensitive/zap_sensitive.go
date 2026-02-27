package zap_sensitive

import "go.uber.org/zap"

func examples() {
	// Rule 4: log messages must not contain sensitive data

	logger := zap.NewNop()
	sugar := logger.Sugar()

	password := "secret123"
	apiKey := "abc123"
	token := "xyz"

	// String literals with sensitive keywords - Logger methods
	logger.Info("user password: secret123")       // want `log message may contain sensitive data \("password"\)`
	logger.Debug("api_key=abc123")                // want `log message may contain sensitive data \("api_key"\)`
	logger.Info("token: xyz")                     // want `log message may contain sensitive data \("token"\)`
	logger.Info("secret value here")              // want `log message may contain sensitive data \("secret"\)`
	logger.Info("private_key loaded")             // want `log message may contain sensitive data \("private_key"\)`
	logger.Info("access_key set")                 // want `log message may contain sensitive data \("access_key"\)`
	logger.Info("api_secret configured")          // want `log message may contain sensitive data \("api_secret"\)`
	logger.Info("user credential found")          // want `log message may contain sensitive data \("credential"\)`
	logger.Info("auth header set")                // want `log message may contain sensitive data \("auth"\)`
	logger.Info("passwd changed")                 // want `log message may contain sensitive data \("passwd"\)`
	logger.Info("apikey sent to server")          // want `log message may contain sensitive data \("apikey"\)`

	// Case insensitive
	logger.Info("PASSWORD: abc")   // want `log message may contain sensitive data \("password"\)`
	logger.Info("Token is valid")  // want `log message may contain sensitive data \("token"\)`

	// String concatenation - Logger
	logger.Info("user password: " + password)  // want `log message may contain sensitive data \("password"\)`
	logger.Debug("api_key=" + apiKey)          // want `log message may contain sensitive data \("api_key"\)`
	logger.Info("token: " + token)             // want `log message may contain sensitive data \("token"\)`

	// SugaredLogger w-methods
	sugar.Infow("user password: secret", "user", "bob")    // want `log message may contain sensitive data \("password"\)`
	sugar.Errorw("token expired", "user", "alice")         // want `log message may contain sensitive data \("token"\)`
	sugar.Debugw("api_key loaded", "key", "k1")            // want `log message may contain sensitive data \("api_key"\)`

	// SugaredLogger f-methods
	sugar.Infof("password is %s", password)   // want `log message may contain sensitive data \("password"\)`
	sugar.Debugf("api_key: %s", apiKey)       // want `log message may contain sensitive data \("api_key"\)`

	// Correct - word boundary checks
	logger.Info("user authenticated successfully")
	logger.Debug("api request completed")
	logger.Info("authorization granted")
	logger.Info("validation complete")
	sugar.Infow("user authenticated", "id", 42)
	sugar.Debugw("api request done", "status", 200)
}
