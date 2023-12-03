package program

import (
	"context"
	"errors"
	"fmt"
	"github.com/jwdev42/logger"
	"github.com/jwdev42/xtagger/internal/cli"
	"github.com/jwdev42/xtagger/internal/io/filesystem"
	"github.com/jwdev42/xtagger/internal/io/printer"
	"github.com/jwdev42/xtagger/internal/softerrors"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

// Program entry point called by main().
func Run() error {
	var err error
	//Parse command line
	commandLine, err = cli.ParseCommandLine()
	if err != nil {
		return fmt.Errorf("Command line error: %s", err)
	}
	//Setup printer
	printMe = printer.NewPrinter(os.Stdout)
	//Update Logger
	logger.Default().SetLevel(commandLine.FlagLogLevel())
	//Set soft error behaviour
	if commandLine.FlagQuitOnSoftError() {
		softerrors.StopOnSoftError()
	}
	//Execute command-specific branch
	switch command := commandLine.Command(); command {
	case cli.CommandTag:
		return runWithOptionalMP(createContext(true), tagFile)
	case cli.CommandPrint:
		return run(createContext(false), printFile)
	case cli.CommandUntag:
		return run(createContext(true), untagFile)
	case cli.CommandInvalidate:
		return run(createContext(true), invalidateFile)
	case cli.CommandRevalidate:
		return run(createContext(true), revalidateFile)
	default:
		return fmt.Errorf("Unknown command \"%s\"", command)
	}
	return nil
}

// Runs fileFunc multithreaded if the corresponding flag was set.
func runWithOptionalMP(opts *filesystem.Context, fileFunc filesystem.FileExaminer) error {
	if commandLine.FlagMultiThread() {
		return runMP(opts, fileFunc)
	}
	return run(opts, fileFunc)
}

// Main runner for fileFunc, singlethreaded by default, can be wrapped by runMP for multithreading.
func run(opts *filesystem.Context, fileFunc filesystem.FileExaminer) error {
	for _, path := range commandLine.Paths() {
		info, err := os.Lstat(path)
		if err != nil {
			if softerrors.Consume(err) == nil {
				continue
			}
			return err
		}
		if info.IsDir() {
			if commandLine.ForbidRecursion() {
				if err := softerrors.Errorf("Recursion is forbidden, cannot descend in directory %s", path); err == nil {
					continue
				} else {
					return err
				}
			}
			err = filesystem.WalkDir(path, opts, fileFunc)
		} else {
			err = filesystem.ExamineFile(filepath.Dir(path), info, opts, fileFunc)
		}
		if err != nil {
			if errors.Is(err, fs.SkipAll) {
				logger.Default().Debugf("run: %s", err)
				return nil
			}
			return err
		}
	}
	return nil
}

// Wrapper for run that runs fileFunc in parallel.
func runMP(opts *filesystem.Context, fileFunc filesystem.FileExaminer) error {
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
			logger.Default().Error(err)
		}
		logger.Default().Debug("runMP: Error callback goroutine exits...")
	}()
	return run(opts, wrapFileExaminer(ctx, cancel, waitForExaminers, errs, fileFunc))
}
