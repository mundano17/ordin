package executor

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"syscall"
)

func copyCommand(srcPath string, destPath string) (err error) {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destDir := filepath.Dir(destPath)
	baseName := filepath.Base(destPath)

	// make a temp file
	tempFile, err := os.CreateTemp(destDir, baseName+".tmp-*")
	if err != nil {
		return err
	}

	tempPath := tempFile.Name()

	defer func() {
		_ = tempFile.Close()

		if err != nil {
			_ = os.Remove(tempPath)
		}
	}()

	// copy file to tempFile , if err -> copy has failed.
	_, err = io.Copy(tempFile, srcFile)
	if err != nil {
		return err
	}

	// flush failed
	err = tempFile.Sync()
	if err != nil {
		return err
	}

	err = os.Rename(tempPath, destPath)
	if err != nil {
		// renaming tempFile to destFile has failed. Essentially every err in copy function leads to failure.
		return err
	}

	return nil
}

func moveCommand(srcPath string, destPath string) error {
	err := os.Rename(srcPath, destPath)

	if errors.Is(err, syscall.EXDEV) {
		err = copyCommand(srcPath, destPath)
		if err != nil {
			// if err here, well srcFile still exists, just clean up the destPath
			err = os.Remove(destPath)
			return err
		}
		err = os.Remove(srcPath)
		if err != nil {
			// if removing the srcPath is well not working, the destPath already has the srcFile
			// essentially cleanup failed, which is fine compared to losing the actual data.
			return &CleanUpErr{err}
		}
	}

	return nil
}
