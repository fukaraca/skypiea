package main

import (
	"fmt"
	"github.com/fukaraca/skypiea/internal/config"
	logg "github.com/fukaraca/skypiea/pkg/log"
	"github.com/spf13/cobra"
	_ "github.com/spf13/cobra"
	"log"
	"log/slog"
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

func runService(mode string) error {
	log.Printf("mode: %s", mode)
	return nil
}

func initialize(mode string) error {
	cfg, err := config.LoadConfig(configName, "./configs")
	if err != nil {
		return err
	}
	logg.New(cfg.Log).Info("config loaded", slog.Any("config", cfg))
	runService(mode)
	return nil
}

func loadConfig() *cobra.Command {
	return &cobra.Command{
		Use:   "load-config",
		Short: "load config",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadConfig(configName, "./configs")
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
