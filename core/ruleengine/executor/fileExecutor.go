package executor

import (
	"errors"
	"io"
	"path/filepath"
	"sync/atomic"
)

func fileExecutioner(file Options, srcPath string, deleteDirPath string, operartionID *atomic.Uint64, successfulWriter io.Writer, errorWriter io.Writer, cleanUpWriter io.Writer, res *RunRes) error {
	fileName := filepath.Base(srcPath)
	if len(file.CopyPaths) != 0 {
		// execute copy command for each directory in copyPaths
		for _, baseDir := range file.CopyPaths {
			destPath, err := destPathMaker(baseDir, fileName)
			if err != nil {
				return err
			}
			err = copyCommand(srcPath, destPath)
			// if any error, log the error and increment error counter
			if err != nil {
				err = writeErrorLogEntry(ErrLogEntry{
					Delete,
					srcPath,
				}, errorWriter)
				incrementer(&res.Error)
				if err != nil {
					return &LogErr{err}
				}
				continue
			}
			// if successful, write the log entry and increment successful counter
			err = writeLogEntry(LogEntry{
				incrementer(operartionID),
				Copy,
				srcPath,
				destPath,
			}, successfulWriter)
			incrementer(&res.Successful)
			if err != nil {
				return &LogErr{err}
			}

		}
	}

	// ignore movePaths if len > 1, which btw should be impossible or illegal for now
	if len(file.MovePaths) == 1 {
		// execute move command for that one file
		baseDir := file.MovePaths[0]
		destPath, err := destPathMaker(baseDir, fileName)
		if err != nil {
			return err
		}
		err = moveCommand(srcPath, destPath)
		// if Clean Up error, log and increment success counter
		if _, ok := errors.AsType[*CleanUpErr](err); ok {
			err = writeCleanUpErrLogEntry(CleanUpErrLogEntry{SrcPath: srcPath}, cleanUpWriter)
			incrementer(&res.Successful)
			if err != nil {
				return &LogErr{err}
			}
			return nil
		}
		if err != nil {
			err = writeErrorLogEntry(ErrLogEntry{
				Delete,
				srcPath,
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
		err = moveCommand(srcPath, destPath)
		// if Clean Up error, log and increment success counter
		if _, ok := errors.AsType[*CleanUpErr](err); ok {
			err = writeCleanUpErrLogEntry(CleanUpErrLogEntry{SrcPath: srcPath}, cleanUpWriter)
			incrementer(&res.Successful)
			if err != nil {
				return &LogErr{err}
			}
			return nil
		}
		if err != nil {
			err = writeErrorLogEntry(ErrLogEntry{
				Delete,
				srcPath,
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
