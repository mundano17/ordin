package section

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type AmbigiousSection struct {
	FileSections []fileSection
	cursor       int
	Collapse     bool
	Focus        bool
}

func AmbigiousSectionInitializer(SrcPaths map[string]string, DestPaths map[string][]string, defaultsel bool) AmbigiousSection {

	innerSections := []fileSection{}
	for key, val := range DestPaths {
		Rows := []row{}
		for _, DestPath := range val {
			y := row{
				RowSelected: defaultsel,
				FileName:    key,
				SrcPath:     SrcPaths[key],
				DestPath:    DestPath,
			}
			Rows = append(Rows, y)
		}
		innerSections = append(innerSections, fileSection{Rows: Rows, FileName: key})
	}

	m := AmbigiousSection{
		FileSections: innerSections,
		cursor:       0,
		Focus:        false,
	}
	return m

}

func (m AmbigiousSection) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m AmbigiousSection) moveUp() {
	if m.cursor > 0 {
		m.cursor--
	}
}

func (m AmbigiousSection) moveDown() {
	if m.cursor < len(m.FileSections)-1 {
		m.cursor++
	}
}

func (m AmbigiousSection) toggleRowSelection() {
	m.FileSections[m.cursor].FileFocus = !m.FileSections[m.cursor].FileFocus

}

func (m AmbigiousSection) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if len(m.FileSections) == 0 {
		return m, nil
	}

	if m.FileSections[m.cursor].FileFocus {
		updated, cmd := m.FileSections[m.cursor].Update(msg)
		m.FileSections[m.cursor] = updated.(fileSection)
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			if msg.String() == "esc" {
				m.FileSections[m.cursor].FileFocus = false
			}
		}
		return m, cmd
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

func (m AmbigiousSection) View() tea.View {
	var s strings.Builder
	for i, fileSection := range m.FileSections {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		section_name := fileSection.FileName
		fmt.Fprintf(&s, "%s%s\n%s\n", cursor, section_name, fileSection.View().Content)
	}
	return tea.NewView(s.String())
}
