package cli

import (
	"fmt"
	"ordin/m/core/cli/section"
	"ordin/m/core/rule_engine"
	"os"

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

// TODO figure out how to display multiple sections, starting with 2
// TODO figure out how to get cli to run from main.go
// TODO figure out v1

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
	scroll   bool
}

func dryRunInit(fileActions rule_engine.FilePaths) model {
	sessionPath := initialize(fileActions)
	return model{
		sections: []section.SectionModel{
			section.SectionInitializer(sessionPath.sourcePaths, sessionPath.copyPaths),
			section.SectionInitializer(sessionPath.sourcePaths, sessionPath.nonConflictMovePaths),
		},
		scroll: true,
		cursor: 0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "tab":
			if !m.scroll && m.cursor > 0 {
				m.cursor--
				if m.cursor == 0 {
					m.scroll = true
				}
			}
			if m.scroll && m.cursor < len(m.sections)-1 {
				m.cursor++
				m.scroll = true
				if m.cursor == len(m.sections) {
					m.scroll = false
				}
			}
		case "c":
			m.sections[m.cursor].Collapse = !m.sections[m.cursor].Collapse
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	s := ""
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
		}

		s += fmt.Sprintf("\n%s %s\n %s\n", cursor, section_name, section.View().Content)
	}
	s += "\nPress q to quit.\n"
	return tea.NewView(s)
}

func DryRun(fileActions rule_engine.FilePaths) {
	p := tea.NewProgram(dryRunInit(fileActions))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
