package ui

import (
	"fmt"
	"onyx/internal/config"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			PaddingLeft(2).
			PaddingRight(2).
			MarginBottom(1)

	subTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A0A0A0")).
			MarginBottom(1)

	statusStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#04B575"))

	offlineStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF0000"))
)

type dashboardModel struct {
	version string
	node    *config.Node
}

// Init is called when the Bubble Tea program starts.
func (m dashboardModel) Init() tea.Cmd {
	return nil
}

// Update handles incoming messages (like keypresses).
func (m dashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the UI based on the current state.
func (m dashboardModel) View() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(titleStyle.Render(fmt.Sprintf(" ONYX SECURITY ADMIN [%s] ", m.version)))
	b.WriteString("\n")

	b.WriteString(fmt.Sprintf("  TARGET NODE: %s\n", m.node.Name))
	b.WriteString(fmt.Sprintf("  ADDRESS:     %s:%d\n", m.node.Address, m.node.Port))
	b.WriteString(subTitleStyle.Render("  ──────────────────────────────────────────"))
	b.WriteString("\n")

	// Placeholder for the actual API fetch we'll build next
	b.WriteString(fmt.Sprintf("  STATUS:      %s\n", offlineStyle.Render("[Offline - Awaiting /status endpoint implementation]")))

	b.WriteString("\n\n  (Press 'q' or 'esc' to return to menu)\n")

	return b.String()
}

// StartDashboard launches the single-node status view.
func StartDashboard(version string, node *config.Node) error {
	m := dashboardModel{
		version: version,
		node:    node,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
