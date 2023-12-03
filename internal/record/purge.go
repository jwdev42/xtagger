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

package record

import (
	"github.com/pkg/xattr"
	"os"
)

// PurgeAttr removes xtagger's extended attribute from the given file.
func PurgeAttr(f *os.File) error {
	attrNames, err := xattr.FList(f)
	if err != nil {
		return err
	}
	for _, name := range attrNames {
		if name == attrName {
			return xattr.FRemove(f, attrName)
		}
	}
	return nil
}
