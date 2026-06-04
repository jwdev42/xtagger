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

// Add adds Constraint b to a, then returns the combined constrained.
func (a Constraint) Add(b Constraint) Constraint {
	return a | b
}

// Equals returns true if a == b.
func (a Constraint) Equals(b Constraint) bool {
	return a == b
}

// Has returns true if b is a subset of a.
// Has will always return false on ConstraintNone.
func (a Constraint) Has(b Constraint) bool {
	return a&b != 0
}

// Remove returns the relative complement between a and b.
func (a Constraint) Remove(b Constraint) Constraint {
	return a &^ b
}

// Sdiff returns the symmetric difference between a and b.
func (a Constraint) Sdiff(b Constraint) Constraint {
	return a ^ b
}

// String returns a string representation of the constraint.
// The returned string will always be lowercase.
// If the constraint has no string representation, String will return
// an empty string.
func (a Constraint) String() string {
	switch a {
	case ConstraintNone:
		return "none"
	case ConstraintUntagged:
		return "untagged"
	}
	return ""
}
