package section

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type deleteRow struct {
	RowSelected bool
	FileName    string
	SrcPath     string
}

type DeleteSectionModel struct {
	DeleteRows []deleteRow
	cursor     int
	Collapse   bool
	Focus      bool
}

func DeleteSectionInitializer(SrcPaths map[string]string, nonConflictDelFlag map[string]bool, defaultsel bool) DeleteSectionModel {

	DeleteRows := []deleteRow{}
	for key, val := range SrcPaths {
		if nonConflictDelFlag[key] {
			y := deleteRow{
				RowSelected: defaultsel,
				FileName:    key,
				SrcPath:     val,
			}
			DeleteRows = append(DeleteRows, y)
		}
	}

	m := DeleteSectionModel{
		DeleteRows: DeleteRows,
		cursor:     0,
		Focus:      false,
	}
	return m

}

func (m DeleteSectionModel) Init() tea.Cmd {
	// Just return `nil`, which means "no currentRowNumber/O right now, please."
	return nil
}

func (m DeleteSectionModel) moveUp() {
	if m.cursor > 0 {
		m.cursor--
	}
}

func (m DeleteSectionModel) moveDown() {
	if m.cursor < len(m.DeleteRows)-1 {
		m.cursor++
	}
}

func (m DeleteSectionModel) toggleRowSelection() {
	m.DeleteRows[m.cursor].RowSelected = !m.DeleteRows[m.cursor].RowSelected

}

func (m DeleteSectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if len(m.DeleteRows) == 0 {
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

func (m DeleteSectionModel) View() tea.View {
	var s strings.Builder
	if !m.Collapse {
		for currentRowNumber, row := range m.DeleteRows {
			checkbox := "[ ]"
			cursor := " "
			if m.cursor == currentRowNumber && m.Focus {
				cursor = ">"
			}
			if row.RowSelected {
				checkbox = "[x]"
			}
			fmt.Fprintf(&s,
				"%s%s %-30s %-100s\n",
				cursor,
				checkbox,
				row.FileName,
				row.SrcPath,
			)

		}
	}

	return tea.NewView(s.String())
}
