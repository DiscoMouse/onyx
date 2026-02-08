// Package main provides the administrative console for the Onyx Security Appliance.
// It serves as the primary interface for managing local and remote Onyx instances
// through a Terminal User Interface (TUI).
package main

import (
	"fmt"
	"os"

	"onyx/internal/ui"

	"github.com/spf13/cobra"
)

// version defines the current build version of the admin tool.
var version = "v0.1.6"

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "onyx-admin",
	Short: "Onyx Admin: The management console for the Onyx Security Appliance.",
	Long: `Onyx Admin is a specialized TUI tool designed to monitor and manage 
Onyx security engines. It supports local management via unix sockets 
and remote management via mTLS-secured connections.`,
}

// statusCmd represents the command to launch the interactive dashboard.
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Launch the interactive status dashboard",
	Run: func(cmd *cobra.Command, args []string) {
		// StartTUI is the entry point for the Bubble Tea interface
		if err := ui.StartTUI(version); err != nil {
			fmt.Printf("Error launching console: %v\n", err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func main() {
	rootCmd.AddCommand(statusCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
