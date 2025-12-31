package pkg

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func NewZapLogger() *ZapLogger {

	// Create an AtomicLevel to manage log level dynamically
	atom := zap.NewAtomicLevel()

	// Set the initial log level
	atom.SetLevel(zap.InfoLevel)

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.Lock(os.Stdout),
		atom,
	))
	return &ZapLogger{
		zapLogger: logger,
		atom:      atom,
	}
}

type ZapLogger struct {
	zapLogger *zap.Logger
	atom      zap.AtomicLevel
}

func (l *ZapLogger) SetLevel(level zapcore.Level) {

	// Change the log level dynamically
	l.atom.SetLevel(level)
}

func (l *ZapLogger) Trace(msg string, fields ...interface{}) {
	l.zapLogger.Sugar().Debugw(msg, fields...) // Zap does not have a Trace level, using Debug instead
}

func (l *ZapLogger) Debug(msg string, fields ...interface{}) {
	l.zapLogger.Sugar().Debugw(msg, fields...)
}

func (l *ZapLogger) Info(msg string, fields ...interface{}) {
	l.zapLogger.Sugar().Infow(msg, fields...)
}

func (l *ZapLogger) Warn(msg string, fields ...interface{}) {
	l.zapLogger.Sugar().Warnw(msg, fields...)
}

func (l *ZapLogger) Error(msg string, fields ...interface{}) {
	l.zapLogger.Sugar().Errorw(msg, fields...)
}

func (l *ZapLogger) Fatal(msg string, fields ...interface{}) {
	l.zapLogger.Sugar().Fatalw(msg, fields...)
}
