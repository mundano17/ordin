package cli

import (
	"fmt"
	"log"
	"ordin/m/core/cli/section"
	"ordin/m/core/rule_engine"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
)

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
	sections []section.SectionModel
	cursor   int
}

func dryRunInit(fileActions rule_engine.FilePaths) model {
	sessionPath := initialize(fileActions)
	return model{
		sections: []section.SectionModel{
			section.SectionInitializer(sessionPath.sourcePaths, sessionPath.copyPaths, true),
			section.SectionInitializer(sessionPath.sourcePaths, sessionPath.nonConflictMovePaths, false),
			section.SectionInitializer(sessionPath.sourcePaths, sessionPath.redundanctDelPaths, false),
			section.SectionInitializer(sessionPath.sourcePaths, sessionPath.ambigiousCmds, false),
		},
		cursor: 0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if m.sections[m.cursor].Selected {

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
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "tab":
			if m.cursor < len(m.sections)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}
		case "c":
			m.sections[m.cursor].Collapse = !m.sections[m.cursor].Collapse

		case "ctrl+c", "q":
			return m, tea.Quit

		case "ctrl+s", "s":
			// need to make an check function to make sure user have checked out the correct stuff or
			// make a new section logic for mabigious moves to make sure the UX forces the user to well not choose more than one move statement per file
			// for now calling BuildPaths over the UX forced, correct shit
			m.BuildPaths()
			return m,tea.Quit

		case "enter":
			m.sections[m.cursor].Selected = true

		case "esc":
			m.sections[m.cursor].Selected = false

		}

	}
	return m, nil
}

type finalPaths struct {
	copy map[string][]string
	move map[string][]string
	del  []string
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
	// delete paths
	for _, row := range m.sections[2].Rows {
		if row.Selected {
			delPaths = append(delPaths, row.SrcPath)
		}
	}
	// redundant paths
	// ambigious paths
	return finalPaths{
		copy: copyPaths,
		move: movePaths,
		del: delPaths,
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
		case 3:
			section_name = "ambigious moves"
		case 4:
			section_name = "delete"
		}

		fmt.Fprintf(&s, "%s%s\n%s\n", cursor, section_name, section.View().Content)
	}
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
