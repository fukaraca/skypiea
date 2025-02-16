package service

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/fukaraca/skypiea/internal/config"
	"github.com/fukaraca/skypiea/internal/server/middlewares"
	"github.com/fukaraca/skypiea/internal/service/templater"
)

func NewRouter(cfg *config.Server, logger *slog.Logger, opts ...gin.OptionFunc) *gin.Engine {
	gin.SetMode(cfg.GinMode)
	e := gin.New(opts...)
	templates := templater.New()
	templates.LoadHTMLGlob("./web/templates")
	e.HTMLRender = templates
	e.Use(middlewares.CommonMiddlewares(logger)...)
	return e
}

func (s *Server) bindRoutes() {
	s.engine.NoRoute()
	s.RegisterRoutes(viewRoutes())
	s.RegisterRoutes(commonRoutes())
	s.RegisterRoutes(strictRoutes())
}
