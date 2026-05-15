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
	"github.com/jwdev42/xtagger/internal/config"
	"github.com/jwdev42/xtagger/internal/record"
	"github.com/jwdev42/xtagger/internal/xio/filesystem"
	"os"
)

func printFile(rt *prt, meta *filesystem.Meta) error {
	// Print prints an attribute and respects program settings
	print := func(attr record.Attribute, path string) error {
		if rt.prefs.PrintRecords {
			// Print whole record
			res, err := attr.PrettyPrintWithPath(path)
			if err != nil {
				return err
			}
			rt.printer.Print(res)
			return nil
		}
		// Only print path by default
		rt.printer.Print(path)
		return nil
	}
	constraint := rt.prefs.PrintConstraint
	// Open file
	f, err := os.Open(meta.Path())
	if err != nil {
		return err
	}
	defer f.Close()
	// Load Attributes
	attr, err := record.FLoadAttribute(f)
	if err != nil {
		return err
	}
	// Filter Attributes by name
	attr = attr.FilterByName(rt.prefs.Names...)

	// Handle empty attributes
	if len(attr) < 1 {
		switch constraint {
		case config.PrintConstraintUntagged:
			//Print recordless file if PrintConstraintUntagged is set
			return print(attr, meta.Path())
		}
		//Skip file otherwise
		return nil
	}

	// Handle nonempty attributes
	switch constraint {
	case config.PrintConstraintUntagged:
		return nil //Skip tagged file
	}
	return print(attr, meta.Path()) //Print tagged file by default
	return nil
}
