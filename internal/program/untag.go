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
