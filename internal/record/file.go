package record

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jwdev42/xtagger/internal/hashes"
	iio "github.com/jwdev42/xtagger/internal/io"
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
	return iio.Hash(src, hash)
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
	return iio.HashCopy(dst, src, hash)
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
	//Create hashcache as different records could use different hash algorithms
	hashCache := make(map[hashes.Algo]string)
	validated := make(Attribute)
	//Recalculate and compare hashes for all records
	for name, rec := range r.attr {
		targetSum, ok := hashCache[rec.HashAlgo]
		if !ok {
			//Hash the file if hash not in cache
			hash := rec.HashAlgo.New()
			if err := r.hash(hash); err != nil {
				return nil, err
			}
			//Write hashsum to hashcache
			hashCache[rec.HashAlgo] = fmt.Sprintf("%x", hash.Sum(nil))
			targetSum = hashCache[rec.HashAlgo]
		}
		if rec.Checksum == targetSum {
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
