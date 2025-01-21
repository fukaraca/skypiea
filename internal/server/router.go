package service

import (
	"log/slog"
	"net/http"

	"github.com/fukaraca/skypiea/internal/config"
	logg "github.com/fukaraca/skypiea/pkg/log"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func NewRouter(cfg *config.Server, logger *slog.Logger, opts ...gin.OptionFunc) *gin.Engine {
	gin.SetMode(cfg.GinMode)
	e := gin.New(opts...)
	e.LoadHTMLGlob("./web/templates/**/*.html")
	e.Use(commonMiddlewares(logger)...)
	return e
}

func commonMiddlewares(logger *slog.Logger) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		gin.Recovery(),
		static.Serve("/web/static", static.LocalFile("./web/static", false)),
		logg.GinMiddleware(logger.With("via", "rest")),
	}
}

func (s *Server) bindRoutes() {
	s.engine.NoRoute()

	s.engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", gin.H{
			"Title":   "Home",
			"CSSFile": "index.css",
		})
	})
}
