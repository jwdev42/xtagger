package program

import (
	"github.com/jwdev42/xtagger/internal/global"
	"github.com/jwdev42/xtagger/internal/io/filesystem"
	"github.com/jwdev42/xtagger/internal/record"
	"io/fs"
	"path/filepath"
)

func invalidateFile(parent string, dirEnt fs.DirEntry, opts *filesystem.WalkDirOpts) error {
	path := filepath.Join(parent, dirEnt.Name())
	f, err := record.NewFile(path)
	if err != nil {
		return global.FilterSoftError(err)
	}
	if err := f.InvalidateOutdatedEntries(commandLine.FlagNames(), commandLine.FlagAllowRevalidation()); err != nil {
		return global.FilterSoftError(err)
	}
	if commandLine.FlagPrint0() {
		if _, err := printMe.Print0(path); err != nil {
			return err
		}
	}
	return nil
}
