package main

import (
	"context"
	"fmt"
	"github.com/fukaraca/skypiea/internal/config"
	service "github.com/fukaraca/skypiea/internal/server"
	logg "github.com/fukaraca/skypiea/pkg/log"
	"github.com/spf13/cobra"
	_ "github.com/spf13/cobra"
	"log"
	"log/slog"
	"slices"
)

var configName, mode string

func main() {
	if err := RootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}

func RootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "skypiea",
		Short: "skypiea",
		RunE: func(cmd *cobra.Command, args []string) error {
			return initialize(mode)
		},
	}
	rootCmd.PersistentFlags().StringVar(&configName, "config", "config.example.yml", "config file name in configs folder")
	rootCmd.Flags().StringVar(&mode, "mode", "", "Mode to run the application (server or worker)")
	rootCmd.MarkFlagRequired("mode")

	rootCmd.AddCommand(loadConfig())
	rootCmd.AddCommand(migrate())
	return rootCmd
}

func initialize(mode string) error {
	if !slices.Contains([]string{config.ModeHttpServer, config.ModeBackgroundWorker}, mode) {
		return fmt.Errorf("invalid mode: %s", mode)
	}
	cfg := config.NewConfig()
	err := cfg.Load(configName, "./configs")
	if err != nil {
		return err
	}
	logg.New(cfg.Log).Info("config loaded", slog.Any("config", cfg))
	cfg.RunningMode = mode
	return service.Start(context.Background(), cfg)
}

func loadConfig() *cobra.Command {
	return &cobra.Command{
		Use:   "load-config",
		Short: "load config",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.NewConfig()
			err := cfg.Load(configName, "./configs")
			if err != nil {
				return err
			}
			logg.New(cfg.Log).Info("config loadedddd", slog.Any("config", cfg))
			fmt.Printf("%+v\n", *cfg)
			return nil
		},
	}
}

func migrate() *cobra.Command {
	migrateCmd := &cobra.Command{
		Use:          "migrate",
		SilenceUsage: true,
	}
	up := &cobra.Command{
		Use:          "up",
		RunE:         nil,
		SilenceUsage: true,
	}
	down := &cobra.Command{
		Use:          "down",
		RunE:         nil,
		SilenceUsage: true,
	}
	migrateCmd.AddCommand(up, down)
	return migrateCmd
}
