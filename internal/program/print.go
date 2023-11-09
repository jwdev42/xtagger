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
	printPath := func(path string) error {
		if commandLine.FlagPrint0() {
			_, err := printMe.Print0(path)
			return err
		}
		_, err := fmt.Printf("%s\n", path)
		return err
	}
	path := filepath.Join(parent, dirEnt.Name())
	constraint := commandLine.PrintConstraint()
	//Open file
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	//Load Attributes
	attrs, err := record.FLoadAttribute(f)
	if err != nil {
		return err
	}
	//Filter Attributes by name
	if names := commandLine.Names(); names != nil {
		attrs = attrs.FilterByName(names...)
	}

	if len(attrs) < 1 {
		switch constraint {
		case cli.PrintConstraintUntagged:
			//Print recordless file if PrintConstraintUntagged is set
			return printPath(path)
		}
		//Skip file otherwise
		return nil
	}

	switch constraint {
	case cli.PrintConstraintNone:
		return printPath(path) //Print tagged file if no constraint is set
	case cli.PrintConstraintUntagged:
		return nil //Skip tagged file
	}

	//Iterate through Attributes to check for invalid and valid records
	var hasInvalidEntry bool
	var hasValidEntry bool
	for _, rec := range attrs {
		if !rec.Valid {
			hasInvalidEntry = true
		} else {
			hasValidEntry = true
		}
	}

	switch constraint {
	case cli.PrintConstraintInvalid:
		//Print if all records are invalid
		if !hasValidEntry {
			return printPath(path)
		}
	case cli.PrintConstraintValid:
		//Print if all records are valid
		if !hasInvalidEntry {
			return printPath(path)
		}
	default:
		panic("You're not supposed to be here")
	}
	return nil
}
