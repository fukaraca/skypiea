package worker

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fukaraca/skypiea/internal/config"
)

func Start(cfg *config.Config, logger *slog.Logger) error {
	logger.Info("start worker")

	ticker := time.NewTicker(cfg.Worker.IntervalTicker)
	defer ticker.Stop()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(quit)

	for {
		select {
		case <-ticker.C:
			logger.Info("worker is running at the background and doing very big tasks however it doesn't cost much interestingly")
		case sig := <-quit:
			logger.Warn("received interrupt signal", slog.Any("signal", sig))
			ticker.Stop()
			return nil
		}
	}
}
