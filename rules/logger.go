package rules

import (
	"encoding/json"
	"os"
)

type ActionType int

const (
	Copy ActionType = iota
	Trash
)

type LogType int

const (
	OperationBegin LogType = iota
	OperationCommitted
)

type OperationLog struct {
	LogType     LogType
	OperationID uint64
	SrcPath     string
	DestPath    string
	Action      ActionType
}

type OperationCommit struct {
	LogType     LogType
	operationID uint64
}

func getFileWriter(path string) (*os.File, error) {
	return os.OpenFile(
		path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0o644,
	)
}

func (entry OperationLog) SaveLog(path string) error {
	file, err := getFileWriter(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(entry)
	if err != nil {
		return err
	}
	return file.Sync()
}

func (entry OperationCommit) SaveLog(path string) error {
	file, err := getFileWriter(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(entry)
	if err != nil {
		return err
	}
	return file.Sync()

}
