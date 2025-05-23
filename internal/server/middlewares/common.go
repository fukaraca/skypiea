package middlewares

import (
	"log/slog"

	"github.com/fukaraca/skypiea/internal/config"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func CommonMiddlewares(logger *slog.Logger, cfg *config.Server) []gin.HandlerFunc {
	logger = logger.With("via", "rest")
	return []gin.HandlerFunc{
		gin.Recovery(),
		static.Serve("/web/static", static.LocalFile("./web/static", false)),
		ErrorHandlerMw(),
		LoggerMw(logger),
		ContextMW(logger, cfg),
	}
}
