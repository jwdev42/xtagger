package record

import (
	"encoding/json"
	"errors"
	"fmt"
	iio "github.com/jwdev42/xtagger/internal/io"
	"github.com/pkg/xattr"
	"hash"
	"io"
	"io/fs"
	"os"
)

// Accesses and manipulates the records of a particular file
type File struct {
	path     string
	attr     Attribute
	fileHash hash.Hash
}

func NewFile(path string, hash hash.Hash) (*File, error) {
	attr, err := LoadAttribute(path)
	if err != nil {
		if errors.Is(err, xattr.ENOATTR) {
			attr = make(Attribute)
		} else {
			return nil, err
		}
	}
	return &File{
		path:     path,
		attr:     attr,
		fileHash: hash,
	}, nil
}

// Hashes the file.
func (r *File) hash() error {
	//Open File
	src, err := r.open()
	if err != nil {
		return err
	}
	defer src.Close()
	r.fileHash.Reset()
	return iio.Hash(src, r.fileHash)
}

// Hashes the file and writes a file copy to dst simultaneously.
func (r *File) hashCopy(dst io.Writer) (written int64, err error) {
	//Open File
	src, err := r.open()
	if err != nil {
		return 0, err
	}
	defer src.Close()
	r.fileHash.Reset()
	return iio.HashCopy(dst, src, r.fileHash)
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

// Calculates the file's checksum, then checks each record if it is still valid.
// Updates each record's field "valid". Returns only valid records.
// Returns an empty slice if no valid records were found.
func (r *File) Validate() (Attribute, error) {
	//Hash the file
	if err := r.hash(); err != nil {
		return nil, err
	}
	//Check all existing records for matching hashes
	validated := make(Attribute)
	for name, rec := range r.attr {
		if rec.Checksum == fmt.Sprintf("%x", r.fileHash.Sum(nil)) {
			rec.Valid = true
			validated[name] = rec
		} else {
			rec.Valid = false
		}
	}
	//Store records in xattrs
	if err := r.attr.Store(r.path); err != nil {
		return nil, err
	}
	return validated, nil
}

func (r *File) CreateRecord(name string) error {
	//Check if identifier is already occupied
	if rec := r.attr[name]; rec != nil {
		return &fs.PathError{
			Op:   "Name conflict:",
			Path: r.path,
			Err:  fmt.Errorf("Cannot create record, identifier \"%s\" already exists", name),
		}
	}
	//Hash the file
	if err := r.hash(); err != nil {
		return err
	}
	//Create new record
	rec := NewRecord()
	rec.Checksum = fmt.Sprintf("%x", r.fileHash.Sum(nil))
	rec.Valid = true
	//Append new record to records
	r.attr[name] = rec
	//Store records in xattrs
	if err := r.attr.Store(r.path); err != nil {
		return err
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
