package rules

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

func giveDestPath(basepath string) (string, error) {
	counter := 0
	for {
		path := appendCounter(basepath, counter)
		f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
		if err == nil {
			_ = f.Close()
			return path, nil
		}
		if errors.Is(err, os.ErrExist) {
			counter++
			continue
		}
		return "", err
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

// TODO : Read more about the ways to make a file transaction more atomic and what each part does in the code and fix the below function.
func copyFileToDestinationPath(srcPath string, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer func() { _ = srcFile.Close() }()

	destDir := filepath.Dir(destPath)
	baseName := filepath.Base(destPath)

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
	_, err = io.Copy(tempFile, srcFile)
	if err != nil {
		return err
	}

	err = tempFile.Sync()
	if err != nil {
		return err
	}

	err = os.Rename(tempPath, destPath)
	if err != nil {
		return err
	}
	return nil
}

func moveFileToTrash(srcPath string, trashpath string) error {
	err := copyFileToDestinationPath(srcPath, trashpath)
	if err != nil {
		return err
	}
	err = os.Remove(srcPath)
	if err != nil {
		return err
	}
	return nil
}

type WAL struct {
	mu   sync.Mutex
	path string
}

func (w *WAL) logCommit(operationID uint64, file *os.File) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	commitLog := OperationCommit{LogType: OperationCommitted, OperationID: operationID}
	return commitLog.SaveLog(w.path, file)
}

func (w *WAL) logCopyOperation(srcPath string, destPath string, operationID uint64, file *os.File) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	oplog := OperationLog{
		LogType:     OperationBegin,
		OperationID: operationID,
		SrcPath:     srcPath,
		DestPath:    destPath,
		Action:      Copy,
	}
	return oplog.SaveLog(w.path, file)
}

func (w *WAL) logTrashOperation(srcPath string, destPath string, operationID uint64, file *os.File) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	oplog := OperationLog{
		LogType:     OperationBegin,
		OperationID: operationID,
		SrcPath:     srcPath,
		DestPath:    destPath,
		Action:      Trash,
	}
	return oplog.SaveLog(w.path, file)
}
