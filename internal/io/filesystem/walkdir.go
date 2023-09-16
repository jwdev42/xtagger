package filesystem

import (
	"errors"
	"fmt"
	"github.com/jwdev42/xtagger/internal/data"
	"github.com/jwdev42/xtagger/internal/global"
	"hash"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	OnErrorAbort ErrorBehaviour = iota //Abort on error
	OnErrorLog                         //Log the error and continue
)

const (
	SymlinksRejectAll  SymlinkBehaviour = iota //Reject all symlinks
	SymlinksRejectDirs                         //Do not resolve symlinks to directories, but follow symlinks to files
	SymlinksRejectNone                         //Follow all symlinks
)

type ErrorBehaviour int
type SymlinkBehaviour int
type FileExaminer func(path string, dirEnt fs.DirEntry, opts *WalkDirOpts) error

type WalkDirOpts struct {
	ErrorMode    ErrorBehaviour
	SymlinkMode  SymlinkBehaviour
	DupeDetector data.DupeDetector
	DetectorHash hash.Hash
}

func WrapWalkDirFunc(exec fs.WalkDirFunc, skipOnError bool) fs.WalkDirFunc {
	walkDir := func(path string, d fs.DirEntry, err error) error {
		if err := exec(path, d, err); err != nil {
			return wrapWalkDirFuncEvalErr(err, path, skipOnError)
		}
		return nil
	}
	return walkDir
}

func wrapWalkDirFuncEvalErr(err error, path string, skipOnError bool) error {
	switch err {
	case fs.SkipAll, fs.SkipDir:
		return err
	}
	if pathErr, ok := err.(*fs.PathError); ok {
		pathErr.Path = filepath.Join(path, pathErr.Path)
		err = pathErr
	}
	if skipOnError {
		global.DefaultLogger.Error(err)
		return nil
	}
	return err
}

func WalkDir(path string, opts *WalkDirOpts, fileEx FileExaminer) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("Not a directory: %s", path)
	}
	return walkDir(path, opts, fileEx)
}

func walkDir(path string, opts *WalkDirOpts, fileEx FileExaminer) error {
	//Stat directory
	info, err := os.Lstat(path)
	if err != nil {
		if opts.ErrorMode == OnErrorLog {
			global.DefaultLogger.Error(err)
		} else {
			return err
		}
	}
	//Evaluate Symlink
	if info.Mode()&fs.ModeSymlink != 0 {
		if opts.SymlinkMode != SymlinksRejectNone {
			global.DefaultLogger.Info("Skipping directory symlink %s", path)
			return nil
		}
		realPath, err := filepath.EvalSymlinks(path)
		if err != nil {
			err = fmt.Errorf("Failed to resolve symlinks for path %s: %s", path, err)
			if opts.ErrorMode == OnErrorLog {
				global.DefaultLogger.Error(err)
				return nil
			}
			return err
		}
		if err := opts.DupeDetector.Register(strings.NewReader(realPath), opts.DetectorHash); err != nil {
			if errors.Is(err, data.DupeDetected) {
				global.DefaultLogger.Errorf("Symlink %s points to already evaluated path %s", path, realPath)
				return nil
			} else {
				if opts.ErrorMode == OnErrorLog {
					global.DefaultLogger.Error(err)
					return nil
				}
				return err
			}
		}
	} else {
		if err := opts.DupeDetector.Register(strings.NewReader(path), opts.DetectorHash); err != nil {
			if errors.Is(err, data.DupeDetected) {
				global.DefaultLogger.Errorf("Path already evaluated %s", path)
			} else {
				if opts.ErrorMode == OnErrorLog {
					global.DefaultLogger.Error(err)
					return nil
				}
				return err
			}
		}
	}
	//Read directory entries
	dirEnts, errs := readDirEnts(path)
	if len(errs) > 0 {
		if opts.ErrorMode == OnErrorLog {
			for _, err := range errs {
				global.DefaultLogger.Error(err)
			}
		} else {
			return errs[0]
		}
	}
	for _, dirEnt := range dirEnts {
		newPath := filepath.Join(path, dirEnt.Name())
		if dirEnt.IsDir() || dirEnt.Type()&(fs.ModeDir|fs.ModeSymlink) != 0 {
			if err := walkDir(newPath, opts, fileEx); err != nil {
				return err
			}
		} else {
			if err := fileEx(path, dirEnt, opts); err != nil {
				if errors.Is(err, fs.SkipDir) {
					global.DefaultLogger.Debugf("Skipping rest of directory %s", newPath)
					return nil
				}
				return err
			}
		}
	}
	return nil
}

func readDirEnts(path string) ([]fs.DirEntry, []error) {
	errs := make([]error, 0)
	dirEnts := make([]fs.DirEntry, 0)
	f, err := os.Open(path)
	if err != nil {
		return nil, []error{err}
	}
	defer f.Close()
	for {
		entries, err := f.ReadDir(1024)
		if len(entries) > 0 {
			dirEnts = append(dirEnts, entries...)
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			errs = append(errs, err)
		}
	}
	return dirEnts, errs
}
