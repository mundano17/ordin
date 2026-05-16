package executor

import (
	"encoding/json"
	"io"
	"os"
)

type ActionType int

const (
	Copy ActionType = 1 + iota
	Move
	Delete
)

type LogEntry struct {
	OperationID uint64
	Action      ActionType
	SrcPath     string
	DestPath    string
}

type ErrLogEntry struct {
	Action  ActionType
	SrcPath string
}

type CleanUpErrLogEntry struct {
	SrcPath string
}

func getFileWriter(path string) (*os.File, error) {
	return os.OpenFile(
		path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0o644,
	)
}

func writeLogEntry(entry LogEntry, writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(entry)
}

func writeErrorLogEntry(entry ErrLogEntry, writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(entry)
}

func writeCleanUpErrLogEntry(entry CleanUpErrLogEntry, writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(entry)
}
