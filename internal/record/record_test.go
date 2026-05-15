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
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/jwdev42/xtagger/internal/hashes"
	"github.com/pkg/xattr"
	"os"
	"testing"
	"time"
)

func mustDecodeHex(input string) []byte {
	res, err := hex.DecodeString(input)
	if err != nil {
		panic(err)
	}
	return res
}

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
		if !bytes.Equal(rec.checksum, loaded[name].checksum) {
			return fmt.Errorf("Checksum mismatch on records for key %q", name)
		}
		if rec.hashAlgo != loaded[name].hashAlgo {
			return fmt.Errorf("Algorithm mismatch on records for key %q", name)
		}
		if !rec.timestamp.Equal(loaded[name].timestamp) {
			return fmt.Errorf("Timestamp mismatch on records for key %q", name)
		}
	}
	return nil
}

func TestAttributeStoreAndLoad(t *testing.T) {
	samples := []Attribute{
		{},
		{
			"私はウサギです": &Record{
				checksum:  mustDecodeHex("368b97b0b055910d97d284f834cbf1f8d5dec95b70576c8aedf6361e6a7bbc63"),
				hashAlgo:  hashes.SHA256,
				timestamp: time.Unix(23, 0),
			},
			"TestBackup123": &Record{
				checksum:  mustDecodeHex("1f2946e2fd7d0be6c4295c1ed828f0ff4aec21e89df898f9efbaddbe445c5c7c"),
				hashAlgo:  hashes.SHA256,
				timestamp: time.Unix(1686676137, 0),
			},
		},
	}
	for i, sample := range samples {
		if err := testAttributeStoreAndLoad(t, sample); err != nil {
			t.Errorf("Failed test for sample %d: %s", i, err)
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
		`{"test":{"c":"NouXsLBVkQ2X0oT4NMvx+NXeyVtwV2yK7fY2Hmp7vGM=","h":"invalid","t":0}}`, //"h" has an illegal value
		`{"test":{"c":NouXsLBVkQ2X0oT4NMvx+NXeyVtwV2yK7fY2Hmp7vGM=,"h":"SHA256","t":0}}`,    //"c" is int instead of string
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
