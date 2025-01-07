package main

import (
	"github.com/spf13/cobra"
	_ "github.com/spf13/cobra"
	"log"
)

func main() {
	if err := RootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}

func RootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "skypiea",
		Short: "skypiea",
		RunE:  nil,
	}

	// add commands
	// run service, workers, migration(nil),
	// load config
	return rootCmd
}
