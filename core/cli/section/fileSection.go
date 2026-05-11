package section

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type fileSection struct {
	FileName    string
	Rows        []row
	SelectedRow int
	cursor      int
	FileFocus   bool
}

func (m fileSection) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m fileSection) moveUp() {
	if m.cursor > 0 {
		m.cursor--
	}
}

func (m fileSection) moveDown() {
	if m.cursor < len(m.Rows)-1 {
		m.cursor++
	}
}

func (m fileSection) radioRowSelect() {
	m.SelectedRow = m.cursor

}

func (m fileSection) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up", "k":
			m.moveUp()
		case "down", "j":
			m.moveDown()
		case "enter", "space":
			m.radioRowSelect()
		}
	}
	return m, nil
}

func (m fileSection) View() tea.View {
	var s strings.Builder
	for i, rows := range m.Rows {
		cursor := " "
		radio := "()"
		if m.cursor == i {
			cursor = ">"
		}
		if m.SelectedRow == i {
			radio = "(*)"
		}
		fmt.Fprintf(&s, "%s %s %s %s\n", cursor, radio, rows.SrcPath, rows.DestPath)
	}
	return tea.NewView(s.String())
}
