package record

import (
	"errors"
	"github.com/pkg/xattr"
	"os"
)

// PurgeFile removes all extended attributes by xtagger from the file at path.
func PurgeFile(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return &os.PathError{
			Op:   "PurgeFile",
			Path: path,
			Err:  errors.New("Not a regular file"),
		}
	}
	return xattr.LRemove(path, attrName)
}
