package section

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type row struct {
	RowSelected bool
	FileName    string
	SrcPath     string
	DestPath    string
}

type SectionModel struct {
	Rows     []row
	cursor   int
	Collapse bool
	Focus    bool
}

func SectionInitializer(SrcPaths map[string]string, DestPaths map[string][]string, defaultsel bool) SectionModel {

	Rows := []row{}
	for key, val := range DestPaths {
		for _, DestPath := range val {
			y := row{
				RowSelected: defaultsel,
				FileName:    key,
				SrcPath:     SrcPaths[key],
				DestPath:    DestPath,
			}
			Rows = append(Rows, y)
		}
	}

	m := SectionModel{
		Rows:   Rows,
		cursor: 0,
		Focus:  false,
	}
	return m

}

func (m SectionModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m *SectionModel) moveUp() {
	if m.cursor > 0 {
		m.cursor--
	}
}

func (m *SectionModel) moveDown() {
	if m.cursor < len(m.Rows)-1 {
		m.cursor++
	}
}

func (m *SectionModel) toggleRowSelection() {
	m.Rows[m.cursor].RowSelected = !m.Rows[m.cursor].RowSelected

}

func (m SectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if len(m.Rows) == 0 {
		return m, nil
	}
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up", "k":
			m.moveUp()
		case "down", "j":
			m.moveDown()
		case "enter", "space":
			m.toggleRowSelection()
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m SectionModel) View() tea.View {
	var s strings.Builder
	if !m.Collapse {
		for i, row := range m.Rows {
			checkbox := "[ ]"
			cursor := " "
			if m.cursor == i && m.Focus {
				cursor = ">"
			}
			if row.RowSelected {
				checkbox = "[x]"
			}
			fmt.Fprintf(&s,
				"%s%s %-30s %-100s %-100s\n",
				cursor,
				checkbox,
				row.FileName,
				row.SrcPath,
				row.DestPath,
			)

		}
	}

	return tea.NewView(s.String())
}
