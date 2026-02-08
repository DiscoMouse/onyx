// Package engine implements the core logic for the Onyx remote management service,
// including the temporary pairing listener for administrative bootstrapping.
package engine

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GeneratePairingToken creates a high-entropy, human-readable 8-character token.
// Example format: "ABCD-1234"
func GeneratePairingToken() (string, error) {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // Removed ambiguous chars like 0, O, 1, I

	out := make([]byte, 8)
	for i := range out {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		out[i] = charset[num.Int64()]
	}

	// Format as XXXX-XXXX for better readability
	return fmt.Sprintf("%s-%s", string(out[:4]), string(out[4:])), nil
}

// PairingStatus represents the result of a pairing attempt.
type PairingStatus int

const (
	StatusWaiting PairingStatus = iota
	StatusSuccess
	StatusExpired
	StatusUserCanceled
)
