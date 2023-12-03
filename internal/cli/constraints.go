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

package cli

const (
	TagConstraintNone TagConstraint = iota
	TagConstraintUntagged
	TagConstraintInvalid
)

const (
	PrintConstraintNone PrintConstraint = iota
	PrintConstraintValid
	PrintConstraintInvalid
	PrintConstraintUntagged
)

const (
	UntagConstraintNone UntagConstraint = iota
	UntagConstraintAll
	UntagConstraintInvalid
)

type TagConstraint int
type UntagConstraint int
type PrintConstraint int
