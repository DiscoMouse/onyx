// Package main provides the administrative console for the Onyx Security Appliance.
package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"

	"onyx/internal/config"
	"onyx/internal/crypto"
	"onyx/internal/ui"

	"github.com/spf13/cobra"
)

var version = "dev" // Default for local builds without tags

var rootCmd = &cobra.Command{
	Use:   "onyx-admin",
	Short: "Onyx Admin: The management console for the Onyx Security Appliance.",
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Launch the interactive status dashboard",
	Run: func(cmd *cobra.Command, args []string) {
		home, _ := os.UserHomeDir()
		configPath := filepath.Join(home, ".config", "onyx", "servers.toml")

		conf, err := config.LoadConfig(configPath)
		if err != nil {
			fmt.Printf("Error loading configuration: %v\n", err)
		}

		if err := ui.StartTUI(version, conf); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var pairCmd = &cobra.Command{
	Use:   "pair [ip-address]",
	Short: "Pair this console with a remote Onyx Security Engine",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetIP := args[0]
		token, _ := cmd.Flags().GetString("token")

		if token == "" {
			fmt.Println("Error: A pairing --token is required. Check server logs.")
			os.Exit(1)
		}

		home, _ := os.UserHomeDir()
		baseDir := filepath.Join(home, ".config", "onyx")
		certDir := filepath.Join(baseDir, "certs")
		configPath := filepath.Join(baseDir, "servers.toml")

		os.MkdirAll(certDir, 0700)

		keyPath := filepath.Join(certDir, "client.key")
		certPath := filepath.Join(certDir, "client.crt")

		fmt.Println("Checking local identity...")
		_, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			fmt.Printf("Failed to generate identity: %v\n", err)
			os.Exit(1)
		}

		privPEM, err := crypto.EncodePrivateKey(priv)
		if err != nil {
			fmt.Printf("Failed to encode private key: %v\n", err)
			os.Exit(1)
		}

		if err := crypto.SavePEM(keyPath, privPEM); err != nil {
			fmt.Printf("Failed to save private key: %v\n", err)
			os.Exit(1)
		}

		hostname, _ := os.Hostname()
		commonName := fmt.Sprintf("admin@%s", hostname)

		csrPEM, err := crypto.GenerateCSR(priv, commonName)
		if err != nil {
			fmt.Printf("Failed to create CSR: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Initiating secure pairing with %s...\n", targetIP)
		signedCert, err := crypto.PerformHandshake(targetIP, token, csrPEM)
		if err != nil {
			fmt.Printf("Handshake failed: %v\n", err)
			os.Exit(1)
		}

		if err := os.WriteFile(certPath, signedCert, 0644); err != nil {
			fmt.Printf("Failed to save certificate: %v\n", err)
			os.Exit(1)
		}

		// Persist the server details to the configuration file.
		conf, err := config.LoadConfig(configPath)
		if err != nil {
			fmt.Printf("Warning: Failed to load config for update: %v\n", err)
		} else {
			conf.AddServer("Remote Engine", targetIP, 2305)
			if err := conf.SaveConfig(configPath); err != nil {
				fmt.Printf("Warning: Failed to persist server config: %v\n", err)
			}
		}

		fmt.Println("[✓] Pairing complete! Your device is now authorized.")
		fmt.Printf("Identity saved to: %s\n", certDir)
	},
}

func main() {
	pairCmd.Flags().StringP("token", "t", "", "One-time pairing token")
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(pairCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
