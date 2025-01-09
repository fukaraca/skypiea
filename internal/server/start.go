package service

import (
	"log"

	"github.com/fukaraca/skypiea/internal/config"
	logg "github.com/fukaraca/skypiea/pkg/log"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func Start(cfg *config.Config) error {
	logger := logg.New(cfg.Log).With("service mode", cfg.ServiceMode)

	gin.SetMode(cfg.Server.GinMode)
	e := gin.New()
	e.LoadHTMLGlob("./web/templates/*.html")
	e.Use(
		gin.Recovery(),
		static.Serve("/", static.LocalFile("./web/static", false)),
		logg.GinMiddleware(logger.With("via", "rest")),
	)
	server := NewServer(cfg, e, logger)
	server.bindRoutes()

	logger.Info("server started")
	log.Fatal(server.engine.Run(":" + cfg.Server.Port))
	// get db conn
	// get server
	// init server
	return nil
}
