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
	//Execute command-specific branch
	switch command := cmdline.Command(); command {
	case cli.CommandTag:
		return run(cmdline, tagDir, tagFile)
	case cli.CommandPrint:
		return run(cmdline, printDir, printFile)
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
		//Check if symlink
		if info.Mode()&fs.ModeSymlink == fs.ModeSymlink {
			if !cmdline.FlagFollowSymlinks() {
				continue
			} else {
				//if symlinks are allowed, follow them
				info, err = os.Stat(path)
				if err != nil {
					return err
				}
			}
		}
		if info.IsDir() {
			if !cmdline.FlagRecursive() {
				continue
			}
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
