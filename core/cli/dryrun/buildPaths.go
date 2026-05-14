package dryrun

import (
	"encoding/json"
	"os"

	"ordin/m/core/ruleengine"
)

func SavePlan(path string, data ruleengine.FileOptions) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := file.Close()
		if err == nil {
			err = closeErr
		}
	}()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func (m model) BuildPaths() ruleengine.FileOptions {
	fileActions := make(ruleengine.FileOptions)

	getOrCreate := func(path string) ruleengine.Options {
		val, ok := fileActions[path]
		if !ok {
			val = ruleengine.Options{}
		}
		return val
	}

	// copy paths
	for _, row := range m.sections[0].Rows {
		if row.RowSelected {
			val := getOrCreate(row.SrcPath)
			val.CopyPaths = append(val.CopyPaths, row.DestPath)
			fileActions[row.SrcPath] = val
		}
	}

	// move paths
	for _, row := range m.sections[1].Rows {
		if row.RowSelected {
			val := getOrCreate(row.SrcPath)
			val.MovePaths = append(val.MovePaths, row.DestPath)
			fileActions[row.SrcPath] = val
		}
	}

	// redundant paths
	for _, row := range m.sections[2].Rows {
		if row.RowSelected {
			val := getOrCreate(row.SrcPath)
			val.MovePaths = append(val.MovePaths, row.DestPath)
			fileActions[row.SrcPath] = val
		}
	}

	// delete paths
	for _, row := range m.deleteSection.DeleteRows {
		if row.RowSelected {
			val := getOrCreate(row.SrcPath)
			val.DeleteFlag = true
			fileActions[row.SrcPath] = val
		}
	}

	// ambiguous paths
	for _, fileSection := range m.ambgiousSection.FileSections {
		index := fileSection.SelectedRow

		row := fileSection.Rows[index]

		val := getOrCreate(row.SrcPath)
		val.MovePaths = append(val.MovePaths, row.DestPath)
		fileActions[row.SrcPath] = val
	}

	return fileActions
}
