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
	"errors"
	"fmt"
	"github.com/jwdev42/xtagger/internal/hashes"
	"time"
)

const attrName = "user.xtagger"

// Represents a single record within a user.xtagger xattr entry
type Record struct {
	Checksum  string      `json:"c"` //Hex-string of the SHA256sum of the file.
	HashAlgo  hashes.Algo `json:"h"` //Name of the used hashing algorithm.
	Timestamp int64       `json:"t"` //Unix timestamp of the record's creation.
	Valid     bool        `json:"v"` //Record valid if true, invalidated if false.
}

// Returns a new record with the current time as timestamp. All other member fields
// are zero-values.
func NewRecord() *Record {
	return &Record{
		Timestamp: time.Now().Unix(),
	}
}

func (a *Record) Equals(b *Record) bool {
	if *a == *b {
		return true
	}
	return false
}

func (r *Record) Copy() *Record {
	recCpy := *r
	return &recCpy
}

func (r *Record) validate() error {
	//Checks if receiver is nil (can be triggered by writing null in JSON)
	if r == nil {
		return errors.New("Record cannot be null")
	}
	//Checks if the hashing algorithm for the Record is known
	if err := r.HashAlgo.Validate(); err != nil {
		return err
	}
	//Checks if Checksum has the correct length
	var checksumLen int
	switch r.HashAlgo {
	case hashes.RIPEMD160:
		checksumLen = 40
	default:
		checksumLen = 64
	}
	if len(r.Checksum) != checksumLen {
		return fmt.Errorf("Expected a checksum of %d characters for %s", checksumLen, r.HashAlgo)
	}
	//Checks if Checksum is represented as hexadecimal string
	for i, ch := range []rune(r.Checksum) {
		if !(ch >= 48 && ch <= 57 || ch >= 97 && ch <= 102) {
			return fmt.Errorf("Checksum has illegal character at index %d", i)
		}
	}
	return nil
}
