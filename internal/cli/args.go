package cli

import(
	"errors"
	"flag"
	"fmt"
)

//Parses command line arguments and returns the result as map.
//Returnes an error if one occured, then the map wil be nil.
func ParseArgs(args []string) (map[string]any, error) {
	return parseCommand(args)
}

func parseCommand(args []string) (map[string]any, error) {
	if len(args) < 1 {
		return errors.New("No command specified")
	}
	cmd := args[0]
	switch cmd {
		case "show":
			return parseShow(args[1:])
		case "copy":
			return parseCopy(args[1:])
		case "remove":
			return parseRemove(args[1:])
		case "todo":
			return parseTodo(args[1:])
	}
	return nil, fmt.Errorf("\"%s\" is not a valid command", cmd)
}

func parseShow(args []string) (map[string]any, error) {
	//flags
	name := &name("")
	fs := flag.NewFlagSet("show", flag.ContinueOnError)
	fs.Var(name, "name", "specifies the backup identifier to be used")
	fs.Parse(args)
}
