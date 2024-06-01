package logger

import (
	"context"

	"go.uber.org/zap"
)

type contextKey int

const (
	loggerContextKey contextKey = iota
)

// ToContext создает контекст с переданным логгером внутри.
func ToContext(ctx context.Context, l *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, loggerContextKey, l)
}

// FromContext достает логгер из контекста. Если в контексте логгер не
// обнаруживается - возвращает глобальный логгер. В обоих случаях логгер уже
// содержит аннотации в виде trace_id и span_id
func FromContext(ctx context.Context) *zap.SugaredLogger {
	l := getLogger(ctx)

	return l
}

// WithFields создает логгер из уже имеющегося в контексте и устанавливает метаданные,
// используя типизированные поля.
func WithFields(ctx context.Context, fields ...zap.Field) context.Context {
	log := FromContext(ctx).
		Desugar().
		With(fields...).
		Sugar()
	return ToContext(ctx, log)
}

func getLogger(ctx context.Context) *zap.SugaredLogger {
	if logger, ok := ctx.Value(loggerContextKey).(*zap.SugaredLogger); ok {
		return logger
	}

	return global
}
