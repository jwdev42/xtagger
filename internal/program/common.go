//This file is part of xtagger. ©2023 Jörg Walter.
//This program is free software: you can redistribute it and/or modify
//it under the terms of the GNU General Public License as published by
//the Free Software Foundation, either version 3 of the License, or
//(at your option) any later version.
//
//This program is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//GNU General Public License for more details.
//
//You should have received a copy of the GNU General Public License
//along with this program.  If not, see <https://www.gnu.org/licenses/>.

package program

import (
	"context"
	"crypto/sha256"
	"github.com/jwdev42/xtagger/internal/cli"
	"github.com/jwdev42/xtagger/internal/data"
	"github.com/jwdev42/xtagger/internal/xio/filesystem"
	"github.com/jwdev42/xtagger/internal/xio/printer"
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
