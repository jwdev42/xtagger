package program

import (
	"context"
	"crypto/sha256"
	"github.com/jwdev42/xtagger/internal/cli"
	"github.com/jwdev42/xtagger/internal/io/filesystem"
	"github.com/jwdev42/xtagger/internal/record"
	"io/fs"
	"os"
)

func runCreate(cmdline *cli.CommandLine) error {
	for _, path := range cmdline.Paths() {
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}
		if info.Mode()&fs.ModeSymlink == fs.ModeSymlink {

		} else if info.Mode()&fs.ModeDir == fs.ModeDir {
			err = tagDir(cmdline, path)
		} else {
			err = tagFile(cmdline, path)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func tagDir(cmdline *cli.CommandLine, path string) error {
	ctx, cancel := context.WithCancel(context.Background())
	ms := filesystem.NewMultistat(cmdline.FlagRecursive(), cmdline.FlagFollowSymlinks())
	nextInfo, err := ms.Run(ctx, path)
	stopOnError := func(err error) error {
		//Cancel context
		cancel()
		//Drain channel
		for {
			_, ok := <-nextInfo
			if !ok {
				break
			}
		}
		return err
	}
	if err != nil {
		return err
	}
	info := <-nextInfo
	for info != nil {
		if info.Err() != nil {
			return stopOnError(info.Err())
		}
		if err := tagFile(cmdline, info.Path()); err != nil {
			return stopOnError(err)
		}
		info = <-nextInfo
	}
	return nil
}

func tagFile(cmdline *cli.CommandLine, path string) error {
	f, err := record.NewFile(path, sha256.New())
	if err != nil {
		return err
	}
	if err := f.CreateRecord(cmdline.FlagName()); err != nil {
		return err
	}
	return nil
}
