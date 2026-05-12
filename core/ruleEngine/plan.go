package ruleengine

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type filterStruct struct {
	fileAction *FilePaths
	filePaths  []string
	copyFlag   bool
	deleteFlag bool
	moveFlag   bool
	copyDest   string
	moveDest   string
	min        int
	max        int
	equal      int
	enable     bool
}

func dataSizeChecker(filePath string, min int, max int, equal int, enable bool) int {
	info, err := os.Stat(filePath)
	if err != nil {
		return -1
	}

	if !enable {
		return 1
	}
	size := info.Size()

	// max  = 0 , max is disabled
	// min = 0, min is disabled
	// equal = 0, equal is disabled
	if (size > int64(min) && size < int64(max) && max != 0) || (size == int64(equal) && equal != 0) || (max == 0 && size > int64(min)) {
		return 1
	} else {
		return 0
	}
}

func planMaker(parameter filterStruct) error {
	/*
		conflict -- 2 types:
			redundant: mv + del (len(val.movePaths) >= 1 && val.deleteFlag)
			ambigious: mv, mv, mv, .... len(val.movePaths) > 1
	*/
	if parameter.moveFlag && parameter.deleteFlag || parameter.deleteFlag && parameter.copyFlag {
		return fmt.Errorf("delete and move/copy isn't allowed")
	}

	for _, path := range parameter.filePaths {
		res := dataSizeChecker(path, parameter.min, parameter.max, parameter.equal, parameter.enable)
		if res == 0 {
			continue
		}
		val, ok := (*parameter.fileAction)[path]
		if !ok {
			val = Paths{}
			val.Srcpath = path
		}
		if res == -1 {
			val.NoAccessFlag = true
			continue
		}
		if parameter.deleteFlag {
			val.DeleteFlag = true
		}

		if parameter.copyFlag {
			val.CopyPaths = append(val.CopyPaths, parameter.copyDest)
		}

		if parameter.moveFlag {
			val.MovePaths = append(val.MovePaths, parameter.moveDest)
		}

		(*parameter.fileAction)[path] = val
	}
	return nil
}

func Plan(sortedRules []rule, path string) (FilePaths, error) {
	_, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("dir path read failed")
	}
	fileAction := make(FilePaths)
	for _, val := range sortedRules {
		// SEARCH

		searchTerm := fmt.Sprintf("%s.*%s$", val.ruleData.Name, val.ruleData.Extension)
		cmd := exec.Command("fd", "-0", searchTerm, path)

		out, err := cmd.Output()
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				if exitError.ExitCode() == 1 {
					return nil, nil
				}
			}
			return nil, err
		}

		filePaths := strings.Split(string(out), "\x00")

		if len(filePaths) > 0 && filePaths[len(filePaths)-1] == "" {
			filePaths = filePaths[:len(filePaths)-1]
		}

		// ACTION + FILTERING

		param := filterStruct{
			fileAction: &fileAction,
			filePaths:  filePaths,
			copyFlag:   (val.ruleData.Action.Copy != ""),
			deleteFlag: val.ruleData.Action.Delete,
			moveFlag:   (val.ruleData.Action.Move != ""),
			copyDest:   val.ruleData.Action.Copy,
			moveDest:   val.ruleData.Action.Move,
			min:        val.ruleData.DataSize.Min,
			max:        val.ruleData.DataSize.Max,
			equal:      val.ruleData.DataSize.Equal,
			enable:     val.ruleData.DataSize.Enable,
		}

		err = planMaker(param)
		if err != nil {
			return nil, err
		}

	}
	return fileAction, nil
}
