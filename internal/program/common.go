package program

import (
	"context"
	"crypto/sha256"
	"github.com/jwdev42/xtagger/internal/cli"
	"github.com/jwdev42/xtagger/internal/data"
	"github.com/jwdev42/xtagger/internal/io/filesystem"
	"github.com/jwdev42/xtagger/internal/io/printer"
	"io/fs"
	"sync"
)

var commandLine *cli.CommandLine
var printMe *printer.Printer

func createContext(detectProcessedFiles bool) *filesystem.Context {
	var opts = new(filesystem.Context)
	if commandLine.FlagFollowSymlinks() {
		opts.SymlinkMode = filesystem.SymlinksRejectNone
	}
	if quota := commandLine.SizeQuota(); quota > 0 {
		if commandLine.FlagQuotaContinue() {
			opts.SetQuota(filesystem.QuotaSkip, quota)
		}
		opts.SetQuota(filesystem.QuotaCutoff, quota)
	}
	if detectProcessedFiles {
		opts.DupeDetector = make(data.DupeDetector)
		opts.DetectorHash = sha256.New()
	}
	return opts
}

func wrapFileExaminer(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup, errs chan<- error, payload filesystem.FileExaminer) filesystem.FileExaminer {
	return func(parent string, info fs.FileInfo) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := payload(parent, info); err != nil {
				cancel()
				errs <- err
			}
		}()
		return nil
	}
}
