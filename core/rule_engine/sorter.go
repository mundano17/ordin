package rule_engine

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type fileCommands map[string]commands
type commands struct {
	moveCmds     []string
	copyCmds     []string
	delCmds      []string
	conflictFlag bool
	noAccessFlag bool
}

type filterStruct struct {
	fileAction *fileCommands
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

func dataSizeChecker(file_path string, min int, max int, equal int, enable bool) int {
	info, err := os.Stat(file_path)
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

func cmdMaker(parameter filterStruct) error {
	/*
		conflict -- 2 types:
			redundant: mv + del
			ambigious: mv, mv, mv, ....
	*/
	if parameter.moveFlag && parameter.deleteFlag || parameter.deleteFlag && parameter.copyFlag {
		return fmt.Errorf("delete and move/copy isn't allowed.")
	}

	for _, path := range parameter.filePaths {
		var res int = dataSizeChecker(path, parameter.min, parameter.max, parameter.equal, parameter.enable)
		if res == 0 {
			continue
		}
		val, ok := (*parameter.fileAction)[path]
		if !ok {
			val = commands{}
		}
		if res == -1 {
			val.noAccessFlag = true
			continue
		}
		if parameter.deleteFlag {
			k := fmt.Sprintf("rm '%s'", path)
			val.delCmds = append(val.delCmds, k)
		}

		if parameter.copyFlag {
			if !val.conflictFlag && (len(val.moveCmds) > 1 || (len(val.moveCmds) >= 1 && len(val.delCmds) >= 1)) {
				val.conflictFlag = true
			}
			k := fmt.Sprintf("cp '%s' '%s'", path, parameter.copyDest)
			val.copyCmds = append(val.copyCmds, k)
		}

		if parameter.moveFlag {
			if !val.conflictFlag && (len(val.moveCmds) > 1 || (len(val.moveCmds) >= 1 && len(val.delCmds) >= 1)) {
				val.conflictFlag = true
			}
			k := fmt.Sprintf("mv '%s' '%s'", path, parameter.moveDest)
			val.moveCmds = append(val.moveCmds, k)
		}

		(*parameter.fileAction)[path] = val
	}
	return nil
}

func plan(sortedRules []rule, path string) (fileCommands, error) {

	_, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var fileAction fileCommands
	for _, val := range sortedRules {
		// SEARCH

		searchTerm := fmt.Sprintf("%s.*%s$", val.ruleData.Name, val.ruleData.Extension)
		cmd := exec.Command("rg", "--files", "-0", "-g", searchTerm, path)
		out, err := cmd.Output()
		if err != nil {
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
			deleteFlag: (val.ruleData.Action.Delete),
			moveFlag:   (val.ruleData.Action.Move != ""),
			copyDest:   val.ruleData.Action.Copy,
			moveDest:   val.ruleData.Action.Move,
			min:        val.ruleData.Data_size.Min,
			max:        val.ruleData.Data_size.Max,
			equal:      val.ruleData.Data_size.Equal,
			enable:     val.ruleData.Data_size.Enable}

		err = cmdMaker(param)

	}
	return fileAction, nil
}
