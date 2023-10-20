package record

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jwdev42/xtagger/internal/hashes"
	"github.com/pkg/xattr"
	"hash"
	"io"
	"io/fs"
	"os"
)

// Accesses and manipulates the records of a particular file
type File struct {
	path string
	attr Attribute
}

func NewFile(path string) (*File, error) {
	attr, err := LoadAttribute(path)
	if err != nil {
		if errors.Is(err, xattr.ENOATTR) {
			attr = make(Attribute)
		} else {
			return nil, err
		}
	}
	return &File{
		path: path,
		attr: attr,
	}, nil
}

// Hashes the file.
func (r *File) hash(hash hash.Hash) error {
	//Open File
	src, err := r.open()
	if err != nil {
		return err
	}
	defer src.Close()
	//hashing
	hash.Reset()
	return hashes.Hash(src, hash)
}

// Hashes the file and writes a file copy to dst simultaneously.
func (r *File) hashCopy(dst io.Writer, hash hash.Hash) (written int64, err error) {
	//Open File
	src, err := r.open()
	if err != nil {
		return 0, err
	}
	defer src.Close()
	//hashing
	hash.Reset()
	return hashes.HashCopy(dst, src, hash)
}

// Wraps os.Open()
func (r *File) open() (*os.File, error) {
	return os.Open(r.path)
}

// Returns all records for the file. The received map is shared with
// the object.
func (r *File) Attributes() Attribute {
	return r.attr
}

func (r *File) InvalidateOutdatedEntries() error {
	//Nothing to do if attribute is empty
	if len(r.attr) < 1 {
		return nil
	}
	//Initialize
	hashMap := make(map[hashes.Algo]hash.Hash)
	//Open File
	src, err := r.open()
	if err != nil {
		return err
	}
	//Fill hashMap with all necessary hashing algorithms
	for _, rec := range r.attr {
		if hashMap[rec.HashAlgo] == nil {
			hashMap[rec.HashAlgo] = rec.HashAlgo.New()
		}
	}
	//Calculate hashes for file
	if err := hashes.MultiHash(src, hashMap); err != nil {
		return err
	}
	//Invalidate entries with hash sums that do not match
	for _, rec := range r.attr {
		if rec.Checksum != fmt.Sprintf("%x", hashMap[rec.HashAlgo].Sum(nil)) {
			rec.Valid = false
		}
	}
	//Store records in xattrs
	if err := r.attr.FStore(src); err != nil {
		return err
	}
	return nil
}

func (r *File) CreateRecord(name string, hashAlgo hashes.Algo) error {
	//Check if identifier is already occupied
	if rec := r.attr[name]; rec != nil {
		return &fs.PathError{
			Op:   "Name conflict:",
			Path: r.path,
			Err:  fmt.Errorf("Cannot create record, identifier \"%s\" already exists", name),
		}
	}
	//Hash the file
	hash := hashAlgo.New()
	if err := r.hash(hash); err != nil {
		return fmt.Errorf("Could not hash file: %s", err)
	}
	//Create new record
	rec := NewRecord()
	rec.Checksum = fmt.Sprintf("%x", hash.Sum(nil))
	rec.HashAlgo = hashAlgo
	rec.Valid = true
	//Append new record to records
	r.attr[name] = rec
	//Store records in xattrs
	if err := r.attr.Store(r.path); err != nil {
		return fmt.Errorf("Could not store xattrs: %s", err)
	}
	return nil
}

func (r *File) String() string {
	entity := struct {
		Path    string
		Records Attribute
	}{
		Path:    r.path,
		Records: r.attr,
	}
	payload, err := json.MarshalIndent(&entity, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("Error marshalling json: %s", err))
	}
	return fmt.Sprintf("%s\n", payload)
}
