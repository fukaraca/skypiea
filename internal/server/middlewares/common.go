package middlewares

import (
	"log/slog"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func CommonMiddlewares(logger *slog.Logger) []gin.HandlerFunc {
	logger = logger.With("via", "rest")
	return []gin.HandlerFunc{
		gin.Recovery(),
		static.Serve("/web/static", static.LocalFile("./web/static", false)),
		LoggerMw(logger),
		ErrorHandlerMw(),
	}
}
