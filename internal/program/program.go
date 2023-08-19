package program

import (
	"fmt"
	"github.com/jwdev42/xtagger/internal/cli"
	"io/fs"
	"os"
)

func Run() error {
	//Parse command line
	cmdline, err := cli.ParseCommandLine()
	if err != nil {
		return fmt.Errorf("Command line error: %s", err)
	}
	switch command := cmdline.Command(); command {
	case cli.CommandCreate:
		return run(cmdline, tagDir, tagFile)
	default:
		return fmt.Errorf("Unknown command \"%s\"", command)
	}
	return nil
}

func run(cmdline *cli.CommandLine, dirFunc, fileFunc func(*cli.CommandLine, string) error) error {
	for _, path := range cmdline.Paths() {
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}
		if info.Mode()&fs.ModeSymlink == fs.ModeSymlink {
			if !cmdline.FlagFollowSymlinks() {
				continue
			}
			err = dirFunc(cmdline, path)
		} else if info.Mode()&fs.ModeDir == fs.ModeDir {
			err = dirFunc(cmdline, path)
		} else {
			err = fileFunc(cmdline, path)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
