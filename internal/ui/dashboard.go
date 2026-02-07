package ui

import (
	"fmt"
	"onyx/internal/state" // Import our new check package
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// tickMsg is sent to the Update function to trigger a re-check
type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type model struct {
	version string
	system  state.SystemState
}

func InitialModel(v string) model {
	return model{
		version: v,
		system:  state.CheckHeartbeat(),
	}
}

func (m model) Init() tea.Cmd {
	return tick() // Start the heartbeat immediately
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tickMsg:
		m.system = state.CheckHeartbeat() // Refresh state
		return m, tick()                  // Schedule next tick
	}
	return m, nil
}

func (m model) View() string {
	var s strings.Builder

	// Lipgloss styles for "Tidy" feedback
	green := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	red := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

	s.WriteString(fmt.Sprintf("  ONYX ORCHESTRATOR [%s]\n", m.version))
	s.WriteString("  ──────────────────────────────────────\n\n")

	// Dynamic status based on real state
	s.WriteString("  CONFIG:  " + renderStatus(m.system.ConfigValid, "VALID", "MISSING", green, red) + "\n")
	s.WriteString("  WAF:     " + renderStatus(m.system.WAFRulesReady, "READY", "NO RULES", green, red) + "\n")
	s.WriteString(fmt.Sprintf("  VOLUMES: [%d/3] Paths Detected\n", m.system.PathsFound))

	s.WriteString("\n  (Press 'q' to exit)\n")
	return s.String()
}

func renderStatus(val bool, pos, neg string, g, r lipgloss.Style) string {
	if val {
		return g.Render(pos)
	}
	return r.Render(neg)
}

func StartTUI(v string) error {
	p := tea.NewProgram(InitialModel(v))
	_, err := p.Run()
	return err
}
