package slog_english

import (
	"context"
	"log/slog"
)

func examples() {
	// Rule 2: log messages must be in English only

	// Cyrillic
	slog.Info("запуск сервера")                     // want `log message must be in English only`
	slog.Error("ошибка подключения к базе данных")  // want `log message must be in English only`
	slog.Warn("предупреждение системы")             // want `log message must be in English only`
	slog.Debug("отладочная информация")             // want `log message must be in English only`

	// Chinese
	slog.Info("启动服务器")   // want `log message must be in English only`
	slog.Error("连接失败")   // want `log message must be in English only`

	// Japanese
	slog.Info("サーバー起動") // want `log message must be in English only`

	// Mixed: English + non-English
	slog.Info("server запуск")         // want `log message must be in English only`
	slog.Error("failed подключение")   // want `log message must be in English only`

	// Accented characters
	slog.Info("café started")  // want `log message must be in English only`

	// Context variants
	ctx := context.Background()
	slog.InfoContext(ctx, "запуск с контекстом") // want `log message must be in English only`

	// Logger instance
	logger := slog.Default()
	logger.Info("ошибка системы") // want `log message must be in English only`

	// Correct - English only
	slog.Info("starting server")
	slog.Error("failed to connect to database")
	slog.Warn("something went wrong")
	slog.Debug("debug information")
	slog.Info("server started successfully")

	// Numbers and punctuation are ok
	slog.Info("processing 42 items")
	slog.Info("server at 127.0.0.1")
	slog.Info("retry 3 of 5")
}
