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
	"context"
	"sync/atomic"
)

// ErrorHandler manages asynchronous error processing by queueing errors
// and passing them to a designated callback function.
type ErrorHandler struct {
	errChan  chan error
	counter  *atomic.Int64
	callback func(context.Context, error)
}

// NewErrorHandler initializes a new ErrorHandler and starts a background goroutine
// to process incoming errors. It returns the handler instance and a closure function
// that should be called to gracefully shut down the error processing.
//
// The bufsize parameter determines the capacity of the underlying error channel.
func NewErrorHandler(ctx context.Context, bufsize int, callback func(context.Context, error)) (eh *ErrorHandler, closeFunc func()) {
	errChan := make(chan error, bufsize)
	doneChan := make(chan struct{}) // Closed when consumer is done

	eh = &ErrorHandler{
		errChan:  errChan,
		counter:  &atomic.Int64{},
		callback: callback,
	}

	// Start consumer
	go func() {
		eh.listen(ctx)
		close(doneChan)
	}()

	// Define canceler
	closeFunc = func() {
		close(errChan)
		<-doneChan
	}

	return
}

// listen continuously reads from the error channel and triggers the callback
// for each received error. It terminates when the error channel is closed.
func (eh *ErrorHandler) listen(ctx context.Context) {
	for err := range eh.errChan {
		eh.callback(ctx, err)
	}
}

// Error enqueues an error to be processed by the ErrorHandler.
// If the provided error is nil, it performs no action and returns false.
// If the provided error is not nil, it returns true.
func (eh *ErrorHandler) Error(err error) bool {
	if err == nil {
		return false
	}
	eh.errChan <- err
	eh.counter.Add(1)
	return true
}

// Errors returns the amount of non-nil errors processed by the handler.
func (eh *ErrorHandler) Errors() int64 {
	return eh.counter.Load()
}
