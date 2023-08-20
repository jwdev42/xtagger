package cli

import (
	"fmt"
)

const (
	CommandInvalid Command = ""
	CommandPrint           = "print"
	CommandTag             = "tag"
	CommandUntag           = "untag"
)

type Command string

func parseCommand(input string) (Command, error) {
	switch cmd := Command(input); cmd {
	case CommandPrint, CommandTag, CommandUntag:
		return cmd, nil
	}
	return CommandInvalid, fmt.Errorf("Unknown command: %s", input)
}
