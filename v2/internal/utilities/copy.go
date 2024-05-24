package utilities

import (
	"bytes"
	"io"
)

func CopyAndClose(done chan struct{}, buf *bytes.Buffer, r io.ReadCloser) {
	defer close(done)
	defer r.Close()
	io.Copy(buf, r)
}
