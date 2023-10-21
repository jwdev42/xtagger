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
