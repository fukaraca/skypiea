package server

import (
	"log/slog"

	"github.com/fukaraca/skypiea/internal/server/handlers"
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
	e.Use(middlewares.CommonMiddlewares(logger, cfg)...)
	return e
}

func (s *Server) bindRoutes() {
	s.engine.NoRoute(middlewares.CommonAuthMw(), handlers.NoRoute404)
	// v1 := s.engine.Group(V1)

	common := handlers.NewCommonHandler(s.Config)
	s.RegisterRoutes(nil, viewRoutes(s, common), middlewares.CommonAuthMw())
	s.RegisterRoutes(nil, openRoutes(s, common), middlewares.CommonAuthMw())
	s.RegisterRoutes(nil, strictRoutes(s, common), middlewares.StrictAuthMw())
}

//
