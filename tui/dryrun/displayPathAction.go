package dryrun

import (
	"fmt"
	"ordin/rules"
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
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
	keys      displayPathActionKeyMap
	help      help.Model
}

type displayPathActionKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Space  key.Binding
	Help   key.Binding
	Escape key.Binding
	Delete key.Binding
}

var displayPathActionKeys = displayPathActionKeyMap{
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
		key.WithHelp("Space", "select file"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", " help "),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("Esc", "return to file selection"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete file"),
	),
}

func (k displayPathActionKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Escape}
}

func (k displayPathActionKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Space, k.Delete}, // first column
		{k.Help, k.Escape},                // second column
	}
}

func getDisplayPathAction(pathAction rules.PathAction) displayPathAction {
	return displayPathAction{
		destPaths: getRows(pathAction.DestPaths),
		toDelete:  pathAction.ToDelete,
		fileName:  pathAction.PathInfo.FileName,
		cursor: func() int {
			if len(pathAction.DestPaths) > 0 {
				return 0
			}
			return -1
		}(),
		focused: false,
		Srcpath: pathAction.PathInfo.SrcPath,
		help:    help.Model{},
		keys:    displayPathActionKeys,
	}
}

func (m displayPathAction) Init() tea.Cmd { return nil }

func (m displayPathAction) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.destPaths)-1 {
				m.cursor++
			}
		case key.Matches(msg, m.keys.Space):
			m.destPaths[m.cursor].selected = !m.destPaths[m.cursor].selected
		case key.Matches(msg, m.keys.Delete):
			m.toDelete = !m.toDelete
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Escape):
			m.focused = false
		}
	}
	return m, nil
}

func (m displayPathAction) View() tea.View {
	var s strings.Builder
	//if m.toDelete {
	//	fmt.Fprintf(&s, "%s", "DELETE ENABLED")
	//}
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
