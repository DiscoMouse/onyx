package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle // Assignment creates a copy automatically
	noStyle      = lipgloss.NewStyle()
	helpStyle    = blurredStyle // Assignment creates a copy automatically
	cursorMode   = cursor.CursorBlink

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type FormResult struct {
	Address string
	Port    string
	Token   string
}

type formModel struct {
	focusIndex int
	inputs     []textinput.Model
	submitted  bool
	canceled   bool
}

func initialFormModel() formModel {
	m := formModel{
		inputs: make([]textinput.Model, 3),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 64

		switch i {
		case 0:
			t.Placeholder = "10.0.0.1"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "2305"
			t.CharLimit = 5
		case 2:
			t.Placeholder = "oy_..."
			t.EchoMode = textinput.EchoNormal // Set to EchoPassword if you want to mask it
		}

		m.inputs[i] = t
	}

	return m
}

func (m formModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m formModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.canceled = true
			return m, tea.Quit

		// Change Focus
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press Enter on the Submit button?
			if s == "enter" && m.focusIndex == len(m.inputs) {
				m.submitted = true
				return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
				} else {
					// Remove focused state
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = noStyle
					m.inputs[i].TextStyle = noStyle
				}
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *formModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only update inputs if they are focused
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m formModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(" PAIR NEW ENGINE "))
	b.WriteString("\n\n")

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render(" (esc to cancel)"))

	return b.String()
}

// StartPairingForm launches the input form.
// Returns the collected data and a boolean (true = submitted, false = canceled).
func StartPairingForm() (FormResult, bool) {
	p := tea.NewProgram(initialFormModel())
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Error running form: %v\n", err)
		return FormResult{}, false
	}

	finalForm := m.(formModel)
	if finalForm.canceled || !finalForm.submitted {
		return FormResult{}, false
	}

	return FormResult{
		Address: finalForm.inputs[0].Value(),
		Port:    finalForm.inputs[1].Value(),
		Token:   finalForm.inputs[2].Value(),
	}, true
}
