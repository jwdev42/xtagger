package cli

import (
	"flag"
	"fmt"
	"github.com/jwdev42/logger"
	"github.com/jwdev42/xtagger/internal/hashes"
	"os"
	"slices"
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
	forbidRecursion     bool
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

func (r *CommandLine) ForbidRecursion() bool {
	return r.forbidRecursion
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

// Parses and validates command line arguments.
func ParseCommandLine() (*CommandLine, error) {
	//Stage 1: Parse flags
	var cmd = new(CommandLine)
	cmd.flagHash = hashes.SHA256 //Default hash algorithm
	var logLevel = logger.LevelFlag(logger.LevelError)
	main := flag.NewFlagSet("main", flag.ContinueOnError)
	main.Var(&logLevel, "ll", "Set the loglevel")
	main.BoolVar(&cmd.flagFollowSymlinks, "symlinks", false, "Program follows symlinks if true")
	main.Func("hash", "Specify the hashing algorithm", func(input string) error {
		hash, err := hashes.ParseAlgo(input)
		if err != nil {
			return err
		}
		cmd.flagHash = hash
		return nil
	})
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
