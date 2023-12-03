package softerrors

import (
	"fmt"
	"github.com/jwdev42/logger"
	"github.com/jwdev42/xtagger/internal/global"
)

var stopOnSoftError bool

func StopOnSoftError() {
	stopOnSoftError = true
}

func Consume(err error) error {
	if stopOnSoftError {
		return err
	}
	//Consume soft error
	if err != nil {
		logger.Default().Error(err)
		global.ExitCode = global.ExitSoftError
	}
	return nil
}

func Errorf(format string, a ...any) error {
	if !stopOnSoftError {
		//Consume soft error
		logger.Default().Errorf(format, a...)
		global.ExitCode = global.ExitSoftError
		return nil
	}
	return fmt.Errorf(format, a...)
}
