// Package crypto provides primitives for generating and managing Ed25519
// key pairs and X.509 certificates for mTLS authentication.
package crypto

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// NewMTLSClient creates an HTTP client configured with the user's Onyx certificates.
// This client is used for all administrative communication with remote engines.
func NewMTLSClient() (*http.Client, error) {
	home, _ := os.UserHomeDir()
	certDir := filepath.Join(home, ".config", "onyx", "certs")

	certFile := filepath.Join(certDir, "client.crt")
	keyFile := filepath.Join(certDir, "client.key")

	// 1. Load the client's certificate and private key
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load client identity: %w", err)
	}

	// 2. Setup the TLS configuration
	// In a production v1, we would also verify the Server's CA here.
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true, // Temporary: Server is using a self-signed session key
	}

	// 3. Create a transport with the TLS config
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   5 * time.Second,
	}, nil
}
