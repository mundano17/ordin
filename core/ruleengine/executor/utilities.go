package executor

import (
	"errors"
	"fmt"
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

func appendCounter(basepath string, counter int) string {
	path := basepath
	if counter != 0 {
		ext := filepath.Ext(basepath)
		path = basepath[:len(basepath)-len(ext)]
		path = fmt.Sprintf("%s_%03d%s", path, counter, ext)
	}
	return path
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
