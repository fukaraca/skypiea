package service

import (
	"github.com/fukaraca/skypiea/internal/config"
	logg "github.com/fukaraca/skypiea/pkg/log"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"log/slog"
	"net/http"
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
	s.engine.GET("/ping", func(c *gin.Context) {
		s.Logger.Info("ping ponged", "id", sloggin.GetRequestID(c))
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	s.engine.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "base.html", gin.H{
			"title": "Welcome to Aquila AI",
		})
	})
	s.engine.POST("/api/query", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"response": "test reponse.",
		})
	})
}
