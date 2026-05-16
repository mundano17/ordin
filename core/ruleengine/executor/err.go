package executor

import "fmt"

type LogErr struct {
	Err error
}

func (e *LogErr) Error() string {
	return fmt.Sprintf("%v", e.Err)
}

type CleanUpErr struct {
	Err error
}

func (e *CleanUpErr) Error() string {
	return fmt.Sprintf("%v", e.Err)
}
