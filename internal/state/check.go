package state

import (
	"os"
	"os/exec"
	"strings"
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

	// 1. Check if the Caddyfile exists
	if _, err := os.Stat("/etc/onyx/Caddyfile"); err == nil {
		s.ConfigValid = true
	}

	// 2. Check for Coraza rules
	if _, err := os.Stat("/var/lib/onyx/rules/crs.conf"); err == nil {
		s.WAFRulesReady = true
	}

	// 3. Check Systemd Service Status
	cmd := exec.Command("systemctl", "is-active", "onyx")
	out, _ := cmd.Output()
	if strings.TrimSpace(string(out)) == "active" {
		s.ProxyActive = true
	}

	// 4. Volume Detect
	paths := []string{"/etc/onyx", "/var/lib/onyx", "/var/log/onyx"}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			s.PathsFound++
		}
	}

	return s
}
