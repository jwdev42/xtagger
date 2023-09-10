package cli

import (
	"errors"
	"fmt"
	"github.com/integrii/flaggy"
	"github.com/jwdev42/logger"
	"github.com/jwdev42/xtagger/internal/hashes"
)

const (
	shortName           = "n"
	longName            = "name"
	shortHash           = "H"
	longHash            = "hash"
	shortRecursive      = "r"
	longRecursive       = "recursive"
	shortFollowSymlinks = "L"
	longFollowSymlinks  = "follow-symlinks"
	shortPath           = "p"
	longPath            = "path"
	shortLogLevel       = "ll"
	longLogLevel        = "loglevel"
)

// Represents a parsed command line argument set.
type CommandLine struct {
	command              Command //Specified command
	paths                []string
	flagLogLevel         logger.Level
	flagRecursive        bool
	flagFollowSymlinks   bool
	flagNames            []string
	flagHash             hashes.Algo
	flagBackupTargetPath string
	flagOmitEmpty        bool
}

// Parses and validates command line arguments.
func ParseCommandLine() (*CommandLine, error) {
	cl := new(CommandLine)
	parser := flaggy.NewParser("xtagger")
	//Intermediates for flags with custom types
	var flagHash string
	flagLogLevel := "error"
	//Command tag
	tag := flaggy.NewSubcommand(string(CommandTag))
	tag.StringSlice(&cl.flagNames, shortName, longName, "Name for the new record")
	tag.String(&flagHash, shortHash, longHash, "Hashing algorithm")
	tag.Bool(&cl.flagRecursive, shortRecursive, longRecursive, "Recurse into subdirectories")
	tag.Bool(&cl.flagFollowSymlinks, shortFollowSymlinks, longFollowSymlinks, "Follow symlinks")
	tag.String(&cl.flagBackupTargetPath, "b", "backup", "Backup target path")
	tag.StringSlice(&cl.paths, shortPath, longPath, "Source path, can be specified multiple times")
	tag.String(&flagLogLevel, shortLogLevel, longLogLevel, "Desired log level, default is Error")
	parser.AttachSubcommand(tag, 1)
	//Command untag
	untag := flaggy.NewSubcommand(string(CommandUntag))
	untag.Bool(&cl.flagRecursive, shortRecursive, longRecursive, "Recurse into subdirectories")
	untag.Bool(&cl.flagFollowSymlinks, shortFollowSymlinks, longFollowSymlinks, "Follow symlinks")
	untag.StringSlice(&cl.flagNames, shortName, longName, "Name of the record to be deleted")
	untag.StringSlice(&cl.paths, shortPath, longPath, "Source path, can be specified multiple times")
	untag.String(&flagLogLevel, shortLogLevel, longLogLevel, "Desired log level, default is Error")
	parser.AttachSubcommand(untag, 1)
	//Command print
	print := flaggy.NewSubcommand(string(CommandPrint))
	print.Bool(&cl.flagRecursive, shortRecursive, longRecursive, "Recurse into subdirectories")
	print.Bool(&cl.flagFollowSymlinks, shortFollowSymlinks, longFollowSymlinks, "Follow symlinks")
	print.StringSlice(&cl.flagNames, shortName, longName, "Only print records matching name")
	print.StringSlice(&cl.paths, shortPath, longPath, "Source path, can be specified multiple times")
	print.String(&flagLogLevel, shortLogLevel, longLogLevel, "Desired log level, default is Error")
	parser.AttachSubcommand(print, 1)
	//Parse
	if err := parser.Parse(); err != nil {
		return nil, err
	}
	//Parse command name
	cl.command = Command(parser.TrailingSubcommand().Name)
	//Process custom types that flaggy doesn't support directly
	cl.flagHash = hashes.Algo(flagHash)
	logLevel, err := logger.ParseLevel(flagLogLevel)
	if err != nil {
		return nil, err
	}
	cl.flagLogLevel = logLevel
	//Validate CommandLine
	if err := cl.validate(); err != nil {
		return nil, err
	}
	return cl, nil
}

func (r *CommandLine) Command() Command {
	return r.command
}

func (r *CommandLine) Paths() []string {
	return r.paths
}

func (r *CommandLine) FlagLogLevel() logger.Level {
	return r.flagLogLevel
}

func (r *CommandLine) FlagNames() []string {
	return r.flagNames
}

func (r *CommandLine) FlagHash() hashes.Algo {
	return r.flagHash
}

func (r *CommandLine) FlagRecursive() bool {
	return r.flagRecursive
}

func (r *CommandLine) FlagFollowSymlinks() bool {
	return r.flagFollowSymlinks
}

func (r *CommandLine) FlagOmitEmpty() bool {
	return r.flagOmitEmpty
}

// Checks if all mandatory command line arguments are set dependent on the command.
func (r *CommandLine) validate() error {
	if r.paths == nil || len(r.paths) < 1 {
		return errors.New("No path specified")
	}
	switch r.command {
	case CommandInvalid:
		return errors.New("No command specified")
	case CommandTag:
		//Check if flag longName is present exactly once
		if names := r.flagNames; names == nil {
			return fmt.Errorf("Command %q: Flag %q is mandatory", r.command, longName)
		} else if len(names) != 1 {
			return fmt.Errorf("Command %q: Flag %q can be set just once", r.command, longName)
		}
		if err := validateName(r.flagNames[0]); err != nil {
			return fmt.Errorf("Flag %q: %s", longName, err)
		}
		//Check if hashing algorithm is valid
		algo, err := hashes.ParseAlgo(string(r.flagHash))
		if err != nil {
			return fmt.Errorf("Flag %q: %s", longHash, err)
		}
		r.flagHash = algo
	}
	return nil
}
