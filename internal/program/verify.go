// This file is part of xtagger. ©2023-2026 Jörg Walter.
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package program

import (
	"github.com/jwdev42/xtagger/internal/record"
	"github.com/jwdev42/xtagger/internal/xio/xfs"
	"os"
	"strings"
)

func verifyFile(rt *prt, meta *xfs.Meta) error {
	// Open file
	f, err := os.Open(meta.Path())
	if err != nil {
		return err
	}
	defer f.Close()
	// Load attribute
	attr, err := record.FLoadAttribute(f)
	if err != nil || len(attr) == 0 {
		return err
	}
	updated, err := attr.Update(f)
	if err != nil {
		return err
	}
	builder := &strings.Builder{}
	cmp := record.VerifyAttribute(attr.Prettify(), updated.Prettify())
	if err := cmp.TomlWithPath(builder, meta.Path()); err != nil {
		return err
	}
	rt.printer.Print(builder.String())
	return nil
}
