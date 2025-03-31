package server

import (
	"log"
	"net"
	"time"

	"github.com/fukaraca/skypiea/internal/config"
	"github.com/fukaraca/skypiea/pkg/gwt"
	logg "github.com/fukaraca/skypiea/pkg/log"
	"github.com/fukaraca/skypiea/pkg/session"
)

func Start(cfg *config.Config) error {
	logger := logg.New(cfg.Log).With("service mode", cfg.ServiceMode)
	router := NewRouter(cfg.Server, logger)
	router.SetTrustedProxies(nil)
	db, err := cfg.Database.GetDBConn()
	if err != nil {
		return err
	}
	session.Cache = session.NewManager(&gwt.Config{Secret: []byte("secret")}, db, time.Minute*10)
	server := NewServer(cfg, router, db, logger)
	server.bindRoutes()

	logger.Info("server started")
	log.Fatal(server.engine.Run(net.JoinHostPort(cfg.Server.Address, cfg.Server.Port)))
	return nil
}
