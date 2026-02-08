package main

import (
	"fmt"
	"os"
	"os/user" // Added for user identification

	"github.com/caddyserver/caddy/v2"
	caddycmd "github.com/caddyserver/caddy/v2/cmd"
	"github.com/spf13/cobra"

	"onyx/internal/ui"

	_ "github.com/caddy-dns/ovh"
	_ "github.com/caddyserver/caddy/v2/modules/standard"
	_ "github.com/corazawaf/coraza-caddy/v2"
)

var version = "dev"

// isRestrictedUser returns true if the current user is the system 'onyx' account
func isRestrictedUser() bool {
	u, err := user.Current()
	if err != nil {
		return true // Default to restricted if lookup fails
	}
	return u.Username == "onyx"
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "onyx",
		Short: "Onyx: Virtual Infrastructure Orchestrator",
	}

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print Onyx and dependency versions",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Onyx Version:  %s\n", version)
			_, full := caddy.Version()
			fmt.Printf("Caddy Version: %s\n", full)
		},
	})

	var proxyCmd = &cobra.Command{
		Use:                "proxy",
		Short:              "Sub-commands for the underlying Caddy engine",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			os.Args = append([]string{"caddy"}, args...)
			caddycmd.Main()
		},
	}
	rootCmd.AddCommand(proxyCmd)

	rootCmd.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Launch the interactive TUI dashboard",
		Run: func(cmd *cobra.Command, args []string) {
			// Explicitly block the background system user from the TUI
			if isRestrictedUser() {
				fmt.Println("Error: Access Denied. The 'onyx' system user cannot access administrative tools.")
				os.Exit(1)
			}

			if err := ui.StartTUI(version); err != nil {
				fmt.Printf("Error starting TUI: %v\n", err)
				os.Exit(1)
			}
		},
	})

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
