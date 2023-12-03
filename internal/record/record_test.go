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
	"fmt"
	"github.com/jwdev42/xtagger/internal/hashes"
	"github.com/pkg/xattr"
	"os"
	"testing"
)

func testAttributeStoreAndLoad(t *testing.T, sample Attribute) error {
	const path = "testAttributeStoreAndLoad.temp"
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return fmt.Errorf("Could not create temp file: %s", err)
	}
	f.Close()
	defer os.Remove(path)

	if err := sample.Store(path); err != nil {
		return fmt.Errorf("Failed to create xattr: %s", err)
	}

	loaded, err := LoadAttribute(path)
	if err != nil {
		return fmt.Errorf("Failed to load xattr: %s", err)
	}
	if len(loaded) != len(sample) {
		return fmt.Errorf("Attributes have different amounts of records (first: %d, second: %d)", len(sample), len(loaded))
	}
	for name, rec := range sample {
		if !rec.Equals(loaded[name]) {
			return fmt.Errorf("Records for key \"%s\" do not match:", name)
		}
	}
	return nil
}

func TestAttributeStoreAndLoad(t *testing.T) {
	samples := []Attribute{
		{},
		{
			"私はウサギです": &Record{
				Checksum:  "368b97b0b055910d97d284f834cbf1f8d5dec95b70576c8aedf6361e6a7bbc63",
				HashAlgo:  hashes.SHA256,
				Timestamp: 23,
				Valid:     false,
			},
			"TestBackup123": &Record{
				Checksum:  "1f2946e2fd7d0be6c4295c1ed828f0ff4aec21e89df898f9efbaddbe445c5c7c",
				HashAlgo:  hashes.SHA256,
				Timestamp: 1686676137,
				Valid:     true,
			},
		},
	}
	for i, sample := range samples {
		if err := testAttributeStoreAndLoad(t, sample); err != nil {
			t.Errorf("Failed test for sample %d: %s", i, err)
		}
	}
}

func TestEquals(t *testing.T) {
	//0 and 1 are equal, the others are different
	sample := []*Record{
		{
			Checksum:  "368b97b0b055910d97d284f834cbf1f8d5dec95b70576c8aedf6361e6a7bbc63",
			Timestamp: 23,
		},
		{
			Checksum:  "368b97b0b055910d97d284f834cbf1f8d5dec95b70576c8aedf6361e6a7bbc63",
			Timestamp: 23,
		},
		{
			Checksum:  "1f2946e2fd7d0be6c4295c1ed828f0ff4aec21e89df898f9efbaddbe445c5c7c",
			Timestamp: 23,
		},
		{
			Checksum:  "368b97b0b055910d97d284f834cbf1f8d5dec95b70576c8aedf6361e6a7bbc63",
			Timestamp: 24,
		},
	}
	for i, rec := range sample {
		if i < 2 {
			if !rec.Equals(sample[0]) {
				t.Errorf("Record at index %d was supposed to match record at index 0", i)
			}
		} else {
			if rec.Equals(sample[0]) {
				t.Errorf("Record at index %d wasn't supposed to match record at index 0", i)
			}
		}
	}
}

func TestNegativeLoadAttribute(t *testing.T) {
	const path = "TestNegativeLoadAttribute.temp"
	samples := []string{
		"",
		"null",
		`"test"`,
		"2",
		`{"test":null}`,
		`{"test":{"c":"368b97b0b055910d97d284f834cbf1f8d5dec95b70576c8aedf6361e6a7bbc63","h":"invalid","t":0,"v":false}}`,  //"h" has an illegal value
		`{"test":{"c":"368b97b0b055910d97d284f834cbf1f8d5dec95b70576c8aedf6361e6a7bbc63","h":"SHA256","t":0,"v":"false"}}`, //"v" is string instead of bool
		`{"test":{"c":368b97b0b055910d97d284f834cbf1f8d5dec95b70576c8aedf6361e6a7bbc63,"h":"SHA256","t":0,"v":false}}`,     //"c" is int instead of string
		`{"test":{}}`, //Record "test" is empty

	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		t.Fatalf("Could not create temp file: %s", err)
	}
	f.Close()
	defer os.Remove(path)

	for i, sample := range samples {
		//Write sample as extended attribute
		if err := xattr.Set(path, attrName, []byte(sample)); err != nil {
			t.Fatalf("Unexpected error: Failed to write extended attribute: %s", err)
		}
		_, err := LoadAttribute(path)
		if err == nil {
			t.Errorf("Index %d: Expected error for sample \"%s\"", i, sample)
		} else {
			t.Logf("[Negative test] index %d: %s", i, err)
		}
	}
}
