package program

import (
	"fmt"
	"github.com/jwdev42/xtagger/internal/cli"
	"github.com/jwdev42/xtagger/internal/record"
	"io/fs"
	"os"
	"path/filepath"
)

func printFile(cmdline *cli.CommandLine, path string) error {
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
	examine := func(name string, d fs.DirEntry, err error) error {
		path := filepath.Join(path, name)
		if d.IsDir() {
			if !cmdline.FlagRecursive() {
				return fs.SkipDir
			}
			return nil
		}
		return printFile(cmdline, path)
	}
	return fs.WalkDir(os.DirFS(path), ".", examine)
}
