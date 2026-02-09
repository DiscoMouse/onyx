// Package engine implements the core logic for the Onyx security service,
// including the temporary pairing listener for administrative bootstrapping.
package engine

import (
	"bufio"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"onyx/internal/crypto"
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
func StartPairingMode(token string) {
	reader := bufio.NewReader(os.Stdin)

	// Generate a temporary CA key for this pairing session
	_, caPriv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Printf("Critical error: failed to generate session key: %v\n", err)
		return
	}

	for {
		fmt.Printf("\n[PAIRING MODE ACTIVE]\n")
		fmt.Printf("Token: %s\n", token)
		fmt.Printf("Port:  2305\n")
		fmt.Printf("Window: 5 Minutes\n\n")

		resultChan := make(chan bool)
		pairingCtx, pairingCancel := context.WithTimeout(context.Background(), 5*time.Minute)

		mux := http.NewServeMux()
		mux.HandleFunc("/pair", func(w http.ResponseWriter, r *http.Request) {
			// 1. Verify the Token
			if r.Header.Get("X-Onyx-Token") != token {
				http.Error(w, "Invalid pairing token", http.StatusUnauthorized)
				return
			}

			// 2. Read the CSR from the body
			csrBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to read CSR", http.StatusBadRequest)
				return
			}

			// 3. Sign the CSR
			certPEM, err := crypto.SignCSR(csrBytes, caPriv)
			if err != nil {
				http.Error(w, fmt.Sprintf("Signing failed: %v", err), http.StatusInternalServerError)
				return
			}

			// 4. Extract identity and save the "Public Key" (the cert) for future mTLS
			cert, _ := crypto.ParseCertificate(certPEM)
			clientID := cert.Subject.CommonName

			// Ensure the auth directory exists
			authDir := "/var/lib/onyx/auth/clients/"
			os.MkdirAll(authDir, 0755)

			certPath := filepath.Join(authDir, fmt.Sprintf("%s.crt", clientID))
			if err := os.WriteFile(certPath, certPEM, 0644); err != nil {
				http.Error(w, "Failed to persist authorization", http.StatusInternalServerError)
				return
			}

			// 5. Send the signed cert back to the client
			w.Header().Set("Content-Type", "application/x-pem-file")
			w.Write(certPEM)

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
			fmt.Println("[✓] Device paired successfully. Certificate saved.")
			success = true
		case <-pairingCtx.Done():
			fmt.Println("\n[!] Pairing window expired.")
		}

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 2*time.Second)
		srv.Shutdown(shutdownCtx)
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
