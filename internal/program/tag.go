package program

import (
	"github.com/jwdev42/xtagger/internal/global"
	"github.com/jwdev42/xtagger/internal/io/filesystem"
	"github.com/jwdev42/xtagger/internal/record"
	"io/fs"
	"path/filepath"
)

func tagFile(parent string, dirEnt fs.DirEntry, opts *filesystem.WalkDirOpts) error {
	path := filepath.Join(parent, dirEnt.Name())
	f, err := record.NewFile(path)
	if err != nil {
		return global.FilterSoftError(err)
	}
	if err := f.CreateRecord(global.CommandLine.FlagNames()[0], global.CommandLine.FlagHash()); err != nil {
		return global.FilterSoftError(err)
	}
	return nil
}
