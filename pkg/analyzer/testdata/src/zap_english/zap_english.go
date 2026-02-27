package zap_english

import "go.uber.org/zap"

func examples() {
	// Rule 2: log messages must be in English only

	logger := zap.NewNop()
	sugar := logger.Sugar()

	// Logger methods with Cyrillic
	logger.Info("запуск сервера")                    // want `log message must be in English only`
	logger.Error("ошибка подключения к базе данных") // want `log message must be in English only`
	logger.Warn("предупреждение системы")            // want `log message must be in English only`
	logger.Debug("отладочная информация")            // want `log message must be in English only`

	// Chinese
	logger.Info("启动服务器") // want `log message must be in English only`

	// Japanese
	logger.Info("サーバー起動") // want `log message must be in English only`

	// Mixed
	logger.Info("server запуск failed") // want `log message must be in English only`

	// SugaredLogger w-methods
	sugar.Infow("запуск сервера", "port", 8080) // want `log message must be in English only`
	sugar.Errorw("ошибка", "key", "val")        // want `log message must be in English only`
	sugar.Warnw("предупреждение")               // want `log message must be in English only`

	// SugaredLogger f-methods
	sugar.Infof("запуск %d", 8080)   // want `log message must be in English only`
	sugar.Errorf("ошибка: %s", "db") // want `log message must be in English only`

	// Correct - English only
	logger.Info("starting server")
	logger.Error("failed to connect to database")
	logger.Warn("something went wrong")
	logger.Debug("debug information")
	sugar.Infow("starting server", "port", 8080)
	sugar.Infof("starting on port %d", 8080)

	// Numbers and punctuation are ok
	logger.Info("processing 42 items")
	logger.Info("server at 127.0.0.1")
}
