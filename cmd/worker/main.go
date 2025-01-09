package main

import (
	"context"
	"github.com/fukaraca/skypiea/internal/config"
	service "github.com/fukaraca/skypiea/internal/server"
	logg "github.com/fukaraca/skypiea/pkg/log"
	"github.com/spf13/cobra"
	_ "github.com/spf13/cobra"
	"log"
	"log/slog"
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
	logg.New(cfg.Log).Info("worker initialized", slog.Any("config", cfg))
	cfg.RunningMode = config.ModeBackgroundWorker
	return service.Start(context.Background(), cfg)
}
