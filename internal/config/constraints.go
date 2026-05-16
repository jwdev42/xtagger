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

package config

const ConstraintNone Constraint = 0 // No constraints

const (
	ConstraintUntagged Constraint = 1 << iota // Only process files without user.xtagger xattr entry
)

type Constraint uint64

type Constraints struct {
	storage Constraint
}

func (cs *Constraints) Add(c Constraint) {
	cs.storage |= c
}

func (cs *Constraints) Has(c Constraint) bool {
	return cs.storage&c != 0
}

func (cs *Constraints) Remove(c Constraint) {
	cs.storage &^= c
}

func (cs *Constraints) Toggle(c Constraint) {
	cs.storage ^= c
}

func (cs *Constraints) Union(c Constraints) {
	cs.Add(c.storage)
}
