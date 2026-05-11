package section

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

/*
type row struct {

	Selected bool
	FileName string
	SrcPath  string
	DestPath string

}
*/

type innerSection struct {
	FileName    string
	Rows        []row
	SelectedRow int
	cursor      int
	Selected    bool
}

type AmbigiousSection struct {
	Section  []innerSection
	cursor   int
	Collapse bool
	Selected bool
}

func AmbigiousSectionInitializer(SrcPaths map[string]string, DestPaths map[string][]string, defaultsel bool) AmbigiousSection {

	innerSections := []innerSection{}
	for key, val := range DestPaths {
		Rows := []row{}
		for _, DestPath := range val {
			y := row{
				Selected: defaultsel,
				FileName: key,
				SrcPath:  SrcPaths[key],
				DestPath: DestPath,
			}
			Rows = append(Rows, y)
		}
		innerSections = append(innerSections, innerSection{Rows: Rows, FileName: key})
	}

	m := AmbigiousSection{
		Section:  innerSections,
		cursor:   0,
		Selected: false,
	}
	return m

}

func (m AmbigiousSection) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m innerSection) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m AmbigiousSection) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if len(m.Section) == 0 {
		return m, nil
	}

	if m.Section[m.cursor].Selected {
		updated, cmd := m.Section[m.cursor].Update(msg)
		m.Section[m.cursor] = updated.(innerSection)
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			if msg.String() == "esc" {
				m.Section[m.cursor].Selected = false
			}
		}
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.Section)-1 {
				m.cursor++
			}
		case "enter", "space":
			m.Section[m.cursor].Selected = !m.Section[m.cursor].Selected

		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m innerSection) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			m.SelectedRow = m.cursor
		}
	}
	return m, nil
}

func (m AmbigiousSection) View() tea.View {
	var s strings.Builder
	for i, section := range m.Section {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		section_name := section.FileName
		fmt.Fprintf(&s, "%s%s\n%s\n", cursor, section_name, section.View().Content)
	}
	return tea.NewView(s.String())
}

func (m innerSection) View() tea.View {
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
