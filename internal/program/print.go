package program

import (
	"fmt"
	"github.com/jwdev42/xtagger/internal/cli"
	"github.com/jwdev42/xtagger/internal/global"
	"github.com/jwdev42/xtagger/internal/record"
	"io/fs"
	"os"
	"path/filepath"
)

func printFile(parent string, info fs.FileInfo) error {
	print := func(attr record.Attribute, path string) error {
		if commandLine.FlagPrint0() {
			_, err := printMe.Print0(path)
			return err
		}
		if commandLine.FlagPrintRecords() {
			_, err := attr.FprintRecordsWithPath(os.Stdout, path)
			return err
		}
		_, err := fmt.Printf("%s\n", path)
		return err
	}
	path := filepath.Join(parent, info.Name())
	constraint := commandLine.PrintConstraint()
	//Open file
	f, err := os.Open(path)
	if err != nil {
		return global.FilterSoftError(err)
	}
	defer f.Close()
	//Load Attributes
	attr, err := record.FLoadAttribute(f)
	if err != nil {
		return global.FilterSoftError(err)
	}
	//Filter Attributes by name
	if names := commandLine.Names(); names != nil {
		attr = attr.FilterByName(names...)
	}

	if len(attr) < 1 {
		switch constraint {
		case cli.PrintConstraintUntagged:
			//Print recordless file if PrintConstraintUntagged is set
			return global.FilterSoftError(print(attr, path))
		}
		//Skip file otherwise
		return nil
	}

	switch constraint {
	case cli.PrintConstraintNone:
		return global.FilterSoftError(print(attr, path)) //Print tagged file if no constraint is set
	case cli.PrintConstraintUntagged:
		return nil //Skip tagged file
	}

	//Iterate through Attributes to check for invalid and valid records
	var hasInvalidEntry bool
	var hasValidEntry bool
	for _, rec := range attr {
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
			return global.FilterSoftError(print(attr, path))
		}
	case cli.PrintConstraintValid:
		//Print if all records are valid
		if !hasInvalidEntry {
			return global.FilterSoftError(print(attr, path))
		}
	default:
		panic("You're not supposed to be here")
	}
	return nil
}
