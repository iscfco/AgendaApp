package logs

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// GinZapMiddleware is a middleware that injects a logger into the context
func GinZapMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Create a new logger with the request ID
		reqLogger := logger.With(zap.String("request_id", requestID))

		// Inject the logger into the context
		ctx := ContextWithLogger(c.Request.Context(), reqLogger)
		c.Request = c.Request.WithContext(ctx)

		c.Header("X-Request-ID", requestID)

		c.Next()

		// Retrieve the logger from the context and log the request
		latency := time.Since(start)
		finalLogger := GetLoggerFromCtx(c.Request.Context())
		finalLogger.Info("Petición HTTP completada",
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Duration("latency", latency),
		)
	}
}
