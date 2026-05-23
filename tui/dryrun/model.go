package dryrun

import (
	"fmt"
	"strings"

	"ordin/rules"

	tea "charm.land/bubbletea/v2"
)

func initializeModel(pathActions map[string]*rules.PathAction) model {
	if len(pathActions) == 0 {
		fmt.Println("no path actions")
		return model{}
	}
	var displayPathActions []displayPathAction
	for _, pathAction := range pathActions {
		displayPathActions = append(displayPathActions, getDisplayPathAction(*pathAction))
	}
	return model{displayPathActions: displayPathActions, cursor: 0}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.displayPathActions[m.cursor].focused {
		updated, cmd := m.displayPathActions[m.cursor].Update(msg)

		m.displayPathActions[m.cursor] = updated.(displayPathAction)

		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			if msg.String() == "esc" {
				m.displayPathActions[m.cursor].focused = false
			}
		}

		return m, cmd

	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up", "j":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "k":
			if m.cursor < len(m.displayPathActions)-1 {
				m.cursor++
			}
		case "space":
			m.displayPathActions[m.cursor].focused = true
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+s", "s":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() tea.View {
	s := strings.Builder{}
	for i, displaypa := range m.displayPathActions {
		cursor := " "
		if i == m.cursor {
			cursor = ">>"
		}
		delStatus := "DISABLED"
		if displaypa.toDelete {
			delStatus = "ENABLED"
		}
		fmt.Fprintf(&s, "%s %s %s\n%s\n", cursor, displaypa.Srcpath, delStatus, displaypa.View().Content)
	}
	return tea.NewView(s.String())
}
