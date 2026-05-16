package executor

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

type FileOptions map[string]Options

type Options struct {
	MovePaths  []string
	CopyPaths  []string
	DeleteFlag bool
}

func sessionInit() (string, error) {
	now := time.Now()
	timestamp := now.Format("2006-01-02_15-04-05")
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dirString := filepath.Join(base, "ordin", "session", timestamp)
	err = os.MkdirAll(dirString, 0o750)
	if err != nil {
		return "", err
	}
	deleteDirPath := filepath.Join(dirString, "delete")
	err = os.MkdirAll(deleteDirPath, 0o750)
	if err != nil {
		return "", err
	}
	logDirPath := filepath.Join(dirString, "log")
	err = os.MkdirAll(logDirPath, 0o750)
	if err != nil {
		return "", err
	}
	return timestamp, err
}

type RunRes struct {
	Successful atomic.Uint64
	Error      atomic.Uint64
}

func incrementer(a *atomic.Uint64) uint64 {
	return a.Add(1)
}

// Executor NOTE: isn't concurrent yet
func Executor(allFiles FileOptions) (*RunRes, error) {
	timestamp, err := sessionInit()
	res := &RunRes{}
	if err != nil {
		return nil, err
	}
	base, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	deleteDirPath := filepath.Join(base, "ordin", "session", timestamp, "delete")
	logFilePath := filepath.Join(base, "ordin", "session", timestamp, "log", "successful.ndjson")
	errorLogFilePath := filepath.Join(base, "ordin", "session", timestamp, "log", "error.ndjson")
	cleanUpLogFilePath := filepath.Join(base, "ordin", "session", timestamp, "log", "cleanup.ndjson")
	successfulWriter, err := getFileWriter(logFilePath)
	if err != nil {
		_ = successfulWriter.Close()
		return nil, err
	}
	errorWriter, err := getFileWriter(errorLogFilePath)
	if err != nil {
		_ = errorWriter.Close()
		return nil, err
	}
	cleanUpWriter, err := getFileWriter(cleanUpLogFilePath)
	if err != nil {
		_ = cleanUpWriter.Close()
		return nil, err
	}
	defer func() {
		_ = successfulWriter.Close()
		_ = errorWriter.Close()
		_ = cleanUpWriter.Close()
	}()
	var operartionID atomic.Uint64
	for srcPath, options := range allFiles {
		err := fileExecutioner(options, srcPath, deleteDirPath, &operartionID, successfulWriter, errorWriter, cleanUpWriter, res)
		if _, ok := errors.AsType[*LogErr](err); ok {
			// rollback to be implemented here incase of log failure
			return nil, fmt.Errorf("run terminated: Log Error")
		}
	}
	return res, err
}
