package program

import (
	"fmt"
	"github.com/jwdev42/xtagger/internal/cli"
)

func Run() error {
	//Parse command line
	cmdline, err := cli.ParseCommandLine()
	if err != nil {
		return fmt.Errorf("Command line error: %s", err)
	}
	switch command := cmdline.Command(); command {
	case cli.CommandCreate:
		return runCreate(cmdline)
	default:
		return fmt.Errorf("Unknown command \"%s\"", command)
	}
	return nil
}
