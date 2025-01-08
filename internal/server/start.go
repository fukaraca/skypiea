package service

import (
	"context"
	"github.com/fukaraca/skypiea/internal/config"
	logg "github.com/fukaraca/skypiea/pkg/log"
)

func Start(ctx context.Context, cfg *config.Config) error {
	logger := logg.New(cfg.Log).With("mode", cfg.RunningMode)
	logger.Info("start server")
	// get db conn
	// get server
	// init server
	return nil
}
