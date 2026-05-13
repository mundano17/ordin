package ruleengine

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func destCheck(basepath string, counter int) (int, error) {
	path := basepath
	if counter != 0 {
		ext := filepath.Ext(basepath)
		path = basepath[:len(basepath)-len(ext)]
		path = fmt.Sprintf("%s_%03d%s", path, counter, ext)
	}

	if _, err := os.Stat(path); err == nil {
		return destCheck(basepath, counter+1)
	} else if errors.Is(err, os.ErrNotExist) {
		return counter, nil
	} else {
		return -1, err
	}
}

func copyCommand(srcPath string, destPath string) error {
	// #nosec G304 -- paths are validated by planner
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}

	defer func() {
		_ = srcFile.Close()
	}()

	// #nosec G304 -- paths are validated by planner
	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}

	defer func() {
		_ = destFile.Close()
		if err != nil {
			_ = os.Remove(destPath)
		}
	}()

	_, err = io.Copy(destFile, srcFile)
	return err
}

// NOTE: this wonderful 2 line function can not do a cross disk transfer and nor will the delete function.
func moveCommand(srcPath string, destPath string) error {
	err := os.Rename(srcPath, destPath)
	return err
}

func appendCounter(basepath string, counter int) string {
	path := basepath
	if counter != 0 {
		ext := filepath.Ext(basepath)
		path = basepath[:len(basepath)-len(ext)]
		path = fmt.Sprintf("%s_%03d%s", path, counter, ext)
	}
	return path
}

func delCommand(srcPath string, destPath string) error {
	err := os.Rename(srcPath, destPath)
	return err
}

func destPathMaker(baseDir string, fileName string) (string, error) {
	destPath := filepath.Join(baseDir, fileName)
	counter, err := destCheck(destPath, 0)
	if err != nil {
		return "", err
	}
	destPath = appendCounter(destPath, counter)
	return destPath, nil
}
