package program

import (
	"github.com/jwdev42/xtagger/internal/cli"
	"github.com/jwdev42/xtagger/internal/global"
	"github.com/jwdev42/xtagger/internal/io/filesystem"
	"github.com/jwdev42/xtagger/internal/record"
	"io/fs"
	"path/filepath"
)

func tagDir(cmdline *cli.CommandLine, path string) error {
	examine := func(path string, d fs.DirEntry, opts *filesystem.WalkDirOpts) error {
		return tagFile(cmdline, filepath.Join(path, d.Name()))
	}
	return filesystem.WalkDir(path, createWalkDirOpts(cmdline, true), examine)
}

func tagFile(cmdline *cli.CommandLine, path string) error {
	f, err := record.NewFile(path)
	if err != nil {
		return global.FilterSoftError(err)
	}
	if err := f.CreateRecord(cmdline.FlagNames()[0], cmdline.FlagHash()); err != nil {
		return global.FilterSoftError(err)
	}
	return nil
}
