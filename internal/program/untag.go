//This file is part of xtagger. ©2023 Jörg Walter.
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
	"github.com/jwdev42/xtagger/internal/record"
	"github.com/jwdev42/xtagger/internal/xio/filesystem"
	"os"
)

func untagFile(rt *payloadRuntime, meta *filesystem.Meta) error {
	//Open file
	f, err := os.Open(meta.Path())
	if err != nil {
		return err
	}
	defer f.Close()
	if len(rt.prefs.Names) > 0 {
		attr, err := record.FLoadAttribute(f)
		if err != nil {
			return err
		}
		initialLength := len(attr)
		for _, name := range rt.prefs.Names {
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
			return err
		}
	}
	if rt.prefs.UsePrint0 {
		if _, err := printMe.Print0(meta.Path()); err != nil {
			return err
		}
	}
	return nil
}
