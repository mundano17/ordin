package rules

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// If Operation has Begin but no commit -> recover pipeline
// worker pool system for recovery

func identifyLatestSession() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	var newestFile string
	var newestTime int64 = 0
	dirString := filepath.Join(base, "ordin", "session")
	files, err := os.ReadDir(dirString)
	for _, f := range files {
		fi, err := os.Stat(filepath.Join(dirString, f.Name()))
		if err != nil {
			return "", err
		}
		currTime := fi.ModTime().Unix()
		if currTime > newestTime {
			newestTime = currTime
			newestFile = filepath.Join(dirString, f.Name())
		}
	}
	return newestFile, nil
}

type BaseLog struct {
	LogType LogType
}

type x struct {
	isCommitted bool
	log         OperationLog
}

func getLogMap() (map[uint64]x, error) {
	latestSession, err := identifyLatestSession()
	if err != nil {
		return nil, err
	}
	logFilePath := filepath.Join(latestSession, "log", "log.ndjson")
	f, err := os.Open(logFilePath)

	if err != nil {
		return nil, err
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	y := map[uint64]x{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Bytes()
		var baseLog BaseLog
		err := json.Unmarshal(line, &baseLog)
		if err != nil {
			return nil, err
		}

		switch baseLog.LogType {
		case OperationBegin:
			var opLog OperationLog
			err := json.Unmarshal(line, &opLog)
			if err != nil {
				return nil, err
			}
			_, ok := y[(opLog.OperationID)]
			if !ok {
				y[(opLog.OperationID)] = x{false, opLog}
			}

		case OperationCommitted:
			var opCom OperationCommit
			err := json.Unmarshal(line, &opCom)
			if err != nil {
				return nil, err
			}

			l := y[(opCom.OperationID)]
			l.isCommitted = true
			y[(opCom.OperationID)] = l

		}
	}
	return y, nil
}

func returnFailedLogs(logs map[uint64]x) map[uint64]x {
	for opID, l := range logs {
		if l.isCommitted {
			delete(logs, opID)
		}
	}
	return logs
}

// real men don't use concurrency

func Recover() (*RunRes, error) {
	y, err := getLogMap()
	if err != nil {
		return nil, err
	}
	if len(y) == 0 {
		fmt.Println(`Recovery Not Required`)
		return nil, nil
	}
	logs := returnFailedLogs(y)
	res := new(RunRes)
	for _, l := range logs {
		srcPath := l.log.SrcPath
		destPath := l.log.DestPath
		if l.log.Action == Copy {
			err = copyFileToDestinationPath(srcPath, destPath)
			if err != nil {
				res.incrementerError()
				continue
			}
		}

		if l.log.Action == Trash {
			err = copyFileToDestinationPath(srcPath, destPath)
			if err != nil {
				res.incrementerError()
				continue
			}
			err := os.Remove(srcPath)
			if err != nil {
				res.incrementerError()
				continue
			}
		}

		res.incrementerSuccessful()
	}
	return res, nil
}
