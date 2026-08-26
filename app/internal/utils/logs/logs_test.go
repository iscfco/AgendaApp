package logs

import (
	"context"
	"sync"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func TestLogger(t *testing.T) {
	logger := InitLogger()
	logger.Info("Starting test Test")

	// Set the logger first time in the context
	ctx := ContextWithLogger(context.Background(), logger)
	Logger(ctx).Info("Logger stored first time in context")

	// Add field 1
	requestID := uuid.New().String()
	ctx = LoggerWithFields(ctx, zap.String("request_id", requestID))

	wg := &sync.WaitGroup{}

	// Validate Concurrency
	wg.Add(1)
	go func() {
		defer wg.Done()

		ctx := LoggerWithFields(ctx, zap.Int("funcID", 1))
		for i := 0; i < 3; i++ {
			Logger(ctx).Info("Logging from gorutine", zap.Int("i", i))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		ctx := LoggerWithFields(ctx, zap.Int("funcID", 2))
		for i := 0; i < 3; i++ {
			Logger(ctx).Info("Logging from gorutine", zap.Int("j", i))
		}
	}()

	func1(ctx)
	func3(ctx)
	Logger(ctx).Info("After func1 and func3")

	// Validate Concurrency passing the context
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()

		ctx = LoggerWithFields(ctx, zap.Int("funcWithCtxID", 1))
		for i := 0; i < 3; i++ {
			Logger(ctx).Info("Logging from gorutine with context", zap.Int("i", i))
		}
	}(ctx)

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()

		ctx = LoggerWithFields(ctx, zap.Int("funcWithCtxID", 2))
		for i := 0; i < 3; i++ {
			Logger(ctx).Info("Logging from gorutine with context", zap.Int("j", i))
		}
	}(ctx)

	wg.Wait()
}

func func1(ctx context.Context) {
	ctx = LoggerWithFields(ctx, zap.Int("transaction_id", 1001))
	Logger(ctx).Info("Logging from func1")

	func2(ctx)
}

func func2(ctx context.Context) {
	ctx = LoggerWithFields(ctx, zap.String("user_id", "USR-0001"))
	Logger(ctx).Info("Logging from func2")
}

func func3(ctx context.Context) {
	ctx = LoggerWithFields(ctx, zap.String("session_id", "SESS-0001"))
	Logger(ctx).Info("Logging from func3")
}
