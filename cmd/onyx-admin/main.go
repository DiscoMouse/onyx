// Package main provides the administrative console for the Onyx Security Appliance.
package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"onyx/internal/config"
	"onyx/internal/crypto"
	"onyx/internal/ui"

	"github.com/spf13/cobra"
)

var version = "dev" // Injected at build time

var rootCmd = &cobra.Command{
	Use:   "onyx-admin [target-ip]",
	Short: "Onyx Admin: The management console for the Onyx Security Appliance.",
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Load Configuration
		home, _ := os.UserHomeDir()
		configPath := filepath.Join(home, ".config", "onyx", "config.toml")

		conf, err := config.LoadConfig(configPath)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		// 2. Direct Connect Mode (onyx-admin 10.0.0.1)
		if len(args) > 0 {
			targetIP := args[0]
			port, _ := cmd.Flags().GetInt("port")

			// Check if this node is already known (to get its name/stats)
			var targetNode *config.Node
			for i := range conf.Nodes {
				if conf.Nodes[i].Address == targetIP {
					targetNode = &conf.Nodes[i]
					// If the user provided a port flag, override the saved port
					if cmd.Flags().Changed("port") {
						targetNode.Port = port
					}
					break
				}
			}

			// If unknown, create a temporary ad-hoc node
			if targetNode == nil {
				targetNode = &config.Node{
					Name:     "Ad-Hoc Session",
					Address:  targetIP,
					Port:     port,
					AddedAt:  time.Now(),
					LastSeen: time.Now(),
				}
			}

			// Launch directly into Dashboard
			if err := ui.StartDashboard(version, targetNode); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			return
		}

		// 3. Interactive Menu Mode (Default)
		for {
			action, node := ui.StartMenu(version, conf)

			switch action {
			case ui.ActionQuit:
				fmt.Println("Bye!")
				return

			case ui.ActionConnect:
				// Launch Dashboard
				if err := ui.StartDashboard(version, node); err != nil {
					fmt.Printf("Dashboard Error: %v\n", err)
					// Pause so user can see error before returning to menu
					time.Sleep(2 * time.Second)
				}
				// Loop continues back to menu after dashboard exits

			case ui.ActionPair:
				// Ideally, we would launch a TUI form here.
				// For now, guide the user to the CLI command.
				fmt.Println("\nTo pair a new engine, please run:")
				fmt.Println("  onyx-admin pair <ip-address> --token <token>")
				fmt.Println("\n(Press Enter to return to menu)")
				fmt.Scanln()
			}
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
		port, _ := cmd.Flags().GetInt("port")

		if token == "" {
			fmt.Println("Error: A pairing --token is required. Check server logs.")
			os.Exit(1)
		}

		// 1. Setup local paths
		home, _ := os.UserHomeDir()
		baseDir := filepath.Join(home, ".config", "onyx")
		certDir := filepath.Join(baseDir, "certs")
		configPath := filepath.Join(baseDir, "config.toml")

		if err := os.MkdirAll(certDir, 0700); err != nil {
			fmt.Printf("Failed to create config directory: %v\n", err)
			os.Exit(1)
		}

		keyPath := filepath.Join(certDir, "client.key")
		certPath := filepath.Join(certDir, "client.crt")

		// 2. Generate or load local identity
		// Check if key already exists to avoid overwriting identity
		var priv ed25519.PrivateKey
		if _, err := os.Stat(keyPath); os.IsNotExist(err) {
			fmt.Println("Generating new local identity...")
			_, newPriv, err := ed25519.GenerateKey(rand.Reader)
			if err != nil {
				fmt.Printf("Failed to generate identity: %v\n", err)
				os.Exit(1)
			}
			priv = newPriv

			privPEM, err := crypto.EncodePrivateKey(priv)
			if err != nil {
				fmt.Printf("Failed to encode private key: %v\n", err)
				os.Exit(1)
			}

			if err := crypto.SavePEM(keyPath, privPEM); err != nil {
				fmt.Printf("Failed to save private key: %v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Println("Loading existing identity...")
			loadedPriv, err := crypto.LoadPrivateKey(keyPath)
			if err != nil {
				fmt.Printf("Failed to load private key: %v\n", err)
				os.Exit(1)
			}
			priv = loadedPriv
		}

		// 3. Create CSR
		hostname, _ := os.Hostname()
		commonName := fmt.Sprintf("admin@%s", hostname)

		csrPEM, err := crypto.GenerateCSR(priv, commonName)
		if err != nil {
			fmt.Printf("Failed to create CSR: %v\n", err)
			os.Exit(1)
		}

		// 4. Perform Handshake
		fmt.Printf("Initiating secure pairing with %s:%d...\n", targetIP, port)

		// Construct the full address (IP:Port)
		targetAddr := fmt.Sprintf("%s:%d", targetIP, port)

		// PASS targetAddr INSTEAD OF targetIP
		signedCert, err := crypto.PerformHandshake(targetAddr, token, csrPEM)
		if err != nil {
			fmt.Printf("Handshake failed: %v\n", err)
			os.Exit(1)
		}

		// 5. Save the signed certificate
		if err := os.WriteFile(certPath, signedCert, 0644); err != nil {
			fmt.Printf("Failed to save certificate: %v\n", err)
			os.Exit(1)
		}

		// 6. PERSISTENCE: Save the server to config.toml
		conf, err := config.LoadConfig(configPath)
		if err != nil {
			fmt.Printf("Warning: Failed to load config for update: %v\n", err)
			// Proceed with empty config if load fails
			conf = &config.AdminConfig{}
		}

		// AddNode handles duplication checks automatically
		conf.AddNode("Onyx Engine", targetIP, port)

		if err := conf.SaveConfig(configPath); err != nil {
			fmt.Printf("Warning: Failed to save server to config: %v\n", err)
		}

		fmt.Println("[✓] Pairing complete! Your device is now authorized.")
		fmt.Printf("Identity saved to: %s\n", certDir)
		fmt.Printf("Server added to: %s\n", configPath)
	},
}

func main() {
	// Global Flags
	rootCmd.PersistentFlags().IntP("port", "p", 2305, "Target port")

	// Pair specific flags
	pairCmd.Flags().StringP("token", "t", "", "One-time pairing token")

	rootCmd.AddCommand(pairCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
