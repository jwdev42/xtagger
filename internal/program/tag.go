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
	global.DefaultLogger.Infof("Hashing file %s", path)
	if err := f.CreateRecord(commandLine.Names()[0], commandLine.FlagHash()); err != nil {
		return global.FilterSoftError(err)
	}
	global.DefaultLogger.Infof("Successfully hashed file %s", path)
	if commandLine.FlagPrint0() {
		if _, err := printMe.Print0(path); err != nil {
			return err
		}
	}
	return nil
}
