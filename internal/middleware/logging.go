package middleware

import (
	"time"

	"github.com/DimKa163/goseller/internal/shared/sellerlog"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		localLogger := sellerlog.FromContext(ctx, logger)
		localLogger = localLogger.With(zap.String("method", ctx.Request.Method)).With(zap.String("path", ctx.Request.URL.Path))
		localLogger.Info("incoming request")
		ctx.Request = ctx.Request.WithContext(sellerlog.WithLogger(ctx.Request.Context(), localLogger))
		ctx.Set("logger", localLogger)
		startTime := time.Now()
		ctx.Next()
		elapsedTime := time.Since(startTime)
		localLogger.Info("request processed", zap.Int("status", ctx.Writer.Status()), zap.Duration("elapsed", elapsedTime))
	}
}
