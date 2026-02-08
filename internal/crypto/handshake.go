// Package crypto provides primitives for generating and managing Ed25519
// key pairs and X.509 certificates for mTLS authentication.
package crypto

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

// PerformHandshake sends a CSR and token to the remote Onyx engine.
// It returns the signed certificate bytes or an error if the pairing fails.
func PerformHandshake(targetIP string, token string, csrPEM []byte) ([]byte, error) {
	// Construct the URL for the pairing endpoint
	url := fmt.Sprintf("http://%s:2305/pair", targetIP)

	// We use a custom client with a reasonable timeout for the handshake
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(csrPEM))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Attach the token in a header for verification
	req.Header.Set("X-Onyx-Token", token)
	req.Header.Set("Content-Type", "application/x-pem-file")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w (is the server in pairing mode?)", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("pairing rejected (status %d): %s", resp.StatusCode, string(body))
	}

	// The server returns the signed certificate in the body
	certBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate from response: %w", err)
	}

	return certBytes, nil
}
