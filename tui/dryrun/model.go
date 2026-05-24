package dryrun

import (
	"fmt"
	"strings"

	"ordin/rules"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"

	tea "charm.land/bubbletea/v2"
)

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Space  key.Binding
	Save   key.Binding
	Help   key.Binding
	Quit   key.Binding
	Escape key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Space: key.NewBinding(
		key.WithKeys("space"),
		key.WithHelp("space", "select file"),
	),
	Save: key.NewBinding(
		key.WithKeys("ctrl+s", "s"),
		key.WithHelp("s", "save/run"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help "),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "return to file selection"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Space, k.Save}, // first column
		{k.Help, k.Quit},                // second column
	}
}

func initializeModel(pathActions map[string]*rules.PathAction) model {
	if len(pathActions) == 0 {
		fmt.Println("no path actions")
		return model{}
	}
	var displayPathActions []displayPathAction
	for _, pathAction := range pathActions {
		displayPathActions = append(displayPathActions, getDisplayPathAction(*pathAction))
	}
	return model{displayPathActions: displayPathActions, cursor: 0, keys: keys, help: help.Model{}}
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
		switch {
		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.displayPathActions)-1 {
				m.cursor++
			}
		case key.Matches(msg, m.keys.Space):
			m.displayPathActions[m.cursor].focused = true
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Save):
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() tea.View {

	helpView := m.help.View(m.keys)
	s := strings.Builder{}
	if len(m.displayPathActions) > 0 {
		for i, displaypa := range m.displayPathActions {
			cursor := " "
			if i == m.cursor {
				cursor = ">>"
			}
			delStatus := "DELETE DISABLED"
			if displaypa.toDelete {
				delStatus = "DELETE ENABLED"
			}
			if i == m.cursor && m.displayPathActions[m.cursor].focused {
				helpView = m.displayPathActions[m.cursor].help.View(m.displayPathActions[m.cursor].keys)
			}
			_, _ = fmt.Fprintf(&s, "%s %s %s\n%s\n", cursor, displaypa.Srcpath, delStatus, displaypa.View().Content)
		}
		return tea.NewView(s.String() + "\n" + helpView)
	}
	return tea.NewView("No path actions available, press q to quit")
}
