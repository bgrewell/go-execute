package execute

// Logger defines the interface for logging within the execute package.
// The package is designed to not output any logging by default.
// Users can integrate their preferred logging system by implementing this interface.
type Logger interface {
	Trace(msg string, fields ...any)
	Debug(msg string, fields ...any)
	Info(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Error(msg string, fields ...any)
}

// logger is the package-level logger instance
var logger Logger = &NoOpLogger{}

// SetLogger allows users to integrate their logging system with the execute package
func SetLogger(l Logger) {
	logger = l
}

// GetLogger returns the current logger instance
func GetLogger() Logger {
	return logger
}

// NoOpLogger is the default logger that discards all log messages
type NoOpLogger struct{}

func (n *NoOpLogger) Trace(msg string, fields ...any) {}
func (n *NoOpLogger) Debug(msg string, fields ...any) {}
func (n *NoOpLogger) Info(msg string, fields ...any)  {}
func (n *NoOpLogger) Warn(msg string, fields ...any)  {}
func (n *NoOpLogger) Error(msg string, fields ...any) {}
