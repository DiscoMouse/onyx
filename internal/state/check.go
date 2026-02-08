// Package state provides functions to monitor the operational health
// of the Onyx Security Engine and its environment.
package state

import (
	"os"
	"os/exec"
	"strings"
)

// SystemState holds the current health and status of the security perimeters.
type SystemState struct {
	ProxyActive   bool
	WAFRulesReady bool
	ConfigValid   bool
	PathsFound    int
}

// CheckHeartbeat runs a diagnostic check on the system to verify the
// environment and the status of the Onyx service.
func CheckHeartbeat() SystemState {
	s := SystemState{}

	// 1. Verify if the Caddyfile configuration exists
	if _, err := os.Stat("/etc/onyx/Caddyfile"); err == nil {
		s.ConfigValid = true
	}

	// 2. Check for the presence of Coraza WAF rules
	if _, err := os.Stat("/var/lib/onyx/rules/crs.conf"); err == nil {
		s.WAFRulesReady = true
	}

	// 3. Query systemd to check if the onyx service is currently active
	cmd := exec.Command("systemctl", "is-active", "onyx")
	out, _ := cmd.Output()
	if strings.TrimSpace(string(out)) == "active" {
		s.ProxyActive = true
	}

	// 4. Verify existence of critical system volumes and log paths
	paths := []string{"/etc/onyx", "/var/lib/onyx", "/var/log/onyx"}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			s.PathsFound++
		}
	}

	return s
}
