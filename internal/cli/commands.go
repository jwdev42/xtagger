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
	switch c := Command(input); c {
	case CommandPrint, CommandTag, CommandUntag:
		return c, nil
	}
	return CommandInvalid, fmt.Errorf("Unknown command \"%s\"", input)
}
