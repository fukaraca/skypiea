package middlewares

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/fukaraca/skypiea/internal/config"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

type ctxKey string

const (
	LoggerCtx ctxKey = "logger"
	GinCtx    ctxKey = "GinCtx"
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
		c.Set(string(LoggerCtx), logger)
		c.Next()
	}
}

func LoggerMw(logger *slog.Logger) gin.HandlerFunc {
	return sloggin.NewWithConfig(logger, sloggin.Config{
		WithUserAgent: true,
		WithRequestID: true,
		Filters:       []sloggin.Filter{filterToSkipLog},
	})
}

func GetLoggerFromContext(ctx context.Context) *slog.Logger {
	if v, ok := ctx.Value(string(LoggerCtx)).(*slog.Logger); ok {
		return v
	}
	return slog.Default()
}

func GetGinCtxFromContext(ctx context.Context) *gin.Context {
	return ctx.Value(GinCtx).(*gin.Context)
}

var filterToSkipLog sloggin.Filter = func(ctx *gin.Context) bool {
	if ctx.FullPath() != "/healthz" || ctx.Errors != nil {
		return true
	}
	return false
}
