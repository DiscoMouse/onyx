// Package crypto provides primitives for generating and managing Ed25519
// key pairs and X.509 certificates for mTLS authentication.
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
		return nil, nil, fmt.Errorf("failed to generate ed25519 key: %w", err)
	}

	privPEM, err := EncodePrivateKey(priv)
	if err != nil {
		return nil, nil, err
	}

	return privPEM, pub, nil
}

// EncodePrivateKey converts an Ed25519 private key into a PEM-encoded PKCS#8 block.
func EncodePrivateKey(priv ed25519.PrivateKey) ([]byte, error) {
	pkcs8Key, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal private key: %w", err)
	}

	pemBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pkcs8Key,
	}

	return pem.EncodeToMemory(pemBlock), nil
}

// SavePEM writes PEM-encoded data to a file with strict 0600 permissions.
func SavePEM(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0600)
}

// LoadPrivateKey reads a PEM-encoded Ed25519 private key from disk.
func LoadPrivateKey(filename string) (ed25519.PrivateKey, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	switch k := key.(type) {
	case ed25519.PrivateKey:
		return k, nil
	default:
		return nil, fmt.Errorf("not an ed25519 private key")
	}
}
