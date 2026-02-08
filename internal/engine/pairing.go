// Package engine implements the core logic for the Onyx security service,
// including the temporary pairing listener for administrative bootstrapping.
package engine

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"
)

// GeneratePairingToken creates a high-entropy, human-readable 8-character token.
func GeneratePairingToken() (string, error) {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	out := make([]byte, 8)
	for i := range out {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		out[i] = charset[num.Int64()]
	}
	return fmt.Sprintf("%s-%s", string(out[:4]), string(out[4:])), nil
}

// StartPairingMode opens a 5-minute window for a new admin to pair.
// It will shut down immediately upon success or prompt to retry on timeout.
func StartPairingMode(token string) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("\n[PAIRING MODE ACTIVE]\n")
		fmt.Printf("Token: %s\n", token)
		fmt.Printf("Port:  2305\n")
		fmt.Printf("Window: 5 Minutes\n\n")

		// resultChan allows the HTTP handler to tell the main loop we're done.
		resultChan := make(chan bool)
		pairingCtx, pairingCancel := context.WithTimeout(context.Background(), 5*time.Minute)

		mux := http.NewServeMux()
		mux.HandleFunc("/pair", func(w http.ResponseWriter, r *http.Request) {
			// TODO: Verify token and sign CSR
			fmt.Println("\n[+] Pairing successful! Finalizing...")
			w.WriteHeader(http.StatusOK)
			resultChan <- true
		})

		srv := &http.Server{Addr: ":2305", Handler: mux}

		go func() {
			if err := srv.ListenAndServe(); err != http.ErrServerClosed {
				fmt.Printf("Error: %v\n", err)
			}
		}()

		success := false
		select {
		case <-resultChan:
			fmt.Println("[✓] Device paired successfully.")
			success = true
		case <-pairingCtx.Done():
			fmt.Println("\n[!] Pairing window expired.")
		}

		// Shutdown the server gracefully with a short 2-second deadline
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 2*time.Second)
		srv.Shutdown(shutdownCtx)

		// CRITICAL: Call cancel functions to release resources/timers immediately
		shutdownCancel()
		pairingCancel()

		if success {
			break
		}

		fmt.Print("No onyx-admin client paired. Try again? (y/N): ")
		answer, _ := reader.ReadString('\n')
		if strings.TrimSpace(strings.ToLower(answer)) != "y" {
			fmt.Println("Exiting pairing mode.")
			break
		}
	}
}
