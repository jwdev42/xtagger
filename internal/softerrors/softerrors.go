//This file is part of xtagger. ©2023 Jörg Walter.
//This program is free software: you can redistribute it and/or modify
//it under the terms of the GNU General Public License as published by
//the Free Software Foundation, either version 3 of the License, or
//(at your option) any later version.
//
//This program is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//GNU General Public License for more details.
//
//You should have received a copy of the GNU General Public License
//along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Softerrors is a helper package that enables the programmer to log
// specific errors instead of exiting the program.
package softerrors

import (
	"fmt"
	"github.com/jwdev42/logger"
	"github.com/jwdev42/xtagger/internal/global"
)

var stopOnSoftError bool

// Disables soft errors.
func StopOnSoftError() {
	stopOnSoftError = true
}

// Logs err as soft error if soft errors are enabled, then returns nil.
// Returns err if soft errors are disabled.
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

// Creates a new formatted soft error, logs it, then returns nil.
// Returns a new formatted error if soft errors are disabled.
func Errorf(format string, a ...any) error {
	if !stopOnSoftError {
		//Consume soft error
		logger.Default().Errorf(format, a...)
		global.ExitCode = global.ExitSoftError
		return nil
	}
	return fmt.Errorf(format, a...)
}
