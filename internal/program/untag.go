package program

import (
	"github.com/jwdev42/xtagger/internal/global"
	"github.com/jwdev42/xtagger/internal/io/filesystem"
	"github.com/jwdev42/xtagger/internal/record"
	"io/fs"
	"os"
	"path/filepath"
)

func untagFile(parent string, dirEnt fs.DirEntry, opts *filesystem.WalkDirOpts) error {
	path := filepath.Join(parent, dirEnt.Name())
	if commandLine.Names() != nil {
		//Open file
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		attr, err := record.FLoadAttribute(f)
		if err != nil {
			return err
		}
		initialLength := len(attr)
		for _, name := range commandLine.Names() {
			delete(attr, name)
		}
		if initialLength == len(attr) {
			//Return if attr didn't change
			return nil
		}
		if err := attr.FStore(f); err != nil {
			return err
		}
	} else {
		if err := record.PurgeFile(path); err != nil {
			return global.FilterSoftError(err)
		}
	}
	if commandLine.FlagPrint0() {
		if _, err := printMe.Print0(path); err != nil {
			return err
		}
	}
	return nil
}
