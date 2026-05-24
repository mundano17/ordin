package rules

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"
)

// TODO concurrent executor
// TODO Fix the copyFileToDestination Function

func executeCopyOperation(srcPath string, destPath string, operationID uint64, w *WAL, file *os.File) (error, error) {
	err := w.logCopyOperation(srcPath, destPath, operationID, file)
	if err != nil {
		return err, nil
	}
	err = copyFileToDestinationPath(srcPath, destPath)
	if err != nil {
		return nil, err
	}
	err = w.logCommit(operationID, file)
	return err, nil
}

func executeTrashOperation(srcPath string, destPath string, operationID uint64, w *WAL, file *os.File) (error, error) {
	err := w.logTrashOperation(srcPath, destPath, operationID, file)
	if err != nil {
		return err, nil
	}
	err = moveFileToTrash(srcPath, destPath)
	if err != nil {
		return nil, err
	}
	err = w.logCommit(operationID, file)
	return err, nil
}

func (p PathAction) execute(sesPath *sessionPath, res *RunRes, opID *atomic.Uint64, file *os.File) error {
	copyErr := false
	srcPath, err := filepath.Abs(p.PathInfo.SrcPath)
	if err != nil {
		return err
	}
	for _, dir := range p.DestPaths {
		dpath := filepath.Join(dir, p.PathInfo.FileName)
		dpath, err := filepath.Abs(dpath)
		if err != nil {
			copyErr = true
			res.incrementerError()
			continue
		}
		destPath, err := giveDestPath(dpath)
		currentOperationId := opID.Add(1)
		if err != nil {
			copyErr = true
			res.incrementerError()
			continue
		}
		logError, operationError := executeCopyOperation(srcPath, destPath, currentOperationId, &sesPath.wal, file)
		if logError != nil {
			return logError
		}
		if operationError != nil {
			copyErr = true
			res.incrementerError()
			continue
		}
		res.incrementerSuccessful()
	}

	if p.ToDelete && !copyErr {
		dpath := filepath.Join(sesPath.trashpath, p.PathInfo.FileName)
		dpath, err := filepath.Abs(dpath)
		if err != nil {
			copyErr = true
			res.incrementerError()

		}
		destPath, err := giveDestPath(dpath)
		if err != nil {
			res.incrementerError()
			return nil
		}
		currentOperationId := opID.Add(1)
		logError, operationError := executeTrashOperation(srcPath, destPath, currentOperationId, &sesPath.wal, file)
		if logError != nil {
			return logError
		}
		if operationError != nil {
			copyErr = true
			res.incrementerError()
			return nil
		}
		res.incrementerSuccessful()
	}
	return nil
}

type sessionPath struct {
	trashpath string
	wal       WAL
}

func sessionInit() (sessionPath, error) {
	now := time.Now()
	timestamp := now.Format("2006-01-02_15-04-05")
	base, err := os.UserConfigDir()
	if err != nil {
		return sessionPath{}, err
	}
	dirString := filepath.Join(base, "ordin", "session", timestamp)
	err = os.MkdirAll(dirString, 0o750)
	if err != nil {
		return sessionPath{}, err
	}
	trashPath := filepath.Join(dirString, "trash")
	err = os.MkdirAll(trashPath, 0o750)
	if err != nil {
		return sessionPath{}, err
	}
	logDirPath := filepath.Join(dirString, "log")
	err = os.MkdirAll(logDirPath, 0o750)
	if err != nil {
		return sessionPath{}, err
	}
	logPath := filepath.Join(logDirPath, "log.ndjson")
	return sessionPath{trashpath: trashPath, wal: WAL{path: logPath, mu: sync.Mutex{}}}, nil
}

type RunRes struct {
	Successful atomic.Uint64
	Error      atomic.Uint64
}

func (r *RunRes) incrementerSuccessful() {
	r.Successful.Add(1)
}

func (r *RunRes) incrementerError() {
	r.Error.Add(1)
}

type Job struct {
	sesPath *sessionPath
	res     *RunRes
	opID    *atomic.Uint64
	pa      PathAction
	file    *os.File
}

func worker(ctx context.Context, jobs chan Job) error {
	for job := range jobs {
		select {
		case <-ctx.Done():
			return nil

		default:
			err := job.pa.execute(job.sesPath, job.res, job.opID, job.file)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Execute the commands for each pathAction and that's it
func Execute(pathActions []PathAction) (*RunRes, error) {
	sesPath, err := sessionInit()
	res := new(RunRes)
	opID := new(atomic.Uint64)
	jobs := make(chan Job, len(pathActions))
	g, ctx := errgroup.WithContext(context.Background())
	if err != nil {
		return res, err
	}
	file, err := getFileWriter(sesPath.wal.path)
	defer func() {
		_ = file.Close()
	}()
	if err != nil {
		return res, err
	}
	numWorkers := 6
	for i := 0; i < numWorkers; i++ {
		g.Go(
			func() error {
				return worker(ctx, jobs)
			})
	}

	for _, pathAction := range pathActions {
		newJob := Job{sesPath: &sesPath, opID: opID, res: res, pa: pathAction, file: file}
		select {
		case jobs <- newJob:
		case <-ctx.Done():
			close(jobs)
			return res, g.Wait()
		}

	}
	close(jobs)
	err = g.Wait()
	return res, err
}
