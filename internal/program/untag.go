package program

import (
	"github.com/jwdev42/xtagger/internal/global"
	"github.com/jwdev42/xtagger/internal/io/filesystem"
	"github.com/jwdev42/xtagger/internal/record"
	"io/fs"
	"path/filepath"
)

func untagFile(parent string, dirEnt fs.DirEntry, opts *filesystem.WalkDirOpts) error {
	path := filepath.Join(parent, dirEnt.Name())
	if err := record.PurgeFile(path); err != nil {
		return global.FilterSoftError(err)
	}
	if commandLine.FlagPrint0() {
		if _, err := printMe.Print0(path); err != nil {
			return err
		}
	}
	return nil
}
