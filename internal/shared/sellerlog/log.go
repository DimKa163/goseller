package sellerlog

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type loggerContextName string

const loggerName loggerContextName = "github.com/DimKa163/goseller_gin_logger"

func FromContext(ctx context.Context, defaultLogger *zap.Logger) *zap.Logger {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ctx = ginCtx.Request.Context()
	}

	if logger, ok := ctx.Value(loggerName).(*zap.Logger); ok {
		return logger
	}
	return defaultLogger
}

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerName, logger)
}
