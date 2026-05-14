package ruleengine

import (
	"errors"
	"fmt"
	"io"
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

type LogErr struct {
	Err error
}

func (e *LogErr) Error() string {
	return fmt.Sprintf("%v", e.Err)
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

func fileExecutioner(file Options, srcPath string, deleteDirPath string, operartionID *atomic.Uint64, successfulWriter io.Writer, errorWriter io.Writer, res *RunRes) error {
	fileName := filepath.Base(srcPath)
	if len(file.CopyPaths) != 0 {
		// execute copy command
		for _, baseDir := range file.CopyPaths {
			destPath, err := destPathMaker(baseDir, fileName)
			if err != nil {
				return err
			}
			err = copyCommand(srcPath, destPath)
			if err != nil {
				err = writeLogEntry(LogEntry{
					incrementer(operartionID),
					Copy,
					srcPath,
					destPath,
				}, errorWriter)
				incrementer(&res.Error)
				if err != nil {
					return &LogErr{err}
				}
				continue
			}
			incrementer(&res.Successful)
			err = writeLogEntry(LogEntry{
				incrementer(operartionID),
				Copy,
				srcPath,
				destPath,
			}, successfulWriter)
			if err != nil {
				return &LogErr{err}
			}

		}
	}

	// ignore movePaths if len > 1, which btw should be impossible or illegal for now
	if len(file.MovePaths) == 1 {
		// execute move command
		baseDir := file.MovePaths[0]
		destPath, err := destPathMaker(baseDir, fileName)
		if err != nil {
			return err
		}
		err = moveCommand(srcPath, destPath)
		if err != nil {
			err = writeLogEntry(LogEntry{
				incrementer(operartionID),
				Move,
				srcPath,
				destPath,
			}, errorWriter)
			incrementer(&res.Error)
			if err != nil {
				return &LogErr{err}
			}
			return nil
		}
		incrementer(&res.Successful)
		err = writeLogEntry(LogEntry{
			incrementer(operartionID),
			Move,
			srcPath,
			destPath,
		}, successfulWriter)
		if err != nil {
			return &LogErr{err}
		}

	}

	if file.DeleteFlag {
		destPath, err := destPathMaker(deleteDirPath, fileName)
		if err != nil {
			return err
		}
		err = delCommand(srcPath, destPath)
		if err != nil {
			err = writeLogEntry(LogEntry{
				incrementer(operartionID),
				Delete,
				srcPath,
				destPath,
			}, errorWriter)
			incrementer(&res.Error)
			if err != nil {
				return &LogErr{err}
			}
			return nil
		}
		incrementer(&res.Successful)
		err = writeLogEntry(LogEntry{
			incrementer(operartionID),
			Delete,
			srcPath,
			destPath,
		}, successfulWriter)
		if err != nil {
			return &LogErr{err}
		}

	}
	return nil
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
	defer func() {
		_ = successfulWriter.Close()
		_ = errorWriter.Close()
	}()
	var operartionID atomic.Uint64
	for srcPath, options := range allFiles {
		err := fileExecutioner(options, srcPath, deleteDirPath, &operartionID, successfulWriter, errorWriter, res)
		if _, ok := errors.AsType[*LogErr](err); ok {
			// rollback
			return nil, fmt.Errorf("run terminated: Log Error")
		}
	}
	return res, err
}
