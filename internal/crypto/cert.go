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
	"math/big"
	"time"
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

// SignCSR takes a raw PEM-encoded CSR and returns a signed X.509 Certificate.
// The certificate is valid for 1 year and is strictly limited to Client Authentication.
func SignCSR(csrPEM []byte, caPriv ed25519.PrivateKey) ([]byte, error) {
	block, _ := pem.Decode(csrPEM)
	if block == nil || block.Type != "CERTIFICATE REQUEST" {
		return nil, fmt.Errorf("invalid CSR PEM")
	}

	csr, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSR: %w", err)
	}

	// Create a certificate template based on the CSR
	serialNumber, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	template := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               csr.Subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0), // Valid for 1 year
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	// The server signs the certificate with its own private key
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, csr.PublicKey, caPriv)
	if err != nil {
		return nil, fmt.Errorf("failed to sign certificate: %w", err)
	}

	certBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}

	return pem.EncodeToMemory(certBlock), nil
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
