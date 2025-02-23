package middlewares

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

const LoggerCtx = "logger"

func LoggerMw(logger *slog.Logger) gin.HandlerFunc {
	lMW := sloggin.NewWithConfig(logger, sloggin.Config{
		WithUserAgent: true,
		WithRequestID: true,
	})
	return func(c *gin.Context) {
		c.Set(LoggerCtx, logger)
		c.Next()
		lMW(c)
	}
}

func GetLoggerFromContext(ctx context.Context) *slog.Logger {
	return ctx.Value(LoggerCtx).(*slog.Logger)
}
