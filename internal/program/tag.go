package program

import (
	"context"
	"github.com/jwdev42/xtagger/internal/cli"
	"github.com/jwdev42/xtagger/internal/io/filesystem"
	"github.com/jwdev42/xtagger/internal/record"
)

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
	f, err := record.NewFile(path)
	if err != nil {
		return err
	}
	if err := f.CreateRecord(cmdline.FlagName(), cmdline.FlagHash()); err != nil {
		return err
	}
	return nil
}
