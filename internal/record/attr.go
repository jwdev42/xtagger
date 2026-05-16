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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pkg/xattr"
	"os"
)

// Attribute represents a user.xtagger xattr entry. It's purpose is to
// retrieve, manipulate and write xtagger xattr entries.
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
	// Purge xattr record if Attribute has no entries
	if len(r) == 0 {
		return PurgeAttr(f)
	}
	// Validation
	if err := r.validate(); err != nil {
		return err
	}
	// Encode JSON
	payload, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("Failed to encode json: %s", err)
	}
	// Write extended attribute
	if err := xattr.FSet(f, attrName, payload); err != nil {
		return fmt.Errorf("Failed to write extended attribute: %s", err)
	}
	return nil
}

// DeleteRecord deletes Records from Attribute by Name.
// Returns true if at least one Record was deleted.
func (r Attribute) DeleteRecord(name ...string) bool {
	before := len(r)
	for _, n := range name {
		delete(r, n)
	}
	return len(r) < before
}

// Returns the newest Record. Returns zero-values if no record was found.
func (r Attribute) MostRecent() (name string, rec *Record) {
	if r == nil || len(r) < 1 {
		return "", nil
	}
	for k, v := range r {
		if rec == nil || v.Time().After(rec.Time()) {
			name = k
			rec = v
		}
	}
	return name, rec
}

// Exists returns true if the receiver holds a record that matches name.
func (r Attribute) Exists(name string) bool {
	if r[name] != nil {
		return true
	}
	return false
}

// FilterByName returns a new Attribute that contains all entries
// from the receiver that match one of the passed names.
// However, if name is empty or nil, the receiver will be returned.
func (r Attribute) FilterByName(name ...string) Attribute {
	if len(name) == 0 {
		return r
	}
	attr := make(Attribute)
	for _, key := range name {
		if rec := r[key]; rec != nil {
			attr[key] = rec
		}
	}
	return attr
}

func (r Attribute) Prettify() (pretty PrettyAttribute) {
	pretty = make(PrettyAttribute)
	for k, v := range r {
		pretty[k] = v.Prettify()
	}
	return
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
