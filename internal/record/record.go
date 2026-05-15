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
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jwdev42/xtagger/internal/hashes"
	"io"
	"time"
)

const attrName = "user.xtagger"

// Record represents a single record within a user.xtagger entry
type Record struct {
	checksum  []byte      // File hash.
	hashAlgo  hashes.Algo // Name of the used hashing algorithm.
	timestamp time.Time   // Timestamp of hashing operation.
}

// CreateRecord hashes src using the given hash algorithm,
// stores the result in a new Record,
// then returns a pointer to that Record.
func CreateRecord(src io.Reader, algo hashes.Algo) (*Record, error) {
	hash := algo.New()
	if err := hashes.Hash(src, hash); err != nil {
		return nil, err
	}
	return &Record{
		checksum:  hash.Sum(nil),
		hashAlgo:  algo,
		timestamp: time.Now(),
	}, nil
}

// Algo returns the Record's hashing algorithm.
func (r *Record) Algo() hashes.Algo {
	return r.hashAlgo
}

// Hex returns the Record's checksum as hex string.
func (r *Record) Hex() string {
	return hex.EncodeToString(r.checksum)
}

// Time returns the time of record creation
func (r *Record) Time() time.Time {
	return r.timestamp
}

// MarshalJSON implements the json.Marshaler interface.
func (r *Record) MarshalJSON() ([]byte, error) {
	// Create anonymous proxy struct with Exported fields
	proxy := struct {
		Checksum  []byte      `json:"c"`
		HashAlgo  hashes.Algo `json:"h"`
		Timestamp int64       `json:"t"`
	}{
		Checksum:  r.checksum,
		HashAlgo:  r.hashAlgo,
		Timestamp: r.timestamp.Unix(),
	}

	return json.Marshal(proxy)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// Calls the receiver's internal validation function implicitly.
func (r *Record) UnmarshalJSON(data []byte) error {
	// Create anonymous proxy struct with Exported fields
	proxy := struct {
		Checksum  []byte      `json:"c"`
		HashAlgo  hashes.Algo `json:"h"`
		Timestamp int64       `json:"t"`
	}{}

	if err := json.Unmarshal(data, &proxy); err != nil {
		return err
	}

	// Map proxy struct fields to unexported struct fields
	r.checksum = proxy.Checksum
	r.hashAlgo = proxy.HashAlgo
	r.timestamp = time.Unix(proxy.Timestamp, 0)
	// Validate input
	return r.validate()
}

func (r *Record) validate() error {
	// Checks if receiver is nil (can be triggered by writing null in JSON)
	if r == nil {
		return errors.New("Record cannot be nil")
	}
	// Checks if the hashing algorithm for the Record is known
	if err := r.hashAlgo.Validate(); err != nil {
		return err
	}
	// Checks if Checksum has the correct length
	var checksumLen int
	switch r.hashAlgo {
	case hashes.RIPEMD160:
		checksumLen = 20
	default:
		checksumLen = 32
	}
	if len(r.checksum) != checksumLen {
		return fmt.Errorf("Expected a checksum of %d bytes for %s", checksumLen, r.hashAlgo)
	}
	return nil
}
