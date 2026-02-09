// Package ui implements the terminal user interface for the Onyx Security Appliance.
// It uses the Bubble Tea framework to provide a real-time monitoring dashboard.
package ui

import (
	"fmt"
	"strings"
	"time"

	"onyx/internal/config"
	"onyx/internal/state"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// tickMsg is an internal signal used to trigger periodic UI refreshes.
type tickMsg time.Time

// tick creates a command that sends a tickMsg every second.
func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// model maintains the internal state of the administration dashboard.
type model struct {
	version string
	system  state.SystemState
	config  *config.AdminConfig // Updated from config.Config
}

// InitialModel prepares the starting state for the TUI, including loaded server configurations.
func InitialModel(v string, cfg *config.AdminConfig) model { // Updated from config.Config
	return model{
		version: v,
		system:  state.CheckHeartbeat(),
		config:  cfg,
	}
}

// Init defines the initial command for the TUI program.
func (m model) Init() tea.Cmd {
	return tick()
}

// Update processes incoming messages and updates the model's state.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tickMsg:
		m.system = state.CheckHeartbeat()
		return m, tick()
	}
	return m, nil
}

// View renders the current state of the dashboard into a string for terminal display.
func (m model) View() string {
	var s strings.Builder

	green := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	red := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

	s.WriteString(fmt.Sprintf("  ONYX SECURITY ADMIN [%s]\n", m.version))
	s.WriteString("  ──────────────────────────────────────────\n\n")

	// Display local engine status
	s.WriteString("  LOCAL ENGINE: " + renderStatus(m.system.ProxyActive, "RUNNING", "STOPPED", green, red) + "\n")
	s.WriteString("  LOCAL CONFIG: " + renderStatus(m.system.ConfigValid, "VALID", "MISSING", green, red) + "\n")

	// Display remote server count from TOML
	remoteCount := 0
	if m.config != nil {
		remoteCount = len(m.config.Servers)
	}
	s.WriteString(fmt.Sprintf("  REMOTE NODES: [%d] Configured in TOML\n", remoteCount))

	s.WriteString("\n  (Press 'q' to exit)\n")
	return s.String()
}

// renderStatus applies color and text based on a boolean condition.
func renderStatus(val bool, pos, neg string, g, r lipgloss.Style) string {
	if val {
		return g.Render(pos)
	}
	return r.Render(neg)
}

// StartTUI initializes and launches the main Bubble Tea program loop with configuration data.
func StartTUI(v string, cfg *config.AdminConfig) error { // Updated from config.Config
	p := tea.NewProgram(InitialModel(v, cfg))
	_, err := p.Run()
	return err
}
