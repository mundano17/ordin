package dryrun

import (
	"fmt"
	"log"
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

type model struct {
	sections        []section.SectionModel
	ambgiousSection section.AmbigiousSection
	deleteSection   section.DeleteSectionModel
	cursor          int
}

func getSessionPaths(fileActions rule_engine.FilePaths) sessionPaths {
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

func DryRunTUIInit(fileActions rule_engine.FilePaths) {
	p := tea.NewProgram(dryRunInit(fileActions))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func dryRunInit(fileActions rule_engine.FilePaths) model {
	sessionPath := getSessionPaths(fileActions)
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
