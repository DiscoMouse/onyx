// Package crypto provides primitives for generating and managing Ed25519
// key pairs and X.509 certificates for mTLS authentication within the Onyx ecosystem.
package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
)

// GenerateCSR creates a Certificate Signing Request (CSR) for a client.
// This is sent to the server during the pairing process to request a signed certificate.
func GenerateCSR(priv ed25519.PrivateKey, commonName string) ([]byte, error) {
	subj := pkix.Name{
		CommonName:   commonName,
		Organization: []string{"Onyx Admin"},
	}

	template := x509.CertificateRequest{
		Subject:            subj,
		SignatureAlgorithm: x509.PureEd25519,
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, priv)
	if err != nil {
		return nil, fmt.Errorf("failed to create CSR: %w", err)
	}

	csrBlock := &pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrBytes,
	}

	return pem.EncodeToMemory(csrBlock), nil
}

// ParseCertificate converts PEM-encoded certificate data into an x509 Certificate object.
func ParseCertificate(certPEM []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(certPEM)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("failed to decode PEM block containing certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	return cert, nil
}
