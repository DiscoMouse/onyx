package crypto

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

// PerformHandshake exchanges the CSR and Token for a signed certificate.
// address should be in the format "ip:port" (e.g., "10.0.0.1:2305").
func PerformHandshake(address, token string, csr []byte) ([]byte, error) {
	// 1. Prepare the request
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/pair", address), bytes.NewReader(csr))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/x-pem-file")

	// 2. Create a client that skips TLS verification (since we are bootstrapping trust)
	// or uses plain HTTP if that's how your pairing endpoint is currently set up.
	// Note: Your current pairing server uses standard HTTP for the initial handshake.
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 3. Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 4. Handle response
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server returned %s: %s", resp.Status, string(body))
	}

	return io.ReadAll(resp.Body)
}
