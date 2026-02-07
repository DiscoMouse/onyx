package main

import (
	"fmt"
	"os"

	"github.com/caddyserver/caddy/v2"
	caddycmd "github.com/caddyserver/caddy/v2/cmd"
	"github.com/spf13/cobra"

	// 1. ADD THIS IMPORT (Make sure it matches your module name in go.mod)
	"onyx/internal/ui"

	_ "github.com/caddy-dns/ovh"
	_ "github.com/caddyserver/caddy/v2/modules/standard"
	_ "github.com/corazawaf/coraza-caddy/v2"
)

var version = "dev"

func main() {
	var rootCmd = &cobra.Command{
		Use:   "onyx",
		Short: "Onyx: Virtual Infrastructure Orchestrator",
	}

	// 1. TOP LEVEL VERSION
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print Onyx and dependency versions",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Onyx Version:  %s\n", version)
			_, full := caddy.Version()
			fmt.Printf("Caddy Version: %s\n", full)
		},
	})

	// 2. PROXY COMMANDS (Passthrough)
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

	// 3. STATUS (The TUI)
	rootCmd.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Launch the interactive TUI dashboard",
		Run: func(cmd *cobra.Command, args []string) {
			// 2. REPLACE THE PRINTLN WITH THIS:
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
