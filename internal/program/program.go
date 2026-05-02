//This file is part of xtagger. ©2023-2026 Jörg Walter.
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
	"errors"
	"fmt"
	"github.com/jwdev42/xtagger/internal/config"
	"github.com/jwdev42/xtagger/internal/logging"
	"github.com/jwdev42/xtagger/internal/xio/filesystem"
	"github.com/jwdev42/xtagger/internal/xio/printer"
	"log/slog"
	"os"
	"sync"
)

// Program entry point called by main().
func Run() error {
	var err error
	// Setup logger
	dynamicLogLevel := setupDefaultLogger()
	// Parse command line
	commandLine, err = config.ParseCommandLine()
	if err != nil {
		return fmt.Errorf("Command line error: %s", err)
	}
	// Setup context
	ctx := context.Background()
	// Adjust log level
	dynamicLogLevel.Set(commandLine.FlagLogLevel())
	// Setup printer
	printMe = printer.NewPrinter(os.Stdout)
	// Generate PushOpts
	pushOpts := pushOptsFromCommandLine(commandLine)
	// Execute command-specific branch
	switch command := commandLine.Command(); command {
	case config.CommandTag:
		execPayload(ctx, pushOpts, commandLine.Threads(), tagFile, commandLine.Paths()...)
	case config.CommandPrint:
		execPayload(ctx, pushOpts, commandLine.Threads(), printFile, commandLine.Paths()...)
	case config.CommandUntag:
		execPayload(ctx, pushOpts, commandLine.Threads(), untagFile, commandLine.Paths()...)
	case config.CommandInvalidate:
		execPayload(ctx, pushOpts, commandLine.Threads(), invalidateFile, commandLine.Paths()...)
	case config.CommandRevalidate:
		execPayload(ctx, pushOpts, commandLine.Threads(), revalidateFile, commandLine.Paths()...)
	case config.CommandLicenses:
		printLicenses()
	default:
		return fmt.Errorf("Unknown command \"%s\"", command)
	}
	return nil
}

// Setup default logger with dynamic leveler, return LevelVar
func setupDefaultLogger() *slog.LevelVar {
	levelSwitch := &slog.LevelVar{} // log level LevelInfo
	defaultLogger := slog.New(slog.NewTextHandler(os.Stderr,
		&slog.HandlerOptions{
			Level:       levelSwitch,
			ReplaceAttr: logging.ReplaceLogLevelNames,
		}))
	slog.SetDefault(defaultLogger)
	return levelSwitch
}

func defaultErrorHandler(ctx context.Context, err error) {
	slog.ErrorContext(ctx, err.Error())
}

func pushOptsFromCommandLine(cmd *config.CommandLine) filesystem.PushOpts {
	return filesystem.PushOpts{
		FollowSymlinks: cmd.FlagFollowSymlinks(),
		Recursive:      !cmd.ForbidRecursion(),
	}
}

func execPayload(
	ctx context.Context,
	opts filesystem.PushOpts,
	threads int,
	payload func(*filesystem.Meta) error,
	paths ...string) error {
	// Setup error handler
	eh, cancelEH := logging.NewErrorHandler(ctx, 10, defaultErrorHandler)
	// Use closure to ensure a finished error handler
	func() {
		defer cancelEH()
		// Setup WaitGroup
		wg := &sync.WaitGroup{}
		// Stat files
		metas := filesystem.PushMetas(ctx, eh, wg, opts, paths...)
		// Setup semaphore
		semaphore := make(chan struct{}, threads)
		// Run payload on files
		for meta := range metas {
			wg.Go(func() {
				semaphore <- struct{}{}
				defer func() { <-semaphore }()
				eh.Error(payload(meta))
			})
		}
		// Wait for payloads to finish
		wg.Wait()
	}()
	// Examinate error count
	errs := eh.Errors()
	switch errs {
	case 0:
		return nil
	case 1:
		return errors.New("An error occured during command execution")
	}
	return fmt.Errorf("%d errors occured during command execution", errs)
}
