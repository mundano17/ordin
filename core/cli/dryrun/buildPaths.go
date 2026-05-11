package dryrun

import (
	"encoding/json"
	"os"
)

type finalPaths struct {
	Copy map[string][]string
	Move map[string][]string
	Del  []string
}

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

func (m model) BuildPaths() finalPaths {
	copyPaths := make(map[string][]string)
	movePaths := make(map[string][]string)
	delPaths := []string{}
	// copy paths
	for _, row := range m.sections[0].Rows {
		if row.RowSelected {
			copyPaths[row.SrcPath] = append(copyPaths[row.SrcPath], row.DestPath)
		}
	}
	// move paths
	for _, row := range m.sections[1].Rows {
		if row.RowSelected {
			movePaths[row.SrcPath] = append(movePaths[row.SrcPath], row.DestPath)
		}
	}
	// redundant paths
	for _, row := range m.sections[2].Rows {
		if row.RowSelected {
			movePaths[row.SrcPath] = append(movePaths[row.SrcPath], row.DestPath)
		}
	}
	// delete paths
	for _, rows := range m.deleteSection.DeleteRows {
		if rows.RowSelected {
			delPaths = append(delPaths, rows.SrcPath)
		}

	}
	// ambigious paths
	for _, fileSection := range m.ambgiousSection.FileSections {
		index := fileSection.SelectedRow
		movePaths[fileSection.Rows[index].SrcPath] = append(movePaths[fileSection.Rows[index].SrcPath], fileSection.Rows[index].DestPath)
	}
	return finalPaths{
		Copy: copyPaths,
		Move: movePaths,
		Del:  delPaths,
	}
}
