package logger

import (
	"context"
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// global глобальный экземпляр логгера.
	global       *zap.SugaredLogger
	defaultLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
)

func init() {
	SetLogger(New(defaultLevel,
		zap.AddStacktrace(zap.FatalLevel),
	))
}

func New(level zapcore.LevelEnabler, options ...zap.Option) *zap.SugaredLogger {
	return NewWithSink(level, os.Stdout, options...)
}

func NewWithSink(level zapcore.LevelEnabler, sink io.Writer, options ...zap.Option) *zap.SugaredLogger {
	if level == nil {
		level = defaultLevel
	}

	core := newZapCore(level, sink)

	return zap.New(core, options...).Sugar()
}

func newZapCore(level zapcore.LevelEnabler, sink io.Writer) zapcore.Core {
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}),
		zapcore.AddSync(sink),
		level,
	)
}

// Level возвращает текущий уровень логгирования глобального логгера.
func Level() zapcore.Level {
	return defaultLevel.Level()
}

// SetLevel устанавливает уровень логгирования глобального логгера.
func SetLevel(l zapcore.Level) {
	defaultLevel.SetLevel(l)
}

// Logger возвращает глобальный логгер.
func Logger() *zap.SugaredLogger {
	return global
}

// SetLogger устанавливает глобальный логгер. Функция непотокобезопасна.
func SetLogger(l *zap.SugaredLogger) {
	global = l
}

func Debug(ctx context.Context, args ...interface{}) {
	if logger := FromContext(ctx); logger.Level().Enabled(zapcore.DebugLevel) {
		logger.Debug(args...)
	}
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	if logger := FromContext(ctx); logger.Level().Enabled(zapcore.DebugLevel) {
		logger.Debugf(format, args...)
	}
}

func DebugKV(ctx context.Context, message string, kvs ...interface{}) {
	if logger := FromContext(ctx); logger.Level().Enabled(zapcore.DebugLevel) {
		logger.Debugw(message, kvs...)
	}
}

func Info(ctx context.Context, args ...interface{}) {
	if logger := FromContext(ctx); logger.Level().Enabled(zapcore.InfoLevel) {
		logger.Info(args...)
	}
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	if logger := FromContext(ctx); logger.Level().Enabled(zapcore.InfoLevel) {
		logger.Infof(format, args...)
	}
}

func InfoKV(ctx context.Context, message string, kvs ...interface{}) {
	if logger := FromContext(ctx); logger.Level().Enabled(zapcore.InfoLevel) {
		logger.Infow(message, kvs...)
	}
}

func Warn(ctx context.Context, args ...interface{}) {
	if logger := FromContext(ctx); logger.Level().Enabled(zapcore.WarnLevel) {
		logger.Warn(args...)
	}
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	if logger := FromContext(ctx); logger.Level().Enabled(zapcore.WarnLevel) {
		logger.Warnf(format, args...)
	}
}

func WarnKV(ctx context.Context, message string, kvs ...interface{}) {
	if logger := FromContext(ctx); logger.Level().Enabled(zapcore.WarnLevel) {
		logger.Warnw(message, kvs...)
	}
}

func Error(ctx context.Context, args ...interface{}) {
	if logger := FromContext(ctx); logger.Level().Enabled(zapcore.ErrorLevel) {
		logger.Error(args...)
	}
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	if logger := FromContext(ctx); logger.Level().Enabled(zapcore.ErrorLevel) {
		logger.Errorf(format, args...)
	}
}

func ErrorKV(ctx context.Context, message string, kvs ...interface{}) {
	if logger := FromContext(ctx); logger.Level().Enabled(zapcore.ErrorLevel) {
		logger.Errorw(message, kvs...)
	}
}

func Fatal(ctx context.Context, args ...interface{}) {
	FromContext(ctx).Fatal(args...)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).Fatalf(format, args...)
}

func FatalKV(ctx context.Context, message string, kvs ...interface{}) {
	FromContext(ctx).Fatalw(message, kvs...)
}

func Panic(ctx context.Context, args ...interface{}) {
	FromContext(ctx).Panic(args...)
}

func Panicf(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).Panicf(format, args...)
}

func PanicKV(ctx context.Context, message string, kvs ...interface{}) {
	FromContext(ctx).Panicw(message, kvs...)
}

func Audit(ctx context.Context, message string, kvs ...interface{}) {
	FromContext(ctx).Errorw(message, kvs...)
}
