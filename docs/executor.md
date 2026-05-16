# Executor

Aims to execute the plan from plan.go -> DRYRUNTUI -> savePaths

``` go
type FileOptions map[string]Options

type Options struct {
 MovePaths  []string
 CopyPaths  []string
 DeleteFlag bool
}

```

FileOptions is just { sourcePath : MovePaths[] , CopyPaths[] , true/false } and then fileExecutioner aims to copy -> move -> delete.

Each copy, move or delete is considered an operation and we log for  each operation.

``` go
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

```

The above snippet gives the structure of the log and recovery can just look at said log and aim to rollback, etc.

If log fails -> rollback and exit

Delete just moves said file to baseConfigDir/.ordin/session/data_time_stamp/.....

Log happens on the same dir as above but /logs/successful.ndjson and error.ndjson and logs a bunch of jsons in there.

Concurrency is pending, so is rollback right now. they will be added once they are started or done ig.
