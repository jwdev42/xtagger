package program

import (
	"fmt"
	"github.com/jwdev42/xtagger/internal/cli"
	"github.com/jwdev42/xtagger/internal/io/filesystem"
	"github.com/jwdev42/xtagger/internal/record"
	"io/fs"
	"os"
	"path/filepath"
)

func printFile(cmdline *cli.CommandLine, path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	//Skip irregular files
	if !info.Mode().IsRegular() {
		return nil
	}
	f, err := record.NewFile(path)
	if err != nil {
		return err
	}
	if cmdline.FlagOmitEmpty() && len(f.Attributes()) == 0 {
		return nil
	}
	fmt.Printf("%s\n", f)
	return nil
}

func printDir(cmdline *cli.CommandLine, path string) error {
	examine := func(path string, d fs.DirEntry, opts *filesystem.WalkDirOpts) error {
		return printFile(cmdline, filepath.Join(path, d.Name()))
	}

	return filesystem.WalkDir(path, createWalkDirOpts(cmdline), examine)
}
