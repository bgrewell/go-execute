package internal

import (
	"bytes"
	"errors"
	"io"
	"sync"
)

type ExecReadWriter struct {
	mu           sync.Mutex
	cond         *sync.Cond
	buffer       bytes.Buffer
	reader       io.ReadCloser
	closed       bool
	readerClosed bool
	readPending  bool
	waiters      int
}

// NewExecReadWriter initializes a new ExecReadWriter with the provided io.ReadCloser.
func NewExecReadWriter(reader io.ReadCloser) *ExecReadWriter {
	rwc := &ExecReadWriter{
		reader: reader,
	}
	rwc.cond = sync.NewCond(&rwc.mu)
	go rwc.readFromReader()
	return rwc
}

// Read blocks until data is available or the io.ReadCloser is closed.
func (rwc *ExecReadWriter) Read(p []byte) (n int, err error) {
	rwc.mu.Lock()
	defer rwc.mu.Unlock()

	for rwc.buffer.Len() == 0 && !rwc.readerClosed {
		rwc.readPending = true
		rwc.cond.Wait()
		rwc.readPending = false
	}

	if rwc.buffer.Len() == 0 && rwc.readerClosed {
		return 0, io.EOF
	}

	return rwc.buffer.Read(p)
}

// Write writes data to the internal buffer.
func (rwc *ExecReadWriter) Write(p []byte) (n int, err error) {
	rwc.mu.Lock()
	defer rwc.mu.Unlock()

	if rwc.closed {
		return 0, errors.New("write to closed writer")
	}

	n, err = rwc.buffer.Write(p)
	if rwc.readPending {
		rwc.cond.Signal()
	}

	return n, err
}

// Close closes the ExecReadWriter.
func (rwc *ExecReadWriter) Close() error {
	rwc.mu.Lock()
	defer rwc.mu.Unlock()

	rwc.closed = true
	rwc.cond.Broadcast()
	return nil
}

// Wait waits until the io.ReadCloser has closed.
func (rwc *ExecReadWriter) Wait() {
	rwc.mu.Lock()
	defer rwc.mu.Unlock()

	for !rwc.readerClosed {
		rwc.waiters++
		rwc.cond.Wait()
		rwc.waiters--
	}
}

// readFromReader reads from the io.ReadCloser and writes to the internal buffer.
func (rwc *ExecReadWriter) readFromReader() {
	buf := make([]byte, 4096)
	for {
		n, err := rwc.reader.Read(buf)
		if n > 0 {
			rwc.mu.Lock()
			rwc.buffer.Write(buf[:n])
			if rwc.readPending {
				rwc.cond.Signal()
			}
			rwc.mu.Unlock()
		}
		if err != nil {
			rwc.mu.Lock()
			rwc.readerClosed = true
			rwc.cond.Broadcast()
			rwc.mu.Unlock()
			break
		}
	}
}
