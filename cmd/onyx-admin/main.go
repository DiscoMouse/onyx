// Package main provides the administrative console for the Onyx Security Appliance.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"onyx/internal/config"
	"onyx/internal/ui"

	"github.com/spf13/cobra"
)

var version = "v0.1.6"

var rootCmd = &cobra.Command{
	Use:   "onyx-admin",
	Short: "Onyx Admin: The management console for the Onyx Security Appliance.",
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Launch the interactive status dashboard",
	Run: func(cmd *cobra.Command, args []string) {
		// Define the default configuration path
		home, _ := os.UserHomeDir()
		configPath := filepath.Join(home, ".config", "onyx", "servers.toml")

		// Attempt to load the server configuration
		conf, err := config.LoadConfig(configPath)
		if err != nil {
			// If config is missing, we proceed in 'Local Only' mode for now
			fmt.Printf("Note: Running in local mode (%v)\n", err)
		} else {
			fmt.Printf("Loaded %d remote server(s) from configuration.\n", len(conf.Servers))
		}

		// Launch the TUI
		if err := ui.StartTUI(version); err != nil {
			fmt.Printf("Error launching console: %v\n", err)
			os.Exit(1)
		}
	},
}

func main() {
	rootCmd.AddCommand(statusCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
