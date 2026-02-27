// Package zap is a stub for go.uber.org/zap used in analysis tests.
package zap

// Field is a marshaling operation used to add a key-value pair to a logger's context.
type Field struct{}

// Logger is a fast, leveled, structured logger.
type Logger struct{}

// NewNop returns a no-op Logger.
func NewNop() *Logger { return &Logger{} }

// Sugar wraps the Logger to provide a sugared logger.
func (l *Logger) Sugar() *SugaredLogger { return &SugaredLogger{} }

// Info logs a message at Info level.
func (l *Logger) Info(msg string, fields ...Field) {}

// Warn logs a message at Warn level.
func (l *Logger) Warn(msg string, fields ...Field) {}

// Error logs a message at Error level.
func (l *Logger) Error(msg string, fields ...Field) {}

// Debug logs a message at Debug level.
func (l *Logger) Debug(msg string, fields ...Field) {}

// Fatal logs a message at Fatal level, then calls os.Exit(1).
func (l *Logger) Fatal(msg string, fields ...Field) {}

// Panic logs a message at Panic level, then panics.
func (l *Logger) Panic(msg string, fields ...Field) {}

// DPanic logs a message at DPanic level.
func (l *Logger) DPanic(msg string, fields ...Field) {}

// Log logs a message at the given level.
func (l *Logger) Log(lvl int, msg string, fields ...Field) {}

// SugaredLogger wraps the base Logger functionality in a slower, but less verbose, API.
type SugaredLogger struct{}

// Infow logs a message with some additional context.
func (s *SugaredLogger) Infow(msg string, keysAndValues ...interface{}) {}

// Warnw logs a message with some additional context.
func (s *SugaredLogger) Warnw(msg string, keysAndValues ...interface{}) {}

// Errorw logs a message with some additional context.
func (s *SugaredLogger) Errorw(msg string, keysAndValues ...interface{}) {}

// Debugw logs a message with some additional context.
func (s *SugaredLogger) Debugw(msg string, keysAndValues ...interface{}) {}

// Fatalw logs a message with some additional context, then calls os.Exit(1).
func (s *SugaredLogger) Fatalw(msg string, keysAndValues ...interface{}) {}

// Panicw logs a message with some additional context, then panics.
func (s *SugaredLogger) Panicw(msg string, keysAndValues ...interface{}) {}

// DPanicw logs a message with some additional context.
func (s *SugaredLogger) DPanicw(msg string, keysAndValues ...interface{}) {}

// Infof uses fmt.Sprintf to log a templated message.
func (s *SugaredLogger) Infof(template string, args ...interface{}) {}

// Warnf uses fmt.Sprintf to log a templated message.
func (s *SugaredLogger) Warnf(template string, args ...interface{}) {}

// Errorf uses fmt.Sprintf to log a templated message.
func (s *SugaredLogger) Errorf(template string, args ...interface{}) {}

// Debugf uses fmt.Sprintf to log a templated message.
func (s *SugaredLogger) Debugf(template string, args ...interface{}) {}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit(1).
func (s *SugaredLogger) Fatalf(template string, args ...interface{}) {}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func (s *SugaredLogger) Panicf(template string, args ...interface{}) {}

// DPanicf uses fmt.Sprintf to log a templated message.
func (s *SugaredLogger) DPanicf(template string, args ...interface{}) {}

// Infoln uses fmt.Sprintln to log a message.
func (s *SugaredLogger) Infoln(args ...interface{}) {}

// Warnln uses fmt.Sprintln to log a message.
func (s *SugaredLogger) Warnln(args ...interface{}) {}

// Errorln uses fmt.Sprintln to log a message.
func (s *SugaredLogger) Errorln(args ...interface{}) {}

// Debugln uses fmt.Sprintln to log a message.
func (s *SugaredLogger) Debugln(args ...interface{}) {}

// Fatalln uses fmt.Sprintln to log a message, then calls os.Exit(1).
func (s *SugaredLogger) Fatalln(args ...interface{}) {}

// Panicln uses fmt.Sprintln to log a message, then panics.
func (s *SugaredLogger) Panicln(args ...interface{}) {}
