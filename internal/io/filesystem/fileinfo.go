package filesystem

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

type FileInfo struct {
	err   error       //Stores an error of one occures during stat
	path  string      //Absolute path
	size  int64       //The file's size in bytes, system-specific value if it's not a regular file
	mode  fs.FileMode //File mode bits
	mtime time.Time   //Modification time
}

// Calls stat on file name and returns a FileInfo object.
// If an error occurs, it can be detected through the second return value and the
// FileInfo's Error() method.
func Stat(name string) (*FileInfo, error) {
	return stat(name, true)
}

// Like Stat, but doesn't folow symlinks.
func Lstat(name string) (*FileInfo, error) {
	return stat(name, false)
}

func stat(name string, followSymlinks bool) (*FileInfo, error) {
	onErr := func(err error) (*FileInfo, error) {
		return &FileInfo{err: err}, err
	}
	//Input checks
	if name == "" {
		return onErr(&fs.PathError{Op: "stat", Path: "", Err: errors.New("Cannot stat an empty path")})
	}
	//Construct absolute path
	path, err := filepath.Abs(name)
	if err != nil {
		return onErr(err)
	}
	//Execute stat syscall
	var osInfo fs.FileInfo
	if followSymlinks {
		osInfo, err = os.Stat(path)
	} else {
		osInfo, err = os.Lstat(path)
	}
	if err != nil {
		return onErr(err)
	}
	//Build FileInfo
	return &FileInfo{
		path:  path,
		size:  osInfo.Size(),
		mode:  osInfo.Mode(),
		mtime: osInfo.ModTime(),
	}, nil
}

// Returns the base name.
func (r *FileInfo) Name() string {
	r.panicOnError()
	return filepath.Base(r.path)
}

// Returns the parent directory.
func (r *FileInfo) Dir() string {
	r.panicOnError()
	return filepath.Dir(r.path)
}

// Returns the absolute file path.
func (r *FileInfo) Path() string {
	r.panicOnError()
	return r.path
}

// Returns the file's size in bytes if it is a regular file.
// Returns a system-specific value otherwise.
func (r *FileInfo) Size() int64 {
	r.panicOnError()
	return r.size
}

// Returns the file's mode bits.
func (r *FileInfo) Mode() fs.FileMode {
	r.panicOnError()
	return r.mode
}

// Returns the file's mtime.
func (r *FileInfo) ModTime() time.Time {
	r.panicOnError()
	return r.mtime
}

// Returns true if the underlying file is a directory.
func (r *FileInfo) IsDir() bool {
	r.panicOnError()
	return r.mode.IsDir()
}

// Always returns nil, method must be there to implement fs.FileInfo.
func (r *FileInfo) Sys() any {
	r.panicOnError()
	return nil
}

// Returns an error if one occured during stat, returns nil otherwise.
// Always check for errors on a newly generated FileInfo!
func (r *FileInfo) Err() error {
	return r.err
}

func (r *FileInfo) panicOnError() {
	if r.err != nil {
		panic("FileInfo cannot be used if it had an error. Please use method Error() to check for errors.")
	}
}
