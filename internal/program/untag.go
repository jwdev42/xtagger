package program

import (
	"github.com/jwdev42/xtagger/internal/cli"
	"github.com/jwdev42/xtagger/internal/io/filesystem"
	"github.com/jwdev42/xtagger/internal/record"
	"io/fs"
	"path/filepath"
)

func untagFile(cmdline *cli.CommandLine, path string) error {
	return record.PurgeFile(path)
}

func untagDir(cmdline *cli.CommandLine, path string) error {
	examine := func(path string, d fs.DirEntry, opts *filesystem.WalkDirOpts) error {
		return untagFile(cmdline, filepath.Join(path, d.Name()))
	}
	return filesystem.WalkDir(path, createWalkDirOpts(cmdline, false), examine)
}
