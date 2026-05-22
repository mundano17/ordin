// Package dryrun TUI stuff
package dryrun

import (
	"fmt"
	"os"

	"ordin/rules"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	displayPathActions []displayPathAction
	cursor             int
}

func getRows(destPaths []string) []row {
	rows := []row{}
	for _, path := range destPaths {
		rows = append(rows, row{path: path, selected: false})
	}
	return rows
}

func getDisplayPathAction(pathAction rules.PathAction) displayPathAction {
	return displayPathAction{
		destPaths: getRows(pathAction.DestPaths),
		toDelete:  pathAction.ToDelete,
		fileName:  pathAction.PathInfo.FileName,
		cursor:    0,
		focused:   false,
		Srcpath:   pathAction.PathInfo.SrcPath,
	}
}

func (r rows) returnChosenRows() []string {
	destPaths := []string{}
	for _, row := range r {
		if row.selected {
			destPaths = append(destPaths, row.path)
		}
	}
	return destPaths
}

func (m model) returnChosenPathActions() []rules.PathAction {
	chosenPathActions := []rules.PathAction{}
	for _, pa := range m.displayPathActions {
		x := rules.PathInfo{SrcPath: pa.Srcpath, IsSymLink: false, FileName: pa.fileName}
		y := rules.PathAction{DestPaths: pa.destPaths.returnChosenRows(), ToDelete: pa.toDelete, PathInfo: x}
		chosenPathActions = append(chosenPathActions, y)
	}
	return chosenPathActions
}

// InitializeDryRun - creates and displays the dry run program and then returns the chosen paths
func InitializeDryRun(pathActions map[string]*rules.PathAction) []rules.PathAction {
	p := tea.NewProgram(initializeModel(pathActions))
	tmodel, err := p.Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	m, ok := tmodel.(model)
	if !ok {
		return nil
	}
	return m.returnChosenPathActions()
}
