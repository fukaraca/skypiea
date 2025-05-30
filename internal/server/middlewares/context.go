package middlewares

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/fukaraca/skypiea/internal/config"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

const (
	LoggerCtx = "logger"
	GinCtx    = "GinCtx"
)

// ContextMW sets up and propagates the request context (adding logger, Gin context, and timeout) for downstream handlers.
func ContextMW(logger *slog.Logger, cfg *config.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, LoggerCtx, logger)
		ctx = context.WithValue(ctx, GinCtx, c)
		ctx, cancel := context.WithTimeout(ctx, cfg.DefaultRequestTimeout)
		defer func() {
			cancel()
			if ctx.Err() == context.DeadlineExceeded {
				c.Writer.WriteHeader(http.StatusGatewayTimeout) // we can delegate this to nginx? slow endpoints?
			}
		}()
		c.Request = c.Request.WithContext(ctx)
		c.Set(LoggerCtx, logger)
		c.Next()
	}
}

func LoggerMw(logger *slog.Logger) gin.HandlerFunc {
	return sloggin.NewWithConfig(logger, sloggin.Config{
		WithUserAgent: true,
		WithRequestID: true,
	})
}

func GetLoggerFromContext(ctx context.Context) *slog.Logger {
	return ctx.Value(LoggerCtx).(*slog.Logger)
}

func GetGinCtxFromContext(ctx context.Context) *gin.Context {
	return ctx.Value(GinCtx).(*gin.Context)
}
