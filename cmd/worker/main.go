package main

import (
	"log"
	"log/slog"

	"github.com/fukaraca/skypiea/internal/config"
	"github.com/fukaraca/skypiea/internal/worker"
	logg "github.com/fukaraca/skypiea/pkg/log"
	"github.com/spf13/cobra"
)

var (
	Version    string = "dev"
	configName string
)

func main() {
	if err := RootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}

func RootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "skypiea-worker",
		Short: "skypiea-worker",
		RunE: func(cmd *cobra.Command, args []string) error {
			return initialize()
		},
		Version: Version,
	}
	rootCmd.PersistentFlags().StringVar(&configName, "config", "config.example.yml", "config file name in configs folder")
	return rootCmd
}

func initialize() error {
	cfg := config.NewConfig()
	err := cfg.Load(configName, "./configs")
	if err != nil {
		return err
	}
	logger := logg.New(cfg.Log)

	cfg.Worker.Version = Version
	cfg.ServiceMode = config.ModeBackgroundWorker

	logger.Info("worker initialized", slog.Any("Version", cfg.Worker.Version))
	return worker.Start(cfg, logger)
}
