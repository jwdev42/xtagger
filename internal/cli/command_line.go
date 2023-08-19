package cli

import (
	"errors"
	"flag"
	"fmt"
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

func (r *CommandLine) Paths() []string {
	return r.args[ArgKeyInput].([]string)
}

func (r *CommandLine) FlagName() string {
	return r.args[ArgKeyName].(string)
}

func (r *CommandLine) FlagRecursive() bool {
	return r.args[ArgKeyRecursive].(bool)
}

func (r *CommandLine) FlagFollowSymlinks() bool {
	return r.args[ArgKeyFollowSymlinks].(bool)
}

func parseCommandArgs(command Command, args []string) (map[ArgKey]any, error) {
	switch command {
	case CommandShow, CommandTag, CommandUntag:
		return parseArgsCommon(args)
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

func parseArgsCommon(args []string) (map[ArgKey]any, error) {
	parsedArgs := make(map[ArgKey]any)
	flagSet := flag.NewFlagSet("common", flag.ContinueOnError)
	flagName := new(name)
	flagSet.Var(flagName, "name", "The xtag's name")
	flagRecursive := flagSet.Bool("recursive", false, "Recurse into subdirectories if true")
	flagFollowSymlinks := flagSet.Bool("follow-symlinks", false, "Follows symbolic links if true")
	flagBackup := flagSet.String("backup", "", "Takes a file path as argument, activates backup mode if set")
	if err := flagSet.Parse(args); err != nil {
		return nil, err
	}
	parsedArgs[ArgKeyName] = flagName.Get()
	parsedArgs[ArgKeyRecursive] = *flagRecursive
	parsedArgs[ArgKeyFollowSymlinks] = *flagFollowSymlinks
	parsedArgs[ArgKeyBackup] = *flagBackup
	if paths, err := parsePaths(flagSet.Args()); err != nil {
		return nil, err
	} else {
		parsedArgs[ArgKeyInput] = paths
	}
	return parsedArgs, nil
}
