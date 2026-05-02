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
	// Setup context
	ctx := context.Background()
	// Create preferences
	prefs := config.DefaultPreferences()
	// Setup logger
	setupDefaultLogger(prefs.LogLevel)
	// Parse command line
	if err := config.ParseCommandLine(prefs); err != nil {
		return fmt.Errorf("Command line error: %s", err)
	}
	// Setup printer
	printMe = printer.NewPrinter(os.Stdout)
	// Execute command-specific branch
	switch prefs.Command {
	case config.CommandTag:
		execPayload(ctx, prefs, tagFile)
	case config.CommandPrint:
		execPayload(ctx, prefs, printFile)
	case config.CommandUntag:
		execPayload(ctx, prefs, untagFile)
	case config.CommandInvalidate:
		execPayload(ctx, prefs, invalidateFile)
	case config.CommandRevalidate:
		execPayload(ctx, prefs, revalidateFile)
	case config.CommandLicenses:
		printLicenses()
	default:
		return fmt.Errorf("Unknown command \"%s\"", prefs.Command)
	}
	return nil
}

// Setup default logger with dynamic leveler level
func setupDefaultLogger(level *slog.LevelVar) {
	levelSwitch := &slog.LevelVar{} // log level LevelInfo
	defaultLogger := slog.New(slog.NewTextHandler(os.Stderr,
		&slog.HandlerOptions{
			Level:       levelSwitch,
			ReplaceAttr: logging.ReplaceLogLevelNames,
		}))
	slog.SetDefault(defaultLogger)
}

func defaultErrorHandler(ctx context.Context, err error) {
	slog.ErrorContext(ctx, err.Error())
}

func pushOptsFromPrefs(prefs *config.Preferences) filesystem.PushOpts {
	return filesystem.PushOpts{
		FollowSymlinks: prefs.FollowSymlinks,
		Recursive:      prefs.UseRecursion,
	}
}

type payloadRuntime struct {
	ctx   context.Context
	eh    *logging.ErrorHandler
	prefs *config.Preferences
}

type payloadFunc func(*payloadRuntime, *filesystem.Meta) error

func execPayload(ctx context.Context, prefs *config.Preferences, payload payloadFunc) error {
	// Setup error handler
	eh, cancelEH := logging.NewErrorHandler(ctx, 10, defaultErrorHandler)
	// Use closure to ensure a finished error handler
	func() {
		defer cancelEH()
		// Create runtime object for payload
		rt := &payloadRuntime{
			ctx:   ctx,
			eh:    eh,
			prefs: prefs,
		}
		// Setup WaitGroup
		wg := &sync.WaitGroup{}
		// Stat files
		metas := filesystem.PushMetas(ctx, eh, wg, pushOptsFromPrefs(prefs), prefs.Paths...)
		// Setup semaphore
		semaphore := make(chan struct{}, prefs.Threads)
		// Run payload on files
		for meta := range metas {
			wg.Go(func() {
				semaphore <- struct{}{}
				defer func() { <-semaphore }()
				eh.Error(payload(rt, meta))
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
