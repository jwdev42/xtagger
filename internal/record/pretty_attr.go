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
	toml "github.com/pelletier/go-toml/v2"
	"io"
)

// PrettyAttribute is a human-readable version of Attribute.
// Use it for printing attribute data.
type PrettyAttribute map[string]*PrettyRecord

// TomlWithPath writes path and attribute as a toml entry to the writer.
func (r PrettyAttribute) TomlWithPath(wr io.Writer, path string) error {
	container := make(map[string][]*NamedPrettyRecord)
	container[path] = r.Slice()
	enc := toml.NewEncoder(wr)
	return enc.Encode(container)
}

// Slice returns PrettyAttribute's data reorganized as slice of NamedPrettyRecords
func (r PrettyAttribute) Slice() []*NamedPrettyRecord {
	recs := make([]*NamedPrettyRecord, len(r))
	i := 0
	for k, v := range r {
		rec := v.WithName(k)
		recs[i] = &rec
		i++
	}
	return recs
}
