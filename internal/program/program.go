package program

import (
	"context"
	"fmt"
	"github.com/jwdev42/xtagger/internal/cli"
	"github.com/jwdev42/xtagger/internal/global"
	"github.com/jwdev42/xtagger/internal/io/filesystem"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

func Run() error {
	var err error
	//Parse command line
	commandLine, err = cli.ParseCommandLine()
	if err != nil {
		return fmt.Errorf("Command line error: %s", err)
	}
	//Update Logger
	global.DefaultLogger.SetLevel(commandLine.FlagLogLevel())
	//Set soft error behaviour
	if commandLine.FlagQuitOnSoftError() {
		global.StopOnSoftError()
	}
	//Execute command-specific branch
	switch command := commandLine.Command(); command {
	case cli.CommandTag:
		return run(createWalkDirOpts(true), tagFile)
	case cli.CommandPrint:
		return run(createWalkDirOpts(false), printFile)
	case cli.CommandUntag:
		return run(createWalkDirOpts(true), untagFile)
	case cli.CommandInvalidate:
		return run(createWalkDirOpts(true), invalidateFile)
	default:
		return fmt.Errorf("Unknown command \"%s\"", command)
	}
	return nil
}

func run(opts *filesystem.WalkDirOpts, fileFunc filesystem.FileExaminer) error {
	if commandLine.FlagMultithreaded() {
		return runMP(opts, fileFunc)
	}
	return runSP(opts, fileFunc)
}

func runSP(opts *filesystem.WalkDirOpts, fileFunc filesystem.FileExaminer) error {
	for _, path := range commandLine.Paths() {
		info, err := os.Lstat(path)
		if err != nil {
			if global.FilterSoftError(err) == nil {
				continue
			}
			return err
		}
		if info.IsDir() {
			if !commandLine.FlagRecursive() {
				if err := global.SoftErrorf("Recursive mode is not set and path is a directory: %s", path); err == nil {
					continue
				} else {
					return err
				}
			}
			err = filesystem.WalkDir(path, opts, fileFunc)
		} else {
			err = fileFunc(filepath.Dir(path), fs.FileInfoToDirEntry(info), opts)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func runMP(opts *filesystem.WalkDirOpts, fileFunc filesystem.FileExaminer) error {
	ctx, cancel := context.WithCancel(context.Background())
	errs := make(chan error)
	waitForErrorCollector := make(chan struct{})
	waitForExaminers := new(sync.WaitGroup)
	defer func() { <-waitForErrorCollector }()
	defer close(errs)
	defer waitForExaminers.Wait()
	go func() {
		defer close(waitForErrorCollector)
		for err := <-errs; err != nil; err = <-errs {
			global.DefaultLogger.Error(err)
		}
		global.DefaultLogger.Debug("runMP: Error callback goroutine exits...")
	}()
	return runSP(opts, wrapFileExaminer(ctx, cancel, waitForExaminers, errs, fileFunc))
}
