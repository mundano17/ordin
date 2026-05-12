package dryrun

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

func (m model) View() tea.View {
	var s strings.Builder
	for i, section := range m.sections {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		sectionName := ""
		switch i {
		case 0:
			sectionName = "copy"
		case 1:
			sectionName = "move"
		case 2:
			sectionName = "redundant deletes"
		}
		fmt.Fprintf(&s, "%s%s\n%s\n", cursor, sectionName, section.View().Content)
	}
	// ambigious section View
	cursor := " "
	if m.cursor == len(m.sections) {
		cursor = ">"
	}
	fmt.Fprintf(&s, "%s%s\n%s\n", cursor, "Ambigious Moves", m.ambgiousSection.View().Content)
	cursor = " "
	if m.cursor == len(m.sections)+1 {
		cursor = ">"
	}
	fmt.Fprintf(&s, "%s%s\n%s\n", cursor, "Delete Section", m.deleteSection.View().Content)

	s.WriteString("\nPress q to quit.\n")
	return tea.NewView(s.String())
}
