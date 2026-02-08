// Package crypto provides primitives for generating and managing Ed25519
// key pairs and X.509 certificates for mTLS authentication within the Onyx ecosystem.
package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// GenerateKeyPair creates a new Ed25519 private/public key pair.
// It returns the PEM-encoded private key and the raw public key.
func GenerateKeyPair() ([]byte, ed25519.PublicKey, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate ed25519 keys: %w", err)
	}

	// Marshal the private key into PKCS#8 format
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal private key: %w", err)
	}

	// Encode the private key into PEM format for storage
	privBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privBytes,
	}

	return pem.EncodeToMemory(privBlock), pub, nil
}

// SavePEM writes PEM-encoded data to the filesystem with restrictive 0600 permissions.
// This is critical for protecting private keys on multi-user Linux systems.
func SavePEM(path string, data []byte) error {
	return os.WriteFile(path, data, 0600)
}

// LoadPrivateKey reads a PEM-encoded Ed25519 private key from the disk.
func LoadPrivateKey(path string) (ed25519.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	key, ok := priv.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not an ed25519 private key")
	}

	return key, nil
}
