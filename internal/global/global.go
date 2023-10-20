package global

import (
	"fmt"
	"github.com/jwdev42/logger"
)

const (
	ExitSuccess ProgramExitCode = iota
	ExitHardError
	ExitSoftError
)

const BufSize = 1048576 //Default buffer size is 1 MiB

type ProgramExitCode int

var DefaultLogger *logger.Logger //The program's default logger
var ExitCode ProgramExitCode = ExitSuccess
var stopOnSoftError bool

func StopOnSoftError() {
	stopOnSoftError = true
}

func FilterSoftError(err error) error {
	if stopOnSoftError {
		return err
	}
	//Consume soft error
	if err != nil {
		DefaultLogger.Error(err)
		ExitCode = ExitSoftError
	}
	return nil
}

func SoftErrorf(format string, a ...any) error {
	if !stopOnSoftError {
		//Consume soft error
		DefaultLogger.Errorf(format, a...)
		ExitCode = ExitSoftError
		return nil
	}
	return fmt.Errorf(format, a...)
}
