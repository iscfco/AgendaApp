package logs

// TODO: Implement log rotation

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerCtxKey struct{}

func InitLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := config.Build()
	if err != nil {
		panic("Error at initializing logger: " + err.Error())
	}
	return logger
}

func Logger(ctx context.Context) *zap.Logger {
	return GetLoggerFromCtx(ctx)
}

// LoggerWithFields adds fields to the logger in the context and returns a new context
func LoggerWithFields(ctx context.Context, fields ...zap.Field) context.Context {
	currentLogger := GetLoggerFromCtx(ctx)
	newLogger := currentLogger.With(fields...)
	return ContextWithLogger(ctx, newLogger)
}
