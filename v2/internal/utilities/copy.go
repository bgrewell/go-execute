package utilities

import (
	"bytes"
	"io"
	"sync"
)

// CopyAndClose copies the contents of the provided io.ReadCloser to the provided bytes.Buffer and closes the io.ReadCloser.
func CopyAndClose(done chan struct{}, buf *bytes.Buffer, r io.ReadCloser, ready *sync.WaitGroup) {
	defer close(done)
	defer r.Close()
	ready.Done()
	io.Copy(buf, r)
}
