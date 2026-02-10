package ui

import (
	"fmt"
	"onyx/internal/config"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

const (
	ActionQuit    = "quit"
	ActionConnect = "connect"
	ActionPair    = "pair"
)

// item represents a list entry (either a Node or a Command).
type item struct {
	title  string
	desc   string
	node   *config.Node
	action string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type menuModel struct {
	list     list.Model
	choice   *item
	quitting bool
}

func (m menuModel) Init() tea.Cmd {
	return nil
}

func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = &i
			}
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m menuModel) View() string {
	if m.quitting {
		return ""
	}
	return "\n" + m.list.View()
}

// StartMenu displays the interactive main menu.
// It returns the selected action (connect/pair/quit) and the target node (if connect).
func StartMenu(version string, conf *config.AdminConfig) (string, *config.Node) {
	items := []list.Item{}

	// 1. Add Saved Nodes
	for i := range conf.Nodes {
		n := &conf.Nodes[i] // Use pointer to reference actual config data
		desc := fmt.Sprintf("%s:%d • Last seen: %s", n.Address, n.Port, n.LastSeen.Format("Jan 02 15:04"))
		items = append(items, item{
			title:  n.Name,
			desc:   desc,
			node:   n,
			action: ActionConnect,
		})
	}

	// 2. Add System Actions
	items = append(items, item{
		title:  "＋ Pair New Engine",
		desc:   "Connect a new Onyx instance to this console",
		action: ActionPair,
	})

	// 3. Setup List Styling
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))

	l := list.New(items, list.NewDefaultDelegate(), width, height-4)
	l.Title = fmt.Sprintf("ONYX ADMIN [%s]", version)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)

	m := menuModel{list: l}

	// 4. Run the Program
	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Println("Error running menu:", err)
		return ActionQuit, nil
	}

	finalMenu := finalModel.(menuModel)
	if finalMenu.choice != nil {
		return finalMenu.choice.action, finalMenu.choice.node
	}

	return ActionQuit, nil
}
