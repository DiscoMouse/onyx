// Package main is the entry point for the Onyx Security Engine.
// It handles the proxy data plane and administrative control plane.
package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/caddy-dns/ovh"
	_ "github.com/caddyserver/caddy/v2/modules/standard"
	_ "github.com/corazawaf/coraza-caddy/v2"

	"onyx/internal/engine"

	"github.com/spf13/cobra"
)

var pairMode bool
var version = "dev" // Default for local builds without tags

var rootCmd = &cobra.Command{
	Use:   "onyx",
	Short: "Onyx: A high-performance, secure mTLS proxy.",
	Long: `Onyx is a security appliance designed to provide encrypted 
access to internal services using Mutual TLS (mTLS).`,
	Run: func(cmd *cobra.Command, args []string) {
		if pairMode {
			runPairing()
			return
		}

		// Normal engine startup would go here
		fmt.Println("Onyx Engine starting in proxy mode...")
	},
}

// runPairing handles the secure bootstrapping of a new admin client.
func runPairing() {
	token, err := engine.GeneratePairingToken()
	if err != nil {
		log.Fatalf("Failed to generate secure token: %v", err)
	}

	fmt.Println("--------------------------------------------------")
	fmt.Println("ONYX BOOTSTRAP MODE")
	fmt.Println("--------------------------------------------------")
	fmt.Println("Use this mode to pair a new admin console via SSH.")

	engine.StartPairingMode(token)
}

func main() {
	rootCmd.Flags().BoolVarP(&pairMode, "pair", "p", false, "Enable temporary pairing mode for new admin consoles")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
