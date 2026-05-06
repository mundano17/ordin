package section

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

type row struct {
	selected bool
	fileName string
	srcPath  string
	destPath string
}

type SectionModel struct {
	rows     []row
	cursor   int
	Collapse bool
}

func SectionInitializer(srcPaths map[string]string, destPaths map[string][]string) SectionModel {

	rows := []row{}
	for key, val := range destPaths {
		for _, destPath := range val {
			y := row{
				selected: false,
				fileName: key,
				srcPath:  srcPaths[key],
				destPath: destPath,
			}
			rows = append(rows, y)
		}
	}

	m := SectionModel{
		rows:   rows,
		cursor: 0,
	}
	return m

}

func (m SectionModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m SectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.rows)-1 {
				m.cursor++
			}
		case "enter", "space":
			m.rows[m.cursor].selected = !m.rows[m.cursor].selected
		}
	}
	return m, nil
}

func (m SectionModel) View() tea.View {
	s := " "
	if !m.Collapse {
		for i, row := range m.rows {
			checkbox := "[ ]"
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			if row.selected {
				checkbox = "[x]"
			}
			s += fmt.Sprintf(
				"%s %s %-10s %-10s %-10s\n",
				cursor,
				checkbox,
				row.fileName,
				row.srcPath,
				row.destPath,
			)

		}
	}

	return tea.NewView(s)
}
