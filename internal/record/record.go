package record

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jwdev42/xtagger/internal/hashes"
	"github.com/pkg/xattr"
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
	//Validation
	if err := attr.validate(); err != nil {
		return nil, err
	}
	return attr, nil
}

// Stores the xtagger extended attribute in path's inode.
func (r Attribute) Store(path string) error {
	if r == nil {
		panic("BUG: Calling Store() with a nil receiver is prohibited")
	}
	//Validation
	if err := r.validate(); err != nil {
		return err
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

func (r Attribute) validate() error {
	if r == nil {
		return fmt.Errorf("Attribute %s cannot be null", attrName)
	}
	for name, rec := range r {
		if err := rec.validate(); err != nil {
			return fmt.Errorf("Validation for record \"%s\" failed: %s", name, err)
		}
	}
	return nil
}
