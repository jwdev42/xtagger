package program

import (
	"fmt"
	"github.com/jwdev42/xtagger/internal/cli"
	"github.com/jwdev42/xtagger/internal/global"
	"github.com/jwdev42/xtagger/internal/io/filesystem"
	"io/fs"
	"os"
	"path/filepath"
)

func Run() error {
	//Parse command line
	cmdline, err := cli.ParseCommandLine()
	if err != nil {
		return fmt.Errorf("Command line error: %s", err)
	}
	global.CommandLine = cmdline
	//Update Logger
	global.DefaultLogger.SetLevel(cmdline.FlagLogLevel())
	//Set soft error behaviour
	if cmdline.FlagQuitOnSoftError() {
		global.StopOnSoftError()
	}
	//Execute command-specific branch
	switch command := cmdline.Command(); command {
	case cli.CommandTag:
		return run(createWalkDirOpts(true), tagFile)
	case cli.CommandPrint:
		return run(createWalkDirOpts(false), printFile)
	case cli.CommandUntag:
		return run(createWalkDirOpts(true), untagFile)
	default:
		return fmt.Errorf("Unknown command \"%s\"", command)
	}
	return nil
}

func run(opts *filesystem.WalkDirOpts, fileFunc filesystem.FileExaminer) error {
	for _, path := range global.CommandLine.Paths() {
		info, err := os.Lstat(path)
		if err != nil {
			if global.FilterSoftError(err) == nil {
				continue
			}
			return err
		}
		if info.IsDir() {
			if !global.CommandLine.FlagRecursive() {
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
