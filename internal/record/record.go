package record

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/xattr"
	"time"
)

const attrName = "user.xtagger"

// Represents a single record within a user.xtagger xattr entry
type Record struct {
	Checksum  string `json:"c"` //Hex-string of the SHA256sum of the file.
	Timestamp int64  `json:"t"` //Unix timestamp of the record's creation.
	Valid     bool   `json:"v"` //Record valid if true, invalidated if false.
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

// Represents the whole content of a user.xtagger xattr entry
type Attribute map[string]*Record

// Loads the xtagger extended attribute for path.
func LoadAttribute(path string) (Attribute, error) {
	//Read extended attribute
	payload, err := xattr.Get(path, attrName)
	if err != nil {
		//Wrap error to be able to catch ENOATTR
		return nil, fmt.Errorf("Failed to read extended attribute: %w", err)
	}
	//Decode JSON
	attr := make(Attribute)
	if err := json.Unmarshal([]byte(payload), &attr); err != nil {
		return nil, fmt.Errorf("Failed to decode json: %s", err)
	}
	return attr, nil
}

// Stores the xtagger extended attribute in path's inode.
func (r Attribute) Store(path string) error {
	if r == nil {
		panic("BUG: Calling Store() with a nil receiver is prohibited")
	}
	//Encode JSON
	payload, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("Failed to encode json: %s", err)
	}
	//Write extended attribute
	if err := xattr.Set(path, attrName, payload); err != nil {
		return fmt.Errorf("Failed to write extended attribute: %s", err)
	}
	return nil
}

// Returns the newest Record. Returns zero-values if no record was found.
func (r Attribute) MostRecent() (name string, rec *Record) {
	if r == nil || len(r) < 1 {
		return "", nil
	}
	for k, v := range r {
		if rec == nil || v.Timestamp > rec.Timestamp {
			name = k
			rec = v
		}
	}
	return name, rec
}
