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

package record

import (
	"errors"
	toml "github.com/pelletier/go-toml/v2"
	"io"
	"time"
)

const (
	RCStateUndefined RCState = "UNDEFINED" // Indicates a not yet defined state.
	RCStateFailed            = "FAILED"    // Verification could not be performed.
	RCStateOK                = "OK"        // Verification succeeded, new checksum equals original checksum.
	RCStateMismatch          = "MISMATCH"  // Verification mismatch, checksums differ.
)

type RCState string // Record comparison state

type RecordVerification struct {
	Result       RCState       `toml:"result"`              // Verification result.
	Interval     time.Duration `toml:"interval"`            // Time interval between original record and verification.
	Error        error         `toml:"error,omitempty"`     // Optional error message. Non-nil if Result is not RCStateOK.
	Original     PrettyRecord  `toml:"original_record"`     // Original record data.
	Verification PrettyRecord  `toml:"verification_record"` // Verification record data.
}

func VerifyRecord(orig, ver PrettyRecord) RecordVerification {
	res := RecordVerification{
		Original:     orig,
		Verification: ver,
		Interval:     ver.Timestamp.Sub(orig.Timestamp),
	}
	if orig.Algorithm != ver.Algorithm {
		res.Result = RCStateMismatch
		res.Error = errors.New("Algorithm mismatch")
		return res
	}
	if orig.Checksum != ver.Checksum {
		res.Result = RCStateMismatch
		res.Error = errors.New("Checksum mismatch")
	}
	res.Result = RCStateOK
	return res
}

type AttributeVerification map[string]RecordVerification

func VerifyAttribute(orig, ver PrettyAttribute) AttributeVerification {
	res := make(AttributeVerification)
	for name, rec := range orig {
		res[name] = VerifyRecord(rec, ver[name])
	}
	return res
}

// TomlWithPath writes path and AttributeVerification as a toml entry to the writer.
func (r AttributeVerification) TomlWithPath(wr io.Writer, path string) error {
	container := make(map[string]AttributeVerification)
	container[path] = r
	enc := toml.NewEncoder(wr)
	return enc.Encode(container)
}
