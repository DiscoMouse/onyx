// Package main provides the administrative console for the Onyx Security Appliance.
package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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

			// Check if this node is already known
			var targetNode *config.Node
			for i := range conf.Nodes {
				if conf.Nodes[i].Address == targetIP {
					targetNode = &conf.Nodes[i]
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

			if err := ui.StartDashboard(version, targetNode); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			return
		}

		// 3. Interactive Menu Mode (Default)
		for {
			// Reload config on every loop so new pairings show up immediately
			conf, _ = config.LoadConfig(configPath)

			action, node := ui.StartMenu(version, conf)

			switch action {
			case ui.ActionQuit:
				fmt.Println("Bye!")
				return

			case ui.ActionConnect:
				if err := ui.StartDashboard(version, node); err != nil {
					fmt.Printf("Dashboard Error: %v\n", err)
					time.Sleep(2 * time.Second)
				}

			case ui.ActionPair:
				// Launch the TUI Form
				result, submitted := ui.StartPairingForm()
				if !submitted {
					continue // User cancelled, go back to menu
				}

				// Convert port string to int
				portInt, err := strconv.Atoi(result.Port)
				if err != nil {
					fmt.Printf("\nError: Invalid port number '%s'\n", result.Port)
					time.Sleep(2 * time.Second)
					continue
				}

				// Use default port if 0 or empty (though form usually catches this)
				if portInt == 0 {
					portInt = 2305
				}

				// Run the pairing logic
				fmt.Println("\nConnecting to engine...")
				if err := performPairing(result.Address, portInt, result.Token); err != nil {
					fmt.Printf("\nPairing Failed: %v\n", err)
					fmt.Println("(Press Enter to return to menu)")
					fmt.Scanln()
				} else {
					// Success! Loop will reload config and show the new node.
					fmt.Println("\nSuccess! Returning to menu...")
					time.Sleep(1 * time.Second)
				}
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
			fmt.Println("Error: A pairing --token is required.")
			os.Exit(1)
		}

		if err := performPairing(targetIP, port, token); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// performPairing handles the core identity generation, handshake, and config persistence.
// This is used by both the CLI 'pair' command and the TUI Form.
func performPairing(targetIP string, port int, token string) error {
	// 1. Setup local paths
	home, _ := os.UserHomeDir()
	baseDir := filepath.Join(home, ".config", "onyx")
	certDir := filepath.Join(baseDir, "certs")
	configPath := filepath.Join(baseDir, "config.toml")

	if err := os.MkdirAll(certDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	keyPath := filepath.Join(certDir, "client.key")
	certPath := filepath.Join(certDir, "client.crt")

	// 2. Generate or load local identity
	var priv ed25519.PrivateKey
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		fmt.Println("Generating new local identity...")
		_, newPriv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return fmt.Errorf("failed to generate identity: %w", err)
		}
		priv = newPriv

		privPEM, err := crypto.EncodePrivateKey(priv)
		if err != nil {
			return fmt.Errorf("failed to encode private key: %w", err)
		}

		if err := crypto.SavePEM(keyPath, privPEM); err != nil {
			return fmt.Errorf("failed to save private key: %w", err)
		}
	} else {
		fmt.Println("Loading existing identity...")
		loadedPriv, err := crypto.LoadPrivateKey(keyPath)
		if err != nil {
			return fmt.Errorf("failed to load private key: %w", err)
		}
		priv = loadedPriv
	}

	// 3. Create CSR
	hostname, _ := os.Hostname()
	commonName := fmt.Sprintf("admin@%s", hostname)

	csrPEM, err := crypto.GenerateCSR(priv, commonName)
	if err != nil {
		return fmt.Errorf("failed to create CSR: %w", err)
	}

	// 4. Perform Handshake
	targetAddr := fmt.Sprintf("%s:%d", targetIP, port)
	fmt.Printf("Initiating secure handshake with %s...\n", targetAddr)

	signedCert, err := crypto.PerformHandshake(targetAddr, token, csrPEM)
	if err != nil {
		return fmt.Errorf("handshake failed: %w", err)
	}

	// 5. Save the signed certificate
	if err := os.WriteFile(certPath, signedCert, 0644); err != nil {
		return fmt.Errorf("failed to save certificate: %w", err)
	}

	// 6. PERSISTENCE: Save the server to config.toml
	conf, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("Warning: Failed to load config for update: %v\n", err)
		conf = &config.AdminConfig{}
	}

	// AddNode handles duplication checks automatically
	conf.AddNode("Onyx Engine", targetIP, port)

	if err := conf.SaveConfig(configPath); err != nil {
		return fmt.Errorf("failed to save server to config: %w", err)
	}

	fmt.Println("[✓] Pairing complete! Identity authorized.")
	return nil
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
