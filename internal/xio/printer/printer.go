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

package printer

import (
	"fmt"
	"io"
	"sync"
)

type Printer struct {
	mu *sync.Mutex
	wr io.Writer
}

func NewPrinter(wr io.Writer) *Printer {
	return &Printer{
		mu: new(sync.Mutex),
		wr: wr,
	}
}

func (r *Printer) Print0(message string) (n int, err error) {
	defer r.mu.Unlock()
	r.mu.Lock()
	return fmt.Fprintf(r.wr, "%s%c", message, 0)
}
