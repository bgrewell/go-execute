package execute

var (
	logger Logger = &NoOpLogger{}
)

// Logger provides the logging interface for integrating your logging with the go-execute module
type Logger interface {
	Trace(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
}

func SetLogger(custom Logger) {
	logger = custom
}

// NoOpLogger is the default logging implementation which doesn't do anything with any of the logging messages
// since the module is designed to not output any logging by default
type NoOpLogger struct{}

func (l *NoOpLogger) Trace(msg string, fields ...interface{}) {}
func (l *NoOpLogger) Debug(msg string, fields ...interface{}) {}
func (l *NoOpLogger) Info(msg string, fields ...interface{})  {}
func (l *NoOpLogger) Warn(msg string, fields ...interface{})  {}
func (l *NoOpLogger) Error(msg string, fields ...interface{}) {}
func (l *NoOpLogger) Fatal(msg string, fields ...interface{}) {}
