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
	e.Use(middlewares.CommonMiddlewares(logger)...)
	return e
}

func (s *Server) bindRoutes() {
	s.engine.NoRoute(middlewares.CommonAuthMw(), handlers.NoRoute404)
	// v1 := s.engine.Group(V1)

	s.RegisterRoutes(nil, viewRoutes(s), middlewares.CommonAuthMw())
	s.RegisterRoutes(nil, openRoutes(s), middlewares.CommonAuthMw())
	s.RegisterRoutes(nil, strictRoutes(s), middlewares.StrictAuthMw())
}

//
