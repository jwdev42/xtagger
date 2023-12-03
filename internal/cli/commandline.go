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

package cli

import (
	"flag"
	"fmt"
	"github.com/jwdev42/logger"
	"github.com/jwdev42/xtagger/internal/hashes"
	"os"
	"slices"
	"strconv"
)

// Represents a parsed command line argument set.
type CommandLine struct {
	command             Command //Specified command
	paths               []string
	names               []string
	flagLogLevel        logger.Level //parsed loglevel
	flagFollowSymlinks  bool
	flagHash            hashes.Algo
	flagQuitOnSoftError bool
	flagMultiThread     bool
	flagPrint0          bool
	printRecords        bool
	forbidRecursion     bool
	quota               int64
	quotaContinue       bool
	tagConstraint       TagConstraint
	untagConstraint     UntagConstraint
	printConstraint     PrintConstraint
}

func (r *CommandLine) Command() Command {
	return r.command
}

func (r *CommandLine) Paths() []string {
	return r.paths
}

func (r *CommandLine) Names() []string {
	return r.names
}

func (r *CommandLine) FlagFollowSymlinks() bool {
	return r.flagFollowSymlinks
}

func (r *CommandLine) FlagLogLevel() logger.Level {
	return r.flagLogLevel
}

func (r *CommandLine) FlagHash() hashes.Algo {
	return r.flagHash
}

func (r *CommandLine) FlagQuitOnSoftError() bool {
	return r.flagQuitOnSoftError
}

func (r *CommandLine) FlagMultiThread() bool {
	return r.flagMultiThread
}

func (r *CommandLine) FlagPrint0() bool {
	return r.flagPrint0
}

func (r *CommandLine) FlagPrintRecords() bool {
	return r.printRecords
}

func (r *CommandLine) ForbidRecursion() bool {
	return r.forbidRecursion
}

func (r *CommandLine) SizeQuota() int64 {
	return r.quota
}

func (r *CommandLine) FlagQuotaContinue() bool {
	return r.quotaContinue
}

func (r *CommandLine) TagConstraint() TagConstraint {
	return r.tagConstraint
}

func (r *CommandLine) UntagConstraint() UntagConstraint {
	return r.untagConstraint
}

func (r *CommandLine) PrintConstraint() PrintConstraint {
	return r.printConstraint
}

func (r *CommandLine) parseHashAlgo(input string) error {
	hash, err := hashes.ParseAlgo(input)
	if err != nil {
		return err
	}
	r.flagHash = hash
	return nil
}

func (r *CommandLine) parseSizeStatement(input string) error {
	var base = make([]rune, len(input))
	var suffix string
	//Parse size limit integer
	for i, ch := range input {
		if !(ch >= 0x30 && ch <= 0x39) {
			base = base[:i]
			suffix = input[i:]
			break
		}
		base[i] = ch
	}
	sizeLimit, err := strconv.ParseInt(string(base), 10, 64)
	if err != nil {
		return fmt.Errorf("Could not parse size statement: %s", err)
	}
	//Parse optional size suffix
	const kib = 1024
	const mib = kib * 1024
	const gib = mib * 1024
	const tib = gib * 1024
	switch suffix {
	case "":
		r.quota = sizeLimit
	case "K":
		r.quota = sizeLimit * kib
	case "M":
		r.quota = sizeLimit * mib
	case "G":
		r.quota = sizeLimit * gib
	case "T":
		r.quota = sizeLimit * tib
	default:
		return fmt.Errorf("Could not parse size statement: Unknown suffix: \"%s\"", suffix)
	}
	return nil
}

// Parses and validates command line arguments.
func ParseCommandLine() (*CommandLine, error) {
	//Stage 1: Parse flags
	var cmd = new(CommandLine)
	cmd.flagHash = hashes.SHA256 //Default hash algorithm
	var logLevel = logger.LevelFlag(logger.LevelError)
	main := flag.NewFlagSet("main", flag.ContinueOnError)
	main.Var(&logLevel, "ll", "Set the loglevel")
	main.BoolVar(&cmd.flagFollowSymlinks, "symlinks", false, "Program follows symlinks if true")
	main.Func("hash", "Specify the hashing algorithm", cmd.parseHashAlgo)
	main.Func("limit", "Specify the size limit", cmd.parseSizeStatement)
	main.BoolVar(&cmd.flagQuitOnSoftError, "hard", false, "Quit on every error if true")
	main.BoolVar(&cmd.flagMultiThread, "mt", false, "Enable multithreading on supported subroutines")
	main.BoolVar(&cmd.flagPrint0, "print0", false, "Print processed file paths null-terminated")
	if err := main.Parse(os.Args[1:]); err != nil {
		return nil, err
	}
	cmd.flagLogLevel = logger.Level(logLevel)
	//Stage 2: Parse command
	p := &parser{
		tokens:      main.Args(),
		commandLine: cmd,
	}
	if err := p.parseCommand(); err != nil {
		return nil, err
	}
	return cmd, nil
}

// Returns an error of b don't holds the same data as a.
// This is a debug function used by unit tests.
func (a *CommandLine) mustEqual(b *CommandLine) error {
	differs := func(field string, a, b any) error {
		return fmt.Errorf("Field %s differs: A: %v, B: %v", field, a, b)
	}
	if a.command != b.command {
		return differs("command", a.command, b.command)
	}
	if slices.Compare(a.paths, b.paths) != 0 {
		return differs("paths", a.paths, b.paths)
	}
	if slices.Compare(a.names, b.names) != 0 {
		return differs("names", a.names, b.names)
	}
	if a.flagLogLevel != b.flagLogLevel {
		return differs("flagLogLevel", a.flagLogLevel, b.flagLogLevel)
	}
	if a.flagFollowSymlinks != b.flagFollowSymlinks {
		return differs("flagFollowSymlinks", a.flagFollowSymlinks, b.flagFollowSymlinks)
	}
	if a.flagHash != b.flagHash {
		return differs("flagHash", a.flagHash, b.flagHash)
	}
	if a.flagQuitOnSoftError != b.flagQuitOnSoftError {
		return differs("flagQuitOnSoftError", a.flagQuitOnSoftError, b.flagQuitOnSoftError)
	}
	if a.flagMultiThread != b.flagMultiThread {
		return differs("flagMultiThread", a.flagMultiThread, b.flagMultiThread)
	}
	if a.flagPrint0 != b.flagPrint0 {
		return differs("flagPrint0", a.flagPrint0, b.flagPrint0)
	}
	if a.tagConstraint != b.tagConstraint {
		return differs("tagConstraint", a.tagConstraint, b.tagConstraint)
	}
	if a.untagConstraint != b.untagConstraint {
		return differs("untagConstraint", a.untagConstraint, b.untagConstraint)
	}
	if a.flagLogLevel != b.flagLogLevel {
		return differs("printConstraint", a.printConstraint, b.printConstraint)
	}
	return nil
}
