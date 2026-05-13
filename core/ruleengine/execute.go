package ruleengine

import (
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

func nextOperationID(id *atomic.Uint64) uint64 {
	return id.Add(1)
}

func fileExecutioner(file Options, srcPath string, deleteDirPath string, operartionID *atomic.Uint64, writer io.Writer) error {
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
				return err
			}
			err = writeLogEntry(LogEntry{
				nextOperationID(operartionID),
				Copy,
				srcPath,
				destPath,
			}, writer)
			if err != nil {
				return err
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
			return err
		}
		err = writeLogEntry(LogEntry{
			nextOperationID(operartionID),
			Move,
			srcPath,
			destPath,
		}, writer)
		if err != nil {
			return err
		}

	}

	if file.DeleteFlag {
		destPath, err := destPathMaker(deleteDirPath, fileName)
		if err != nil {
			return err
		}
		err = delCommand(srcPath, destPath)
		if err != nil {
			return err
		}
		err = writeLogEntry(LogEntry{
			nextOperationID(operartionID),
			Delete,
			srcPath,
			destPath,
		}, writer)
		if err != nil {
			return err
		}

	}
	return nil
}

func Executior(allFiles FileOptions) error {
	timestamp, err := sessionInit()
	if err != nil {
		return err
	}
	base, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	deleteDirPath := filepath.Join(base, "ordin", "session", timestamp, "delete")
	logFilePath := filepath.Join(base, "ordin", "session", timestamp, "log", "journal.ndjson")
	writer, err := getFileWriter(logFilePath)
	if err != nil {
		_ = writer.Close()
		return err
	}
	var operartionID atomic.Uint64
	for srcPath, options := range allFiles {
		err := fileExecutioner(options, srcPath, deleteDirPath, &operartionID, writer)
		if err != nil {
			return err
		}
	}
	return nil
}
