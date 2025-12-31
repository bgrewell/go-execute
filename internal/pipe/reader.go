// Package pipe provides utilities for handling command output pipes
// with proper synchronization to avoid timing issues.
package pipe

import (
	"bytes"
	"io"
	"sync"
)

// Reader wraps an io.ReadCloser and buffers its output to handle
// the timing issues that can occur with command execution pipes.
//
// The problem: When executing commands asynchronously, the stdout/stderr
// pipes from exec.Cmd can be read too early (before data arrives) or
// the underlying reader can be closed before all data is consumed.
// Using bytes.Buffer directly doesn't work because it returns EOF
// immediately when empty, even if more data is coming.
//
// This implementation:
// - Spawns a goroutine to continuously read from the underlying pipe
// - Buffers all data internally
// - Blocks Read() calls until data is available or the pipe closes
// - Provides Wait() to block until the underlying pipe is fully consumed
type Reader struct {
	mu           sync.Mutex
	cond         *sync.Cond
	buffer       bytes.Buffer
	reader       io.ReadCloser
	closed       bool
	readerClosed bool
	readPending  bool
}

// NewReader creates a new Reader that wraps the given io.ReadCloser.
// It immediately starts a background goroutine to read from the pipe.
func NewReader(r io.ReadCloser) *Reader {
	pr := &Reader{
		reader: r,
	}
	pr.cond = sync.NewCond(&pr.mu)
	go pr.readLoop()
	return pr
}

// Read implements io.Reader. It blocks until data is available
// or the underlying pipe is closed.
func (pr *Reader) Read(p []byte) (n int, err error) {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	// Wait for data or pipe closure
	for pr.buffer.Len() == 0 && !pr.readerClosed {
		pr.readPending = true
		pr.cond.Wait()
		pr.readPending = false
	}

	// If buffer is empty and pipe is closed, we're done
	if pr.buffer.Len() == 0 && pr.readerClosed {
		return 0, io.EOF
	}

	return pr.buffer.Read(p)
}

// Close marks the Reader as closed. This does not close the underlying
// pipe - that happens naturally when the command exits.
func (pr *Reader) Close() error {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	pr.closed = true
	pr.cond.Broadcast()
	return nil
}

// Wait blocks until the underlying pipe has been fully read and closed.
// This is useful for coordinating with command completion.
func (pr *Reader) Wait() {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	for !pr.readerClosed {
		pr.cond.Wait()
	}
}

// readLoop continuously reads from the underlying pipe into the buffer.
// It runs in a background goroutine.
func (pr *Reader) readLoop() {
	buf := make([]byte, 4096)
	for {
		n, err := pr.reader.Read(buf)
		if n > 0 {
			pr.mu.Lock()
			pr.buffer.Write(buf[:n])
			if pr.readPending {
				pr.cond.Signal()
			}
			pr.mu.Unlock()
		}
		if err != nil {
			pr.mu.Lock()
			pr.readerClosed = true
			pr.cond.Broadcast()
			pr.mu.Unlock()
			return
		}
	}
}

// Bytes returns all buffered data without consuming it.
// This should only be called after Wait() returns.
func (pr *Reader) Bytes() []byte {
	pr.mu.Lock()
	defer pr.mu.Unlock()
	return pr.buffer.Bytes()
}

// String returns all buffered data as a string without consuming it.
// This should only be called after Wait() returns.
func (pr *Reader) String() string {
	pr.mu.Lock()
	defer pr.mu.Unlock()
	return pr.buffer.String()
}
