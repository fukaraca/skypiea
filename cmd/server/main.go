package main

import (
	"log"
	"log/slog"

	"github.com/fukaraca/skypiea/internal/config"
	service "github.com/fukaraca/skypiea/internal/server"
	"github.com/fukaraca/skypiea/internal/storage/migration"
	logg "github.com/fukaraca/skypiea/pkg/log"
	_ "github.com/lib/pq"
	"github.com/rubenv/sql-migrate"
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
		Use:   "skypiea-server",
		Short: "skypiea-server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return initialize()
		},
		Version: Version,
	}
	rootCmd.PersistentFlags().StringVar(&configName, "config", "config.example.yml", "config file name in configs folder")

	rootCmd.AddCommand(loadConfig())
	rootCmd.AddCommand(migrateCmd())
	return rootCmd
}

func initialize() error {
	cfg := config.NewConfig()
	err := cfg.Load(configName, "./configs")
	if err != nil {
		return err
	}
	cfg.Server.Version = Version
	cfg.ServiceMode = config.ModeHttpServer

	logg.New(cfg.Log).Info("server initialized", slog.Any("version", cfg.Server.Version))
	return service.Start(cfg)
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
			return nil
		},
	}
}

func migrateCmd() *cobra.Command {
	migrateCmd := &cobra.Command{
		Use:          "migration",
		SilenceUsage: true,
	}
	up := &cobra.Command{
		Use:          "up",
		RunE:         migrateDB(migrate.Up),
		SilenceUsage: true,
	}
	down := &cobra.Command{
		Use:          "down",
		RunE:         migrateDB(migrate.Down),
		SilenceUsage: true,
	}
	migrateCmd.AddCommand(up, down)
	return migrateCmd
}

func migrateDB(direction migrate.MigrationDirection) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cfg := config.NewConfig()
		err := cfg.Load(configName, "./configs")
		if err != nil {
			return err
		}
		return migration.RunMigration(cfg.Database, direction)
	}
}
