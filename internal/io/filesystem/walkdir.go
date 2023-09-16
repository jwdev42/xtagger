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
	OnErrorAbort ErrorBehaviour = iota //Abort on soft errors
	OnErrorLog                         //Log the soft error and continue
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
	ErrorMode      ErrorBehaviour
	SymlinkMode    SymlinkBehaviour
	DupeDetector   data.DupeDetector
	DetectorHash   hash.Hash
	symlinkCounter int
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
	return walkDir(path, opts, fileEx)
}

func walkDir(path string, opts *WalkDirOpts, fileEx FileExaminer) error {
	softErr := func(err error) error {
		if err == nil {
			return err
		}
		if opts.ErrorMode == OnErrorLog {
			global.DefaultLogger.Error(err)
		} else {
			return err
		}
		return nil
	}
	//Stat directory
	info, err := os.Lstat(path)
	if softErr(err) != nil {
		return err
	}
	//Check if path is a directory
	if !(info.IsDir() || info.Mode()&(fs.ModeDir|fs.ModeSymlink) != 0) {
		return fmt.Errorf("Not a directory: %s", path)
	}
	//Evaluate Symlink
	if info.Mode()&fs.ModeSymlink != 0 {
		//Check if symlinks are to follow
		if opts.SymlinkMode != SymlinksRejectNone {
			global.DefaultLogger.Infof("Skipping directory symlink: %s", path)
			return nil
		}
		//Symlink counter
		if opts.symlinkCounter >= 40 {
			return errors.New("Symlink limit reached")
		}
		opts.symlinkCounter++
		defer func() {
			opts.symlinkCounter--
			global.DefaultLogger.Debugf("Symlink counter: %02d", opts.symlinkCounter)
		}()
		global.DefaultLogger.Debugf("Symlink counter: %02d", opts.symlinkCounter)
	}
	//Read directory entries
	dirEnts, errs := readDirEnts(path)
	if len(errs) > 0 {
		for i, err := range errs {
			if len(errs)-i > 1 {
				global.DefaultLogger.Error(err)
				continue
			}
			if softErr(err) != nil {
				return err
			}
		}
	}
	//Loop directory entries
	for _, dirEnt := range dirEnts {
		newPath := filepath.Join(path, dirEnt.Name())
		if dirEnt.IsDir() || dirEnt.Type()&(fs.ModeDir|fs.ModeSymlink) != 0 {
			//Descend into darkness
			if err := walkDir(newPath, opts, fileEx); err != nil {
				return err
			}
		} else {
			//Use DupeDetector for files if available
			if opts.DupeDetector != nil {
				realPath, err := filepath.EvalSymlinks(newPath)
				if err != nil {
					return err
				}
				if err := opts.DupeDetector.Register(strings.NewReader(realPath), opts.DetectorHash); err != nil {
					global.DefaultLogger.Debugf("DupeDetector: Skipping already processed file: %s", newPath)
					continue
				}
			}
			//Call file executor function
			if err := fileEx(path, dirEnt, opts); err != nil {
				if errors.Is(err, fs.SkipDir) {
					global.DefaultLogger.Debugf("File executor returned fs.SkipDir, skipping rest of directory: %s", path)
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
