package ruleengine

import (
	"encoding/json"
	"io"
	"os"
)

type actionType int

const (
	Copy actionType = 1 + iota
	Move
	Delete
)

type LogEntry struct {
	OperationID uint64
	Action      actionType
	SrcPath     string
	DestPath    string
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
