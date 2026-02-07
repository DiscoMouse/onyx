package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Style definitions using Lip Gloss
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00ADD8")).
			MarginLeft(2)
	statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))
)

type model struct {
	version string
	proxy   string // State awareness: e.g., "Caddy v2.7.6"
	waf     string // State awareness: e.g., "Coraza + OWASP CRS"
}

func InitialModel(v string) model {
	return model{
		version: v,
		proxy:   "Caddy (Active)",
		waf:     "Coraza (Loaded)",
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render(fmt.Sprintf("ONYX ORCHESTRATOR [%s]", m.version)) + "\n\n")
	s.WriteString(fmt.Sprintf("  Proxy: %s\n", statusStyle.Render(m.proxy)))
	s.WriteString(fmt.Sprintf("  WAF:   %s\n", statusStyle.Render(m.waf)))
	s.WriteString("\n  (Press 'q' to exit)\n")

	return s.String()
}

func StartTUI(v string) error {
	p := tea.NewProgram(InitialModel(v))
	_, err := p.Run()
	return err
}
