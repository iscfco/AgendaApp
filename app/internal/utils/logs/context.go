package logs

import (
	"context"

	"go.uber.org/zap"
)

// ContextWithLogger injects a logger into the context
func ContextWithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey{}, logger)
}

// GetLoggerFromCtx gets the logger from the context. If the context is nil, it returns the global logger
func GetLoggerFromCtx(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return zap.L()
	}
	if logger, ok := ctx.Value(loggerCtxKey{}).(*zap.Logger); ok {
		return logger
	}
	return zap.L()
}
