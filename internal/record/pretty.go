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

package record

import (
	"github.com/jwdev42/xtagger/internal/hashes"
	"time"
)

// PrettyRecord represents a single record within a user.xtagger entry
type PrettyRecord struct {
	Checksum  string      `json:"checksum"`  // Hex encoded hash.
	Algorithm hashes.Algo `json:"algorithm"` // Name of the hashing algorithm.
	Timestamp time.Time   `json:"timestamp"` // Timestamp of hashing operation.
}

// PrettyAttribute is a human-readable version of Attribute.
// Use it for printing attribute data.
type PrettyAttribute struct {
	Path    string                  `json:"path"`    // Attribute file path
	Records map[string]PrettyRecord `json:"records"` // Attribute records
}
