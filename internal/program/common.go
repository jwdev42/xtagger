package program

import (
	"crypto/sha256"
	"github.com/jwdev42/xtagger/internal/cli"
	"github.com/jwdev42/xtagger/internal/data"
	"github.com/jwdev42/xtagger/internal/io/filesystem"
)

func createWalkDirOpts(cmdline *cli.CommandLine) *filesystem.WalkDirOpts {
	var opts = &filesystem.WalkDirOpts{
		DupeDetector: make(data.DupeDetector),
		DetectorHash: sha256.New(),
	}
	if cmdline.FlagFollowSymlinks() {
		opts.SymlinkMode = filesystem.SymlinksRejectNone
	}
	return opts
}
