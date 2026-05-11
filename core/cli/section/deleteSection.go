package section

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type deleteRows struct {
	Selected bool
	FileName string
	SrcPath  string
}

type DeleteSectionModel struct {
	DeleteRows []deleteRows
	cursor     int
	Collapse   bool
	Selected   bool
}

func DeleteSectionInitializer(SrcPaths map[string]string, nonConflictDelFlag map[string]bool, defaultsel bool) DeleteSectionModel {

	DeleteRows := []deleteRows{}
	for key, val := range SrcPaths {
		if nonConflictDelFlag[key] {
			y := deleteRows{
				Selected: defaultsel,
				FileName: key,
				SrcPath:  val,
			}
			DeleteRows = append(DeleteRows, y)
		}
	}

	m := DeleteSectionModel{
		DeleteRows: DeleteRows,
		cursor:     0,
		Selected:   false,
	}
	return m

}

func (m DeleteSectionModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m DeleteSectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if len(m.DeleteRows) == 0 {
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
			if m.cursor < len(m.DeleteRows)-1 {
				m.cursor++
			}
		case "enter", "space":
			m.DeleteRows[m.cursor].Selected = !m.DeleteRows[m.cursor].Selected

		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m DeleteSectionModel) View() tea.View {
	var s strings.Builder
	if !m.Collapse {
		for i, deleteRows := range m.DeleteRows {
			checkbox := "[ ]"
			cursor := " "
			if m.cursor == i && m.Selected {
				cursor = ">"
			}
			if deleteRows.Selected {
				checkbox = "[x]"
			}
			fmt.Fprintf(&s,
				"%s%s %-30s %-100s\n",
				cursor,
				checkbox,
				deleteRows.FileName,
				deleteRows.SrcPath,
			)

		}
	}

	return tea.NewView(s.String())
}
