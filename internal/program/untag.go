package program

import (
	"github.com/jwdev42/xtagger/internal/record"
	"github.com/jwdev42/xtagger/internal/softerrors"
	"io/fs"
	"os"
	"path/filepath"
)

func untagFile(parent string, info fs.FileInfo) error {
	path := filepath.Join(parent, info.Name())
	//Open file
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if commandLine.Names() != nil {
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
		if err := record.PurgeAttr(f); err != nil {
			return softerrors.Consume(err)
		}
	}
	if commandLine.FlagPrint0() {
		if _, err := printMe.Print0(path); err != nil {
			return err
		}
	}
	return nil
}
