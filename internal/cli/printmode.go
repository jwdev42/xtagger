package cli

import (
	"fmt"
)

const (
	PrintModeInvalid   = ""
	PrintModeAll       = "all"
	PrintModeValidOnly = "valid"
	PrintModeOmitEmpty = "nonempty"
)

type PrintMode string

func ParsePrintMode(input string) (PrintMode, error) {
	switch mode := PrintMode(input); mode {
	case PrintModeAll, PrintModeValidOnly, PrintModeOmitEmpty:
		return mode, nil
	}
	return PrintModeInvalid, fmt.Errorf("Unknown print mode: %s", input)
}
