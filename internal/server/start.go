package service

import (
	"log"

	"github.com/fukaraca/skypiea/internal/config"
	logg "github.com/fukaraca/skypiea/pkg/log"
)

func Start(cfg *config.Config) error {
	logger := logg.New(cfg.Log).With("service mode", cfg.ServiceMode)
	router := NewRouter(cfg.Server, logger)
	router.SetTrustedProxies(nil)
	server := NewServer(cfg, router, logger)
	server.bindRoutes()

	logger.Info("server started")
	log.Fatal(server.engine.Run(":" + cfg.Server.Port))
	// get db conn
	// get server
	// init server
	return nil
}
