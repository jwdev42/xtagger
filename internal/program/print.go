package program

import (
	"fmt"
	"github.com/jwdev42/xtagger/internal/cli"
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
	attrs := f.Attributes()
	//Filter empty entries if PrintModeOmitEmpty is set
	if mode := commandLine.FlagPrintMode(); len(attrs) == 0 && (mode == cli.PrintModeOmitEmpty || mode == cli.PrintModeValidOnly) {
		return nil
	}
	//Filter by name if at least one name restriction is set
	if names := commandLine.FlagNames(); names != nil {
		var found bool
		for _, name := range names {
			if rec := attrs[name]; rec != nil {
				//Skip invalid records if mode is PrintModeValidOnly
				if commandLine.FlagPrintMode() == cli.PrintModeValidOnly && !rec.Valid {
					continue
				}
				found = true
				break
			}
		}
		if !found {
			return nil
		}
	} else {
		//skip entries that have no valid record if mode is PrintModeValidOnly
		if commandLine.FlagPrintMode() == cli.PrintModeValidOnly {
			var hasValidRecord bool
			for _, rec := range attrs {
				if rec.Valid {
					hasValidRecord = true
				}
			}
			if !hasValidRecord {
				return nil
			}
		}
	}
	if commandLine.FlagPrint0() {
		_, err := printMe.Print0(path)
		return err
	}
	fmt.Printf("%s\n", f)
	return nil
}
