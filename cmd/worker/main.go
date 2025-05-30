package main

import (
	"log"
	"log/slog"

	"github.com/fukaraca/skypiea/internal/config"
	"github.com/fukaraca/skypiea/internal/worker"
	logg "github.com/fukaraca/skypiea/pkg/log"
	"github.com/spf13/cobra"
)

var configName string

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
	logger.Info("worker initialized", slog.Any("config", cfg))
	cfg.ServiceMode = config.ModeBackgroundWorker
	return worker.Start(cfg, logger)
}
