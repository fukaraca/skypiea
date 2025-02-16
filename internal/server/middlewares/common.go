package middlewares

import (
	"log/slog"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	logg "github.com/fukaraca/skypiea/pkg/log"
)

func CommonMiddlewares(logger *slog.Logger) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		gin.Recovery(),
		static.Serve("/web/static", static.LocalFile("./web/static", false)),
		logg.GinMiddleware(logger.With("via", "rest")),
	}
}
