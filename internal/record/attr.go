package record

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/xattr"
	"os"
)

// Represents the whole content of a user.xtagger xattr entry
type Attribute map[string]*Record

// Loads the xtagger extended attribute for path.
func LoadAttribute(path string) (Attribute, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return FLoadAttribute(f)
}

// Loads the xtagger extended attribute for File f.
func FLoadAttribute(f *os.File) (Attribute, error) {
	//Read extended attribute
	payload, err := xattr.FGet(f, attrName)
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
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return r.FStore(f)
}

// Stores the xtagger extended attribute in File's inode.
func (r Attribute) FStore(f *os.File) error {
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
	if err := xattr.FSet(f, attrName, payload); err != nil {
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
