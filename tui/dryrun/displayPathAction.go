package dryrun

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type row struct {
	path     string
	selected bool
}

type rows []row

type displayPathAction struct {
	destPaths rows
	toDelete  bool // glow in red if file will be deleted from the cwd
	fileName  string
	Srcpath   string
	cursor    int
	focused   bool
}

func (m displayPathAction) Init() tea.Cmd { return nil }

func (m displayPathAction) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up", "j":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "k":
			if m.cursor < len(m.destPaths)-1 {
				m.cursor--
			}
		case "space":
			m.destPaths[m.cursor].selected = true
		case "Esc":
			m.focused = false
		}
	}
	return m, nil
}

func (m displayPathAction) View() tea.View {
	var s strings.Builder
	if m.toDelete {
		fmt.Fprintf(&s, "%s", "DELETE ENABLED")
	}
	fmt.Fprintf(&s, "%s\n", "")
	for i, row := range m.destPaths {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		checkbox := "[ ]"
		if m.destPaths[i].selected {
			checkbox = "[x]"
		}
		fmt.Fprintf(&s, "%s %s %s\n", cursor, checkbox, row.path)
	}

	return tea.NewView(s.String())
}
