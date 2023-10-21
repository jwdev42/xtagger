package program

import (
	"fmt"
	"github.com/jwdev42/xtagger/internal/io/filesystem"
	"github.com/jwdev42/xtagger/internal/record"
	"io/fs"
	"os"
	"path/filepath"
)

func printFile(parent string, dirEnt fs.DirEntry, opts *filesystem.WalkDirOpts) error {
	path := filepath.Join(parent, dirEnt.Name())
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	//Skip irregular files
	if !info.Mode().IsRegular() {
		return nil
	}
	f, err := record.NewFile(path)
	if err != nil {
		return err
	}
	if commandLine.FlagOmitEmpty() && len(f.Attributes()) == 0 {
		return nil
	}
	if commandLine.FlagPrint0() {
		_, err := printMe.Print0(path)
		return err
	}
	fmt.Printf("%s\n", f)
	return nil
}
