package section

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type row struct {
	Selected bool
	FileName string
	SrcPath  string
	DestPath string
}

type SectionModel struct {
	Rows     []row
	cursor   int
	Collapse bool
	Selected bool
}

func SectionInitializer(SrcPaths map[string]string, DestPaths map[string][]string, defaultsel bool) SectionModel {

	Rows := []row{}
	for key, val := range DestPaths {
		for _, DestPath := range val {
			y := row{
				Selected: defaultsel,
				FileName: key,
				SrcPath:  SrcPaths[key],
				DestPath: DestPath,
			}
			Rows = append(Rows, y)
		}
	}

	m := SectionModel{
		Rows:     Rows,
		cursor:   0,
		Selected: false,
	}
	return m

}

func (m SectionModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m SectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if len(m.Rows) == 0 {
		return m, nil
	}
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.Rows)-1 {
				m.cursor++
			}
		case "enter", "space":
			m.Rows[m.cursor].Selected = !m.Rows[m.cursor].Selected

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
			if m.cursor == i && m.Selected {
				cursor = ">"
			}
			if row.Selected {
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
