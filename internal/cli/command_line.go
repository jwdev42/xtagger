package cli

import (
	"errors"
	"flag"
	"fmt"
	"github.com/jwdev42/xtagger/internal/hashes"
	"os"
	"path/filepath"
)

// Represents a parsed command line argument set
type CommandLine struct {
	command Command        //Specified command
	args    map[ArgKey]any //Holds the parsed command-specific arguments
}

func ParseCommandLine() (*CommandLine, error) {
	if len(os.Args) < 2 {
		return nil, errors.New("No command specified")
	}
	command, err := parseCommand(os.Args[1])
	if err != nil {
		return nil, fmt.Errorf("Error parsing command: %s", err)
	}
	var commandArgs []string
	if len(os.Args) > 2 {
		commandArgs = os.Args[2:]
	}
	args, err := parseCommandArgs(command, commandArgs)
	if err != nil {
		return nil, fmt.Errorf("Error parsing arguments for command \"%s\": %s", command, err)
	}
	return &CommandLine{
		command: command,
		args:    args,
	}, nil
}

func (r *CommandLine) Command() Command {
	return r.command
}

func (r *CommandLine) Arg(key ArgKey) (any, bool) {
	val, ok := r.args[key]
	return val, ok
}

// Checks for problems not detected by the parser such as missing
// mandatory arguments. Returns an error if it finds a problem.
func (r *CommandLine) Check() error {
	if err := r.checkMandatoryName(); err != nil {
		return err
	}
	return nil
}

func (r *CommandLine) checkMandatoryName() error {
	switch r.command {
	case CommandTag, CommandUntag:
		val, ok := r.Arg(ArgKeyName)
		if !ok || val.(string) == "" {
			return fmt.Errorf("Flag \"-name\" is mandatory for command \"%s\"", r.command)
		}
	}
	return nil
}

func (r *CommandLine) Paths() []string {
	return r.args[ArgKeyInput].([]string)
}

func (r *CommandLine) FlagName() string {
	return r.args[ArgKeyName].(string)
}

func (r *CommandLine) FlagHash() hashes.Algo {
	return r.args[ArgKeyHashAlgo].(hashes.Algo)
}

func (r *CommandLine) FlagRecursive() bool {
	return r.args[ArgKeyRecursive].(bool)
}

func (r *CommandLine) FlagFollowSymlinks() bool {
	return r.args[ArgKeyFollowSymlinks].(bool)
}

func (r *CommandLine) FlagOmitEmpty() bool {
	return r.args[ArgKeyOmitEmpty].(bool)
}

func parseCommandArgs(command Command, args []string) (map[ArgKey]any, error) {
	switch command {
	case CommandPrint, CommandTag, CommandUntag:
		return parseArgs(args)
	case CommandInvalid:
		panic("BUG: Zero-value trap CommandInvalid triggered")
	}
	panic("BUG: You're not supposed to be here")
}

func parsePaths(args []string) ([]string, error) {
	if args == nil || len(args) < 1 {
		return nil, errors.New("No path in input")
	}
	paths := make([]string, len(args)) //stores the parsed paths
	for i, arg := range args {
		var err error
		paths[i], err = parsePath(arg)
		if err != nil {
			return nil, err
		}
	}
	return paths, nil
}

func parsePath(path string) (string, error) {
	if path == "" {
		return "", errors.New("Path cannot be empty")
	}
	parsedPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return parsedPath, nil
}

func parseArgs(args []string) (map[ArgKey]any, error) {
	parsedArgs := make(map[ArgKey]any)
	flagSet := flag.NewFlagSet("common", flag.ContinueOnError)
	flagName := new(name)
	flagSet.Var(flagName, "name", "The xtag's name")
	flagHash := new(hash)
	flagSet.Var(flagName, "hash", "The hash algorithm to be used")
	flagRecursive := flagSet.Bool("R", false, "Recurse into subdirectories if true")
	flagFollowSymlinks := flagSet.Bool("L", false, "Follows symbolic links if true")
	flagOmitEmpty := flagSet.Bool("omitempty", false, "Skips empty entries if true")
	flagBackup := flagSet.String("backup", "", "Takes a file path as argument, activates backup mode if set")
	if err := flagSet.Parse(args); err != nil {
		return nil, err
	}
	parsedArgs[ArgKeyName] = flagName.Get()
	parsedArgs[ArgKeyHashAlgo] = flagHash.Get()
	parsedArgs[ArgKeyRecursive] = *flagRecursive
	parsedArgs[ArgKeyFollowSymlinks] = *flagFollowSymlinks
	parsedArgs[ArgKeyOmitEmpty] = *flagOmitEmpty
	parsedArgs[ArgKeyBackup] = *flagBackup
	if paths, err := parsePaths(flagSet.Args()); err != nil {
		return nil, err
	} else {
		parsedArgs[ArgKeyInput] = paths
	}
	return parsedArgs, nil
}
