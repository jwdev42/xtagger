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

// Internal Record type for marshaling and unmarshaling
type mRecord struct {
	Checksum  []byte      `json:"c"`
	HashAlgo  hashes.Algo `json:"h"`
	Timestamp int64       `json:"t"`
}

// Update Record rec with mRecord's data.
func (m *mRecord) update(rec *Record) {
	rec.checksum = m.Checksum
	rec.hashAlgo = m.HashAlgo
	rec.timestamp = time.Unix(m.Timestamp, 0)
}
