//This file is part of xtagger. ©2023-2026 Jörg Walter.
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

package logging

import "fmt"

// NonFatalErrors is intended to report the error count upon program
// termination. Only errors that didn't lead to premature program
// termination should be reported.
type NonFatalErrors int64

// ReportNonFatalErrors returns a new error of type NonFatalErrors if
// count is > 0. It returns nil otherwise.
func ReportNonFatalErrors(count int64) error {
	if count > 0 {
		return NonFatalErrors(count)
	}
	return nil
}

func (pe NonFatalErrors) Error() string {
	switch pe {
	case 1:
		return fmt.Sprintf("Program finished with %d error", pe)
	}
	return fmt.Sprintf("Program finished with %d errors", pe)
}
