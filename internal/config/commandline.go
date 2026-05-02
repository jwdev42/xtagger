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

package config

import (
	"flag"
	"fmt"
	"github.com/jwdev42/xtagger/internal/hashes"
	"os"
	"strconv"
)

// ParseCommandLine parses and validates command line arguments.
// ParseCommandLine fills argument prefs with the parsed values.
func ParseCommandLine(prefs *Preferences) error {
	// Stage 1: Parse flags
	var hashAlgo hashes.Algo
	var logLevel = &flagLogLevel{}
	var quotaSize int64
	var threads int
	var followSymlinks, usePrint0, useRecursion bool
	main := flag.NewFlagSet("main", flag.ContinueOnError)
	main.Var(logLevel, "loglevel", "Set the loglevel")
	main.Func("hash", "Specify the hashing algorithm", parseHashAlgoFunc(&hashAlgo))
	main.Func("quota", "Specify a quota", parseSizeFunc(&quotaSize))
	main.IntVar(&threads, "threads", prefs.Threads, "Number of threads, set this to 1 on HDDs")
	main.BoolVar(&followSymlinks, "symlinks", prefs.FollowSymlinks, "Program follows symlinks if true")
	main.BoolVar(&usePrint0, "print0", prefs.UsePrint0, "Print processed file paths null-terminated")
	main.BoolVar(&useRecursion, "recursive", prefs.UseRecursion, "Recurse into subdirectories if true")
	if err := main.Parse(os.Args[1:]); err != nil {
		return err
	}
	// Stage 2: Parse command
	commandArgs, err := newCommandParser(main.Args()).start()
	if err != nil {
		return err
	}
	// Stage 3: Apply parsed vars to prefs
	prefs.Command = commandArgs.command
	prefs.Paths = commandArgs.paths
	prefs.Names = commandArgs.names
	if hashAlgo != hashes.INVALID {
		prefs.UseHash = hashAlgo
	}
	if logLevel.Used() {
		prefs.LogLevel.Set(logLevel.Level())
	}
	if quotaSize > 0 {
		prefs.Quota = quotaSize
	}
	prefs.Threads = threads
	prefs.FollowSymlinks = followSymlinks
	prefs.UsePrint0 = usePrint0
	prefs.UseRecursion = useRecursion
	prefs.TagConstraint = commandArgs.tagConstraint
	prefs.PrintConstraint = commandArgs.printConstraint
	prefs.UntagConstraint = commandArgs.untagConstraint
	return nil
}

func parseHashAlgoFunc(storage *hashes.Algo) func(string) error {
	return func(input string) error {
		hash, err := hashes.ParseAlgo(input)
		if err != nil {
			return err
		}
		*storage = hash
		return nil
	}
}

func parseSizeFunc(size *int64) func(string) error {
	return func(input string) error {
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
		const (
			kib = 1024
			mib = kib * kib
			gib = mib * kib
			tib = gib * kib
		)
		switch suffix {
		case "":
		case "K":
			sizeLimit *= kib
		case "M":
			sizeLimit *= mib
		case "G":
			sizeLimit *= gib
		case "T":
			sizeLimit *= tib
		default:
			return fmt.Errorf("Could not parse size statement: Unknown suffix: \"%s\"", suffix)
		}
		*size = sizeLimit
		return nil
	}
}
