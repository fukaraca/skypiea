package server

import (
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/fukaraca/skypiea/internal/config"
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
	server, err := NewServer(cfg, router, db, logger)
	if err != nil {
		return err
	}
	server.bindRoutes()
	session.Cache = session.NewManager(cfg.JWT, server.Repo.Users, cfg.Server.SessionTimeout)

	logger.Info("Server started")
	httpServer := &http.Server{
		Addr:              net.JoinHostPort(cfg.Server.Address, cfg.Server.Port),
		Handler:           server.engine,
		ReadHeaderTimeout: time.Second * 5,
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		logger.Warn("received interrupt signal")
		if errInner := httpServer.Close(); errInner != nil {
			log.Fatal("Server Close:", errInner)
		}
	}()

	if err = httpServer.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			logger.Warn("Server closed under request")
		} else {
			log.Fatal("Server closed unexpectedly")
		}
	}

	logger.Info("Server shutting down")
	return nil
}
