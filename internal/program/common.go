package program

import (
	"crypto/sha256"
	"github.com/jwdev42/xtagger/internal/cli"
	"github.com/jwdev42/xtagger/internal/data"
	"github.com/jwdev42/xtagger/internal/io/filesystem"
)

func createWalkDirOpts(cmdline *cli.CommandLine, detectProcessedFiles bool) *filesystem.WalkDirOpts {
	var opts = new(filesystem.WalkDirOpts)
	if cmdline.FlagFollowSymlinks() {
		opts.SymlinkMode = filesystem.SymlinksRejectNone
	}
	if detectProcessedFiles {
		opts.DupeDetector = make(data.DupeDetector)
		opts.DetectorHash = sha256.New()
	}
	return opts
}
