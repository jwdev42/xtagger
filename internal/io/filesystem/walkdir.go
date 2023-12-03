package filesystem

import (
	"errors"
	"fmt"
	"github.com/jwdev42/logger"
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
	SymlinksRejectAll  SymlinkBehaviour = iota //Reject all symlinks.
	SymlinksRejectDirs                         //Do not resolve symlinks to directories, but follow symlinks to files.
	SymlinksRejectNone                         //Follow all symlinks.
)

const (
	QuotaDisabled QuotaMode = iota //Disable the quota.
	QuotaCutoff                    //Enable the quota, will stop the current operation if the threshold is exeeded.
	QuotaSkip                      //Enable the quota, will skip the current file if the threshold is exeeded.
)

type ErrorBehaviour int
type SymlinkBehaviour int
type QuotaMode int
type FileExaminer func(parent string, info fs.FileInfo) error

type Context struct {
	SymlinkMode    SymlinkBehaviour
	DupeDetector   data.DupeDetector
	DetectorHash   hash.Hash
	quotaMode      QuotaMode
	quota          int64 //Quota left in bytes
	symlinkCounter int
}

func (r *Context) SetQuota(mode QuotaMode, quota int64) {
	r.quotaMode = mode
	r.quota = quota
}

func WalkDir(path string, opts *Context, fileEx FileExaminer) error {
	//Stat directory
	info, err := os.Lstat(path)
	if err != nil {
		return global.FilterSoftError(err)
	}
	//Check if path is a directory
	if !(info.IsDir() || info.Mode()&(fs.ModeDir|fs.ModeSymlink) != 0) {
		return fmt.Errorf("Not a directory: %s", path)
	}
	//Evaluate Symlink
	if info.Mode()&fs.ModeSymlink != 0 {
		//Check if symlinks are to follow
		if opts.SymlinkMode != SymlinksRejectNone {
			logger.Default().Infof("Skipping directory symlink: %s", path)
			return nil
		}
		//Symlink counter
		if opts.symlinkCounter >= 40 {
			return errors.New("Symlink limit reached")
		}
		opts.symlinkCounter++
		defer func() {
			opts.symlinkCounter--
			logger.Default().Debugf("Symlink counter: %02d", opts.symlinkCounter)
		}()
		logger.Default().Debugf("Symlink counter: %02d", opts.symlinkCounter)
	}
	//Read directory entries
	dirEnts, errs := readDirEnts(path)
	if len(errs) > 0 {
		for i, err := range errs {
			if len(errs)-i > 1 {
				logger.Default().Error(err)
				continue
			}
			if global.FilterSoftError(err) != nil {
				return err
			}
		}
	}
	//Loop directory entries
	for _, dirEnt := range dirEnts {
		if dirEnt.IsDir() || dirEnt.Type()&(fs.ModeDir|fs.ModeSymlink) != 0 {
			//Recurse into subdirectory
			if err := WalkDir(filepath.Join(path, dirEnt.Name()), opts, fileEx); err != nil {
				return err
			}
		} else {
			//examine file
			info, err := dirEnt.Info()
			if err != nil {
				if err := global.SoftErrorf("Could not read FileInfo: %s", err); err != nil {
					return err
				}
				continue
			}
			if err := examineFile(path, info, opts, fileEx); err != nil {
				if errors.Is(err, fs.SkipDir) {
					logger.Default().Debugf("walkDir: File executor returned fs.SkipDir, skipping rest of directory: %s", path)
					return nil
				}
				return err
			}
		}
	}
	return nil
}

func ExamineFile(parent string, info fs.FileInfo, opts *Context, fileEx FileExaminer) error {
	err := examineFile(parent, info, opts, fileEx)
	if errors.Is(err, fs.SkipDir) {
		return nil
	}
	return err
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

func examineFile(parent string, info fs.FileInfo, opts *Context, fileEx FileExaminer) error {
	path := filepath.Join(parent, info.Name())
	//Use DupeDetector for files if available
	if opts.DupeDetector != nil {
		realPath, err := filepath.EvalSymlinks(path)
		if err != nil {
			return global.FilterSoftError(err)
		}
		if err := opts.DupeDetector.Register(strings.NewReader(realPath), opts.DetectorHash); err != nil {
			logger.Default().Debugf("examineFile: DupeDetector detected already processed file, skipping: %s", path)
			return nil
		}
	}
	//Check quota on regular files
	if opts.quotaMode != QuotaDisabled && info.Mode().IsRegular() {
		opts.quota = opts.quota - info.Size()
		if opts.quota < 0 {
			switch opts.quotaMode {
			case QuotaCutoff:
				logger.Default().Debugf("examineFile: File exceeds quota in mode QuotaCutoff, aborting: %s", path)
				return fs.SkipAll
			case QuotaSkip:
				logger.Default().Debugf("examineFile: File exceeds quota in mode QuotaSkip, skipping: %s", path)
				return nil
			default:
				panic(fmt.Errorf("examineFile: Unknown QuotaMode: %d", opts.quotaMode))
			}
		}
	}
	//Call file executor function
	return fileEx(parent, info)
}
