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

import (
	"bufio"
	"fmt"
	"io"
)

type Printer struct {
	wr       *bufio.Writer
	messages chan string
}

func NewPrinter(w io.Writer, eh *ErrorHandler, bufsize int, separator string) (printer *Printer, closeFunc func()) {
	printer = &Printer{
		wr:       bufio.NewWriter(w),
		messages: make(chan string, bufsize),
	}
	done := make(chan struct{})

	go func() {
		defer close(done)
		for message := range printer.messages {
			_, err := printer.wr.WriteString(message)
			eh.Error(err)
			_, err = printer.wr.WriteString(separator)
			eh.Error(err)
		}
		eh.Error(printer.wr.Flush())
	}()

	closeFunc = func() {
		close(printer.messages)
		<-done
	}

	return
}

func (r *Printer) Printf(format string, a ...any) {
	r.messages <- fmt.Sprintf(format, a...)
}

func (r *Printer) Print(a ...any) {
	r.messages <- fmt.Sprint(a...)
}
