package state

import (
	"os"
)

// SystemState holds the "Health" of your perimeters
type SystemState struct {
	ProxyActive   bool
	WAFRulesReady bool
	ConfigValid   bool
	PathsFound    int
}

// CheckHeartbeat runs a sanity check on the expected environment
func CheckHeartbeat() SystemState {
	s := SystemState{}

	// 1. Check if the Caddyfile exists (Config Sanity)
	if _, err := os.Stat("/etc/onyx/Caddyfile"); err == nil {
		s.ConfigValid = true
	}

	// 2. Check for Coraza rules (Signature Sanity)
	// For now, we look for a placeholder rule file
	if _, err := os.Stat("/var/lib/onyx/rules/crs.conf"); err == nil {
		s.WAFRulesReady = true
	}

	// 3. Simple path count for visual feedback
	paths := []string{"/etc/onyx", "/var/lib/onyx", "/var/log/onyx"}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			s.PathsFound++
		}
	}

	return s
}
