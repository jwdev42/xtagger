package cli

import (
	"fmt"
)

const (
	CommandInvalid Command = ""
	CommandShow            = "show"
	CommandCreate          = "create"
	CommandRemove          = "remove"
)

type Command string

func parseCommand(input string) (Command, error) {
	switch cmd := Command(input); cmd {
	case CommandShow, CommandCreate, CommandRemove:
		return cmd, nil
	}
	return CommandInvalid, fmt.Errorf("Unknown command: %s", input)
}
