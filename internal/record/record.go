package record

import (
	"encoding/json"
	"fmt"
	"github.com/dop251/scsu"
	"github.com/pkg/xattr"
)

const attrName = "user.xbackup"

// Represents a single record within a user.xbackup xattr entry
type Record struct {
	Checksum   string `json:"c"` //Hex-string of the SHA256sum of the file
	Identifier string `json:"i"` //Name of the backup job that saved the file with the given checksum at the given time
	Timestamp  int64  `json:"t"` //Unix timestamp
}

func (a *Record) Equals(b *Record) bool {
	if *a == *b {
		return true
	}
	return false
}

// Represents the whole content of a user.xbackup xattr entry
type Attribute []*Record

func LoadAttribute(path string) (Attribute, error) {
	//Read extended attribute
	payload, err := xattr.Get(path, attrName)
	if err != nil {
		return nil, fmt.Errorf("Failed to read extended attribute: %s", err)
	}
	//Decompress payload
	jsonText, err := scsu.Decode(payload)
	if err != nil {
		return nil, fmt.Errorf("Failed to decompress payload: %s", err)
	}
	//Decode JSON
	attr := make(Attribute, 0)
	attrp := &attr
	if err := json.Unmarshal([]byte(jsonText), attrp); err != nil {
		return nil, fmt.Errorf("Failed to decode json: %s", err)
	}
	return *attrp, nil
}

func (r Attribute) Store(path string) error {
	//Encode JSON
	jsonText, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("Failed to encode json: %s", err)
	}
	//Compress payload
	payload, err := scsu.EncodeStrict(string(jsonText), nil)
	if err != nil {
		return fmt.Errorf("Failed to compress text: %s", err)
	}
	//Write extended attribute
	if err := xattr.Set(path, attrName, payload); err != nil {
		return fmt.Errorf("Failed to write extended attribute: %s", err)
	}
	return nil
}
