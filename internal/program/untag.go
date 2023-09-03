package program

import (
	"github.com/jwdev42/xtagger/internal/cli"
	"github.com/jwdev42/xtagger/internal/record"
	"io/fs"
	"os"
	"path/filepath"
)

func untagFile(cmdline *cli.CommandLine, path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	//Skip irregular files
	if !info.Mode().IsRegular() {
		return nil
	}
	return record.PurgeFile(path)
}

func untagDir(cmdline *cli.CommandLine, path string) error {
	examine := func(name string, d fs.DirEntry, err error) error {
		path := filepath.Join(path, name)
		if d.IsDir() {
			return nil
		}
		return untagFile(cmdline, path)
	}
	return fs.WalkDir(os.DirFS(path), ".", examine)
}
