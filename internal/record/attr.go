package record

import (
	"encoding/json"
	"errors"
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
// Returns an empty Attribute if the file does not have an extended attribute.
func FLoadAttribute(f *os.File) (Attribute, error) {
	//Read extended attribute
	payload, err := xattr.FGet(f, attrName)
	if errors.Is(err, xattr.ENOATTR) {
		//Create a new Attribute if file doesn't have one yet
		return make(Attribute), nil
	} else if err != nil {
		return nil, fmt.Errorf("Failed to read extended attribute: %s", err)
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

func (r Attribute) Exists(name string) bool {
	if r[name] != nil {
		return true
	}
	return false
}

func (r Attribute) FilterByName(name ...string) Attribute {
	attr := make(Attribute)
	for _, key := range name {
		if rec := r[key]; rec != nil {
			attr[key] = rec
		}
	}
	return attr
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
