package program

import (
	"github.com/jwdev42/xtagger/internal/global"
	"github.com/jwdev42/xtagger/internal/io/filesystem"
	"github.com/jwdev42/xtagger/internal/record"
	"io/fs"
	"path/filepath"
)

func untagFile(parent string, dirEnt fs.DirEntry, opts *filesystem.WalkDirOpts) error {
	return global.FilterSoftError(record.PurgeFile(filepath.Join(parent, dirEnt.Name())))
}
