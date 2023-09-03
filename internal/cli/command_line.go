package cli

import (
	"errors"
	"fmt"
	"github.com/integrii/flaggy"
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
)

// Represents a parsed command line argument set.
type CommandLine struct {
	command              Command //Specified command
	paths                []string
	flagRecursive        bool
	flagFollowSymlinks   bool
	flagName             string
	flagHash             hashes.Algo
	flagBackupTargetPath string
	flagOmitEmpty        bool
}

// Parses and validates command line arguments.
func ParseCommandLine() (*CommandLine, error) {
	cl := new(CommandLine)
	parser := flaggy.NewParser("xtagger")
	var flagHash string
	//Command tag
	tag := flaggy.NewSubcommand(string(CommandTag))
	tag.String(&cl.flagName, shortName, longName, "Name for the new record")
	tag.String(&flagHash, shortHash, longHash, "Hashing algorithm")
	tag.Bool(&cl.flagRecursive, shortRecursive, longRecursive, "Recurse into subdirectories")
	tag.Bool(&cl.flagFollowSymlinks, shortFollowSymlinks, longFollowSymlinks, "Follow symlinks")
	tag.String(&cl.flagBackupTargetPath, "b", "backup", "Backup target path")
	tag.StringSlice(&cl.paths, shortPath, longPath, "Source path, can be specified multiple times")
	parser.AttachSubcommand(tag, 1)
	//Command untag
	untag := flaggy.NewSubcommand(string(CommandUntag))
	untag.Bool(&cl.flagRecursive, shortRecursive, longRecursive, "Recurse into subdirectories")
	untag.Bool(&cl.flagFollowSymlinks, shortFollowSymlinks, longFollowSymlinks, "Follow symlinks")
	untag.String(&cl.flagName, shortName, longName, "Name of the record to be deleted")
	untag.StringSlice(&cl.paths, shortPath, longPath, "Source path, can be specified multiple times")
	parser.AttachSubcommand(untag, 1)
	//Command print
	print := flaggy.NewSubcommand(string(CommandPrint))
	print.Bool(&cl.flagRecursive, shortRecursive, longRecursive, "Recurse into subdirectories")
	print.Bool(&cl.flagFollowSymlinks, shortFollowSymlinks, longFollowSymlinks, "Follow symlinks")
	print.String(&cl.flagName, shortName, longName, "Only print records matching name")
	print.StringSlice(&cl.paths, shortPath, longPath, "Source path, can be specified multiple times")
	parser.AttachSubcommand(print, 1)
	//Parse
	if err := parser.Parse(); err != nil {
		return nil, err
	}
	//Parse command name
	cl.command = Command(parser.TrailingSubcommand().Name)
	//Parse specific oddities
	switch cl.command {
	case CommandTag:
		algo, err := hashes.ParseAlgo(flagHash)
		if err != nil {
			return nil, err
		}
		cl.flagHash = algo
	}
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

func (r *CommandLine) FlagName() string {
	return r.flagName
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
		if err := validateName(r.flagName); err != nil {
			return fmt.Errorf("Flag %q: %s", longName, err)
		}
	}
	return nil
}
