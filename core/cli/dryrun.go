package cli

import (
	"fmt"
	"log"
	"ordin/m/core/cli/section"
	"ordin/m/core/rule_engine"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"

	"encoding/json"
)

func SavePlan(path string, data finalPaths) error {
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()
    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    return encoder.Encode(data)
}

type sessionPaths struct {
	cursor               int
	copyPaths            map[string][]string
	nonConflictMovePaths map[string][]string
	nonConflictDelFlag   map[string]bool
	redundanctDelPaths   map[string][]string
	ambigiousCmds        map[string][]string
	noAccessFiles        []string
	sourcePaths          map[string]string
}

func initialize(fileActions rule_engine.FilePaths) sessionPaths {
	m := sessionPaths{
		copyPaths:            make(map[string][]string),
		nonConflictMovePaths: make(map[string][]string),
		nonConflictDelFlag:   make(map[string]bool),
		redundanctDelPaths:   make(map[string][]string),
		ambigiousCmds:        make(map[string][]string),
		noAccessFiles:        make([]string, 0),
		sourcePaths:          make(map[string]string),
	}

	for key, value := range fileActions {
		moveLen := len(value.MovePaths)
		m.sourcePaths[key] = value.Srcpath
		log.Println(value.Srcpath)
		m.copyPaths[key] = value.CopyPaths
		if value.NoAccessFlag {
			m.noAccessFiles = append(m.noAccessFiles, key)
			continue
		}
		if moveLen <= 1 {
			m.nonConflictMovePaths[key] = value.MovePaths
		}
		if moveLen == 0 && value.DeleteFlag {
			m.nonConflictDelFlag[key] = true
		}
		if moveLen >= 1 && value.DeleteFlag {
			m.redundanctDelPaths[key] = value.MovePaths
		}
		if moveLen > 1 && !value.DeleteFlag {
			m.ambigiousCmds[key] = value.MovePaths
		}

	}
	return m
}

type model struct {
	sections        []section.SectionModel
	ambgiousSection section.AmbigiousSection
	deleteSection   section.DeleteSectionModel
	cursor          int
}

func dryRunInit(fileActions rule_engine.FilePaths) model {
	sessionPath := initialize(fileActions)
	return model{
		sections: []section.SectionModel{
			section.SectionInitializer(sessionPath.sourcePaths, sessionPath.copyPaths, true),
			section.SectionInitializer(sessionPath.sourcePaths, sessionPath.nonConflictMovePaths, false),
			section.SectionInitializer(sessionPath.sourcePaths, sessionPath.redundanctDelPaths, false),
		},
		ambgiousSection: section.AmbigiousSectionInitializer(sessionPath.sourcePaths, sessionPath.ambigiousCmds, false),
		deleteSection:   section.DeleteSectionInitializer(sessionPath.sourcePaths, sessionPath.nonConflictDelFlag, false),
		cursor:          0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if m.cursor < len(m.sections) && m.sections[m.cursor].Selected {

		updated, cmd :=
			m.sections[m.cursor].Update(msg)

		m.sections[m.cursor] =
			updated.(section.SectionModel)

		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			if msg.String() == "esc" {
				m.sections[m.cursor].Selected = false
			}
		}

		return m, cmd
	} else if m.cursor == len(m.sections) && m.ambgiousSection.Selected {

		updated, cmd :=
			m.ambgiousSection.Update(msg)

		m.ambgiousSection =
			updated.(section.AmbigiousSection)

		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			if msg.String() == "esc" {
				m.ambgiousSection.Selected = false
			}
		}

		return m, cmd
	} else if m.cursor == len(m.sections)+1 && m.deleteSection.Selected {

		updated, cmd :=
			m.deleteSection.Update(msg)

		m.deleteSection =
			updated.(section.DeleteSectionModel)

		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			if msg.String() == "esc" {
				m.deleteSection.Selected = false
			}
		}

		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "tab":
			if m.cursor < (len(m.sections)-1)+2 {
				m.cursor++
			} else {
				m.cursor = 0
			}
		case "c":
			if m.cursor < len(m.sections) {
				m.sections[m.cursor].Collapse = !m.sections[m.cursor].Collapse
			} else if m.cursor == len(m.sections) {
				m.ambgiousSection.Collapse = !m.ambgiousSection.Collapse
			} else if m.cursor == len(m.sections)+1 {
				m.deleteSection.Collapse = !m.deleteSection.Collapse
			}
		case "ctrl+c", "q":
			return m, tea.Quit

		case "ctrl+s", "s":
			finalPaths := m.BuildPaths()
			SavePlan("plan.json",finalPaths)
			return m, tea.Quit

		case "enter":
			if m.cursor < len(m.sections) {
				m.sections[m.cursor].Selected = true
			} else if m.cursor == len(m.sections) {
				m.ambgiousSection.Selected = true
			}
			if m.cursor == len(m.sections)+1 {
				m.deleteSection.Selected = true
			}

		case "esc":
			if m.cursor < len(m.sections) {
				m.sections[m.cursor].Selected = false
			} else if m.cursor == len(m.sections) {
				m.ambgiousSection.Selected = false
			}
			if m.cursor == len(m.sections)+1 {
				m.deleteSection.Selected = false
			}

		}

	}
	return m, nil
}

type finalPaths struct {
	Copy map[string][]string
	Move map[string][]string
	Del  []string
}

func (m model) BuildPaths() finalPaths {
	copyPaths := make(map[string][]string)
	movePaths := make(map[string][]string)
	delPaths := []string{}
	// copy paths
	for _, row := range m.sections[0].Rows {
		if row.Selected {
			copyPaths[row.SrcPath] = append(copyPaths[row.SrcPath], row.DestPath)
		}
	}
	// move paths
	for _, row := range m.sections[1].Rows {
		if row.Selected {
			movePaths[row.SrcPath] = append(movePaths[row.SrcPath], row.DestPath)
		}
	}
	// redundant paths
	for _, row := range m.sections[2].Rows {
		if row.Selected {
			movePaths[row.SrcPath] = append(movePaths[row.SrcPath], row.DestPath)
		}
	}
	// delete paths
	for _, rows := range m.deleteSection.DeleteRows {
		if rows.Selected {
			delPaths = append(delPaths, rows.SrcPath)
		}

	}
	// ambigious paths
	for _, fileSection := range m.ambgiousSection.Section {
		index := fileSection.SelectedRow
		movePaths[fileSection.Rows[index].SrcPath] = append(movePaths[fileSection.Rows[index].SrcPath], fileSection.Rows[index].DestPath)
	}
	return finalPaths{
		Copy: copyPaths,
		Move: movePaths,
		Del:  delPaths,
	}
}

func (m model) View() tea.View {
	var s strings.Builder
	for i, section := range m.sections {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		section_name := ""
		switch i {
		case 0:
			section_name = "copy"
		case 1:
			section_name = "move"
		case 2:
			section_name = "redundant deletes"
		}
		fmt.Fprintf(&s, "%s%s\n%s\n", cursor, section_name, section.View().Content)
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

func DryRun(fileActions rule_engine.FilePaths) {
	p := tea.NewProgram(dryRunInit(fileActions))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
