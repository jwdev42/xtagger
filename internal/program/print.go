//This file is part of xtagger. ©2023-2026 Jörg Walter.
//This program is free software: you can redistribute it and/or modify
//it under the terms of the GNU General Public License as published by
//the Free Software Foundation, either version 3 of the License, or
//(at your option) any later version.
//
//This program is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//GNU General Public License for more details.
//
//You should have received a copy of the GNU General Public License
//along with this program.  If not, see <https://www.gnu.org/licenses/>.

package program

import (
	"fmt"
	"github.com/jwdev42/xtagger/internal/config"
	"github.com/jwdev42/xtagger/internal/record"
	"github.com/jwdev42/xtagger/internal/softerrors"
	"github.com/jwdev42/xtagger/internal/xio/filesystem"
	"os"
)

func printFile(rt *payloadRuntime, meta *filesystem.Meta) error {
	print := func(attr record.Attribute, path string) error {
		if rt.prefs.UsePrint0 {
			_, err := printMe.Print0(path)
			return err
		}
		if rt.prefs.PrintRecords {
			_, err := attr.FprintRecordsWithPath(os.Stdout, path)
			return err
		}
		_, err := fmt.Printf("%s\n", path)
		return err
	}
	constraint := rt.prefs.PrintConstraint
	//Open file
	f, err := os.Open(meta.Path())
	if err != nil {
		return softerrors.Consume(err)
	}
	defer f.Close()
	//Load Attributes
	attr, err := record.FLoadAttribute(f)
	if err != nil {
		return softerrors.Consume(err)
	}
	//Filter Attributes by name
	if names := rt.prefs.Names; names != nil {
		attr = attr.FilterByName(names...)
	}

	if len(attr) < 1 {
		switch constraint {
		case config.PrintConstraintUntagged:
			//Print recordless file if PrintConstraintUntagged is set
			return softerrors.Consume(print(attr, meta.Path()))
		}
		//Skip file otherwise
		return nil
	}

	switch constraint {
	case config.PrintConstraintNone:
		return softerrors.Consume(print(attr, meta.Path())) //Print tagged file if no constraint is set
	case config.PrintConstraintUntagged:
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
	case config.PrintConstraintInvalid:
		//Print if all records are invalid
		if !hasValidEntry {
			return softerrors.Consume(print(attr, meta.Path()))
		}
	case config.PrintConstraintValid:
		//Print if all records are valid
		if !hasInvalidEntry {
			return softerrors.Consume(print(attr, meta.Path()))
		}
	default:
		panic("You're not supposed to be here")
	}
	return nil
}
