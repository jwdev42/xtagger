package record

import (
	"fmt"
	"os"
	"testing"
)

func testAttributeStoreAndLoad(t *testing.T, sample Attribute) error {
	const path = "xattr_test.temp"
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
	for i, rec := range sample {
		if !rec.Equals(loaded[i]) {
			return fmt.Errorf("Records at index %d do not match:", i)
		}
	}
	return nil
}

func TestAttributeStoreAndLoad(t *testing.T) {
	samples := []Attribute{
		{},
		{
			&Record{
				Checksum:   "368b97b0b055910d97d284f834cbf1f8d5dec95b70576c8aedf6361e6a7bbc63",
				Identifier: "私はウサギです",
				Timestamp:  23,
				Valid:      false,
			},
			&Record{
				Checksum:   "1f2946e2fd7d0be6c4295c1ed828f0ff4aec21e89df898f9efbaddbe445c5c7c",
				Identifier: "TestBackup123",
				Timestamp:  1686676137,
				Valid:      true,
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
	sample := Attribute{
		&Record{
			Checksum:   "368b97b0b055910d97d284f834cbf1f8d5dec95b70576c8aedf6361e6a7bbc63",
			Identifier: "私はウサギです",
			Timestamp:  23,
		},
		&Record{
			Checksum:   "368b97b0b055910d97d284f834cbf1f8d5dec95b70576c8aedf6361e6a7bbc63",
			Identifier: "私はウサギです",
			Timestamp:  23,
		},
		&Record{
			Checksum:   "1f2946e2fd7d0be6c4295c1ed828f0ff4aec21e89df898f9efbaddbe445c5c7c",
			Identifier: "私はウサギです",
			Timestamp:  23,
		},
		&Record{
			Checksum:   "368b97b0b055910d97d284f834cbf1f8d5dec95b70576c8aedf6361e6a7bbc63",
			Identifier: "test",
			Timestamp:  23,
		},
		&Record{
			Checksum:   "368b97b0b055910d97d284f834cbf1f8d5dec95b70576c8aedf6361e6a7bbc63",
			Identifier: "私はウサギです",
			Timestamp:  24,
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
