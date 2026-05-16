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
	"github.com/jwdev42/xtagger/internal/xio/filesystem"
	"log/slog"
	"os"
)

func tagFile(rt *prt, meta *filesystem.Meta) error {
	name := rt.prefs.Names[0]
	// Open file
	f, err := os.Open(meta.Path())
	if err != nil {
		return err
	}
	defer f.Close()
	// Load attribute
	attr, err := record.FLoadAttribute(f)
	if err != nil {
		return err
	}
	// Process untagged constraint
	if rt.prefs.Constraints.Has(config.ConstraintUntagged) && len(attr) > 0 {
		return nil //Skip already tagged files
	}
	// Check if a record with the designated name already exists
	if attr.Exists(name) {
		return fmt.Errorf("Record %q already exists for path %q", name, meta.Path())
	}
	// Create record
	slog.Debug("Create new tag record", "path", meta.Path())
	rec, err := record.CreateRecord(f, rt.prefs.UseHash)
	if err != nil {
		return err
	}
	// Add record to attribute
	attr[name] = rec
	// Save attribute
	if err := attr.FStore(f); err != nil {
		return err
	}
	// Send info log
	slog.Info("Tagged file", "path", meta.Path(), "checksum", rec.Hex(), "algorithm", rec.Algo())
	return nil
}
